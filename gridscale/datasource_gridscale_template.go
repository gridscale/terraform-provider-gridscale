package gridscale

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go/v3"
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
				Description: "The capacity of a storage/ISO Image/template/snapshot in GB.",
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

func dataSourceGridscaleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	name := d.Get("name").(string)
	errorPrefix := fmt.Sprintf("read template (%s) datasource -", name)

	template, err := client.GetTemplateByName(context.Background(), name)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(template.Properties.ObjectUUID)
	if err = d.Set("location_uuid", template.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", template.Properties.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", template.Properties.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", template.Properties.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("status", template.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("ostype", template.Properties.Ostype); err != nil {
		return fmt.Errorf("%s error setting ostype: %v", errorPrefix, err)
	}
	if err = d.Set("version", template.Properties.Version); err != nil {
		return fmt.Errorf("%s error setting version: %v", errorPrefix, err)
	}
	if err = d.Set("private", template.Properties.Private); err != nil {
		return fmt.Errorf("%s error setting private: %v", errorPrefix, err)
	}
	if err = d.Set("license_product_no", template.Properties.LicenseProductNo); err != nil {
		return fmt.Errorf("%s error setting license_product_no: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", template.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", template.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("distro", template.Properties.Distro); err != nil {
		return fmt.Errorf("%s error setting distro: %v", errorPrefix, err)
	}
	if err = d.Set("description", template.Properties.Description); err != nil {
		return fmt.Errorf("%s error setting description: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", template.Properties.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("capacity", template.Properties.Capacity); err != nil {
		return fmt.Errorf("%s error setting capacity: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", template.Properties.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}

	if err = d.Set("labels", template.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	log.Printf("Found template with key: %v", template.Properties.ObjectUUID)

	return nil
}
