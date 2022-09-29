package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type gridscaleS3Provider struct {
	AccessKey, SecretKey string
}

func (m *gridscaleS3Provider) Retrieve() (credentials.Value, error) {

	return credentials.Value{
		AccessKeyID:     m.AccessKey,
		SecretAccessKey: m.SecretKey,
	}, nil
}

func (m *gridscaleS3Provider) IsExpired() bool { return false }

func resourceGridscaleBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleBucketCreate,
		Read:   resourceGridscaleBucketRead,
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
			"loc_constrain": {
				Type:        schema.TypeString,
				Description: "The Location Constrain. Default: eu",
				Optional:    true,
				ForceNew:    true,
				Default:     "eu",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleBucketRead(d *schema.ResourceData, meta interface{}) error {
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

	loc := d.Get("loc_constrain")
	locStr := loc.(string)
	bucketName := d.Get("bucket_name")
	bucketNameStr := bucketName.(string)
	bucketInput := s3.CreateBucketInput{
		Bucket: &bucketNameStr,
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: &locStr,
		},
	}

	errorPrefix := fmt.Sprintf("Create bucket %s resource at s3host %s-", bucketNameStr, s3HostStr)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	_, err := s3Client.CreateBucketWithContext(ctx, &bucketInput)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	id := fmt.Sprintf("%s/%s", s3HostStr, bucketNameStr)
	d.SetId(id)

	log.Printf("The id for the new bucket has been set to %v", id)
	return nil
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
	_, err := s3Client.DeleteBucketWithContext(ctx, &bucketInput)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

func initS3Client(provider credentials.Provider, s3host string) *s3.S3 {
	forcePathStyle := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:           aws.String("default"),
			Endpoint:         &s3host,
			S3ForcePathStyle: &forcePathStyle,
			Credentials:      credentials.NewCredentials(provider),
		},
	}))
	return s3.New(sess)
}
