package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceGridscalePaaSSecurityZone() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGridscalePaaSSecurityZoneRead,
		Create: resourceGridscalePaaSSecurityZoneCreate,
		Update: resourceGridscalePaaSSecurityZoneUpdate,
		Delete: resourceGridscalePaaSSecurityZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The human-readable name of the object",
				ValidateFunc: validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Helps to identify which datacenter an object belongs to",
			},
			"location_country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the location",
			},
			"location_iata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Uses IATA airport code, which works as a location identifier",
			},
			"location_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the location",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the date and time the object was initially created",
			},
			"change_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the date and time of the last object change",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status indicates the status of the object",
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"relations": {
				Type:        schema.TypeSet,
				Description: "List of PaaS services' UUIDs relating to the security zone",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGridscalePaaSSecurityZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read paas security zone (%s) resource -", d.Id())
	secZone, err := client.GetPaaSSecurityZone(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := secZone.Properties
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", props.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", props.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", props.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", props.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	//Set relations
	var rels []string
	for _, val := range props.Relation.Services {
		rels = append(rels, val.ObjectUUID)
	}
	if err = d.Set("relations", rels); err != nil {
		return fmt.Errorf("%s error setting relations: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscalePaaSSecurityZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.PaaSSecurityZoneCreateRequest{
		Name: d.Get("name").(string),
	}
	response, err := client.CreatePaaSSecurityZone(emptyCtx, requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for security zone %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscalePaaSSecurityZoneRead(d, meta)
}

func resourceGridscalePaaSSecurityZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update paas security zone (%s) resource -", d.Id())
	requestBody := gsclient.PaaSSecurityZoneUpdateRequest{
		Name: d.Get("name").(string),
	}
	err := client.UpdatePaaSSecurityZone(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscalePaaSSecurityZoneRead(d, meta)
}

func resourceGridscalePaaSSecurityZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete paas security zone (%s) resource -", d.Id())
	err := client.DeletePaaSSecurityZone(emptyCtx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
