package gridscale

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/nvthongswansea/gsclient-go"
)

func dataSourceGridscaleTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleTemplateRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the domain",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func dataSourceGridscaleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	name := d.Get("name").(string)

	template, err := client.GetTemplateByName(emptyCtx, name)

	if err == nil {
		d.SetId(template.Properties.ObjectUUID)
		log.Printf("Found template with key: %v", template.Properties.ObjectUUID)
	}

	return err
}
