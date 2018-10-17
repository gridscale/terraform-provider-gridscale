package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:		schema.TypeString,
				Required:	true,
			},
			"token": {
				Type:		schema.TypeString,
				Required:	true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			//"gridscale_server":	dataSourceGridscaleServer(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gridscale_server":	resourceGridscaleServer(),
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
