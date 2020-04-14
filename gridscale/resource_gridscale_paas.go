package gridscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"log"
)

func resourceGridscalePaaS() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscalePaaSServiceCreate,
		Read:   resourceGridscalePaaSServiceRead,
		Delete: resourceGridscalePaaSServiceDelete,
		Update: resourceGridscalePaaSServiceUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for PaaS service",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for PaaS service",
				Computed:    true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "Ports that PaaS service listens to",
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
				Description: "Security zone UUID linked to PaaS service",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "Network UUID containing security zone",
				Computed:    true,
			},
			"service_template_uuid": {
				Type:         schema.TypeString,
				Description:  "Template that PaaS service uses",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"usage_in_minute": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that PaaS service is in use",
				Computed:    true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Current price of PaaS service",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Time of the creation",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of PaaS service",
				Computed:    true,
			},
			"parameter": {
				Type:        schema.TypeSet,
				Description: "Parameter for PaaS service",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								valid := false
								for _, primType := range supportedPrimTypes {
									if v.(string) == primType {
										valid = true
										break
									}
								}
								if !valid {
									errors = append(errors, fmt.Errorf("%v is not a valid primitive type. Valid primitive types are: %v", v.(string), strings.Join(supportedPrimTypes, ",")))
								}
								return
							},
						},
					},
				},
			},
			"resource_limit": {
				Type:        schema.TypeSet,
				Description: "Resource for PaaS service",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
						"limit": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Update: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Delete: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
		},
	}
}

func resourceGridscalePaaSServiceRead(d *schema.ResourceData, meta interface{}) error {
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
	if err = d.Set("usage_in_minute", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minute: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", props.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
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

	//Get parameters
	parameters := make([]interface{}, 0)
	for k, value := range props.Parameters {
		paramValType, err := getInterfaceType(value)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		valueInString, err := convInterfaceToString(paramValType, value)
		param := map[string]interface{}{
			"param": k,
			"value": valueInString,
			"type":  paramValType,
		}
		parameters = append(parameters, param)
	}
	if err = d.Set("parameter", parameters); err != nil {
		return fmt.Errorf("%s error setting parameters: %v", errorPrefix, err)
	}

	//Get resource limits
	resourceLimits := make([]interface{}, 0)
	for _, value := range props.ResourceLimits {
		limit := map[string]interface{}{
			"resource": value.Resource,
			"limit":    value.Limit,
		}
		resourceLimits = append(resourceLimits, limit)
	}
	if err = d.Set("resource_limit", resourceLimits); err != nil {
		return fmt.Errorf("%s error setting resource limits: %v", errorPrefix, err)
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
	//look for a network that the PaaS service is in
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

func resourceGridscalePaaSServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: d.Get("service_template_uuid").(string),
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
		PaaSSecurityZoneUUID:    d.Get("security_zone_uuid").(string),
	}

	params := make(map[string]interface{}, 0)
	for _, value := range d.Get("parameter").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var param string
		var val interface{}
		param = mapVal["param"].(string)
		paramValType := mapVal["type"].(string)
		typedVal, err := convStrToTypeInterface(paramValType, mapVal["value"].(string))
		if err != nil {
			return err
		}
		val = typedVal
		params[param] = val
	}
	requestBody.Parameters = params

	limits := make([]gsclient.ResourceLimit, 0)
	for _, value := range d.Get("resource_limit").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var resLim gsclient.ResourceLimit
		resLim.Resource = mapVal["resource"].(string)
		resLim.Limit = mapVal["limit"].(int)
		limits = append(limits, resLim)
	}
	requestBody.ResourceLimits = limits

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for PaaS service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscalePaaSServiceRead(d, meta)
}

func resourceGridscalePaaSServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update paas (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	params := make(map[string]interface{}, 0)
	for _, value := range d.Get("parameter").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var param string
		var val interface{}
		param = mapVal["param"].(string)
		paramValType := mapVal["type"].(string)
		typedVal, err := convStrToTypeInterface(paramValType, mapVal["value"].(string))
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		val = typedVal
		params[param] = val
	}
	requestBody.Parameters = params

	limits := make([]gsclient.ResourceLimit, 0)
	for _, value := range d.Get("resource_limit").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var resLim gsclient.ResourceLimit
		resLim.Resource = mapVal["resource"].(string)
		resLim.Limit = mapVal["limit"].(int)
		limits = append(limits, resLim)
	}
	requestBody.ResourceLimits = limits

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscalePaaSServiceRead(d, meta)
}

func resourceGridscalePaaSServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete paas (%s) resource -", d.Id())
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := client.DeletePaaSService(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
