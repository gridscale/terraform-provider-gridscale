package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGridscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters",
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "The location this object is placed.",
				Computed:    true,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Description: "The algorithm used to process requests. Accepted values: roundrobin/leastconn.",
				Computed:    true,
			},
			"forwarding_rule": {
				Type:        schema.TypeSet,
				Description: "List of forwarding rules for the Load balancer.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"letsencrypt_ssl": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A valid domain name that points to the loadbalancer's IP address.",
						},
						"certificate_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of a custom certificate.",
						},
						"listen_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the entry port of the load balancer.",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Supports HTTP and TCP mode. Valid values: http, tcp.",
						},
						"target_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the exit port that the load balancer uses to forward the traffic to the backend server.",
						},
					},
				},
			},
			"backend_server": {
				Type:        schema.TypeSet,
				Description: "List of backend servers.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object.",
				Computed:    true,
			},
			"redirect_http_to_https": {
				Type:        schema.TypeBool,
				Description: "Whether the Load balancer is forced to redirect requests from HTTP to HTTPS",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"listen_ipv4_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the IPv4 address the Load balancer will listen to for incoming requests.",
				Computed:    true,
			},
			"listen_ipv6_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the IPv6 address the Load balancer will listen to for incoming requests.",
				Computed:    true,
			},
		},
	}
}

func dataSourceGridscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read loadbalancer (%s) datasource-", id)
	loadbalancer, err := client.GetLoadBalancer(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(loadbalancer.Properties.ObjectUUID)
	if err = d.Set("name", loadbalancer.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("algorithm", loadbalancer.Properties.Algorithm); err != nil {
		return fmt.Errorf("%s error setting algorithm: %v", errorPrefix, err)
	}
	if err = d.Set("status", loadbalancer.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("redirect_http_to_https", loadbalancer.Properties.RedirectHTTPToHTTPS); err != nil {
		return fmt.Errorf("%s error setting redirect_http_to_https: %v", errorPrefix, err)
	}
	if err = d.Set("listen_ipv4_uuid", loadbalancer.Properties.ListenIPv4UUID); err != nil {
		return fmt.Errorf("%s error setting listen_ipv4_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("listen_ipv6_uuid", loadbalancer.Properties.ListenIPv6UUID); err != nil {
		return fmt.Errorf("%s error setting listen_ipv6_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("forwarding_rule", flattenLoadbalancerForwardingRules(loadbalancer.Properties.ForwardingRules)); err != nil {
		return fmt.Errorf("%s error setting forwarding_rule: %v", errorPrefix, err)
	}
	if err = d.Set("backend_server", flattenLoadbalancerBackendServers(loadbalancer.Properties.BackendServers)); err != nil {
		return fmt.Errorf("%s error setting BackendServers: %v", errorPrefix, err)
	}
	if err = d.Set("labels", loadbalancer.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting Labels: %v", errorPrefix, err)
	}

	return nil
}
