package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscalePaaS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscalePaaSRead,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
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
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "Network UUID containing security zone",
				Computed:    true,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "Template that PaaS service uses",
				Computed:    true,
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
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource_limit": {
				Type:        schema.TypeSet,
				Description: "Resource for PaaS service",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscalePaaSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	paas, err := client.GetPaaSService(emptyCtx, id)
	if err != nil {
		return err
	}
	props := paas.Properties
	creds := props.Credentials
	d.SetId(props.ObjectUUID)
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
