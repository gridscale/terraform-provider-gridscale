package gridscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"log"
)

const filesystemTemplateFlavourName = "filesystem"

func resourceGridscaleFilesystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleFilesystemCreate,
		Read:   resourceGridscaleFilesystemRead,
		Delete: resourceGridscaleFilesystemDelete,
		Update: resourceGridscaleFilesystemUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateChange("release", func(ctx context.Context, old, new, meta interface{}) error {
				client := meta.(*gsclient.Client)
				newReleaseVal := new.(string)
				paasTemplates, err := client.GetPaaSTemplateList(ctx)
				if err != nil {
					return err
				}
				var isReleaseValid bool
				var releaseList []string
			TEMPLATELOOP:
				for _, template := range paasTemplates {
					if template.Properties.Flavour == filesystemTemplateFlavourName {
						// check if release already presents in the release list.
						// If so, ignore it.
						for _, release := range releaseList {
							if release == template.Properties.Release {
								continue TEMPLATELOOP
							}
						}
						releaseList = append(releaseList, template.Properties.Release)
						if template.Properties.Release == newReleaseVal {
							isReleaseValid = true
						}
					}
				}
				if !isReleaseValid {
					return fmt.Errorf("%v is not a valid Filesystem service release. Valid releases are: %v", newReleaseVal, strings.Join(releaseList, ", "))
				}
				return nil
			}),
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"release": {
				Type: schema.TypeString,
				Description: `The Filesystem service release of this instance.\n
				For convenience, please use gscloud https://github.com/gridscale/gscloud to get the list of available Filesystem service releases.`,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"performance_class": {
				Type:        schema.TypeString,
				Description: "Performance class of Filesystem service.",
				Required:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, class := range filesystemPerformanceClasses {
						if v.(string) == class {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid Filesystem performance class. Valid values are: %v", v.(string), strings.Join(postgreSQLPerformanceClasses, ",")))
					}
					return
				},
			},
			"listen_port": {
				Type:        schema.TypeSet,
				Description: "The port numbers where this Filesystem service accepts connections.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
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
				Description: "Security zone UUID linked to Filesystem service.",
				Deprecated:  "Security zone is deprecated for gridSQL, gridStore, and gridFs. Please consider to use private network instead.",
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"root_squash": {
				Type:        schema.TypeBool,
				Description: "Map root user/group ownership to anon_uid/anon_gid",
				Optional:    true,
			},
			"allowed_ip_ranges": {
				Type:        schema.TypeSet,
				Description: "Allowed CIDR block or IP address in CIDR notation.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"anon_uid": {
				Type:        schema.TypeInt,
				Description: "Target user id when root squash is active.",
				Optional:    true,
			},
			"anon_gid": {
				Type:        schema.TypeInt,
				Description: "Target group id when root squash is active.",
				Optional:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "The UUID of the network that the service is attached to.",
				Optional:    true,
				Computed:    true,
			},
			"service_template_uuid": {
				Type:        schema.TypeString,
				Description: "PaaS service template that Filesystem service uses.",
				Computed:    true,
			},
			"service_template_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template service's category used to create the service.",
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Description: "Number of minutes that Filesystem service is in use.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Time of the last change.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Date time this service has been created.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of Filesystem service.",
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

func resourceGridscaleFilesystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read gridFs (%s) resource -", d.Id())
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
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("security_zone_uuid", props.SecurityZoneUUID); err != nil {
		return fmt.Errorf("%s error setting security_zone_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("network_uuid", props.NetworkUUID); err != nil {
		return fmt.Errorf("%s error setting network_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_uuid", props.ServiceTemplateUUID); err != nil {
		return fmt.Errorf("%s error setting service_template_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("service_template_category", props.ServiceTemplateCategory); err != nil {
		return fmt.Errorf("%s error setting service_template_category: %v", errorPrefix, err)
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

	// Set Filesystem parameters
	if err = d.Set("root_squash", props.Parameters["root_squash"]); err != nil {
		return fmt.Errorf("%s error setting root_squash: %v", errorPrefix, err)
	}
	if err = d.Set("allowed_ip_ranges", props.Parameters["allowed_ip_ranges"]); err != nil {
		return fmt.Errorf("%s error setting allowed_ip_ranges: %v", errorPrefix, err)
	}
	if err = d.Set("anon_uid", props.Parameters["anon_uid"]); err != nil {
		return fmt.Errorf("%s error setting anon_uid: %v", errorPrefix, err)
	}
	if err = d.Set("anon_gid", props.Parameters["anon_gid"]); err != nil {
		return fmt.Errorf("%s error setting anon_gid: %v", errorPrefix, err)
	}

	//Get listen ports
	listenPorts := make([]interface{}, 0)
	for host, value := range props.ListenPorts {
		for k, portValue := range value {
			port := map[string]interface{}{
				"name": k,
				"host": host,
				"port": portValue,
			}
			listenPorts = append(listenPorts, port)
		}
	}
	if err = d.Set("listen_port", listenPorts); err != nil {
		return fmt.Errorf("%s error setting listen ports: %v", errorPrefix, err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	// Look for security zone's network that the PaaS service is connected to
	// (if the paas is connected to security zone. O.w skip)
	if props.SecurityZoneUUID == "" {
		return nil
	}
	networks, err := client.GetNetworkList(context.Background())
	if err != nil {
		return fmt.Errorf("%s error getting networks: %v", errorPrefix, err)
	}
	//look for a network that the Filesystem service is in
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

func resourceGridscaleFilesystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("create gridFs (%s) resource -", d.Id())

	release := d.Get("release").(string)
	performanceClass := d.Get("performance_class").(string)
	// Get filesystem template UUID
	templateUUID, err := getFilesystemTemplateUUID(client, release, performanceClass)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	requestBody := gsclient.PaaSServiceCreateRequest{
		Name:                    d.Get("name").(string),
		PaaSServiceTemplateUUID: templateUUID,
		Labels:                  convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	networkUUIDInf, isNetworkSet := d.GetOk("network_uuid")
	if isNetworkSet {
		requestBody.NetworkUUID = networkUUIDInf.(string)
	}
	// If network_uuid is set, skip setting security_zone_uuid.
	if secZoneUUIDInf, ok := d.GetOk("security_zone_uuid"); ok && !isNetworkSet {
		requestBody.PaaSSecurityZoneUUID = secZoneUUIDInf.(string)
	}
	params := make(map[string]interface{})
	if rootSquash, ok := d.GetOk("root_squash"); ok {
		params["root_squash"] = rootSquash
	}
	if allowedIPRanges, ok := d.GetOk("allowed_ip_ranges"); ok {
		params["allowed_ip_ranges"] = convSOStrings(allowedIPRanges.(*schema.Set).List())
	}
	if anonUID, ok := d.GetOk("anon_uid"); ok {
		params["anon_uid"] = anonUID
	}
	if anonGID, ok := d.GetOk("anon_gid"); ok {
		params["anon_gid"] = anonGID
	}
	requestBody.Parameters = params
	log.Printf("-------%v\n", requestBody)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreatePaaSService(ctx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for Filesystem service %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleFilesystemRead(d, meta)
}

func resourceGridscaleFilesystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update gridFs (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.PaaSServiceUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}
	if d.HasChange("network_uuid") {
		requestBody.NetworkUUID = d.Get("network_uuid").(string)
	}
	params := make(map[string]interface{})
	if rootSquash, ok := d.GetOk("root_squash"); ok {
		params["root_squash"] = rootSquash
	}
	if allowedIPRanges, ok := d.GetOk("allowed_ip_ranges"); ok {
		params["allowed_ip_ranges"] = convSOStrings(allowedIPRanges.(*schema.Set).List())
	}
	if anonUID, ok := d.GetOk("anon_uid"); ok {
		params["anon_uid"] = anonUID
	}
	if anonGID, ok := d.GetOk("anon_gid"); ok {
		params["anon_gid"] = anonGID
	}
	requestBody.Parameters = params

	// Only update templateUUID, when `release` or `performance_class` is changed
	if d.HasChange("release") || d.HasChange("performance_class") {
		// Get postgres template UUID
		release := d.Get("release").(string)
		performanceClass := d.Get("performance_class").(string)
		templateUUID, err := getFilesystemTemplateUUID(client, release, performanceClass)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
		requestBody.PaaSServiceTemplateUUID = templateUUID
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdatePaaSService(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleFilesystemRead(d, meta)
}

func resourceGridscaleFilesystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete paas (%s) resource -", d.Id())

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

// getFilesystemTemplateUUID returns the UUID of the filesystem service template.
func getFilesystemTemplateUUID(client *gsclient.Client, release, performanceClass string) (string, error) {
	paasTemplates, err := client.GetPaaSTemplateList(context.Background())
	if err != nil {
		return "", err
	}
	var isReleaseValid bool
	var releases []string
	var uTemplate gsclient.PaaSTemplate
	for _, template := range paasTemplates {
		if template.Properties.Flavour == filesystemTemplateFlavourName {
			releases = append(releases, template.Properties.Release)
			if template.Properties.Release == release && template.Properties.PerformanceClass == performanceClass {
				isReleaseValid = true
				uTemplate = template
			}
		}
	}
	if !isReleaseValid {
		return "", fmt.Errorf("%v is not a valid Filesystem service release. Valid releases are: %v", release, strings.Join(releases, ", "))
	}

	return uTemplate.Properties.ObjectUUID, nil
}
