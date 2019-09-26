package gridscale

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/nvthongswansea/gsclient-go"
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
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, alg := range loadbalancerAlgs {
						if v.(string) == alg {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid loadbalancer algorithm. Valid loadbalancer algorithms are: %v", v.(string), strings.Join(loadbalancerAlgs, ",")))
					}
					return
				},
			},
			"forwarding_rule": {
				Type:        schema.TypeSet,
				Description: "List of forwarding rules for the Load balancer.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"letsencrypt_ssl": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  nil,
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
			"backend_server": {
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
				Default:     "active",
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

func resourceGridscaleLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.LoadBalancerCreateRequest{
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		RedirectHTTPToHTTPS: d.Get("redirect_http_to_https").(bool),
		ListenIPv4UUID:      d.Get("listen_ipv4_uuid").(string),
		ListenIPv6UUID:      d.Get("listen_ipv6_uuid").(string),
		Labels:              convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	if d.Get("algorithm").(string) == "roundrobin" {
		requestBody.Algorithm = gsclient.LoadbalancerRoundrobinAlg
	} else if d.Get("algorithm").(string) == "leastconn" {
		requestBody.Algorithm = gsclient.LoadbalancerLeastConnAlg
	}

	if backendServers, ok := d.GetOk("backend_server"); ok {
		requestBody.BackendServers = expandLoadbalancerBackendServers(backendServers)
	}
	if forwardingRules, ok := d.GetOk("forwarding_rule"); ok {
		requestBody.ForwardingRules = expandLoadbalancerForwardingRules(forwardingRules)
	}
	response, err := client.CreateLoadBalancer(emptyCtx, requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for loadbalancer (%s) to be created: %s", requestBody.Name, err)
	}
	d.SetId(response.ObjectUUID)
	return resourceGridscaleLoadBalancerRead(d, meta)
}

func resourceGridscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	loadbalancer, err := client.GetLoadBalancer(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", loadbalancer.Properties.Name)
	d.Set("algorithm", loadbalancer.Properties.Algorithm)
	d.Set("status", loadbalancer.Properties.Status)
	d.Set("redirect_http_to_https", loadbalancer.Properties.RedirectHTTPToHTTPS)
	d.Set("listen_ipv4_uuid", loadbalancer.Properties.ListenIPv4UUID)
	d.Set("listen_ipv6_uuid", loadbalancer.Properties.ListenIPv6UUID)

	if err = d.Set("forwarding_rule", flattenLoadbalancerForwardingRules(loadbalancer.Properties.ForwardingRules)); err != nil {
		return fmt.Errorf("Error setting ForwardingRules: %v", err)
	}

	if err = d.Set("backend_server", flattenLoadbalancerBackendServers(loadbalancer.Properties.BackendServers)); err != nil {
		return fmt.Errorf("Error setting BackendServers: %v", err)
	}

	if err = d.Set("labels", loadbalancer.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.LoadBalancerUpdateRequest{
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		RedirectHTTPToHTTPS: d.Get("redirect_http_to_https").(bool),
		ListenIPv4UUID:      d.Get("listen_ipv4_uuid").(string),
		ListenIPv6UUID:      d.Get("listen_ipv6_uuid").(string),
		Labels:              convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	if d.Get("algorithm").(string) == "roundrobin" {
		requestBody.Algorithm = gsclient.LoadbalancerRoundrobinAlg
	} else if d.Get("algorithm").(string) == "leastconn" {
		requestBody.Algorithm = gsclient.LoadbalancerLeastConnAlg
	}

	if backendServers, ok := d.GetOk("backend_server"); ok {
		requestBody.BackendServers = expandLoadbalancerBackendServers(backendServers)
	}
	if forwardingRules, ok := d.GetOk("forwarding_rule"); ok {
		requestBody.ForwardingRules = expandLoadbalancerForwardingRules(forwardingRules)
	}
	err := client.UpdateLoadBalancer(emptyCtx, d.Id(), requestBody)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for loadbalancer (%s) to be updated: %s", requestBody.Name, err)
	}
	return resourceGridscaleLoadBalancerRead(d, meta)
}

func resourceGridscaleLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Id()
	return client.DeleteLoadBalancer(emptyCtx, id)
}

func expandLoadbalancerBackendServers(backendServers interface{}) []gsclient.BackendServer {
	tempBackendServers := []gsclient.BackendServer{}
	for _, value := range backendServers.(*schema.Set).List() {
		server := value.(map[string]interface{})
		backendServer := gsclient.BackendServer{
			Weight: server["weight"].(int),
			Host:   server["host"].(string),
		}
		tempBackendServers = append(tempBackendServers, backendServer)
	}
	return tempBackendServers
}

func expandLoadbalancerForwardingRules(forwardingRules interface{}) []gsclient.ForwardingRule {
	tempForwardingRules := []gsclient.ForwardingRule{}
	for _, value := range forwardingRules.(*schema.Set).List() {
		rule := value.(map[string]interface{})
		forwardingRule := gsclient.ForwardingRule{
			LetsencryptSSL: nil,
			ListenPort:     rule["listen_port"].(int),
			Mode:           rule["mode"].(string),
			TargetPort:     rule["target_port"].(int),
		}
		tempForwardingRules = append(tempForwardingRules, forwardingRule)
	}
	return tempForwardingRules
}

func flattenLoadbalancerForwardingRules(forwardingRules []gsclient.ForwardingRule) []interface{} {
	tempForwardingRules := make([]interface{}, 0)

	if forwardingRules != nil {
		for _, value := range forwardingRules {
			forwardingRule := map[string]interface{}{
				"letsencrypt_ssl": value.LetsencryptSSL,
				"listen_port":     value.ListenPort,
				"mode":            value.Mode,
				"target_port":     value.TargetPort,
			}
			tempForwardingRules = append(tempForwardingRules, forwardingRule)
		}
	}
	return tempForwardingRules
}

func flattenLoadbalancerBackendServers(backendServers []gsclient.BackendServer) []interface{} {
	tempBackendServers := make([]interface{}, 0)
	if backendServers != nil {
		for _, value := range backendServers {
			backendServer := map[string]interface{}{
				"weight": value.Weight,
				"host":   value.Host,
			}
			tempBackendServers = append(tempBackendServers, backendServer)
		}
	}
	return tempBackendServers
}
