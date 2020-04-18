package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gridscale/gsclient-go/v2"
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
			Create: schema.DefaultTimeout(0 * time.Second),
			Update: schema.DefaultTimeout(0 * time.Second),
			Delete: schema.DefaultTimeout(0 * time.Second),
		},
	}
}

func resourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read network (%s) resource -", d.Id())
	network, err := client.GetNetwork(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	if err = d.Set("name", network.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", network.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("l2security", network.Properties.L2Security); err != nil {
		return fmt.Errorf("%s error setting l2security: %v", errorPrefix, err)
	}
	if err = d.Set("status", network.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("network_type", network.Properties.NetworkType); err != nil {
		return fmt.Errorf("%s error setting network_type: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", network.Properties.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", network.Properties.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", network.Properties.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("delete_block", network.Properties.DeleteBlock); err != nil {
		return fmt.Errorf("%s error setting delete_block: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", network.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", network.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}

	if err = d.Set("labels", network.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}

func resourceGridscaleNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update network (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.NetworkUpdateRequest{
		Name:       d.Get("name").(string),
		L2Security: d.Get("l2security").(bool),
		Labels:     &labels,
	}

	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutUpdate) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
		defer cancel()
	}
	err := client.UpdateNetwork(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.NetworkCreateRequest{
		Name:       d.Get("name").(string),
		L2Security: d.Get("l2security").(bool),
		Labels:     convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutCreate) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
		defer cancel()
	}
	response, err := client.CreateNetwork(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for network %v has been set to %v", requestBody.Name, response.ObjectUUID)

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete network (%s) resource -", d.Id())
	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutDelete) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
		defer cancel()
	}
	net, err := client.GetNetwork(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	//Stop all servers relating to this network address if there is one
	for _, server := range net.Properties.Relations.Servers {
		unlinkNetAction := func(ctx context.Context) error {
			err := client.UnlinkNetwork(ctx, server.ObjectUUID, d.Id())
			return err
		}
		//UnlinkNetwork requires the server to be off
		err = globalServerStatusList.runActionRequireServerOff(ctx, client, server.ObjectUUID, false, unlinkNetAction)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	err = client.DeleteNetwork(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
