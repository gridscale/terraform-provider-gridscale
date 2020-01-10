package gridscale

import (
	"errors"
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceGridscaleFirewall() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGridscaleFirewallRead,
		Create: resourceGridscaleFirewallCreate,
		Update: resourceGridscaleFirewallUpdate,
		Delete: resourceGridscaleFirewallDelete,
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
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object",
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
			"private": {
				Type:        schema.TypeBool,
				Description: "The object is private, the value will be true. Otherwise the value will be false.",
				Computed:    true,
			},
			"network": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
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
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the Firewall.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
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

func resourceGridscaleFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	template, err := client.GetFirewall(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := template.Properties
	d.Set("name", props.Name)
	d.Set("location_name", props.LocationName)
	d.Set("status", props.Status)
	d.Set("private", props.Private)
	d.Set("create_time", props.CreateTime)
	d.Set("change_time", props.ChangeTime)
	d.Set("description", props.Description)

	//Get network relating to this firewall
	networks := make([]interface{}, 0)
	for _, value := range props.Relations.Networks {
		rule := map[string]interface{}{
			"network_uuid": value.NetworkUUID,
			"object_uuid":  value.ObjectUUID,
			"network_name": value.NetworkName,
			"object_name":  value.ObjectName,
			"create_time":  value.CreateTime,
		}
		networks = append(networks, rule)
	}
	if err = d.Set("network", networks); err != nil {
		return fmt.Errorf("Error setting network: %v", err)
	}

	//Get rules_v4_in
	rulesV4In := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV4In)
	if err = d.Set("rules_v4_in", rulesV4In); err != nil {
		return fmt.Errorf("Error setting rules_v4_in: %v", err)
	}

	//Get rules_v4_out
	rulesV4Out := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV4Out)
	if err = d.Set("rules_v4_out", rulesV4Out); err != nil {
		return fmt.Errorf("Error setting rules_v4_out: %v", err)
	}

	//Get rules_v6_in
	rulesV6In := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV6In)
	if err = d.Set("rules_v6_in", rulesV6In); err != nil {
		return fmt.Errorf("Error setting rules_v6_in: %v", err)
	}

	//Get rules_v6_out
	rulesV6Out := convFirewallRuleSliceToInterfaceSlice(props.Rules.RulesV6In)
	if err = d.Set("rules_v6_out", rulesV6Out); err != nil {
		return fmt.Errorf("Error setting rules_v6_out: %v", err)
	}

	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	var rulesV4In, rulesV4Out, rulesV6In, rulesV6Out []gsclient.FirewallRuleProperties
	//Get firewall rules from schema
	if attr, ok := d.GetOk("rules_v4_in"); ok {
		rulesV4In = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v4_out"); ok {
		rulesV4Out = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v6_in"); ok {
		rulesV6In = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v6_out"); ok {
		rulesV6Out = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	//at least one rules in firewall create request
	if len(rulesV4In) == 0 && len(rulesV4Out) == 0 && len(rulesV6In) == 0 && len(rulesV6Out) == 0 {
		return errors.New("at least 1 firewall rule in create request")
	}
	requestBody := gsclient.FirewallCreateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
		Rules: gsclient.FirewallRules{
			RulesV6In:  rulesV6In,
			RulesV6Out: rulesV6Out,
			RulesV4In:  rulesV4In,
			RulesV4Out: rulesV4Out,
		},
	}

	response, err := client.CreateFirewall(emptyCtx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new firewall has been set to %v", response.ObjectUUID)

	return resourceGridscaleFirewallRead(d, meta)
}

func resourceGridscaleFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	var rulesV4In, rulesV4Out, rulesV6In, rulesV6Out []gsclient.FirewallRuleProperties
	if attr, ok := d.GetOk("rules_v4_in"); ok {
		rulesV4In = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v4_out"); ok {
		rulesV4Out = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v6_in"); ok {
		rulesV6In = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("rules_v6_out"); ok {
		rulesV6Out = convInterfaceSliceToFirewallRulesSlice(attr.([]interface{}))
	}
	requestBody := gsclient.FirewallUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	if len(rulesV4In) == 0 && len(rulesV4Out) == 0 && len(rulesV6In) == 0 && len(rulesV6Out) == 0 {
		requestBody.Rules = nil
	} else {
		firewallRules := &gsclient.FirewallRules{
			RulesV6In:  rulesV6In,
			RulesV6Out: rulesV6Out,
			RulesV4In:  rulesV4In,
			RulesV4Out: rulesV4Out,
		}
		requestBody.Rules = firewallRules
	}
	err := client.UpdateFirewall(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleFirewallRead(d, meta)
}

func resourceGridscaleFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	return client.DeleteFirewall(emptyCtx, d.Id())
}

//convFirewallRuleSliceToInterfaceSlice converts slice of firewall rules to slice of interface
func convFirewallRuleSliceToInterfaceSlice(rules []gsclient.FirewallRuleProperties) []interface{} {
	res := make([]interface{}, 0)
	for _, value := range rules {
		rule := map[string]interface{}{
			"order":    value.Order,
			"action":   value.Action,
			"protocol": value.Protocol,
			"dst_port": value.DstPort,
			"src_port": value.SrcPort,
			"src_cidr": value.SrcCidr,
			"dst_cidr": value.DstCidr,
			"comment":  value.Comment,
		}
		res = append(res, rule)
	}
	return res
}

//convInterfaceSliceToFirewallRulesSlice converts slice of interface to slice of firewall rules
func convInterfaceSliceToFirewallRulesSlice(interfaceRules []interface{}) []gsclient.FirewallRuleProperties {
	var firewallRules []gsclient.FirewallRuleProperties
	for _, value := range interfaceRules {
		rule := value.(map[string]interface{})
		firewallRules = append(firewallRules, gsclient.FirewallRuleProperties{
			Protocol: rule["protocol"].(string),
			DstPort:  rule["dst_port"].(string),
			SrcPort:  rule["src_port"].(string),
			SrcCidr:  rule["src_cidr"].(string),
			Action:   rule["action"].(string),
			Comment:  rule["comment"].(string),
			DstCidr:  rule["dst_cidr"].(string),
			Order:    rule["order"].(int),
		})
	}
	return firewallRules
}
