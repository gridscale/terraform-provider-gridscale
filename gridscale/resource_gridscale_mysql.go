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

const mysqlTemplateFlavourName = "mysql"

func resourceGridscaleMySQL() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleMySQLCreate,
		Read:   resourceGridscaleMySQLRead,
		Delete: resourceGridscaleMySQLDelete,
		Update: resourceGridscaleMySQLUpdate,
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
				if template.Properties.Flavour == mysqlTemplateFlavourName {
					perfClasses := releaseWPerfClasess[template.Properties.Release]
					releaseWPerfClasess[template.Properties.Release] = append(perfClasses, template.Properties.PerformanceClass)
					if template.Properties.Release == releaseVal && template.Properties.PerformanceClass == perfClassVal {
						isReleasePerfClassValid = true
						chosenTemplate = template
					}
				}
			}
			if !isReleasePerfClassValid {
				errMess := fmt.Sprintf("release %v with performance class %s is not a valid MySQL release/performance class. Valid releases with corresponding performance classes are:\n\t", releaseVal, perfClassVal)
				for release, perfClasses := range releaseWPerfClasess {
					errMess += fmt.Sprintf("release %s has following perfomance classes: %s\n\t", release, strings.Join(perfClasses, ", "))
				}
				return errors.New(errMess)
			}
			return validateMySQLParameters(d, chosenTemplate)
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
				Description: `The MySQL release of this instance.\n
				For convenience, please use gscloud https://github.com/gridscale/gscloud to get the list of available MySQL service releases.`,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"performance_class": {
				Type:         schema.TypeString,
				Description:  "Performance class of MySQL service.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"mysql_log_bin": {
				Type:        schema.TypeBool,
				Description: "Binary Logging.",
				Optional:    true,
				Default:     false,
			},
			"mysql_sql_mode": {
				Type:        schema.TypeString,
				Description: "SQL Mode.",
				Optional:    true,
				Default:     "ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION",
			},
			"mysql_server_id": {
				Type:        schema.TypeInt,
				Description: "Server Id.",
				Optional:    true,
				Default:     1,
			},
			"mysql_query_cache": {
				Type:        schema.TypeBool,
				Description: "Enable query cache.",
				Optional:    true,
				Default:     true,
			},
			"mysql_binlog_format": {
				Type:        schema.TypeString,
				Description: "Binary Logging Format.",
				Optional:    true,
				Default:     "ROW",
			},
			"mysql_max_connections": {
				Type:        schema.TypeInt,
				Description: "Max Connections.",
				Optional:    true,
				Default:     4000,
			},
			"mysql_query_cache_size": {
				Type:        schema.TypeString,
				Description: "Query Cache Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)).",
				Optional:    true,
				Default:     "128M",
			},
			"mysql_default_time_zone": {
				Type:        schema.TypeString,
				Description: "Server Timezone.",
				Optional:    true,
				Default:     "UTC",
			},
			"mysql_query_cache_limit": {
				Type:        schema.TypeString,
				Description: "Query Cache Limit. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)).",
				Optional:    true,
				Default:     "1M",
			},
			"mysql_max_allowed_packet": {
				Type:        schema.TypeString,
				Description: "Max Allowed Packet Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)).",
				Optional:    true,
				Default:     "64M",
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for MySQL service. It is used to connect to the MySQL instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for MySQL service. It is used to connect to the MySQL instance.",
				Computed:    true,
				Sensitive:   true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "The port numbers where this MySQL service accepts connections.",
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
				Description: "Security zone UUID linked to MySQL service.",
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
				Description: "PaaS service template that MySQL service uses.",
				Computed:    true,
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that MySQL service is in use.",
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
				Description: "Current status of MySQL service.",
				Computed:    true,
			},
			"max_core_count": {
				Type:         schema.TypeInt,
				Description:  "Maximum CPU core count. The MySQL instance's CPU core count will be autoscaled based on the workload. The number of cores stays between 1 and `max_core_count`.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
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

func resourceGridscaleMySQLRead(d *schema.ResourceData, meta interface{}) error {
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

	// Set MySQL parameters
	if err = d.Set("mysql_log_bin", props.Parameters["mysql_log_bin"]); err != nil {
		return fmt.Errorf("%s error setting mysql_log_bin: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_sql_mode", props.Parameters["mysql_sql_mode"]); err != nil {
		return fmt.Errorf("%s error setting mysql_sql_mode: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_server_id", props.Parameters["mysql_server_id"]); err != nil {
		return fmt.Errorf("%s error setting mysql_server_id: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_query_cache", props.Parameters["mysql_query_cache"]); err != nil {
		return fmt.Errorf("%s error setting mysql_query_cache: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_binlog_format", props.Parameters["mysql_binlog_format"]); err != nil {
		return fmt.Errorf("%s error setting mysql_binlog_format: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_max_connections", props.Parameters["mysql_max_connections"]); err != nil {
		return fmt.Errorf("%s error setting mysql_max_connections: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_query_cache_size", props.Parameters["mysql_query_cache_size"]); err != nil {
		return fmt.Errorf("%s error setting mysql_query_cache_size: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_default_time_zone", props.Parameters["mysql_default_time_zone"]); err != nil {
		return fmt.Errorf("%s error setting mysql_default_time_zone: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_query_cache_limit", props.Parameters["mysql_query_cache_limit"]); err != nil {
		return fmt.Errorf("%s error setting mysql_query_cache_limit: %v", errorPrefix, err)
	}
	if err = d.Set("mysql_max_allowed_packet", props.Parameters["mysql_max_allowed_packet"]); err != nil {
		return fmt.Errorf("%s error setting mysql_max_allowed_packet: %v", errorPrefix, err)
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

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	//Get all available networks
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	//look for a network that the MySQL service is in
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

func resourceGridscaleMySQLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	// Get mysql template UUID
	release := d.Get("release").(string)
	performanceClass := d.Get("performance_class").(string)
	templateUUID, err := getMySQLTemplateUUID(client, release, performanceClass)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: templateUUID,
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
		PaaSSecurityZoneUUID:    d.Get("security_zone_uuid").(string),
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

	params := make(map[string]interface{})
	params["mysql_log_bin"] = d.Get("mysql_log_bin")
	params["mysql_sql_mode"] = d.Get("mysql_sql_mode")
	params["mysql_server_id"] = d.Get("mysql_server_id")
	params["mysql_query_cache"] = d.Get("mysql_query_cache")
	params["mysql_binlog_format"] = d.Get("mysql_binlog_format")
	params["mysql_max_connections"] = d.Get("mysql_max_connections")
	params["mysql_query_cache_size"] = d.Get("mysql_query_cache_size")
	params["mysql_default_time_zone"] = d.Get("mysql_default_time_zone")
	params["mysql_query_cache_limit"] = d.Get("mysql_query_cache_limit")
	params["mysql_max_allowed_packet"] = d.Get("mysql_max_allowed_packet")
	requestBody.Parameters = params

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for MySQL service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleMySQLRead(d, meta)
}

func resourceGridscaleMySQLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update k8s (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	// Only update templateUUID, when `release` or `performance_class` is changed
	if d.HasChange("release") || d.HasChange("performance_class") {
		// Get mysql template UUID
		release := d.Get("release").(string)
		performanceClass := d.Get("performance_class").(string)
		templateUUID, err := getMySQLTemplateUUID(client, release, performanceClass)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		requestBody.PaaSServiceTemplateUUID = templateUUID
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

	params := make(map[string]interface{})
	params["mysql_log_bin"] = d.Get("mysql_log_bin")
	params["mysql_sql_mode"] = d.Get("mysql_sql_mode")
	params["mysql_server_id"] = d.Get("mysql_server_id")
	params["mysql_query_cache"] = d.Get("mysql_query_cache")
	params["mysql_binlog_format"] = d.Get("mysql_binlog_format")
	params["mysql_max_connections"] = d.Get("mysql_max_connections")
	params["mysql_query_cache_size"] = d.Get("mysql_query_cache_size")
	params["mysql_default_time_zone"] = d.Get("mysql_default_time_zone")
	params["mysql_query_cache_limit"] = d.Get("mysql_query_cache_limit")
	params["mysql_max_allowed_packet"] = d.Get("mysql_max_allowed_packet")
	requestBody.Parameters = params

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleMySQLRead(d, meta)
}

func resourceGridscaleMySQLDelete(d *schema.ResourceData, meta interface{}) error {
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

// getMySQLTemplateUUID returns the UUID of the mysql service template.
func getMySQLTemplateUUID(client *gsclient.Client, release, performanceClass string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == mysqlTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release && template.Properties.PerformanceClass == performanceClass {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid {
		return "", fmt.Errorf("%v is not a valid MySQL release. Valid releases are: %v\n", release, strings.Join(releases, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}

func validateMySQLParameters(d *schema.ResourceDiff, template gsclient.PaaSTemplate) error {
	var errorMessages []string
	if sqlMode, ok := d.GetOk("mysql_sql_mode"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_sql_mode"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(sqlMode.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_sql_mode' value. Example value: '%s'\n", scheme.Default))
			}
		}
	}
	if serverID, ok := d.GetOk("mysql_server_id"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_server_id"]; ok {
			if scheme.Min > serverID.(int) || serverID.(int) > scheme.Max {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_server_id' value. Value must stays between %d and %d\n", scheme.Min, scheme.Max))
			}
		}
	}
	if binLogFormat, ok := d.GetOk("mysql_binlog_format"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_binlog_format"]; ok {
			var isValidFormat bool
			for _, allowedValue := range scheme.Allowed {
				if binLogFormat.(string) == allowedValue {
					isValidFormat = true
				}
			}
			if !isValidFormat {
				errorMessages = append(errorMessages,
					fmt.Sprintf("Invalid 'mysql_binlog_format' value. Value must be one of these:\n\t%s\n",
						strings.Join(scheme.Allowed, "\n\t"),
					),
				)
			}
		}
	}
	if maxNConn, ok := d.GetOk("mysql_max_connections"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_max_connections"]; ok {
			if scheme.Min > maxNConn.(int) || maxNConn.(int) > scheme.Max {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_max_connections' value. Value must stays between %d and %d\n", scheme.Min, scheme.Max))
			}
		}
	}
	if cacheSize, ok := d.GetOk("mysql_query_cache_size"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_query_cache_size"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(cacheSize.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_query_cache_size' value. Example value: '%s'\n", scheme.Default))
			}
		}
	}
	if defaultTimeZone, ok := d.GetOk("mysql_default_time_zone"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_default_time_zone"]; ok {
			var isValidFormat bool
			for _, allowedValue := range scheme.Allowed {
				if defaultTimeZone.(string) == allowedValue {
					isValidFormat = true
				}
			}
			if !isValidFormat {
				errorMessages = append(errorMessages,
					fmt.Sprintf("Invalid 'mysql_default_time_zone' value. Value must be one of these:\n\t%s",
						strings.Join(scheme.Allowed, "\n\t"),
					),
				)
			}
		}
	}
	if cacheLimit, ok := d.GetOk("mysql_query_cache_limit"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_query_cache_limit"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(cacheLimit.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_query_cache_limit' value. Example value: '%s'\n", scheme.Default))
			}
		}
	}
	if maxAllowedPacket, ok := d.GetOk("mysql_max_allowed_packet"); ok {
		if scheme, ok := template.Properties.ParametersSchema["mysql_max_allowed_packet"]; ok {
			validMode := regexp.MustCompile(scheme.Regex)
			if !validMode.MatchString(maxAllowedPacket.(string)) {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'mysql_max_allowed_packet' value. Example value: '%s'\n", scheme.Default))
			}
		}
	}
	if maxNCore, ok := d.GetOk("max_core_count"); ok {
		autoscalingNCore := template.Properties.Autoscaling.Cores
		if autoscalingNCore.Min > maxNCore.(int) || maxNCore.(int) > autoscalingNCore.Max {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'max_core_count' value. Value must stays between %d and %d\n", autoscalingNCore.Min, autoscalingNCore.Max))
		}
	}
	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, ""))
	}
	return nil
}
