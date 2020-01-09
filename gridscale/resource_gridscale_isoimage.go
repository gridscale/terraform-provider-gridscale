package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceGridscaleISOImage() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGridscaleISOImageRead,
		Create: resourceGridscaleISOImageCreate,
		Update: resourceGridscaleISOImageUpdate,
		Delete: resourceGridscaleISOImageDelete,
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
				Description: "Contains the source URL of the ISO Image that it was originally fetched from.",
				Required:    true,
				ForceNew:    true,
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
				Description: "Description of the ISO Image.",
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

func resourceGridscaleISOImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	iso, err := client.GetISOImage(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := iso.Properties
	d.Set("name", props.Name)
	d.Set("source_url", props.SourceURL)
	d.Set("location_uuid", props.LocationUUID)
	d.Set("location_country", props.LocationCountry)
	d.Set("location_iata", props.LocationIata)
	d.Set("location_name", props.LocationName)
	d.Set("status", props.Status)
	d.Set("version", props.Version)
	d.Set("private", props.Private)
	d.Set("create_time", props.CreateTime)
	d.Set("change_time", props.ChangeTime)
	d.Set("description", props.Description)
	d.Set("usage_in_minutes", props.UsageInMinutes)
	d.Set("capacity", props.Capacity)
	d.Set("current_price", props.CurrentPrice)

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
		return fmt.Errorf("Error setting server-rels: %v", err)
	}

	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleISOImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.ISOImageCreateRequest{
		Name:      d.Get("name").(string),
		SourceURL: d.Get("source_url").(string),
		Labels:    convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	response, err := client.CreateISOImage(emptyCtx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new ISO image has been set to %v", response.ObjectUUID)

	return resourceGridscaleISOImageRead(d, meta)
}

func resourceGridscaleISOImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.ISOImageUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	err := client.UpdateISOImage(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleISOImageRead(d, meta)
}

func resourceGridscaleISOImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	return client.DeleteISOImage(emptyCtx, d.Id())
}
