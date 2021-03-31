package gridscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const postgresTemplateFlavourName = "postgres"

const (
	postgresReleaseValidationOpt = iota
	postgresMaxCoreCountValidationOpt
)

func resourceGridscalePostgresSQL() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscalePostgresSQLCreate,
		Read:   resourceGridscalePostgresSQLRead,
		Delete: resourceGridscalePostgresSQLDelete,
		Update: resourceGridscalePostgresSQLUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"release_no": {
				Type:         schema.TypeString,
				Description:  "Release no. of postgresSQL service.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"performance_class": {
				Type:        schema.TypeString,
				Description: "Performance class of postgresSQL service.",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for PostgresSQL service.",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for PostgresSQL service.",
				Computed:    true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "Ports that PostgresSQL service listens to.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
				Description: "Security zone UUID linked to PostgresSQL service.",
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
				Description: "PaaS service template that postgresSQL service uses.",
				Computed:    true,
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that PostgresSQL service is in use",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Time of the creation.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of PostgresSQL service.",
				Computed:    true,
			},
			"max_core_count": {
				Type:        schema.TypeInt,
				Description: "Maximum CPU core count. It is used for autoscaling.",
				Optional:    true,
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

func resourceGridscalePostgresSQLRead(d *schema.ResourceData, meta interface{}) error {
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
	if creds != nil && len(creds) > 0 {
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
	for _, value := range props.ListenPorts {
		for k, portValue := range value {
			port := map[string]interface{}{
				"name": k,
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
	//look for a network that the PostgresSQL service is in
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

func resourceGridscalePostgresSQLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	// Validate k8s parameters
	templateUUID, err := validatePostgresParameters(client, d,
		postgresReleaseValidationOpt,
		postgresMaxCoreCountValidationOpt,
	)
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

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for PostgresSQL service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscalePostgresSQLRead(d, meta)
}

func resourceGridscalePostgresSQLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update k8s (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	// Validate
	templateUUID, err := validatePostgresParameters(client, d,
		postgresReleaseValidationOpt,
		postgresMaxCoreCountValidationOpt,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	requestBody.PaaSServiceTemplateUUID = templateUUID

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
	err = client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscalePostgresSQLRead(d, meta)
}

func resourceGridscalePostgresSQLDelete(d *schema.ResourceData, meta interface{}) error {
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

// validatePostgresParameters validate postgres resource's selected parameters.
// It returns the UUID of the postgres service template, if the validation is successful.
// Otherwise, an error will be returned.
func validatePostgresParameters(client *gsclient.Client, d *schema.ResourceData, parameters ...int) (string, error) {
	errorMessages := []string{"List of validation errors:\n"}
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	// Check if the postgres release number exists
	release := d.Get("release_no").(string)
	performanceClass := d.Get("performance_class").(string)
	var isReleaseValid bool
	releases := make(map[string][]string)
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == postgresTemplateFlavourName {
			releasePerformanceClasses := releases[template.Properties.Release]
			releases[template.Properties.Release] = append(releasePerformanceClasses, template.Properties.PerformanceClass)
			if template.Properties.Release == release && template.Properties.PerformanceClass == performanceClass {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid && isIntInList(postgresReleaseValidationOpt, parameters) {
		errorMessages = append(errorMessages, fmt.Sprintf("%v/%v is not a valid postgresSQL release no/performance class. Valid release numbers (and corresponding performance classes) are:\n", release, performanceClass))
		for k, v := range releases {
			errorMessages = append(errorMessages, fmt.Sprintf("      Release No.: %s. Performance classes: %s\n", k, strings.Join(v, ", ")))
		}
	}

	// Check max core count are valid
	if val, ok := d.GetOk("max_core_count"); ok {
		coreCount := val.(int)
		minCoreCount := uTemplate.Properties.Autoscaling.Cores.Min
		maxCoreCount := uTemplate.Properties.Autoscaling.Cores.Max
		if (minCoreCount > coreCount || maxCoreCount < coreCount) &&
			isIntInList(postgresMaxCoreCountValidationOpt, parameters) {
			errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid value for \"max_core_count\". Valid value stays between %v and %v\n", coreCount, minCoreCount, maxCoreCount))
		}
	}
	if len(errorMessages) > 1 {
		return "", fmt.Errorf(strings.Join(errorMessages, ""))
	}
	return uTemplate.Properties.ObjectUUID, nil
}
