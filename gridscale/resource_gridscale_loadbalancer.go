package gridscale

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleLoadBalancerCreate,
		Read:   resourceGridscaleLoadBalancerRead,
		Delete: resourceGridscaleLoadBalancerDelete,
		Update: resourceGridscaleLoadBalancerUpdate,
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
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Description: "The algorithm used to process requests. Accepted values: roundrobin/leastconn.",
				Required:    true,
			},
			"forwarding_rules": {
				Type:        schema.TypeSet,
				Description: "List of forwarding rules for the Load balancer.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"letsencrypt_ssl": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"listen_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"backend_servers": {
				Type:        schema.TypeSet,
				Description: "List of backend servers.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  100,
						},
						"host": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status indicates the status of the object.",
				Optional:    true,
			},
			"redirect_http_to_https": {
				Type:        schema.TypeBool,
				Description: "Whether the Load balancer is forced to redirect requests from HTTP to HTTPS",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"listen_ipv4_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the IPv4 address the Load balancer will listen to for incoming requests.",
				Required:    true,
			},
			"listen_ipv6_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the IPv6 address the Load balancer will listen to for incoming requests.",
				Required:    true,
			},
		},
	}
}

func resourceGridscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	loadbalancer, err := client.GetLoadBalancer(d.Id())
	if err != nil {
		if requestError, ok := err.(*gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", loadbalancer.Properties.Name)
	d.Set("algorithm", loadbalancer.Properties.Algorithm)
	d.Set("forwarding_rules", loadbalancer.Properties.ForwardingRules)
	d.Set("backend_servers", loadbalancer.Properties.BackendServers)
	d.Set("status", loadbalancer.Properties.Status)
	d.Set("redirect_http_to_https", loadbalancer.Properties.RedirectHTTPToHTTPS)
	d.Set("listen_ipv4_uuid", loadbalancer.Properties.ListenIPv4Uuid)
	d.Set("listen_ipv6_uuid", loadbalancer.Properties.ListenIPv6Uuid)

	if err = d.Set("labels", loadbalancer.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.LoadBalancerCreateRequest{
		Name:      d.Get("name").(string),
		Algorithm: d.Get("algorithm").(string),
		Status:    d.Get("status").(string),
		//ForwardingRules:     d.Get("forwarding_rules").(*schema.Set),
		//BackendServers:      d.Get("backend_servers").(*schema.Set),
		RedirectHTTPToHTTPS: d.Get("redirect_http_to_https").(bool),
		ListenIPv4Uuid:      d.Get("listen_ipv4_uuid").(string),
		ListenIPv6Uuid:      d.Get("listen_ipv6_uuid").(string),
	}
	if attr, ok := d.GetOk("forwarding_rules"); ok && len(attr.([]interface{})) > 0 {
		requestBody.ForwardingRules = expandLoadbalancerForwardingRules(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("backend_servers"); ok && len(attr.([]interface{})) > 0 {
		requestBody.BackendServers = expandLoadbalancerBackendServers(attr.(*schema.Set).List())
	}
	_, err := client.CreateLoadBalancer(requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for loadbalancer (%s) to be created: %s", requestBody.Name, err)
	}

	return resourceGridscaleLoadBalancerRead(d, meta)
}

func resourceGridscaleLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	return client.DeleteLoadBalancer(d.Id())
}

func resourceGridscaleLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.LoadBalancerUpdateRequest{
		Name:      d.Get("name").(string),
		Algorithm: d.Get("algorithm").(string),
		//ForwardingRules:     d.Get("forwarding_rules").(*schema.Set),
		//BackendServers:      d.Get("backend_servers").(*schema.Set),
		Status:              d.Get("status").(string),
		RedirectHTTPToHTTPS: d.Get("redirect_http_to_https").(bool),
		ListenIPv4Uuid:      d.Get("listen_ipv4_uuid").(string),
		ListenIPv6Uuid:      d.Get("listen_ipv6_uuid").(string),
	}
	if attr, ok := d.GetOk("forwarding_rules"); ok && len(attr.([]interface{})) > 0 {
		requestBody.ForwardingRules = expandLoadbalancerForwardingRules(attr.(*schema.Set).List())
	}

	if attr, ok := d.GetOk("backend_servers"); ok && len(attr.([]interface{})) > 0 {
		requestBody.BackendServers = expandLoadbalancerBackendServers(attr.(*schema.Set).List())
	}
	err := client.UpdateLoadBalancer(d.Id(), requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", requestBody.Name, err)
	}
	return resourceGridscaleLoadBalancerRead(d, meta)
}

func expandLoadbalancerBackendServers(l []interface{}) []gsclient.BackendServer {
	backendServers := []gsclient.BackendServer{}

	for _, mRaw := range l {
		if mRaw == nil {
			continue
		}

		m := mRaw.(map[string]interface{})

		backendServer := gsclient.BackendServer{
			Weight: m["weight"].(int),
			Host:   m["host"].(string),
		}

		backendServers = append(backendServers, backendServer)
	}

	return backendServers
}

func expandLoadbalancerForwardingRules(l []interface{}) []gsclient.ForwardingRule {
	forwardingRules := []gsclient.ForwardingRule{}

	for _, mRaw := range l {
		if mRaw == nil {
			continue
		}

		m := mRaw.(map[string]interface{})

		forwardingRule := gsclient.ForwardingRule{
			LetsencryptSSL: m["letsencrypt_ssl"].(*string),
			ListenPort:     m["listen_port"].(int),
			Mode:           m["mode"].(string),
			TargetPort:     m["target_port"].(int),
		}

		forwardingRules = append(forwardingRules, forwardingRule)
	}

	return forwardingRules
}
