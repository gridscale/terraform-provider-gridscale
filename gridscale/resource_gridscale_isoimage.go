package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"
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
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"source_url": {
				Type:        schema.TypeString,
				Description: "Contains the source URL of the ISO image that it was originally fetched from.",
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
				Description: "The location this object is placed.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Two digit country code (ISO 3166-2) of the location where this object is placed.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Upstream version of the ISO image release",
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
				Description: "Description of the ISO image.",
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
				Description: "The capacity of a storage/ISO image/template/snapshot in GB.",
				Computed:    true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Defines the price for the current period since the last bill.",
				Computed:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleISOImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read ISO-Image (%s) resource -", d.Id())
	iso, err := client.GetISOImage(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := iso.Properties
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("source_url", props.SourceURL); err != nil {
		return fmt.Errorf("%s error setting source_url: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", props.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", props.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", props.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", props.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("version", props.Version); err != nil {
		return fmt.Errorf("%s error setting version: %v", errorPrefix, err)
	}
	if err = d.Set("private", props.Private); err != nil {
		return fmt.Errorf("%s error setting private: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("description", props.Description); err != nil {
		return fmt.Errorf("%s error setting description: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("capacity", props.Capacity); err != nil {
		return fmt.Errorf("%s error setting capacity: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", props.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
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
		return fmt.Errorf("%s error setting server-rels: %v", errorPrefix, err)
	}

	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
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

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateISOImage(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new ISO image has been set to %v", response.ObjectUUID)

	return resourceGridscaleISOImageRead(d, meta)
}

func resourceGridscaleISOImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update ISO-Image (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.ISOImageUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateISOImage(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleISOImageRead(d, meta)
}

func resourceGridscaleISOImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete ISO-Image (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	isoimage, err := client.GetISOImage(ctx, d.Id())
	//In case of 404, don't catch the error
	if errHandler.SuppressHTTPErrorCodes(err, http.StatusNotFound) != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	//Remove all links between this ISO-Image and all servers.
	for _, server := range isoimage.Properties.Relations.Servers {
		//No need to unlink when server returns 409 or 404
		err = errHandler.SuppressHTTPErrorCodes(
			client.UnlinkIsoImage(ctx, server.ObjectUUID, d.Id()),
			http.StatusConflict,
			http.StatusNotFound,
		)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}
	err = errHandler.SuppressHTTPErrorCodes(
		client.DeleteISOImage(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
