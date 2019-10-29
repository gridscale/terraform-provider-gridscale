package gridscale

import (
	"context"
	"fmt"
	relation_manager "github.com/terraform-providers/terraform-provider-gridscale/gridscale/relation-manager"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go"
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
						errors = append(errors, fmt.Errorf("%v is not a valid hardware profile. Valid hardware profiles are: %v", v.(string), strings.Join(hardwareProfiles, ",")))
					}
					return
				},
			},
			"storage": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    8,
				Description: `A list of storages attached to the server. The first storage in the list is always set as the boot storage of the server.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_uuid": {
							Type:     schema.TypeString,
							Required: true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v4_out": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v6_in": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"rules_v6_out": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: getFirewallRuleCommonSchema(),
							},
						},
						"firewall_template_uuid": {
							Type:     schema.TypeString,
							Optional: true,
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

//getFirewallRuleCommonSchema returns schema for custom firewall rules.
//**Note: Every time `getFirewallRuleCommonSchema()` is called,
//all `*schema.Schema` in `map[string]*schema.Schema` are different.
func getFirewallRuleCommonSchema() map[string]*schema.Schema {
	commonSchema := map[string]schema.Schema{
		"order": {
			Type: schema.TypeInt,
			Description: `The order at which the firewall will compare packets against its rules, 
