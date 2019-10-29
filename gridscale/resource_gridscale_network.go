package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleNetworkCreate,
		Read:   resourceGridscaleNetworkRead,
		Delete: resourceGridscaleNetworkDelete,
		Update: resourceGridscaleNetworkUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "MAC spoofing protection - filters layer2 and ARP traffic based on source MAC",
				Optional:    true,
				Default:     false,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object",
				Computed:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "The type of this network, can be mpls, breakout or network.",
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
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
			"delete_block": {
				Type:        schema.TypeBool,
				Description: "If deleting this network is allowed",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 3),
		},
	}
}

func resourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	network, err := client.GetNetwork(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", network.Properties.Name)
	d.Set("location_uuid", network.Properties.LocationUUID)
	d.Set("l2security", network.Properties.L2Security)
	d.Set("status", network.Properties.Status)
	d.Set("network_type", network.Properties.NetworkType)
	d.Set("location_country", network.Properties.LocationCountry)
	d.Set("location_iata", network.Properties.LocationIata)
	d.Set("location_name", network.Properties.LocationName)
	d.Set("delete_block", network.Properties.DeleteBlock)
	d.Set("create_time", network.Properties.CreateTime.String())
	d.Set("change_time", network.Properties.ChangeTime.String())

	if err = d.Set("labels", network.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.NetworkUpdateRequest{
		Name:       d.Get("name").(string),
		L2Security: d.Get("l2security").(bool),
	}

	err := client.UpdateNetwork(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.NetworkCreateRequest{
		Name:         d.Get("name").(string),
		LocationUUID: d.Get("location_uuid").(string),
		L2Security:   d.Get("l2security").(bool),
		Labels:       convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	response, err := client.CreateNetwork(emptyCtx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for network %v has been set to %v", requestBody.Name, response.ObjectUUID)

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	net, err := client.GetNetwork(emptyCtx, d.Id())
	if err != nil {
		return err
	}

	//Stop all servers relating to this network address if there is one
	for _, server := range net.Properties.Relations.Servers {
		unlinkNetAction := func(ctx context.Context) error {
			err = client.UnlinkNetwork(ctx, server.ObjectUUID, d.Id())
			return err
		}
		//UnlinkNetwork requires the server to be off
		err = serverPowerStateList.runActionRequireServerOff(emptyCtx, client, server.ObjectUUID, unlinkNetAction)
		if err != nil {
			return err
		}
	}

	return client.DeleteNetwork(emptyCtx, d.Id())
}
