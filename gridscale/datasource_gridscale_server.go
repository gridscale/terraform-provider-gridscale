package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleServerRead,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
			},
			"memory": {
				Type:        schema.TypeInt,
				Description: "The amount of server memory in GB.",
				Computed:    true,
			},
			"cores": {
				Type:        schema.TypeInt,
				Description: "The number of server cores.",
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Computed:    true,
			},
			"hardware_profile": {
				Type:        schema.TypeString,
				Description: "The number of server cores.",
				Computed:    true,
			},
			"storage": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `A list of storages attached to the server. The first storage in the list is always set as the boot storage of the server.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Computed: true,
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
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootdevice": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"object_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules_v4_in": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v4_out": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v6_in": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v6_out": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"firewall_template_uuid": {
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
				Computed: true,
			},
			"ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isoimage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"power": {
				Type:        schema.TypeBool,
				Description: "The number of server cores.",
				Computed:    true,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"auto_recovery": {
				Type:        schema.TypeBool,
				Description: "If the server should be auto-started in case of a failure (default=true).",
				Computed:    true,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Description: "Defines which Availability-Zone the Server is placed.",
				Computed:    true,
			},
			"console_token": {
				Type:        schema.TypeString,
				Description: "The token used by the panel to open the websocket VNC connection to the server console.",
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"change_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read server (%s) datasource-", id)

	server, err := client.GetServer(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	props := server.Properties
	d.SetId(props.ObjectUUID)

	if err = d.Set("name", server.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("memory", server.Properties.Memory); err != nil {
		return fmt.Errorf("%s error setting memory: %v", errorPrefix, err)
	}
	if err = d.Set("cores", server.Properties.Cores); err != nil {
		return fmt.Errorf("%s error setting cores: %v", errorPrefix, err)
	}
	if err = d.Set("hardware_profile", server.Properties.HardwareProfile); err != nil {
		return fmt.Errorf("%s error setting hardware_profile: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", server.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("power", server.Properties.Power); err != nil {
		return fmt.Errorf("%s error setting power: %v", errorPrefix, err)
	}
	if err = d.Set("status", server.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", server.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", server.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", server.Properties.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}
	if err = d.Set("availability_zone", server.Properties.AvailabilityZone); err != nil {
		return fmt.Errorf("%s error setting availability_zone: %v", errorPrefix, err)
	}
	if err = d.Set("auto_recovery", server.Properties.AutoRecovery); err != nil {
		return fmt.Errorf("%s error setting auto_recovery: %v", errorPrefix, err)
	}
	if err = d.Set("console_token", server.Properties.ConsoleToken); err != nil {
		return fmt.Errorf("%s error setting console_token: %v", errorPrefix, err)
	}
	if err = d.Set("legacy", server.Properties.Legacy); err != nil {
		return fmt.Errorf("%s error setting legacy: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes_memory", server.Properties.UsageInMinutesMemory); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes_memory: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes_cores", server.Properties.UsageInMinutesCores); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes_cores: %v", errorPrefix, err)
	}

	if err = d.Set("labels", server.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	//Get storages
	storages := make([]interface{}, 0)
	for _, value := range server.Properties.Relations.Storages {
		storage := map[string]interface{}{
			"object_uuid":        value.ObjectUUID,
			"bootdevice":         value.BootDevice,
			"create_time":        value.CreateTime.String(),
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

	//Get networks
	networks := readServerNetworkRels(server.Properties.Relations.Networks)
	if err = d.Set("network", networks); err != nil {
		return fmt.Errorf("%s error setting network: %v", errorPrefix, err)
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
	if err = d.Set("ipv4", ipv4); err != nil {
		return fmt.Errorf("%s error setting ipv4: %v", errorPrefix, err)
	}
	if err = d.Set("ipv6", ipv6); err != nil {
		return fmt.Errorf("%s error setting ipv6: %v", errorPrefix, err)
	}

	//Get the ISO image, there can only be one attached to a server but it is in a list anyway
	for _, isoimage := range server.Properties.Relations.IsoImages {
		if err = d.Set("isoimage", isoimage.ObjectUUID); err != nil {
			return fmt.Errorf("%s error setting isoimage: %v", errorPrefix, err)
		}
	}

	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}
