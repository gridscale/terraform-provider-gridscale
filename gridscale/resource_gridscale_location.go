package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"
)

func resourceGridscaleLocation() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleLocationCreate,
		Read:   resourceGridscaleLocationRead,
		Delete: resourceGridscaleLocationDelete,
		Update: resourceGridscaleLocationUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"parent_location_uuid": {
				Type:        schema.TypeString,
				Description: "The location_uuid of an existing public location in which to create the private location.",
				Required:    true,
				ForceNew:    true,
			},
			"cpunode_count": {
				Type:        schema.TypeInt,
				Description: "The number of dedicated cpunodes to assigne to the private location.",
				Required:    true,
			},
			"product_no": {
				Type:        schema.TypeInt,
				Description: "The product number of a valid and available dedicated cpunode article.",
				Required:    true,
				ForceNew:    true,
			},
			"iata": {
				Type:        schema.TypeString,
				Description: "IATA airport code, which works as a location identifier.",
				Computed:    true,
			},
			"country": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Computed:    true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "True if the location is active.",
				Computed:    true,
			},
			"cpunode_count_change_requested": {
				Type:        schema.TypeInt,
				Description: "The requested number of dedicated cpunodes.",
				Computed:    true,
			},
			"product_no_change_requested": {
				Type:        schema.TypeInt,
				Description: "The product number of a valid and available dedicated cpunode article.",
				Computed:    true,
			},
			"parent_location_uuid_change_requested": {
				Type:        schema.TypeInt,
				Description: "The location_uuid of an existing public location in which to create the private location.",
				Computed:    true,
			},
			"public": {
				Type:        schema.TypeBool,
				Description: "True if this location is publicly available or a private location.",
				Computed:    true,
			},
			"certification_list": {
				Type:        schema.TypeString,
				Description: "List of certifications.",
				Computed:    true,
			},
			"city": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Computed:    true,
			},
			"data_protection_agreement": {
				Type:        schema.TypeString,
				Description: "Data protection agreement.",
				Computed:    true,
			},
			"geo_location": {
				Type:        schema.TypeString,
				Description: "Geo location.",
				Computed:    true,
			},
			"green_energy": {
				Type:        schema.TypeString,
				Description: "Green energy.",
				Computed:    true,
			},
			"operator_certification_list": {
				Type:        schema.TypeString,
				Description: "List of operator certifications.",
				Computed:    true,
			},
			"owner": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the owner.",
				Computed:    true,
			},
			"owner_website": {
				Type:        schema.TypeString,
				Description: "The website of the owner.",
				Computed:    true,
			},
			"site_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the website.",
				Computed:    true,
			},
			"hardware_profiles": {
				Type:        schema.TypeString,
				Description: "List of supported hardware profiles.",
				Computed:    true,
			},
			"has_rocket_storage": {
				Type:        schema.TypeString,
				Description: "TRUE if the location supports rocket storage.",
				Computed:    true,
			},
			"has_server_provisioning": {
				Type:        schema.TypeString,
				Description: "TRUE if the location supports server provisioning.",
				Computed:    true,
			},
			"object_storage_region": {
				Type:        schema.TypeString,
				Description: "The region of the object storage.",
				Computed:    true,
			},
			"backup_center_location_uuid": {
				Type:        schema.TypeString,
				Description: "The location_uuid of a backup location.",
				Computed:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleLocationRead(d *schema.ResourceData, meta interface{}) error {
	errorPrefix := fmt.Sprintf("read location (%s) resource -", d.Id())
	client := meta.(*gsclient.Client)
	loc, err := client.GetLocation(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	locProp := loc.Properties
	if err = d.Set("name", locProp.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("iata", locProp.Iata); err != nil {
		return fmt.Errorf("%s error setting iata: %v", errorPrefix, err)
	}
	if err = d.Set("country", locProp.Country); err != nil {
		return fmt.Errorf("%s error setting country: %v", errorPrefix, err)
	}
	if err = d.Set("active", locProp.Active); err != nil {
		return fmt.Errorf("%s error setting active: %v", errorPrefix, err)
	}
	if err = d.Set("cpunode_count_change_requested", locProp.ChangeRequested.CPUNodeCount); err != nil {
		return fmt.Errorf("%s error setting cpunode_count_change_requested: %v", errorPrefix, err)
	}
	if err = d.Set("product_no_change_requested", locProp.ChangeRequested.ProductNo); err != nil {
		return fmt.Errorf("%s error setting product_no_change_requested: %v", errorPrefix, err)
	}
	if err = d.Set("parent_location_uuid_change_requested", locProp.ChangeRequested.ParentLocationUUID); err != nil {
		return fmt.Errorf("%s error setting parent_location_uuid_change_requested: %v", errorPrefix, err)
	}
	if err = d.Set("cpunode_count", locProp.CPUNodeCount); err != nil {
		return fmt.Errorf("%s error setting cpunode_count: %v", errorPrefix, err)
	}
	if err = d.Set("public", locProp.Public); err != nil {
		return fmt.Errorf("%s error setting public: %v", errorPrefix, err)
	}
	if err = d.Set("product_no", locProp.ProductNo); err != nil {
		return fmt.Errorf("%s error setting product_no: %v", errorPrefix, err)
	}
	if err = d.Set("certification_list", locProp.LocationInformation.CertificationList); err != nil {
		return fmt.Errorf("%s error setting certification_list: %v", errorPrefix, err)
	}
	if err = d.Set("city", locProp.LocationInformation.City); err != nil {
		return fmt.Errorf("%s error setting city: %v", errorPrefix, err)
	}
	if err = d.Set("data_protection_agreement", locProp.LocationInformation.DataProtectionAgreement); err != nil {
		return fmt.Errorf("%s error setting data_protection_agreement: %v", errorPrefix, err)
	}
	if err = d.Set("geo_location", locProp.LocationInformation.GeoLocation); err != nil {
		return fmt.Errorf("%s error setting geo_location: %v", errorPrefix, err)
	}
	if err = d.Set("green_energy", locProp.LocationInformation.GreenEnergy); err != nil {
		return fmt.Errorf("%s error setting green_energy: %v", errorPrefix, err)
	}
	if err = d.Set("operator_certification_list", locProp.LocationInformation.OperatorCertificationList); err != nil {
		return fmt.Errorf("%s error setting operator_certification_list: %v", errorPrefix, err)
	}
	if err = d.Set("owner", locProp.LocationInformation.Owner); err != nil {
		return fmt.Errorf("%s error setting owner: %v", errorPrefix, err)
	}
	if err = d.Set("owner_website", locProp.LocationInformation.OwnerWebsite); err != nil {
		return fmt.Errorf("%s error setting owner_website: %v", errorPrefix, err)
	}
	if err = d.Set("site_name", locProp.LocationInformation.SiteName); err != nil {
		return fmt.Errorf("%s error setting site_name: %v", errorPrefix, err)
	}
	if err = d.Set("hardware_profiles", locProp.Features.HardwareProfiles); err != nil {
		return fmt.Errorf("%s error setting hardware_profiles: %v", errorPrefix, err)
	}
	if err = d.Set("has_rocket_storage", locProp.Features.HasRocketStorage); err != nil {
		return fmt.Errorf("%s error setting has_rocket_storage: %v", errorPrefix, err)
	}
	if err = d.Set("has_server_provisioning", locProp.Features.HasServerProvisioning); err != nil {
		return fmt.Errorf("%s error setting has_server_provisioning: %v", errorPrefix, err)
	}
	if err = d.Set("object_storage_region", locProp.Features.ObjectStorageRegion); err != nil {
		return fmt.Errorf("%s error setting object_storage_region: %v", errorPrefix, err)
	}
	if err = d.Set("backup_center_location_uuid", locProp.Features.BackupCenterLocationUUID); err != nil {
		return fmt.Errorf("%s error setting backup_center_location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("labels", locProp.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleLocationCreate(d *schema.ResourceData, meta interface{}) error {
	errorPrefix := "create location resource -"
	client := meta.(*gsclient.Client)
	requestBody := gsclient.LocationCreateRequest{
		Name:               d.Get("name").(string),
		Labels:             convSOStrings(d.Get("labels").(*schema.Set).List()),
		ParentLocationUUID: d.Get("parent_location_uuid").(string),
		CPUNodeCount:       d.Get("cpunode_count").(int),
		ProductNo:          d.Get("product_no").(int),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateLocation(ctx, requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for the new location has been set to %v", response.ObjectUUID)
	return resourceGridscaleLocationRead(d, meta)
}

func resourceGridscaleLocationUpdate(d *schema.ResourceData, meta interface{}) error {
	errorPrefix := fmt.Sprintf("update location (%s) resource -", d.Id())
	client := meta.(*gsclient.Client)
	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	cpunodeCount := d.Get("cpunode_count").(int)
	requestBody := gsclient.LocationUpdateRequest{
		Name:         d.Get("name").(string),
		Labels:       &labels,
		CPUNodeCount: &cpunodeCount,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateLocation(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleLocationRead(d, meta)
}

func resourceGridscaleLocationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete location (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeleteLocation(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
