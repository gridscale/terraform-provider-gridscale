package gridscale

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const (
	msSQLTemplateFlavourName = "mssql"
	defaultBackupServerURL   = "https://gos3.io/"
)

func resourceGridscaleMSSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleMSSQLServerCreate,
		Read:   resourceGridscaleMSSQLServerRead,
		Delete: resourceGridscaleMSSQLServerDelete,
		Update: resourceGridscaleMSSQLServerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			client := meta.(*gsclient.Client)
			paasTemplates, err := client.GetPaaSTemplateList(ctx)
			if err != nil {
				return err
			}

			releaseVal := d.Get("release").(string)
			perfClassVal := d.Get("performance_class").(string)
			var chosenTemplate gsclient.PaaSTemplate
			var isReleasePerfClassValid bool
			releaseWPerfClasess := make(map[string][]string)
			for _, template := range paasTemplates {
				if template.Properties.Flavour == msSQLTemplateFlavourName {
					perfClasses := releaseWPerfClasess[template.Properties.Release]
					releaseWPerfClasess[template.Properties.Release] = append(perfClasses, template.Properties.PerformanceClass)
					if template.Properties.Release == releaseVal && template.Properties.PerformanceClass == perfClassVal {
						isReleasePerfClassValid = true
						chosenTemplate = template
					}
				}
			}
			if !isReleasePerfClassValid {
				errMess := fmt.Sprintf("release %v with performance class %s is not a valid MSSQL release/performance class. Valid releases with corresponding performance classes are:\n\t", releaseVal, perfClassVal)
				for release, perfClasses := range releaseWPerfClasess {
					errMess += fmt.Sprintf("release %s has following perfomance classes: %s\n\t", release, strings.Join(perfClasses, ", "))
				}
				return errors.New(errMess)
			}
			return validateMSSQLParameters(d, chosenTemplate)
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"release": {
				Type: schema.TypeString,
				Description: `The MS SQL Server release of this instance.\n
				For convenience, please use gscloud https://github.com/gridscale/gscloud to get the list of available MS SQL Server releases.`,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"performance_class": {
				Type:        schema.TypeString,
				Description: "Performance class of MS SQL Server.",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, class := range msSQLServerPerformanceClasses {
						if v.(string) == class {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid MS SQL Server performance class. Valid values are: %v", v.(string), strings.Join(postgreSQLPerformanceClasses, ",")))
					}
					return
				},
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for MS SQL Server . It is used to connect to the MS SQL Server instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for MS SQL Server. It is used to connect to the MS SQL Server instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "The port numbers where this MS SQL Server accepts connections.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"s3_backup": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Allow backup/restore MS SQL server to/from a S3 bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Object Storage bucket to upload backups to and restore backups from.",
						},
						"backup_retention": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Retention (in seconds) for local originals of backups. (0 for immediate removal once uploaded to Object Storage (default), higher values for delayed removal after the given time and once uploaded to Object Storage).",
						},
						"backup_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Access key used to authenticate against Object Storage server.",
						},
						"backup_secret_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Secret key used to authenticate against Object Storage server.",
						},
						"backup_server_url": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  defaultBackupServerURL,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								if v.(string) != defaultBackupServerURL {
									errors = append(errors, fmt.Errorf("Currently, only %s is supported", defaultBackupServerURL))
								}
								return
							},
							Description: "Object Storage server URL the bucket is located on.",
						},
					},
				},
			},
			"security_zone_uuid": {
				Type:        schema.TypeString,
				Description: "Security zone UUID linked to MS SQL Server.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "Network UUID containing security zone.",
				Computed:    true,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "PaaS service template that MS SQL Server uses.",
				Computed:    true,
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that MS SQL Server is in use.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Date time this service has been created.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of MS SQL Server.",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceGridscaleMSSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read paas (%s) resource -", d.Id())
	paas, err := client.GetPaaSService(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := paas.Properties
	creds := props.Credentials
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if len(creds) > 0 {
		if err = d.Set("username", creds[0].Username); err != nil {
			return fmt.Errorf("%s error setting username: %v", errorPrefix, err)
		}
		if err = d.Set("password", creds[0].Password); err != nil {
			return fmt.Errorf("%s error setting password: %v", errorPrefix, err)
		}
	}
	if err = d.Set("security_zone_uuid", props.SecurityZoneUUID); err != nil {
		return fmt.Errorf("%s error setting security_zone_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_uuid", props.ServiceTemplateUUID); err != nil {
		return fmt.Errorf("%s error setting service_template_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}

	//Get listen ports
	listenPorts := make([]interface{}, 0)
	for host, value := range props.ListenPorts {
		for k, portValue := range value {
			port := map[string]interface{}{
				"name": k,
				"host": host,
				"port": portValue,
			}
			listenPorts = append(listenPorts, port)
		}
	}
	if err = d.Set("listen_port", listenPorts); err != nil {
		return fmt.Errorf("%s error setting listen ports: %v", errorPrefix, err)
	}

	if props.Parameters["backup_bucket"] != nil {
		var s3Backup []interface{}
		s3Backup = append(s3Backup, map[string]interface{}{
			"backup_bucket":     props.Parameters["backup_bucket"],
			"backup_retention":  props.Parameters["backup_retention"],
			"backup_access_key": props.Parameters["backup_access_key"],
			"backup_secret_key": props.Parameters["backup_secret_key"],
			"backup_server_url": props.Parameters["backup_server_url"],
		})
		if err = d.Set("s3_backup", s3Backup); err != nil {
			return fmt.Errorf("%s error setting s3_backup: %v", errorPrefix, err)
		}
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	//Get all available networks
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	//look for a network that the MSSQLServer service is in
	for _, network := range networks {
		securityZones := network.Properties.Relations.PaaSSecurityZones
		//Each network can contain only one security zone
		if len(securityZones) >= 1 {
			if securityZones[0].ObjectUUID == props.SecurityZoneUUID {
				if err = d.Set("network_uuid", network.Properties.ObjectUUID); err != nil {
					return fmt.Errorf("%s error setting network_uuid: %v", errorPrefix, err)
				}
			}
		}
	}
	return nil
}

func resourceGridscaleMSSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	// get ms sql template UUID
	release := d.Get("release").(string)
	performanceClass := d.Get("performance_class").(string)
	templateUUID, err := getMSSQLTemplateUUID(client, release, performanceClass)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: templateUUID,
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
		PaaSSecurityZoneUUID:    d.Get("security_zone_uuid").(string),
	}

	// If s3_backup is set, attach the backup parameters to the create request.
	if _, ok := d.GetOk("s3_backup"); ok {
		params := make(map[string]interface{})
		params["backup_bucket"] = d.Get("s3_backup.0.backup_bucket")
		params["backup_retention"] = d.Get("s3_backup.0.backup_retention")
		params["backup_access_key"] = d.Get("s3_backup.0.backup_access_key")
		params["backup_secret_key"] = d.Get("s3_backup.0.backup_secret_key")
		params["backup_server_url"] = d.Get("s3_backup.0.backup_server_url")
		requestBody.Parameters = params
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for MSSQLServer service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleMSSQLServerRead(d, meta)
}

func resourceGridscaleMSSQLServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update k8s (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	// Only update templateUUID, when `release` or `performance_class` is changed.
	if d.HasChange("performance_class") || d.HasChange("release") {
		// get ms sql template UUID
		release := d.Get("release").(string)
		performanceClass := d.Get("performance_class").(string)
		templateUUID, err := getMSSQLTemplateUUID(client, release, performanceClass)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		requestBody.PaaSServiceTemplateUUID = templateUUID
	}

	params := make(map[string]interface{})
	if _, ok := d.GetOk("s3_backup"); ok {
		params["backup_bucket"] = d.Get("s3_backup.0.backup_bucket")
		params["backup_retention"] = d.Get("s3_backup.0.backup_retention")
		params["backup_access_key"] = d.Get("s3_backup.0.backup_access_key")
		params["backup_secret_key"] = d.Get("s3_backup.0.backup_secret_key")
		params["backup_server_url"] = d.Get("s3_backup.0.backup_server_url")
	}
	requestBody.Parameters = params

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleMSSQLServerRead(d, meta)
}

func resourceGridscaleMSSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete paas (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeletePaaSService(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

// getMSSQLTemplateUUID returns the UUID of the ms SQL service template.
func getMSSQLTemplateUUID(client *gsclient.Client, release, performanceClass string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == msSQLTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release && template.Properties.PerformanceClass == performanceClass {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid {
		return "", fmt.Errorf("%v is not a valid MS SQL Server release. Valid releases are: %v\n", release, strings.Join(releases, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}

func validateMSSQLParameters(d *schema.ResourceDiff, template gsclient.PaaSTemplate) error {
	var errorMessages []string
	if backupBucket, ok := d.GetOk("s3_backup.0.backup_bucket"); ok {
		if scheme, ok := template.Properties.ParametersSchema["backup_bucket"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(backupBucket.(string)) {
				errorMessages = append(errorMessages, "Invalid 'backup_bucket' value.\n")
			}
		}
	}
	if rentation, ok := d.GetOk("s3_backup.0.backup_retention"); ok {
		if scheme, ok := template.Properties.ParametersSchema["backup_retention"]; ok {
			if scheme.Min > rentation.(int) || rentation.(int) > scheme.Max {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'backup_retention' value. Value must stays between %d and %d\n", scheme.Min, scheme.Max))
			}
		}
	}
	if accessKey, ok := d.GetOk("s3_backup.0.backup_access_key"); ok {
		if scheme, ok := template.Properties.ParametersSchema["backup_access_key"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(accessKey.(string)) {
				errorMessages = append(errorMessages, "Invalid 'backup_access_key' value.\n")
			}
		}
	}
	if secretKey, ok := d.GetOk("s3_backup.0.backup_secret_key"); ok {
		if scheme, ok := template.Properties.ParametersSchema["backup_secret_key"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(secretKey.(string)) {
				errorMessages = append(errorMessages, "Invalid 'backup_secret_key' value.\n")
			}
		}
	}
	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, ""))
	}
	return nil
}
