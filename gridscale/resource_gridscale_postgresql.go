package gridscale

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const postgresTemplateFlavourName = "postgres"

func resourceGridscalePostgreSQL() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscalePostgreSQLCreate,
		Read:   resourceGridscalePostgreSQLRead,
		Delete: resourceGridscalePostgreSQLDelete,
		Update: resourceGridscalePostgreSQLUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			requestedReleaseInterface, isReleaseSet := d.GetOk("release")
			requestedPerformanceClassInterface, isPerformanceClassSet := d.GetOk("performance_class")

			if !isReleaseSet {
				return errors.New("\"release\" has to be defined")
			}
			if !isPerformanceClassSet {
				return errors.New("\"performance_class\" has to be defined")
			}

			requestedRelease := requestedReleaseInterface.(string)
			requestedPerformanceClass := requestedPerformanceClassInterface.(string)
			client := meta.(*gsclient.Client)
			paasTemplates, err := client.GetPaaSTemplateList(ctx)

			if err != nil {
				return err
			}
			var chosenTemplate gsclient.PaaSTemplate
			var isReleasePerformanceClassValid bool
			releaseSupportedPerformanceClasses := make(map[string][]string)
			for _, template := range paasTemplates {
				if template.Properties.Flavour == postgresTemplateFlavourName {
					performanceClasses := releaseSupportedPerformanceClasses[template.Properties.Release]
					releaseSupportedPerformanceClasses[template.Properties.Release] = append(performanceClasses, template.Properties.PerformanceClass)
					if template.Properties.Release == requestedRelease && template.Properties.PerformanceClass == requestedPerformanceClass {
						isReleasePerformanceClassValid = true
						chosenTemplate = template
					}
				}
			}
			if !isReleasePerformanceClassValid {
				errMessage := fmt.Sprintf("release %v with performance class %s is not a valid PostgreSQL release/performance class. Valid releases with corresponding performance classes are:\n\t", requestedRelease, requestedPerformanceClass)
				for release, performanceClasses := range releaseSupportedPerformanceClasses {
					errMessage += fmt.Sprintf("release %s is compatible with following performance classes: %s\n\t", release, strings.Join(performanceClasses, ", "))
				}
				return errors.New(errMessage)
			}
			return validatePostgreSQLParameters(d, chosenTemplate)
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
				Description: `The PostgreSQL release of this instance.\n
				For convenience, please use gscloud https://github.com/gridscale/gscloud to get the list of available PostgreSQL service releases.`,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"performance_class": {
				Type:        schema.TypeString,
				Description: "Performance class of PostgreSQL service.",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, class := range postgreSQLPerformanceClasses {
						if v.(string) == class {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid PostgreSQL performance class. Valid values are: %v", v.(string), strings.Join(postgreSQLPerformanceClasses, ",")))
					}
					return
				},
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for PostgreSQL service. It is used to connect to the PostgreSQL instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for PostgreSQL service. It is used to connect to the PostgreSQL instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "The port numbers where this PostgreSQL service accepts connections.",
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
			"security_zone_uuid": {
				Type:        schema.TypeString,
				Description: "Security zone UUID linked to PostgreSQL service.",
				Deprecated:  "Security zone is deprecated for gridSQL, gridStore, and gridFs. Please consider to use private network instead.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the network that the service is attached to.",
				Optional:    true,
				Computed:    true,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "PaaS service template that PostgreSQL service uses.",
				Computed:    true,
			},
			"service_template_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template service's category used to create the service.",
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that PostgreSQL service is in use.",
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
				Description: "Current status of PostgreSQL service.",
				Computed:    true,
			},
			"max_core_count": {
				Type:        schema.TypeInt,
				Description: "Maximum CPU core count. The PostgreSQL instance's CPU core count will be autoscaled based on the workload. The number of cores stays between 1 and `max_core_count`.",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if 1 > v.(int) || v.(int) > 32 {
						errors = append(errors, fmt.Errorf("%v is not a valid value for number of \"max_core_count\". Valid value should be between 1 and 32", v.(int)))
					}
					return
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pgaudit_log_bucket": {
				Type:        schema.TypeString,
				Description: "Object Storage bucket to upload audit logs to. For pgAudit to be enabled these additional parameters need to be configured: pgaudit_log_server_url, pgaudit_log_access_key, pgaudit_log_secret_key.",
				Optional:    true,
			},
			"pgaudit_log_server_url": {
				Type:        schema.TypeString,
				Description: "Object Storage server URL the bucket is located on.",
				Optional:    true,
			},
			"pgaudit_log_access_key": {
				Type:        schema.TypeString,
				Description: "Access key used to authenticate against Object Storage server.",
				Optional:    true,
			},
			"pgaudit_log_secret_key": {
				Type:        schema.TypeString,
				Description: "Secret key used to authenticate against Object Storage server.",
				Optional:    true,
			},
			"pgaudit_log_rotation_frequency": {
				Type:        schema.TypeInt,
				Description: "Rotation (in minutes) for audit logs. Logs are uploaded to Object Storage once rotated.",
				Optional:    true,
				Default:     5,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceGridscalePostgreSQLRead(d *schema.ResourceData, meta interface{}) error {
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
	if err = d.Set("network_uuid", props.NetworkUUID); err != nil {
		return fmt.Errorf("%s error setting network_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_uuid", props.ServiceTemplateUUID); err != nil {
		return fmt.Errorf("%s error setting service_template_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_category", props.ServiceTemplateCategory); err != nil {
		return fmt.Errorf("%s error setting service_template_category: %v", errorPrefix, err)
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

	//Get core count's limit
	for _, value := range props.ResourceLimits {
		if value.Resource == "cores" {
			if err = d.Set("max_core_count", value.Limit); err != nil {
				return fmt.Errorf("%s error setting max_core_count: %v", errorPrefix, err)
			}
		}
	}

	// Set PostgreSQL parameters
	if err = d.Set("pgaudit_log_bucket", props.Parameters["pgaudit_log_bucket"]); err != nil {
		return fmt.Errorf("%s error setting pgaudit_log_bucket: %v", errorPrefix, err)
	}
	if err = d.Set("pgaudit_log_server_url", props.Parameters["pgaudit_log_server_url"]); err != nil {
		return fmt.Errorf("%s error setting pgaudit_log_server_url: %v", errorPrefix, err)
	}
	if err = d.Set("pgaudit_log_access_key", props.Parameters["pgaudit_log_access_key"]); err != nil {
		return fmt.Errorf("%s error setting pgaudit_log_access_key: %v", errorPrefix, err)
	}
	if err = d.Set("pgaudit_log_secret_key", props.Parameters["pgaudit_log_secret_key"]); err != nil {
		return fmt.Errorf("%s error setting pgaudit_log_secret_key: %v", errorPrefix, err)
	}
	if err = d.Set("pgaudit_log_rotation_frequency", props.Parameters["pgaudit_log_rotation_frequency"]); err != nil {
		return fmt.Errorf("%s error setting pgaudit_log_rotation_frequency: %v", errorPrefix, err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	// Look for security zone's network that the PaaS service is connected to
	// (if the paas is connected to security zone. O.w skip)
	if props.SecurityZoneUUID == "" {
		return nil
	}
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	//look for a network that the PostgreSQL service is in
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

func resourceGridscalePostgreSQLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	// Get postgres template UUID
	release := d.Get("release").(string)
	performanceClass := d.Get("performance_class").(string)
	templateUUID, err := getPostgresTemplateUUID(client, release, performanceClass)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: templateUUID,
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	networkUUIDInf, isNetworkSet := d.GetOk("network_uuid")
	if isNetworkSet {
		requestBody.NetworkUUID = networkUUIDInf.(string)
	}
	// If network_uuid is set, skip setting security_zone_uuid.
	if secZoneUUIDInf, ok := d.GetOk("security_zone_uuid"); ok && !isNetworkSet {
		requestBody.PaaSSecurityZoneUUID = secZoneUUIDInf.(string)
	}
	if val, ok := d.GetOk("max_core_count"); ok {
		limits := []gsclient.ResourceLimit{
			{
				Resource: "cores",
				Limit:    val.(int),
			},
		}
		requestBody.ResourceLimits = limits
	}
	requestBody.Parameters = make(map[string]interface{})
	requestBody.Parameters["pgaudit_log_bucket"] = d.Get("pgaudit_log_bucket")
	requestBody.Parameters["pgaudit_log_server_url"] = d.Get("pgaudit_log_server_url")
	requestBody.Parameters["pgaudit_log_access_key"] = d.Get("pgaudit_log_access_key")
	requestBody.Parameters["pgaudit_log_secret_key"] = d.Get("pgaudit_log_secret_key")
	requestBody.Parameters["pgaudit_log_rotation_frequency"] = d.Get("pgaudit_log_rotation_frequency")

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for PostgreSQL service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscalePostgreSQLRead(d, meta)
}

func resourceGridscalePostgreSQLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update k8s (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}
	if d.HasChange("network_uuid") {
		requestBody.NetworkUUID = d.Get("network_uuid").(string)
	}
	// Only update templateUUID, when `release` or `performance_class` is changed
	if d.HasChange("release") || d.HasChange("performance_class") {
		// Get postgres template UUID
		release := d.Get("release").(string)
		performanceClass := d.Get("performance_class").(string)
		templateUUID, err := getPostgresTemplateUUID(client, release, performanceClass)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		requestBody.PaaSServiceTemplateUUID = templateUUID
	}
	requestBody.Parameters = make(map[string]interface{})
	requestBody.Parameters["pgaudit_log_bucket"] = d.Get("pgaudit_log_bucket")
	requestBody.Parameters["pgaudit_log_server_url"] = d.Get("pgaudit_log_server_url")
	requestBody.Parameters["pgaudit_log_access_key"] = d.Get("pgaudit_log_access_key")
	requestBody.Parameters["pgaudit_log_secret_key"] = d.Get("pgaudit_log_secret_key")
	requestBody.Parameters["pgaudit_log_rotation_frequency"] = d.Get("pgaudit_log_rotation_frequency")

	if val, ok := d.GetOk("max_core_count"); ok {
		limits := []gsclient.ResourceLimit{
			{
				Resource: "cores",
				Limit:    val.(int),
			},
		}
		requestBody.ResourceLimits = limits
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscalePostgreSQLRead(d, meta)
}

func resourceGridscalePostgreSQLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete paas (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.SuppressHTTPErrorCodes(
		client.DeletePaaSService(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

// getPostgresTemplateUUID returns the UUID of the postgres service template.
func getPostgresTemplateUUID(client *gsclient.Client, release, performanceClass string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == postgresTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release && template.Properties.PerformanceClass == performanceClass {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid {
		return "", fmt.Errorf("%v is not a valid PostgreSQL release. Valid releases are: %v", release, strings.Join(releases, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}

// validatePostgreSQLParameters validates PostgreSQL parameter taken from passed input against as well passed PaaS template.
func validatePostgreSQLParameters(d *schema.ResourceDiff, template gsclient.PaaSTemplate) error {
	var errorMessages []string
	if logBucket, ok := d.GetOk("pgaudit_log_bucket"); ok {
		if scheme, ok := template.Properties.ParametersSchema["pgaudit_log_bucket"]; ok {
			validBucket := regexp.MustCompile(scheme.Regex)
			if !validBucket.MatchString(logBucket.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'pgaudit_log_bucket' value. Value needs to match RegEx: '%s'\n", scheme.Regex))
			}
		}
	}
	if logServerURL, ok := d.GetOk("pgaudit_log_server_url"); ok {
		_, err := url.ParseRequestURI(logServerURL.(string))
		if err != nil {
			errorMessages = append(errorMessages, "Invalid 'pgaudit_log_bucket' value, doesn't match URL format")
		}
	}
	if logAccessKey, ok := d.GetOk("pgaudit_log_access_key"); ok {
		if scheme, ok := template.Properties.ParametersSchema["pgaudit_log_access_key"]; ok {
			validAccessKey := regexp.MustCompile(scheme.Regex)
			if !validAccessKey.MatchString(logAccessKey.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'pgaudit_log_access_key' value. Value needs to match RegEx: '%s'\n", scheme.Regex))
			}
		}
	}
	if logSecretKey, ok := d.GetOk("pgaudit_log_secret_key"); ok {
		if scheme, ok := template.Properties.ParametersSchema["pgaudit_log_secret_key"]; ok {
			validSecretKey := regexp.MustCompile(scheme.Regex)
			if !validSecretKey.MatchString(logSecretKey.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'pgaudit_log_secret_key' value. Value needs to match RegEx: '%s'\n", scheme.Regex))
			}
		}
	}
	if logRotationFrequency, ok := d.GetOk("pgaudit_log_rotation_frequency"); ok {
		if scheme, ok := template.Properties.ParametersSchema["pgaudit_log_rotation_frequency"]; ok {
			if scheme.Min > logRotationFrequency.(int) || logRotationFrequency.(int) > scheme.Max {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'pgaudit_log_rotation_frequency' value. Value must stays between %d and %d\n", scheme.Min, scheme.Max))
			}
		}
	}
	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, ""))
	}
	return nil
}
