package gridscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const k8sTemplateFlavourName = "kubernetes"

type k8sValidationOpt int

const (
	k8sReleaseValidationOpt k8sValidationOpt = iota
	k8sNodeCountValidationOpt
	k8sCoreCountValidationOpt
	k8sMemoryValidationOpt
	k8sStorageValidationOpt
)

func resourceGridscaleK8s() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleK8sCreate,
		Read:   resourceGridscaleK8sRead,
		Delete: resourceGridscaleK8sDelete,
		Update: resourceGridscaleK8sUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "K8s config data",
				Computed:    true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "Ports that PaaS service listens to",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"security_zone_uuid": {
				Type:        schema.TypeString,
				Description: "Security zone UUID linked to PaaS service",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "Network UUID containing security zone",
				Computed:    true,
			},
			"k8s_release": {
				Type:         schema.TypeString,
				Description:  "Release number of k8s service",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "PaaS service template that k8s service uses.",
				Computed:    true,
			},
			"node_pool": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Node pool's specification.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of node pool",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Description: "Number of worker nodes",
							Required:    true,
						},
						"cores": {
							Type:        schema.TypeInt,
							Description: "Cores per worker node",
							Required:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory per worker node (in GiB)",
							Required:    true,
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "Storage per worker node (in GiB)",
							Required:    true,
						},
						"storage_type": {
							Type:        schema.TypeString,
							Description: "Storage type (one of storage, storage_high, storage_insane)",
							Required:    true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								valid := false
								for _, stype := range storageTypes {
									if v.(string) == stype {
										valid = true
										break
									}
								}
								if !valid {
									errors = append(errors, fmt.Errorf("%v is not a valid storage type. Valid types are: %v", v.(string), strings.Join(storageTypes, ",")))
								}
								return
							},
						},
					},
				},
			},
			"usage_in_minute": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that PaaS service is in use",
				Computed:    true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Current price of PaaS service",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Time of the creation",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of PaaS service",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
	}
}

