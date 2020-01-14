package gridscale

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go"
)

func dataSourceGridscaleISOImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleISOImageRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"source_url": {
				Type:        schema.TypeString,
				Description: "Contains the source URL of the ISO Image that it was originally fetched from.",
				Computed:    true,
			},
			"server": {
				Type:        schema.TypeSet,
				Description: "The information about servers which are related to this ISO image.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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
				Description: "Upstream version of the ISO Image release",
				Computed:    true,
			},
			"private": {
				Type:        schema.TypeBool,
				Description: "The object is private, the value will be true. Otherwise the value will be false.",
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
				Description: "Description of the ISO Image.",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
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

func dataSourceGridscaleISOImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	isoimage, err := client.GetISOImage(emptyCtx, id)

	if err == nil {
		props := isoimage.Properties
		d.SetId(props.ObjectUUID)
		if err = d.Set("name", props.Name); err != nil {
			return fmt.Errorf("error setting name: %v", err)
		}
		if err = d.Set("source_url", props.SourceURL); err != nil {
			return fmt.Errorf("error setting source_url: %v", err)
		}
		if err = d.Set("location_uuid", props.LocationUUID); err != nil {
			return fmt.Errorf("error setting location_uuid: %v", err)
		}
		if err = d.Set("location_country", props.LocationCountry); err != nil {
			return fmt.Errorf("error setting location_country: %v", err)
		}
		if err = d.Set("location_iata", props.LocationIata); err != nil {
			return fmt.Errorf("error setting location_iata: %v", err)
		}
		if err = d.Set("location_name", props.LocationName); err != nil {
			return fmt.Errorf("error setting location_name: %v", err)
		}
		if err = d.Set("status", props.Status); err != nil {
			return fmt.Errorf("error setting status: %v", err)
		}
		if err = d.Set("version", props.Version); err != nil {
			return fmt.Errorf("error setting version: %v", err)
		}
		if err = d.Set("private", props.Private); err != nil {
			return fmt.Errorf("error setting private: %v", err)
		}
		if err = d.Set("create_time", props.CreateTime.String()); err != nil {
			return fmt.Errorf("error setting create_time: %v", err)
		}
		if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
			return fmt.Errorf("error setting change_time: %v", err)
		}
		if err = d.Set("description", props.Description); err != nil {
			return fmt.Errorf("error setting description: %v", err)
		}
		if err = d.Set("usage_in_minutes", props.UsageInMinutes); err != nil {
			return fmt.Errorf("error setting usage_in_minutes: %v", err)
		}
		if err = d.Set("capacity", props.Capacity); err != nil {
			return fmt.Errorf("error setting capacity: %v", err)
		}
		if err = d.Set("current_price", props.CurrentPrice); err != nil {
			return fmt.Errorf("error setting current_price: %v", err)
		}

		servers := make([]interface{}, 0)
		for _, value := range props.Relations.Servers {
			server := map[string]interface{}{
				"object_uuid": value.ObjectUUID,
				"create_time": value.CreateTime.String(),
				"object_name": value.ObjectName,
				"bootdevice":  value.Bootdevice,
			}
			servers = append(servers, server)
		}
		if err = d.Set("server", servers); err != nil {
			return fmt.Errorf("error setting server-rels: %v", err)
		}

		if err = d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("error setting labels: %v", err)
		}

		log.Printf("Found ISO image with key: %v", props.ObjectUUID)
	}

	return err
}
