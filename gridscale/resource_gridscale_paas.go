package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
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
	}
}

func resourceGridscalePaaSServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	paas, err := client.GetPaaSService(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := paas.Properties
	creds := props.Credentials
	d.Set("name", props.Name)
	if creds != nil && len(creds) > 0 {
		d.Set("username", creds[0].Username)
		d.Set("password", creds[0].Password)
	}
	d.Set("security_zone_uuid", props.SecurityZoneUUID)
	d.Set("service_template_uuid", props.ServiceTemplateUUID)
	d.Set("usage_in_minute", props.UsageInMinutes)
	d.Set("current_price", props.CurrentPrice)
	d.Set("change_time", props.ChangeTime)
	d.Set("create_time", props.CreateTime)
	d.Set("status", props.Status)

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
		return fmt.Errorf("Error setting listen ports: %v", err)
	}

	//Get parameters
	parameters := make([]interface{}, 0)
	for k, value := range props.Parameters {
		param := map[string]interface{}{
			"param": k,
			"value": value,
		}
		parameters = append(parameters, param)
	}
	if err = d.Set("parameter", parameters); err != nil {
		return fmt.Errorf("Error setting parameters: %v", err)
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
		return fmt.Errorf("Error setting resource limits: %v", err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	//Get all available networks
	networks, err := client.GetNetworkList(emptyCtx)
	if err != nil {
		return fmt.Errorf("Error getting networks: %v", err)
	}
	//look for a network that the PaaS service is in
	for _, network := range networks {
		securityZones := network.Properties.Relations.PaaSSecurityZones
		//Each network can contain only one security zone
		if len(securityZones) >= 1 {
			if securityZones[0].ObjectUUID == props.SecurityZoneUUID {
				d.Set("network_uuid", network.Properties.ObjectUUID)
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
		for k, v := range mapVal {
			if k == "param" {
				param = v.(string)
			}
			if k == "value" {
				val = v
			}
		}
		params[param] = val
	}
	requestBody.Parameters = params

	limits := make([]gsclient.ResourceLimit, 0)
	for _, value := range d.Get("resource_limit").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var resLim gsclient.ResourceLimit
		for k, v := range mapVal {
			if k == "resource" {
				resLim.Resource = v.(string)
			}
			if k == "limit" {
				resLim.Limit = v.(int)
			}
		}
		limits = append(limits, resLim)
	}
	requestBody.ResourceLimits = limits

	response, err := client.CreatePaaSService(emptyCtx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for PaaS service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscalePaaSServiceRead(d, meta)
}

func resourceGridscalePaaSServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	params := make(map[string]interface{}, 0)
	for _, value := range d.Get("parameter").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var param string
		var val interface{}
		for k, v := range mapVal {
			if k == "param" {
				param = v.(string)
			}
			if k == "value" {
				val = v
			}
		}
		params[param] = val
	}
	requestBody.Parameters = params

	limits := make([]gsclient.ResourceLimit, 0)
	for _, value := range d.Get("resource_limit").(*schema.Set).List() {
		mapVal := value.(map[string]interface{})
		var resLim gsclient.ResourceLimit
		for k, v := range mapVal {
			if k == "resource" {
				resLim.Resource = v.(string)
			}
			if k == "limit" {
				resLim.Limit = v.(int)
			}
		}
		limits = append(limits, resLim)
	}
	requestBody.ResourceLimits = limits

	err := client.UpdatePaaSService(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}
	return resourceGridscalePaaSServiceRead(d, meta)
}

func resourceGridscalePaaSServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	return client.DeletePaaSService(emptyCtx, d.Id())
}
