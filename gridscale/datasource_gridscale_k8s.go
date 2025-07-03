package gridscale

import (
	"context"
	"fmt"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGridscaleK8s() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleK8sRead,
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
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "K8s config data",
				Computed:    true,
				Sensitive:   true,
			},
			"k8s_private_network_uuid": {
				Type:        schema.TypeString,
				Description: "Private network UUID which k8s nodes are attached to. It can be used to attach other PaaS/VMs.",
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

func dataSourceGridscaleK8sRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	d.SetId(id)
	errorPrefix := fmt.Sprintf("read k8s (%s) resource -", id)
	paas, err := client.GetPaaSService(context.Background(), id)

	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	creds := paas.Properties.Credentials
	if len(creds) > 0 {
		// if expiration_time of kubeconfig is reached, renew it and get new kubeconfig
		if creds[0].ExpirationTime.Before(time.Now()) {
			err = client.RenewK8sCredentials(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("%s error renewing k8s kubeconfig: %v", errorPrefix, err)
			}
			paas, err = client.GetPaaSService(context.Background(), d.Id())
			if err != nil {
				return fmt.Errorf("%s error: %v", errorPrefix, err)
			}
			creds = paas.Properties.Credentials
		}
		if err = d.Set("kubeconfig", creds[0].KubeConfig); err != nil {
			return fmt.Errorf("%s error setting kubeconfig: %v", errorPrefix, err)
		}
	}

	if err = d.Set("name", paas.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("kubeconfig", creds[0].KubeConfig); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("labels", paas.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	// find the network with the label that matches the k8s label
	k8sLabel := fmt.Sprintf("%s%s", k8sLabelPrefix, d.Id())
NETWORK_LOOOP:
	for _, network := range networks {
		for _, label := range network.Properties.Labels {
			if label == k8sLabel {
				if err = d.Set("k8s_private_network_uuid", network.Properties.ObjectUUID); err != nil {
					return fmt.Errorf("%s error setting k8s_private_network_uuid: %v", errorPrefix, err)
				}
				break NETWORK_LOOOP
			}
		}
	}
	return nil
}