a packet will be compared against the first rule, it will either allow it to pass or block it 
and it won t be matched against any other rules. However, if it does no match the rule, 
then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.`,
			Required: true,
		},
		"action": {
			Type:        schema.TypeString,
			Description: "This defines what the firewall will do. Either accept or drop.",
			Required:    true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				valid := false
				for _, action := range firewallActionTypes {
					if v.(string) == action {
						valid = true
						break
					}
				}
				if !valid {
					errors = append(errors, fmt.Errorf("%v is not a valid firewall action. Valid firewall actions are: %v", v.(string), strings.Join(firewallActionTypes, ",")))
				}
				return
			},
		},
		"protocol": {
			Type:        schema.TypeString,
			Description: "Either 'udp' or 'tcp'",
			Optional:    true,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				valid := false
				for _, prot := range firewallRuleProtocols {
					if v.(string) == prot {
						valid = true
						break
					}
				}
				if !valid {
					errors = append(errors, fmt.Errorf("%v is not a valid protocol. Valid protocols are: %v", v.(string), strings.Join(firewallRuleProtocols, ",")))
				}
				return
			},
		},
		"dst_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A Number between 1 and 65535, port ranges are seperated by a colon for FTP",
		},
		"src_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A Number between 1 and 65535, port ranges are seperated by a colon for FTP",
		},
		"src_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.",
		},
		"dst_cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then all IPs have access to this service.",
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comment.",
		},
	}
	//Every time `getFirewallRuleCommonSchema()` is called,
	//all `*schema.Schema` in `map[string]*schema.Schema` have to be different.
	//So that new `*schema.Schema` are created.
	schemaWithPointers := make(map[string]*schema.Schema)
	for k, v := range commonSchema {
		newVal := new(schema.Schema)
		*newVal = v
		schemaWithPointers[k] = newVal
	}
	return schemaWithPointers
}

func resourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	server, err := client.GetServer(emptyCtx, d.Id())
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
	d.Set("availability_zone", server.Properties.AvailabilityZone)
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
	if err = d.Set("storage", storages); err != nil {
		return fmt.Errorf("Error setting storage: %v", err)
	}

	//Get networks
	networks := readServerNetworkRels(server.Properties.Relations.Networks)
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

//readServerNetworkRels extract relationships between server and networks
func readServerNetworkRels(serverNetRels []gsclient.ServerNetworkRelationProperties) []interface{} {
	networks := make([]interface{}, 0)
	for _, rel := range serverNetRels {
		network := map[string]interface{}{
			"object_uuid":            rel.ObjectUUID,
			"bootdevice":             rel.BootDevice,
			"create_time":            rel.CreateTime.String(),
			"mac":                    rel.Mac,
			"firewall_template_uuid": rel.FirewallTemplateUUID,
			"object_name":            rel.ObjectName,
			"network_type":           rel.NetworkType,
			"ordering":               rel.Ordering,
		}
		//Init all types of firewall rule
		v4InRuleProps := make([]interface{}, 0)
		v4OutRuleProps := make([]interface{}, 0)
		v6InRuleProps := make([]interface{}, 0)
		v6OutRuleProps := make([]interface{}, 0)

		//Add rules of type rules_v4_in
		for _, props := range rel.Firewall.RulesV4In {
			v4InRuleProp := flattenFirewallRuleProperties(props)
			v4InRuleProps = append(v4InRuleProps, v4InRuleProp)
		}
		network["rules_v4_in"] = v4InRuleProps

		//Add rules of type rules_v4_out
		for _, props := range rel.Firewall.RulesV4Out {
			v4OutRuleProp := flattenFirewallRuleProperties(props)
			v4OutRuleProps = append(v4OutRuleProps, v4OutRuleProp)
		}
		network["rules_v4_out"] = v4OutRuleProps

		//Add rules of type rules_v6_in
		for _, props := range rel.Firewall.RulesV6In {
			v6InRuleProp := flattenFirewallRuleProperties(props)
			v6InRuleProps = append(v6InRuleProps, v6InRuleProp)
		}
		network["rules_v6_in"] = v6InRuleProps

		//Add rules of type rules_v6_out
		for _, props := range rel.Firewall.RulesV6Out {
			v6OutRuleProp := flattenFirewallRuleProperties(props)
			v6OutRuleProps = append(v6OutRuleProps, v6OutRuleProp)
		}
		network["rules_v6_out"] = v6OutRuleProps

		networks = append(networks, network)
	}
	return networks
}

//flattenFirewallRuleProperties converts variable of type gsclient.FirewallRuleProperties to
//map[string]interface{}
func flattenFirewallRuleProperties(props gsclient.FirewallRuleProperties) map[string]interface{} {
	return map[string]interface{}{
		"order":    props.Order,
		"action":   props.Action,
		"protocol": props.Protocol,
		"dst_port": props.DstPort,
		"src_port": props.SrcPort,
		"src_cidr": props.SrcCidr,
		"dst_cidr": props.DstCidr,
		"comment":  props.Comment,
	}
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	gsc := meta.(*gsclient.Client)
	serverRelMan := relation_manager.NewServerRelationManger(gsc, d)
	requestBody := gsclient.ServerCreateRequest{
		Name:            d.Get("name").(string),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
		LocationUUID:    d.Get("location_uuid").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	profile := d.Get("hardware_profile").(string)
	if profile == "legacy" {
		requestBody.HardwareProfile = gsclient.LegacyServerHardware
	} else if profile == "nested" {
		requestBody.HardwareProfile = gsclient.NestedServerHardware
	} else if profile == "cisco_csr" {
		requestBody.HardwareProfile = gsclient.CiscoCSRServerHardware
	} else if profile == "sophos_utm" {
		requestBody.HardwareProfile = gsclient.SophosUTMServerHardware
	} else if profile == "f5_bigip" {
		requestBody.HardwareProfile = gsclient.F5BigipServerHardware
	} else if profile == "q35" {
		requestBody.HardwareProfile = gsclient.Q35ServerHardware
	} else if profile == "q35_nested" {
		requestBody.HardwareProfile = gsclient.Q35NestedServerHardware
	} else {
		requestBody.HardwareProfile = gsclient.DefaultServerHardware
	}
	response, err := gsc.CreateServer(emptyCtx, requestBody)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", requestBody.Name, err)
	}
	d.SetId(response.ServerUUID)
	log.Printf("[DEBUG] The id for %s has been set to: %v", requestBody.Name, response.ServerUUID)

	//Add server power status to serverPowerStateList
	err = serverPowerStateList.addServer(d.Id())
	if err != nil {
		return err
	}

	//Link storages
	err = serverRelMan.LinkStorages(emptyCtx)
	if err != nil {
		return err
	}

	//Link IPv4
	err = serverRelMan.LinkIPv4(emptyCtx)
	if err != nil {
		return err
	}

	//Link IPv6
	err = serverRelMan.LinkIPv6(emptyCtx)
	if err != nil {
		return err
	}

	//Link ISO-Image
	err = serverRelMan.LinkISOImage(emptyCtx)
	if err != nil {
		return err
	}

	//Link networks
	err = serverRelMan.LinkNetworks(emptyCtx)
	if err != nil {
		return err
	}

	//Set the power state if needed
	power := d.Get("power").(bool)
	if power {
		err = serverPowerStateList.startServerSynchronously(emptyCtx, gsc, d.Id())
		if err != nil {
			return err
		}
	}

	return resourceGridscaleServerRead(d, meta)
}

func resourceGridscaleServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Id()
	//Stop the server synchronously
	err := serverPowerStateList.shutdownServerSynchronously(emptyCtx, client, id)
	if err != nil {
		return err
	}
	//Delete the server
	err = client.DeleteServer(emptyCtx, id)
	if err != nil {
		return err
	}
	//remove server power state from serverPowerStateList
	err = serverPowerStateList.removeServer(d.Id())
	return err
}

func resourceGridscaleServerUpdate(d *schema.ResourceData, meta interface{}) error {
	gsc := meta.(*gsclient.Client)
	serverDepClient := relation_manager.NewServerRelationManger(gsc, d)
	shutdownRequired := serverDepClient.IsShutdownRequired(emptyCtx)
	var err error
	requestBody := gsclient.ServerUpdateRequest{
		Name:            d.Get("name").(string),
		AvailablityZone: d.Get("availability_zone").(string),
		Labels:          convSOStrings(d.Get("labels").(*schema.Set).List()),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
	}

	if shutdownRequired {
		updateSequence := func(ctx context.Context) error {
			//Execute the update request
			err = gsc.UpdateServer(ctx, d.Id(), requestBody)
			if err != nil {
				return err
			}

			//Update relationship between the server and IP addresses
			err = serverDepClient.UpdateIPv4Rel(ctx)
			if err != nil {
				return err
			}
			err = serverDepClient.UpdateIPv6Rel(ctx)
			if err != nil {
				return err
			}

			//Update relationship between the server and networks
			err = serverDepClient.UpdateNetworksRel(ctx)
			if err != nil {
				return err
			}

			//Update relationship between the server and storages
			err = serverDepClient.UpdateStoragesRel(ctx)
			return err
		}
		err = serverPowerStateList.runActionRequireServerOff(emptyCtx, gsc, d.Id(), updateSequence)
		if err != nil {
			return err
		}
	} else {
		//Execute the update request
		err = gsc.UpdateServer(emptyCtx, d.Id(), requestBody)
		if err != nil {
			return err
		}
	}

	//Update relationship between the server and an ISO-image
	err = serverDepClient.UpdateISOImageRel(emptyCtx)
	if err != nil {
		return err
	}

	// Make sure the server in is the expected power state.
	// The StartServer and ShutdownServer functions do a check to see if the server isn't already running, so we don't need to do that here.
	if d.Get("power").(bool) {
		err = serverPowerStateList.startServerSynchronously(emptyCtx, gsc, d.Id())
		if err != nil {
			return err
		}
	} else {
		err = serverPowerStateList.shutdownServerSynchronously(emptyCtx, gsc, d.Id())
		if err != nil {
			return err
		}
	}
	return resourceGridscaleServerRead(d, meta)
}
