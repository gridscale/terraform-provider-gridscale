package gridscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceGridscalePaaS() *schema.Resource {
	return &schema.Resource{
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

