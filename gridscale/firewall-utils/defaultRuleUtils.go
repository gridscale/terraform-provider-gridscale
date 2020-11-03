package fwu

import "github.com/gridscale/gsclient-go/v3"

// AddDefaultFirewallInboundRules adds default fw rules
func AddDefaultFirewallInboundRules(rules []gsclient.FirewallRuleProperties, forIPv6 bool) []gsclient.FirewallRuleProperties {
	if len(rules) == 0 { // If no custom fw rules are added, no need to add default ones
		return rules
	}
	srcCidr := "0.0.0.0/0"
	DHCPDstPort := "67:68"
	DHCPComment := "DHCP IPv4"
	nextOrder := getNextFWRuleOrder(rules)
	if forIPv6 {
		srcCidr = "::/0"
		DHCPDstPort = "546:547"
		DHCPComment = "DHCP IPv6"
	}
	defaultInboundRules := []gsclient.FirewallRuleProperties{
		{
			Protocol: gsclient.UDPTransport,
			DstPort:  DHCPDstPort,
			SrcPort:  "",
			SrcCidr:  srcCidr,
			Action:   "accept",
			Comment:  DHCPComment,
			DstCidr:  "",
			Order:    nextOrder,
		},
		{
			Protocol: gsclient.TCPTransport,
			DstPort:  "32768:65535",
			SrcPort:  "",
			SrcCidr:  srcCidr,
			Action:   "accept",
			Comment:  "Highports TCP",
			DstCidr:  "",
			Order:    nextOrder + 1,
		},
		{
			Protocol: gsclient.UDPTransport,
			DstPort:  "32768:65535",
			SrcPort:  "",
			SrcCidr:  srcCidr,
			Action:   "accept",
			Comment:  "Highports UDP",
			DstCidr:  "",
			Order:    nextOrder + 2,
		},
		{
			Protocol: gsclient.UDPTransport,
			DstPort:  "1:65535",
			SrcPort:  "",
			SrcCidr:  srcCidr,
			Action:   "drop",
			Comment:  "Drop all other UDP",
			DstCidr:  "",
			Order:    nextOrder + 3,
		},
		{
			Protocol: gsclient.TCPTransport,
			DstPort:  "1:65535",
			SrcPort:  "",
			SrcCidr:  srcCidr,
			Action:   "drop",
			Comment:  "Drop all other TCP",
			DstCidr:  "",
			Order:    nextOrder + 4,
		},
	}
	rules = append(rules, defaultInboundRules...)
	return rules
}

func getNextFWRuleOrder(rules []gsclient.FirewallRuleProperties) int {
	var max int
	for _, v := range rules {
		if max < v.Order {
			max = v.Order
		}
	}
	return max + 1
}

// RemoveDefaultFirewallInboundRules removes default fw rules
// It is used when we don't want to display the default fw rules in tf
func RemoveDefaultFirewallInboundRules(rules []gsclient.FirewallRuleProperties) []gsclient.FirewallRuleProperties {
	defaultRulesNames := map[string]int{
		"DHCP IPv4":          0,
		"DHCP IPv6":          0,
		"Highports TCP":      0,
		"Highports UDP":      0,
		"Drop all other UDP": 0,
		"Drop all other TCP": 0,
	}
	for i := 0; i < len(rules); i++ {
		if _, ok := defaultRulesNames[rules[i].Comment]; ok {
			rules = append(rules[:i], rules[i+1:]...)
			i-- // As we've just deleted rules[i], we'd to redo that index
		}
	}
	return rules
}