func resourceGridscaleK8sRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read k8s (%s) resource -", d.Id())
	paas, err := client.GetPaaSService(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := paas.Properties
	creds := props.Credentials
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if creds != nil && len(creds) > 0 {
		if err = d.Set("kubeconfig", creds[0].KubeConfig); err != nil {
			return fmt.Errorf("%s error setting kubeconfig: %v", errorPrefix, err)
		}
	}
	if err = d.Set("security_zone_uuid", props.SecurityZoneUUID); err != nil {
		return fmt.Errorf("%s error setting security_zone_uuid: %v", errorPrefix, err)
	}

	if err = d.Set("usage_in_minute", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minute: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", props.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_uuid", props.ServiceTemplateUUID); err != nil {
		return fmt.Errorf("%s error setting service_template_uuid: %v", errorPrefix, err)
	}

	//Get listen ports
	listenPorts := make([]interface{}, 0)
	for _, value := range props.ListenPorts {
		for k, portValue := range value {
			port := map[string]interface{}{
				"name": k,
				"port": portValue,
			}
			listenPorts = append(listenPorts, port)
		}
	}
	if err = d.Set("listen_port", listenPorts); err != nil {
		return fmt.Errorf("%s error setting listen ports: %v", errorPrefix, err)
	}

	// Get node pool parameters
	nodePoolList := make([]interface{}, 0)
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	nodePool := map[string]interface{}{
		"name":         d.Get("node_pool.0.name"),
		"node_count":   props.Parameters["k8s_worker_node_count"],
		"cores":        props.Parameters["k8s_worker_node_cores"],
		"memory":       props.Parameters["k8s_worker_node_ram"],
		"storage":      props.Parameters["k8s_worker_node_storage"],
		"storage_type": props.Parameters["k8s_worker_node_storage_type"],
	}
	nodePoolList = append(nodePoolList, nodePool)
	if err = d.Set("node_pool", nodePoolList); err != nil {
		return fmt.Errorf("%s error setting node_pool: %v", errorPrefix, err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	//Get all available networks
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	//look for a network that the PaaS service is in
	for _, network := range networks {
		securityZones := network.Properties.Relations.PaaSSecurityZones
		//Each network can contain only one security zone
		if len(securityZones) >= 1 {
			if securityZones[0].ObjectUUID == props.SecurityZoneUUID {
				if err = d.Set("network_uuid", network.Properties.ObjectUUID); err != nil {
					return fmt.Errorf("%s error setting network_uuid: %v", errorPrefix, err)
				}
			}
		}
	}
	return nil
}

func resourceGridscaleK8sCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	// Validate k8s parameters
	templateUUID, err := validateK8sParameters(client, d,
		k8sReleaseValidationOpt,
		k8sNodeCountValidationOpt,
		k8sCoreCountValidationOpt,
		k8sMemoryValidationOpt,
		k8sStorageValidationOpt,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: templateUUID,
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
		PaaSSecurityZoneUUID:    d.Get("security_zone_uuid").(string),
	}

	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	params := make(map[string]interface{})
	params["k8s_worker_node_ram"] = d.Get("node_pool.0.memory")
	params["k8s_worker_node_cores"] = d.Get("node_pool.0.cores")
	params["k8s_worker_node_count"] = d.Get("node_pool.0.node_count")
	params["k8s_worker_node_storage"] = d.Get("node_pool.0.storage")
	params["k8s_worker_node_storage_type"] = d.Get("node_pool.0.storage_type")
	requestBody.Parameters = params

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for PaaS service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleK8sRead(d, meta)
}

func resourceGridscaleK8sUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update k8s (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	// Only update k8s_release, when it is changed
	if d.HasChange("k8s_release") {
		// Check if the k8s release number exists
		templateUUID, err := validateK8sParameters(client, d, k8sReleaseValidationOpt)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		requestBody.PaaSServiceTemplateUUID = templateUUID
	}

	// Validate k8s parameters
	if _, err := validateK8sParameters(client, d,
		k8sNodeCountValidationOpt,
		k8sCoreCountValidationOpt,
		k8sMemoryValidationOpt,
		k8sStorageValidationOpt,
	); err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	params := make(map[string]interface{})
	params["k8s_worker_node_ram"] = d.Get("node_pool.0.memory")
	params["k8s_worker_node_cores"] = d.Get("node_pool.0.cores")
	params["k8s_worker_node_count"] = d.Get("node_pool.0.node_count")
	params["k8s_worker_node_storage"] = d.Get("node_pool.0.storage")
	params["k8s_worker_node_storage_type"] = d.Get("node_pool.0.storage_type")
	requestBody.Parameters = params

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleK8sRead(d, meta)
}

func resourceGridscaleK8sDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete k8s (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeletePaaSService(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

// validateK8sParameters validate k8s resource's selected parameters.
// It returns the UUID of the k8s service template, if the validation is successful.
// Otherwise, an error will be returned.
func validateK8sParameters(client *gsclient.Client, d *schema.ResourceData, parameters ...k8sValidationOpt) (string, error) {
	errorMessages := []string{"List of validation errors:\n"}
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	// Check if the k8s release number exists
	release := d.Get("k8s_release").(string)
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == k8sTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid && isValOptSelected(k8sReleaseValidationOpt, parameters) {
		errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid kubernetes release number. Valid release numbers are: %v\n", release, strings.Join(releases, ", ")))
	}

	// Check if mem, core count, node count, and storage are valid
	if attr, ok := d.GetOk("node_pool"); ok {
		for _, element := range attr.([]interface{}) {
			nodePool := element.(map[string]interface{})

			mem := nodePool["memory"].(int)
			_, ok := uTemplate.Properties.ParametersSchema["k8s_worker_node_ram"]
			minMem := uTemplate.Properties.ParametersSchema["k8s_worker_node_ram"].Min
			maxMem := uTemplate.Properties.ParametersSchema["k8s_worker_node_ram"].Max
			if (minMem > mem || maxMem < mem) &&
				isValOptSelected(k8sMemoryValidationOpt, parameters) && ok {
				errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid value for \"memory\". Valid value stays between %v and %v\n", mem, minMem, maxMem))
			}

			coreCount := nodePool["cores"].(int)
			_, ok = uTemplate.Properties.ParametersSchema["k8s_worker_node_cores"]
			minCoreCount := uTemplate.Properties.ParametersSchema["k8s_worker_node_cores"].Min
			maxCoreCount := uTemplate.Properties.ParametersSchema["k8s_worker_node_cores"].Max
			if (minCoreCount > coreCount || maxCoreCount < coreCount) &&
				isValOptSelected(k8sCoreCountValidationOpt, parameters) && ok {
				errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid value for \"cores\". Valid value stays between %v and %v\n", coreCount, minCoreCount, maxCoreCount))
			}

			nodeCount := nodePool["node_count"].(int)
			_, ok = uTemplate.Properties.ParametersSchema["k8s_worker_node_count"]
			minNodeCount := uTemplate.Properties.ParametersSchema["k8s_worker_node_count"].Min
			maxNodeCount := uTemplate.Properties.ParametersSchema["k8s_worker_node_count"].Max
			if (minNodeCount > nodeCount || maxNodeCount < nodeCount) &&
				isValOptSelected(k8sNodeCountValidationOpt, parameters) && ok {
				errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid value for number of \"node_count\". Valid value stays between %v and %v\n", nodeCount, minNodeCount, maxNodeCount))
			}

			storage := nodePool["storage"].(int)
			_, ok = uTemplate.Properties.ParametersSchema["k8s_worker_node_storage"]
			minStorage := uTemplate.Properties.ParametersSchema["k8s_worker_node_storage"].Min
			maxStorage := uTemplate.Properties.ParametersSchema["k8s_worker_node_storage"].Max
			if (minStorage > storage || maxStorage < storage) &&
				isValOptSelected(k8sStorageValidationOpt, parameters) && ok {
				errorMessages = append(errorMessages, fmt.Sprintf("%v is not a valid value for \"storage\". Valid value stays between %v and %v\n", storage, minStorage, maxStorage))
			}
		}
	}
	if len(errorMessages) > 1 {
		return "", fmt.Errorf(strings.Join(errorMessages, ""))
	}
	return uTemplate.Properties.ObjectUUID, nil
}

// isValOptSelected checks if a k8s validation option presents in a list of k8s validation options.
func isValOptSelected(opt k8sValidationOpt, list []k8sValidationOpt) bool {
	for _, v := range list {
		if v == opt {
			return true
		}
	}
	return false
}
