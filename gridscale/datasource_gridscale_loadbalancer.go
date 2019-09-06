package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func dataSourceGridscaleLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceGridscaleLoadBalancerRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	loadbalancer, err := client.GetLoadBalancer(id)
	if err == nil {
		d.SetId(loadbalancer.Properties.ObjectUUID)
	}
	return err
}
