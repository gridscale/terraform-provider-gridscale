package gridscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type gridscaleS3Provider struct {
	AccessKey, SecretKey string
}

func (m *gridscaleS3Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     m.AccessKey,
		SecretAccessKey: m.SecretKey,
	}, nil
}

func resourceGridscaleBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleBucketCreate,
		Read:   resourceGridscaleBucketRead,
		Update: resourceGridscaleBucketUpdate,
		Delete: resourceGridscaleBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Description: "The object storage secret_key.",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The object storage access_key.",
				Required:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Description: "The name of the bucket.",
				Required:    true,
				ForceNew:    true,
			},
			"s3_host": {
				Type:        schema.TypeString,
				Description: "The S3 host.",
				Optional:    true,
				ForceNew:    true,
				Default:     "gos3.io",
			},
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expiration_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  365,
						},
						"noncurrent_version_expiration_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  365,
						},
						"incomplete_upload_expiration_days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleBucketRead(d *schema.ResourceData, meta interface{}) error {
	s3Host := d.Get("s3_host").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	bucketName := d.Get("bucket_name").(string)

	s3Client := initS3Client(&gridscaleS3Provider{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, s3Host)

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	// Fetch lifecycle configuration
	output, err := s3Client.GetBucketLifecycleConfiguration(ctx, &s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &apiErr) && apiErr.Code == "NoSuchLifecycleConfiguration" {
			// If the error indicates no lifecycle configuration exists, set the lifecycle_rule attribute to nil
			d.Set("lifecycle_rule", nil)
		} else {
			// For any other error, return a formatted error message with context
			return fmt.Errorf("error reading lifecycle configuration for bucket %s: %w", bucketName, err)
		}
	} else {
		rules := []map[string]interface{}{}
		for _, rule := range output.Rules {
			r := map[string]interface{}{
				"id":                                 aws.ToString(rule.ID),
				"enabled":                            rule.Status == types.ExpirationStatusEnabled,
				"expiration_days":                    0,
				"noncurrent_version_expiration_days": 0,
			}
			// Check if the rule has a filter and set the prefix accordingly
			if rule.Filter != nil && rule.Filter.Prefix != nil {
				r["prefix"] = aws.ToString(rule.Filter.Prefix)
			} else {
				r["prefix"] = ""
			}
			// Check if the rule has expiration or noncurrent version expiration days set
			if rule.Expiration != nil && rule.Expiration.Days != nil {
				r["expiration_days"] = aws.ToInt32(rule.Expiration.Days)
			}
			if rule.NoncurrentVersionExpiration != nil && rule.NoncurrentVersionExpiration.NoncurrentDays != nil {
				r["noncurrent_version_expiration_days"] = aws.ToInt32(rule.NoncurrentVersionExpiration.NoncurrentDays)
			}
			// Check if the rule has incomplete upload expiration days set
			if rule.AbortIncompleteMultipartUpload != nil && rule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
				r["incomplete_upload_expiration_days"] = aws.ToInt32(rule.AbortIncompleteMultipartUpload.DaysAfterInitiation)
			}
			rules = append(rules, r)
		}
		d.Set("lifecycle_rule", rules)
	}

	return nil
}

func resourceGridscaleBucketCreate(d *schema.ResourceData, meta interface{}) error {
	s3Host := d.Get("s3_host")
	accessKey := d.Get("access_key")
	secretKey := d.Get("secret_key")

	s3HostStr := s3Host.(string)
	s3Client := initS3Client(&gridscaleS3Provider{
		AccessKey: accessKey.(string),
		SecretKey: secretKey.(string),
	}, s3HostStr)

	bucketName := d.Get("bucket_name")
	bucketNameStr := bucketName.(string)
	bucketInput := s3.CreateBucketInput{
		Bucket: &bucketNameStr,
	}

	errorPrefix := fmt.Sprintf("Create bucket %s resource at s3host %s-", bucketNameStr, s3HostStr)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	_, err := s3Client.CreateBucket(ctx, &bucketInput)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	lifecycleRules := d.Get("lifecycle_rule").([]interface{})
	if len(lifecycleRules) > 0 {
		lifecycleConfig := &types.BucketLifecycleConfiguration{
			Rules: []types.LifecycleRule{},
		}

		for _, rule := range lifecycleRules {
			r := rule.(map[string]interface{})
			lifecycleRule := types.LifecycleRule{
				ID: aws.String(r["id"].(string)),
				Filter: &types.LifecycleRuleFilter{
					Prefix: aws.String(r["prefix"].(string)),
				},
				Status: types.ExpirationStatusEnabled,
			}
			// Check if the rule is enabled
			if !r["enabled"].(bool) {
				lifecycleRule.Status = types.ExpirationStatusDisabled
			}
			// Set expiration days if provided
			if v, ok := r["expiration_days"].(int); ok && v > 0 {
				lifecycleRule.Expiration = &types.LifecycleExpiration{
					Days: aws.Int32(int32(v)),
				}
			}
			// Set noncurrent version expiration days if provided
			if v, ok := r["noncurrent_version_expiration_days"].(int); ok && v > 0 {
				lifecycleRule.NoncurrentVersionExpiration = &types.NoncurrentVersionExpiration{
					NoncurrentDays: aws.Int32(int32(v)),
				}
			}
			// Set incomplete upload expiration days if provided
			if v, ok := r["incomplete_upload_expiration_days"].(int); ok && v > 0 {
				lifecycleRule.AbortIncompleteMultipartUpload = &types.AbortIncompleteMultipartUpload{
					DaysAfterInitiation: aws.Int32(int32(v)),
				}
			}

			lifecycleConfig.Rules = append(lifecycleConfig.Rules, lifecycleRule)
		}

		_, err := s3Client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
			Bucket:                 &bucketNameStr,
			LifecycleConfiguration: lifecycleConfig,
		})
		if err != nil {
			// Delete the bucket if lifecycle configuration fails to set
			return resourceGridscaleBucketDelete(d, meta)
		}
	}

	id := fmt.Sprintf("%s/%s", s3HostStr, bucketNameStr)
	d.SetId(id)

	log.Printf("The id for the new bucket has been set to %v", id)
	return nil
}

func resourceGridscaleBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	s3Host := d.Get("s3_host").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	bucketName := d.Get("bucket_name").(string)

	s3Client := initS3Client(&gridscaleS3Provider{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}, s3Host)

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	if d.HasChange("lifecycle_rule") {
		lifecycleRules := d.Get("lifecycle_rule").([]interface{})

		if len(lifecycleRules) == 0 {
			// If no lifecycle rules are provided, clear the lifecycle configuration
			_, err := s3Client.DeleteBucketLifecycle(ctx, &s3.DeleteBucketLifecycleInput{
				Bucket: aws.String(bucketName),
			})
			if err != nil {
				return fmt.Errorf("error clearing lifecycle configuration for bucket %s using DeleteBucketLifecycle: %v", bucketName, err)
			}
			return resourceGridscaleBucketRead(d, meta)
		} else {
			lifecycleConfig := &types.BucketLifecycleConfiguration{
				Rules: []types.LifecycleRule{},
			}

			for _, rule := range lifecycleRules {
				r := rule.(map[string]interface{})
				lifecycleRule := types.LifecycleRule{
					ID: aws.String(r["id"].(string)),
					Filter: &types.LifecycleRuleFilter{
						Prefix: aws.String(r["prefix"].(string)),
					},
					Status: types.ExpirationStatusEnabled,
				}

				if !r["enabled"].(bool) {
					lifecycleRule.Status = types.ExpirationStatusDisabled
				}

				if v, ok := r["expiration_days"].(int); ok && v > 0 {
					lifecycleRule.Expiration = &types.LifecycleExpiration{
						Days: aws.Int32(int32(v)),
					}
				}

				if v, ok := r["noncurrent_version_expiration_days"].(int); ok && v > 0 {
					lifecycleRule.NoncurrentVersionExpiration = &types.NoncurrentVersionExpiration{
						NoncurrentDays: aws.Int32(int32(v)),
					}
				}

				if v, ok := r["incomplete_upload_expiration_days"].(int); ok && v > 0 {
					lifecycleRule.AbortIncompleteMultipartUpload = &types.AbortIncompleteMultipartUpload{
						DaysAfterInitiation: aws.Int32(int32(v)),
					}
				}

				lifecycleConfig.Rules = append(lifecycleConfig.Rules, lifecycleRule)
			}

			_, err := s3Client.PutBucketLifecycleConfiguration(ctx, &s3.PutBucketLifecycleConfigurationInput{
				Bucket:                 aws.String(bucketName),
				LifecycleConfiguration: lifecycleConfig,
			})
			if err != nil {
				return fmt.Errorf("error updating lifecycle configuration for bucket %s: %v", bucketName, err)
			}
		}
	}

	return resourceGridscaleBucketRead(d, meta)
}

func resourceGridscaleBucketDelete(d *schema.ResourceData, meta interface{}) error {
	s3Host := d.Get("s3_host")
	accessKey := d.Get("access_key")
	secretKey := d.Get("secret_key")

	s3HostStr := s3Host.(string)
	s3Client := initS3Client(&gridscaleS3Provider{
		AccessKey: accessKey.(string),
		SecretKey: secretKey.(string),
	}, s3HostStr)

	bucketName := d.Get("bucket_name")
	bucketNameStr := bucketName.(string)
	bucketInput := s3.DeleteBucketInput{
		Bucket: &bucketNameStr,
	}

	errorPrefix := fmt.Sprintf("delete bucket %s resource at s3host %s-", bucketNameStr, s3HostStr)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	_, err := s3Client.DeleteBucket(ctx, &bucketInput)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

func initS3Client(provider *gridscaleS3Provider, s3host string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			provider.AccessKey,
			provider.SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Configure custom endpoint and path style
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://" + s3host)
		o.UsePathStyle = true
	})

	return client
}
