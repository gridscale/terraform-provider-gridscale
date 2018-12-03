package gridscale

import (
	"bitbucket.org/gridscale/gsclient-go"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleServerCreate,
		Read:   resourceGridscaleServerRead,
		Delete: resourceGridscaleServerDelete,
		Update: resourceGridscaleServerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Required:    true,
			},
			"memory": {
				Type:         schema.TypeInt,
				Description:  "The amount of server memory in GB.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cores": {
				Type:         schema.TypeInt,
				Description:  "The number of server cores.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"hardware_profile": {
				Type:        schema.TypeString,
				Description: "The number of server cores.",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
			},
			"storages": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"networks": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv4": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"power": {
				Type:        schema.TypeBool,
				Description: "The number of server cores.",
				Optional:    true,
				Default:     false,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"auto_recovery": {
				Type:        schema.TypeInt,
				Description: "If the server should be auto-started in case of a failure (default=true).",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Defines which Availability-Zone the Server is placed.",
				Optional:    true,
			},
			"console_token": {
				Type:        schema.TypeString,
				Description: "If the server should be auto-started in case of a failure (default=true).",
				Computed:    true,
			},
			"legacy": {
				Type:        schema.TypeBool,
				Description: "Legacy-Hardware emulation instead of virtio hardware. If enabled, hotplugging cores, memory, storage, network, etc. will not work, but the server will most likely run every x86 compatible operating system. This mode comes with a performance penalty, as emulated hardware does not benefit from the virtio driver infrastructure.",
				Computed:    true,
			},
			"usage_in_minutes_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"usage_in_minutes_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	server, err := client.GetServer(d.Id())
	if err != nil {
		if requestError, ok := err.(*gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", server.Properties.Name)
	d.Set("memory", server.Properties.Memory)
	d.Set("cores", server.Properties.Cores)
	d.Set("hardware_profile", server.Properties.HardwareProfile)
	d.Set("location_uuid", server.Properties.LocationUuid)
	d.Set("power", server.Properties.Power)
	d.Set("current_price", server.Properties.CurrentPrice)
	d.Set("availability_zone", server.Properties.AvailablityZone)
	d.Set("auto_recovery", server.Properties.AutoRecovery)
	d.Set("console_token", server.Properties.ConsoleToken)
	d.Set("legacy", server.Properties.Legacy)
	d.Set("usage_in_minutes_memory", server.Properties.UsageInMinutesMemory)
	d.Set("usage_in_minutes_cores", server.Properties.UsageInMinutesCores)
	d.Set("labels", server.Properties.Labels)

	//Get storages
	storages := make([]interface{}, 0)
	for _, storage := range server.Properties.Relations.Storages {
		storages = append(storages, storage.ObjectUuid)
	}
	d.Set("storages", storages)

	//Get Networks
	networks := make([]interface{}, 0)
	for _, network := range server.Properties.Relations.Networks {
		if !network.PublicNet {
			networks = append(networks, network.NetworkUuid)
		}
	}
	d.Set("networks", networks)

	//Get IP addresses
	var ipv4, ipv6 string
	for _, ip := range server.Properties.Relations.PublicIps {
		if ip.Family == 4 {
			ipv4 = ip.ObjectUuid
		}
		if ip.Family == 6 {
			ipv6 = ip.ObjectUuid
		}
	}
	d.Set("ipv4", ipv4)
	d.Set("ipv6", ipv6)

	log.Printf("Read the following: %v", server)
	return nil
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.ServerCreateRequest{
		Name:            d.Get("name").(string),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
		LocationUuid:    d.Get("location_uuid").(string),
		HardwareProfile: d.Get("hardware_profile").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          d.Get("labels").(*schema.Set).List(),
	}

	requestBody.Relations.IsoImages = []gsclient.ServerIsoImage{}
	requestBody.Relations.Storages = []gsclient.ServerCreateRequestStorage{}
	requestBody.Relations.Networks = []gsclient.ServerCreateRequestNetwork{}
	requestBody.Relations.PublicIps = []gsclient.ServerCreateRequestIp{}

	if attr, ok := d.GetOk("storages"); ok {
		for index, value := range attr.(*schema.Set).List() {
			storage := gsclient.ServerCreateRequestStorage{
				StorageUuid: value.(string),
			}
			if index == 0 {
				storage.BootDevice = true
			}

			requestBody.Relations.Storages = append(requestBody.Relations.Storages, storage)
		}
	}

	if attr, ok := d.GetOk("ipv4"); ok {
		if client.GetIpVersion(attr.(string)) != 4 {
			return fmt.Errorf("The IP address with UUID %v is not version 4", attr.(string))
		}
		ip := gsclient.ServerCreateRequestIp{
			IpaddrUuid: attr.(string),
		}
		requestBody.Relations.PublicIps = append(requestBody.Relations.PublicIps, ip)
	}
	if attr, ok := d.GetOk("ipv6"); ok {
		if client.GetIpVersion(attr.(string)) != 6 {
			return fmt.Errorf("The IP address with UUID %v is not version 6", attr.(string))
		}
		ip := gsclient.ServerCreateRequestIp{
			IpaddrUuid: attr.(string),
		}
		requestBody.Relations.PublicIps = append(requestBody.Relations.PublicIps, ip)
	}

	//Add public network if we have an IP
	if len(requestBody.Relations.PublicIps) > 0 {
		publicNetwork, err := client.GetNetworkPublic()
		if err != nil {
			return err
		}
		network := gsclient.ServerCreateRequestNetwork{
			NetworkUuid: publicNetwork.Properties.ObjectUuid,
		}
		requestBody.Relations.Networks = append(requestBody.Relations.Networks, network)
	}

	if attr, ok := d.GetOk("networks"); ok {
		for index, value := range attr.(*schema.Set).List() {
			network := gsclient.ServerCreateRequestNetwork{
				NetworkUuid: value.(string),
			}
			if index == 0 {
				network.BootDevice = true
			}
			requestBody.Relations.Networks = append(requestBody.Relations.Networks, network)
		}
	}

	response, err := client.CreateServer(requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", requestBody.Name, err)
	}

	d.SetId(response.ServerUuid)

	log.Printf("[DEBUG] The id for %s has been set to: %v", requestBody.Name, response.ServerUuid)

	power := d.Get("power").(bool)
	if power {
		client.StartServer(d.Id())
	}

	return resourceGridscaleServerRead(d, meta)
}

func resourceGridscaleServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Id()
	err := client.StopServer(id)
	if err != nil {
		return err
	}
	err = client.DeleteServer(id)

	return err
}

func resourceGridscaleServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	var err error

	if d.HasChange("power") {
		_, change := d.GetChange("power")
		power := change.(bool)
		if power {
			err = client.StartServer(d.Id())
		} else {
			err = client.ShutdownServer(d.Id())
		}
		if err != nil {
			return err
		}
	}

	requestBody := gsclient.ServerUpdateRequest{
		Name:            d.Get("name").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          d.Get("labels").(*schema.Set).List(),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
	}

	err = client.UpdateServer(d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleServerRead(d, meta)

}
