package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gridscale_storage": dataSourceGridscaleStorage(),
			"gridscale_network": dataSourceGridscaleNetwork(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gridscale_server":  resourceGridscaleServer(),
			"gridscale_storage": resourceGridscaleStorage(),
			"gridscale_network": resourceGridscaleNetwork(),
			"gridscale_ip":      resourceGridscaleIp(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		UserUUID: d.Get("uuid").(string),
		APIToken: d.Get("token").(string),
	}

	return config.Client()
}
