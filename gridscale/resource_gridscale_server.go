package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	resource_dependency_crud "github.com/terraform-providers/terraform-provider-gridscale/gridscale/resource-dependency-crud"
	"github.com/terraform-providers/terraform-provider-gridscale/gridscale/service-query"
	"log"
	"strings"
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
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"memory": {
				Type:         schema.TypeInt,
				Description:  "The amount of server memory in GB.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cores": {
				Type:         schema.TypeInt,
				Description:  "The number of server cores.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
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
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, profile := range hardwareProfiles {
						if v.(string) == profile {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid hardware profile. Valid hardware profiles are: %v", v.(string), strings.Join(storageTypes, ",")))
					}
					return
				},
			},
			"storage": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 8,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"controller": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lun": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"license_product_no": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_used_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 7,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_template_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"partner_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ordering": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ipv4": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isoimage": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"power": {
				Type:        schema.TypeBool,
				Description: "The number of server cores.",
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, profile := range availabilityZones {
						if v.(string) == profile {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid availability zone. Valid availability zones are: %v", v.(string), strings.Join(availabilityZones, ",")))
					}
					return
				},
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
		if requestError, ok := err.(gsclient.RequestError); ok {
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
	d.Set("location_uuid", server.Properties.LocationUUID)
	d.Set("power", server.Properties.Power)
	d.Set("current_price", server.Properties.CurrentPrice)
	d.Set("availability_zone", server.Properties.AvailablityZone)
	d.Set("auto_recovery", server.Properties.AutoRecovery)
	d.Set("console_token", server.Properties.ConsoleToken)
	d.Set("legacy", server.Properties.Legacy)
	d.Set("usage_in_minutes_memory", server.Properties.UsageInMinutesMemory)
	d.Set("usage_in_minutes_cores", server.Properties.UsageInMinutesCores)
	if err = d.Set("labels", server.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}
	//Get storages
	storages := make([]interface{}, 0)
	for _, value := range server.Properties.Relations.Storages {
		storage := map[string]interface{}{
			"object_uuid":        value.ObjectUUID,
			"bootdevice":         value.BootDevice,
			"create_time":        value.CreateTime,
			"controller":         value.Controller,
			"target":             value.Target,
			"lun":                value.Lun,
			"license_product_no": value.LicenseProductNo,
			"bus":                value.Bus,
			"object_name":        value.ObjectName,
			"storage_type":       value.StorageType,
			"last_used_template": value.LastUsedTemplate,
			"capacity":           value.Capacity,
		}
		storages = append(storages, storage)
	}
	if err = d.Set("storage", storages); err != nil {
		return fmt.Errorf("Error setting storage: %v", err)
	}

	//Get networks
	networks := make([]interface{}, 0)
	for _, value := range server.Properties.Relations.Networks {
		if !value.PublicNet {
			network := map[string]interface{}{
				"object_uuid":            value.ObjectUUID,
				"bootdevice":             value.BootDevice,
				"create_time":            value.CreateTime,
				"mac":                    value.Mac,
				"firewall":               value.Firewall,
				"firewall_template_uuid": value.FirewallTemplateUUID,
				"object_name":            value.ObjectName,
				"network_type":           value.NetworkType,
				"ordering":               value.Ordering,
			}
			networks = append(networks, network)
		}
	}
	if err = d.Set("network", networks); err != nil {
		return fmt.Errorf("Error setting network: %v", err)
	}
	//Get IP addresses
	var ipv4, ipv6 string
	for _, ip := range server.Properties.Relations.PublicIPs {
		if ip.Family == 4 {
			ipv4 = ip.ObjectUUID
		}
		if ip.Family == 6 {
			ipv6 = ip.ObjectUUID
		}
	}
	d.Set("ipv4", ipv4)
	d.Set("ipv6", ipv6)
	//Get the ISO image, there can only be one attached to a server but it is in a list anyway
	d.Set("isoimage", "")
	for _, isoimage := range server.Properties.Relations.IsoImages {
		d.Set("isoimage", isoimage.ObjectUUID)
	}
	return nil
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := resource_dependency_crud.NewServerDepClient(meta.(*gsclient.Client), d)
	gsc := client.GetGSClient()
	requestBody := gsclient.ServerCreateRequest{
		Name:            d.Get("name").(string),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
		LocationUUID:    d.Get("location_uuid").(string),
		HardwareProfile: d.Get("hardware_profile").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	response, err := gsc.CreateServer(requestBody)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", requestBody.Name, err)
	}
	d.SetId(response.ObjectUUID)
	log.Printf("[DEBUG] The id for %s has been set to: %v", requestBody.Name, response.ObjectUUID)
	err = client.LinkStorages()
	if err != nil {
		return err
	}
	err = client.LinkIPv4()
	if err != nil {
		return err
	}
	err = client.LinkIPv6()
	if err != nil {
		return err
	}

	//Add public network if we have an IP
	_, hasIPv4 := d.GetOk("ipv4")
	_, hasIPv6 := d.GetOk("ipv6")
	if hasIPv4 || hasIPv6 {
		err = client.LinkNetworks(true)
		if err != nil {
			return err
		}
	}
	err = client.LinkNetworks(false)
	if err != nil {
		return err
	}

	//Set the power state if needed
	power := d.Get("power").(bool)
	if power {
		gsc.StartServer(d.Id())
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
	client := resource_dependency_crud.NewServerDepClient(meta.(*gsclient.Client), d)
	gsc := client.GetGSClient()
	var err error
	//The ShutdownServer command will check if the server is running and shut it down if it is running, so no extra checks are needed here
	if client.IsShutdownRequired() {
		err = gsc.ShutdownServer(d.Id())
		if err != nil {
			return err
		}
	}
	requestBody := gsclient.ServerUpdateRequest{
		Name:            d.Get("name").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          convSOStrings(d.Get("labels").(*schema.Set).List()),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
	}
	//Execute the update request
	err = gsc.UpdateServer(d.Id(), requestBody)
	if err != nil {
		return err
	}

	//Link/unlink isoimages
	err = client.UpdateISOImageRel()
	if err != nil {
		return err
	}
	//Link/Unlink ip addresses
	needsPublicNetwork, err := client.UpdateIPv4Rel()
	if err != nil {
		return err
	}
	needsPublicNetwork, err = client.UpdateIPv6Rel()
	if err != nil {
		return err
	}
	//Disconnect from the public network if there is no longer and IP
	if (d.HasChange("ipv6") || d.HasChange("ipv4")) && d.Get("ipv6").(string) == "" && d.Get("ipv4").(string) == "" {
		err = client.UpdatePublicNetworkRel(false)
		if err != nil {
			return nil
		}
	}
	//Connect to the public network if an IP was added
	if (d.HasChange("ipv6") || d.HasChange("ipv4")) && needsPublicNetwork {
		err = client.UpdatePublicNetworkRel(true)
		if err != nil {
			return nil
		}
	}

	//Link/unlink networks
	err = client.UpdateOtherNetworkRel()
	if err != nil {
		return err
	}

	//Link/unlink storages
	err = client.UpdateStorageRel()
	if err != nil {
		return nil
	}

	// Make sure the server in is the expected power state.
	// The StartServer and ShutdownServer functions do a check to see if the server isn't already running, so we don't need to do that here.
	if d.Get("power").(bool) {
		err = gsc.StartServer(d.Id())
	} else {
		err = gsc.ShutdownServer(d.Id())
	}
	if err != nil {
		return err
	}
	err = service_query.BlockProvisoning(gsc, service_query.ServerService, d.Id())
	if err != nil {
		return err
	}
	return resourceGridscaleServerRead(d, meta)
}
