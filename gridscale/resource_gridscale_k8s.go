package gridscale

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const (
	k8sTemplateFlavourName         = "kubernetes"
	k8sLabelPrefix                 = "#gsk#"
	k8sRocketStorageSupportRelease = "1.26"
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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			client := meta.(*gsclient.Client)
			newReleaseValInf, isReleaseSet := d.GetOk("release")
			newReleaseVal := newReleaseValInf.(string)
			newVersionValInf, isVersionSet := d.GetOk("gsk_version")
			newVersionVal := newVersionValInf.(string)
			if !isReleaseSet && !isVersionSet {
				return errors.New("either \"release\" or \"gsk_version\" has to be defined.")
			}
			if isReleaseSet && isVersionSet {
				return errors.New("\"release\" and \"gsk_version\" cannot be set at the same time. Only one of them is set at a time.")
			}

			paasTemplates, err := client.GetPaaSTemplateList(ctx)
			if err != nil {
				return err
			}
			var isReleaseValid bool
			var isVersionValid bool
			var chosenTemplate gsclient.PaaSTemplate
			var releaseList []string
			var versionList []string
		TEMPLATELOOP:
			for _, template := range paasTemplates {
				if template.Properties.Flavour == k8sTemplateFlavourName {
					// check if release already presents in the release list.
					// If so, ignore it.
					for _, release := range releaseList {
						if release == template.Properties.Release {
							continue TEMPLATELOOP
						}
					}
					releaseList = append(releaseList, template.Properties.Release)
					versionList = append(versionList, template.Properties.Version)

					if template.Properties.Version == newVersionVal && isVersionSet {
						isVersionValid = true
						chosenTemplate = template
					}
					if template.Properties.Release == newReleaseVal && isReleaseSet {
						isReleaseValid = true
						chosenTemplate = template
					}
				}
			}
			if !isReleaseValid && isReleaseSet {
				return fmt.Errorf("%v is an INVALID Kubernetes minor release. Valid releases are: %v\n", newReleaseVal, strings.Join(releaseList, ", "))
			}
			if !isVersionValid && isVersionSet {
				return fmt.Errorf("%v is an INVALID gridscale Kubernetes (GSK) version. Valid GSK versions are: %v\n", newVersionVal, strings.Join(versionList, ", "))
			}
			return validateK8sParameters(d, chosenTemplate)
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: "K8s config data",
				Computed:    true,
				Sensitive:   true,
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "The port number where this k8s service accepts connections.",
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
				Description: "Security zone UUID linked to PaaS service.",
				Deprecated:  "GSK service does not support security zone.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "Network UUID containing security zone",
				Deprecated: `network_uuid of a security zone is no more available for GSK.
					Please consider to use k8s_private_network_uuid for connecting external services to the cluster.`,
				Computed: true,
			},
			"k8s_private_network_uuid": {
				Type:        schema.TypeString,
				Description: "Private network UUID which k8s nodes are attached to. It can be used to attach other PaaS/VMs.",
				Computed:    true,
			},
			"release": {
				Type:        schema.TypeString,
				Description: "The k8s release of this instance.",
				Optional:    true,
			},
			"gsk_version": {
				Type:        schema.TypeString,
				Description: "The gridscale k8s PaaS version (issued by gridscale) of this instance.",
				Optional:    true,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "PaaS service template identifier for this service.",
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
							Description: "Name of node pool.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Description: "Number of worker nodes.",
							Required:    true,
						},
						"cores": {
							Type:        schema.TypeInt,
							Description: "Cores per worker node.",
							Required:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory per worker node (in GiB).",
							Required:    true,
						},
						"storage": {
							Type:        schema.TypeInt,
							Description: "Storage per worker node (in GiB).",
							Required:    true,
						},
						"storage_type": {
							Type:        schema.TypeString,
							Description: "Storage type.",
							Required:    true,
						},
						"rocket_storage": {
							Type:        schema.TypeInt,
							Description: "Rocket storage per worker node (in GiB).",
							Optional:    true,
						},
						"surge_node": {
							Type:        schema.TypeBool,
							Description: "Enable surge node to avoid resources shortage during the cluster upgrade.",
							Optional:    true,
							Default:     true,
						},
						"cluster_cidr": {
							Type:        schema.TypeString,
							Description: "The cluster CIDR that will be used to generate the CIDR of nodes, services, and pods. The allowed CIDR prefix length is /16. If this field is empty, the default value is \"10.244.0.0/16\"",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that PaaS service is in use",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Time this service was created.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the service",
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
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
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
			props = paas.Properties
			creds = props.Credentials
		}
		if err = d.Set("kubeconfig", creds[0].KubeConfig); err != nil {
			return fmt.Errorf("%s error setting kubeconfig: %v", errorPrefix, err)
		}
	}
	if err = d.Set("security_zone_uuid", props.SecurityZoneUUID); err != nil {
		return fmt.Errorf("%s error setting security_zone_uuid: %v", errorPrefix, err)
	}

	if err = d.Set("usage_in_minutes", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
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

	// Set rocket storage if it is set
	if _, isRocketStorageSet := props.Parameters["k8s_worker_node_rocket_storage"]; isRocketStorageSet {
		nodePool["rocket_storage"] = props.Parameters["k8s_worker_node_rocket_storage"]
	}

	// Set cluster CIDR if it is set
	if _, isClusterCIDRSet := props.Parameters["k8s_cluster_cidr"]; isClusterCIDRSet {
		nodePool["cluster_cidr"] = props.Parameters["k8s_cluster_cidr"]
	}

	// Surge node feature is enable if k8s_surge_node_count > 0
	if surgeNodeCount, ok := props.Parameters["k8s_surge_node_count"].(float64); ok {
		nodePool["surge_node"] = surgeNodeCount > 0
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

	k8sLabel := fmt.Sprintf("%s%s", k8sLabelPrefix, d.Id())
	// look for a network having the defined k8sLabel.
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

func resourceGridscaleK8sCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create k8s (%s) resource -", d.Id())

	var templateUUID string
	// Validate k8s release, and get template UUID from release
	if release, isReleaseSet := d.GetOk("release"); isReleaseSet {
		var err error
		templateUUID, err = getK8sTemplateUUIDFromRelease(client, release.(string))
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	// Validate gsk version, and get template UUID from gsk version.
	if version, isVersionSet := d.GetOk("gsk_version"); isVersionSet {
		var err error
		templateUUID, err = getK8sTemplateUUIDFromGSKVersion(client, version.(string))
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
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
	// Set rocket storage if it is set
	if rocketStorage, isRocketStorageSet := d.GetOk("node_pool.0.rocket_storage"); isRocketStorageSet {
		params["k8s_worker_node_rocket_storage"] = rocketStorage
	}
	// Set cluster CIDR if it is set
	if clusterCIDR, isClusterCIDRSet := d.GetOk("node_pool.0.cluster_cidr"); isClusterCIDRSet {
		params["k8s_cluster_cidr"] = clusterCIDR
	}
	isSurgeNodeEnabled := d.Get("node_pool.0.surge_node").(bool)
	if isSurgeNodeEnabled {
		params["k8s_surge_node_count"] = 1
	} else {
		params["k8s_surge_node_count"] = 0
	}
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
	currentTemplateUUID := d.Get("service_template_uuid")
	// Only update release/gsk version, when it is changed
	if releaseValInf, isReleaseSet := d.GetOk("release"); d.HasChange("release") && isReleaseSet {
		// Check if the k8s release number exists
		release := releaseValInf.(string)
		templateUUID, err := getK8sTemplateUUIDFromRelease(client, release)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		// Only add template UUID when it really has been changed.
		if templateUUID != currentTemplateUUID.(string) {
			requestBody.PaaSServiceTemplateUUID = templateUUID
		}
	}
	if versionValInf, isVersionSet := d.GetOk("gsk_version"); d.HasChange("gsk_version") && isVersionSet {
		version := versionValInf.(string)
		templateUUID, err := getK8sTemplateUUIDFromGSKVersion(client, version)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		// Only add template UUID when it really has been changed.
		if templateUUID != currentTemplateUUID.(string) {
			requestBody.PaaSServiceTemplateUUID = templateUUID
		}
	}

	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	params := make(map[string]interface{})
	params["k8s_worker_node_ram"] = d.Get("node_pool.0.memory")
	params["k8s_worker_node_cores"] = d.Get("node_pool.0.cores")
	params["k8s_worker_node_count"] = d.Get("node_pool.0.node_count")
	params["k8s_worker_node_storage"] = d.Get("node_pool.0.storage")
	params["k8s_worker_node_storage_type"] = d.Get("node_pool.0.storage_type")
	// Set rocket storage if it is set
	if rocketStorage, isRocketStorageSet := d.GetOk("node_pool.0.rocket_storage"); isRocketStorageSet {
		params["k8s_worker_node_rocket_storage"] = rocketStorage
	}
	// Set cluster CIDR if it is set
	if clusterCIDR, isClusterCIDRSet := d.GetOk("node_pool.0.cluster_cidr"); isClusterCIDRSet {
		params["k8s_cluster_cidr"] = clusterCIDR
	}
	isSurgeNodeEnabled := d.Get("node_pool.0.surge_node").(bool)
	if isSurgeNodeEnabled {
		params["k8s_surge_node_count"] = 1
	} else {
		params["k8s_surge_node_count"] = 0
	}
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
	err := errHandler.SuppressHTTPErrorCodes(
		client.DeletePaaSService(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}

// getK8sTemplateUUIDFromRelease returns the UUID of the k8s service template from given release.
func getK8sTemplateUUIDFromRelease(client *gsclient.Client, release string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == k8sTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release {
				isReleaseValid = true
				uTemplate = template
				break
			}
		}
	}
	if !isReleaseValid {
		return "", fmt.Errorf("%v is an INVALID Kubernetes minor release. Valid releases are: %v\n", release, strings.Join(releases, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}

// getK8sTemplateUUIDFromGSKVersion returns the UUID of the k8s service template from given GSK version.
func getK8sTemplateUUIDFromGSKVersion(client *gsclient.Client, version string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isVersionValid bool
	var versions []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == k8sTemplateFlavourName {
			versions = append(versions, template.Properties.Version)
			if template.Properties.Version == version {
				isVersionValid = true
				uTemplate = template
				break
			}
		}
	}
	if !isVersionValid {
		return "", fmt.Errorf("%v is an INVALID gridscale Kubernetes (GSK) version. Valid GSK versions are: %v\n", version, strings.Join(versions, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}

func validateK8sParameters(d *schema.ResourceDiff, template gsclient.PaaSTemplate) error {
	var errorMessages []string

	worker_memory_scheme, mem_ok := template.Properties.ParametersSchema["k8s_worker_node_ram"]
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	if memory, ok := d.GetOk("node_pool.0.memory"); ok && mem_ok {
		if memory.(int) < worker_memory_scheme.Min || memory.(int) > worker_memory_scheme.Max {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'node_pool.0.memory' value. Value must stay between %d and %d\n", worker_memory_scheme.Min, worker_memory_scheme.Max))
		}
	}

	worker_core_scheme, core_ok := template.Properties.ParametersSchema["k8s_worker_node_cores"]
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	if core, ok := d.GetOk("node_pool.0.cores"); ok && core_ok {
		if core.(int) < worker_core_scheme.Min || core.(int) > worker_core_scheme.Max {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'node_pool.0.cores' value. Value must stay between %d and %d\n", worker_core_scheme.Min, worker_core_scheme.Max))
		}
	}

	worker_count_scheme, worker_count_ok := template.Properties.ParametersSchema["k8s_worker_node_count"]
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	if node_count, ok := d.GetOk("node_pool.0.node_count"); ok && worker_count_ok {
		if node_count.(int) < worker_count_scheme.Min || node_count.(int) > worker_count_scheme.Max {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'node_pool.0.node_count' value. Value must stay between %d and %d\n", worker_count_scheme.Min, worker_count_scheme.Max))
		}
	}

	worker_storage_scheme, storage_ok := template.Properties.ParametersSchema["k8s_worker_node_storage"]
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	if storage, ok := d.GetOk("node_pool.0.storage"); ok && storage_ok {
		if storage.(int) < worker_storage_scheme.Min || storage.(int) > worker_storage_scheme.Max {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'node_pool.0.storage' value. Value must stay between %d and %d\n", worker_storage_scheme.Min, worker_storage_scheme.Max))
		}
	}

	worker_rocket_storage_scheme, rocket_storage_ok := template.Properties.ParametersSchema["k8s_worker_node_rocket_storage"]
	// TODO: The API scheme will be CHANGED in the future. There will be multiple node pools.
	if rocket_storage, ok := d.GetOk("node_pool.0.rocket_storage"); ok && rocket_storage_ok {
		rocketStorageValidation := true
		featureReleaseCompabilityValidation := true
		supportedRelease, err := NewRelease(k8sRocketStorageSupportRelease)
		if err != nil {
			panic("Something went wrong at backend side parsing of version string expected for support of rocket storage at k8s.")
		}
		requestedRelease, err := NewRelease(template.Properties.Release)
		if err != nil {
			errorMessages = append(errorMessages, "The release doesn't match a valid version string.")
			featureReleaseCompabilityValidation = false
		}
		if featureReleaseCompabilityValidation {
			err := requestedRelease.CheckIfFeatureIsKnown(&Feature{Description: "rocket storage", Release: *supportedRelease})
			if err != nil {
				errorMessages = append(errorMessages, err.Error())
				rocketStorageValidation = false
			}
		}
		if rocketStorageValidation && (rocket_storage.(int) < worker_rocket_storage_scheme.Min || rocket_storage.(int) > worker_rocket_storage_scheme.Max) {
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid 'node_pool.0.rocket_storage' value. Value must stay between %d and %d\n", worker_rocket_storage_scheme.Min, worker_rocket_storage_scheme.Max))
		}
	}

	worker_storage_type_scheme, storage_type_ok := template.Properties.ParametersSchema["k8s_worker_node_storage_type"]
	if storage_type, ok := d.GetOk("node_pool.0.storage_type"); ok && storage_type_ok {
		var isValid bool
		for _, allowedValue := range worker_storage_type_scheme.Allowed {
			if storage_type.(string) == allowedValue {
				isValid = true
			}
		}
		if !isValid {
			errorMessages = append(errorMessages,
				fmt.Sprintf("Invalid 'node_pool.0.storage_type' value. Value must be one of these:\n\t%s",
					strings.Join(worker_storage_type_scheme.Allowed, "\n\t"),
				),
			)
		}
	}

	cluster_cidr_template, cluster_cidr_template_ok := template.Properties.ParametersSchema["k8s_cluster_cidr"]
	if cluster_cidr, ok := d.GetOk("node_pool.0.cluster_cidr"); ok {
		// if the template doesn't support cluster_cidr, return error if it is set
		if !cluster_cidr_template_ok {
			errorMessages = append(errorMessages, "The template doesn't support cluster_cidr. Please remove it from your configuration.\n")
		} else {
			// if the template supports cluster_cidr, validate the value
			if cluster_cidr.(string) != "" {
				_, _, err := net.ParseCIDR(cluster_cidr.(string))
				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for PaaS template release. Value must be a valid CIDR.\n"))
				}
			}
			// if cluster_cidr_template is immutable, return error if it is set during k8s creation
			// and it is changed during k8s update
			if cluster_cidr_template.Immutable {
				oldClusterCIDR, _ := d.GetChange("node_pool.0.cluster_cidr")
				if oldClusterCIDR != "" && d.HasChange("node_pool.0.cluster_cidr") {
					errorMessages = append(errorMessages, "Cannot change parameter cluster_cidr, because it is immutable.\n")
				}
			}
		}
	}

	if len(errorMessages) != 0 {
		return errors.New(strings.Join(errorMessages, ""))
	}
	return nil
}
