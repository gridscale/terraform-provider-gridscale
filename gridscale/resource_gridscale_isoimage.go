package gridscale

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGridscaleISOImage() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"source_url": {
				Type:        schema.TypeString,
				Description: "Contains the source URL of the ISO-Image that it was originally fetched from.",
				Required:    true,
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
			"version": {
				Type:        schema.TypeString,
				Description: "Upstream version of the ISO-Image release",
				Computed:    true,
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
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the ISO-Image.",
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
