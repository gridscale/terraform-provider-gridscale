package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGridscaleTemplate() *schema.Resource {
	return &schema.Resource{
		Read: resourceGridscaleTemplateRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Optional:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object",
				Computed:    true,
			},
			"ostype": {
				Type:        schema.TypeString,
				Description: "The operating system installed in the template",
				Computed:    true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private": {
				Type:        schema.TypeBool,
				Description: "The object is private, the value will be true. Otherwise the value will be false.",
				Computed:    true,
			},
			"license_product_no": {
				Type:        schema.TypeInt,
				Description: "If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change.",
				Computed:    true,
			},
			"distro": {
				Type:        schema.TypeString,
				Description: "The OS distrobution that the Template contains.",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the Template.",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Total minutes the object has been running.",
				Computed:    true,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Description: "The capacity of a storage/ISO-Image/template/snapshot in GB.",
				Computed:    true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Defines the price for the current period since the last bill.",
				Computed:    true,
			},
		},
	}
}

func resourceGridscaleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	template, err := client.GetTemplate(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := template.Properties
	d.Set("name", props.Name)
	d.Set("location_uuid", props.LocationUUID)
	d.Set("location_country", props.LocationCountry)
	d.Set("location_iata", props.LocationIata)
	d.Set("location_name", props.LocationName)
	d.Set("status", props.Status)
	d.Set("ostype", props.Ostype)
	d.Set("version", props.Version)
	d.Set("private", props.Private)
	d.Set("license_product_no", props.LicenseProductNo)
	d.Set("create_time", props.CreateTime)
	d.Set("change_time", props.ChangeTime)
	d.Set("distro", props.Distro)
	d.Set("description", props.Description)
	d.Set("usage_in_minutes", props.UsageInMinutes)
	d.Set("capacity", props.Capacity)
	d.Set("current_price", props.CurrentPrice)

	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}
