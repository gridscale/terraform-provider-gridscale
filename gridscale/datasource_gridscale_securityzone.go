package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscalePaaSSecurityZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscalePaaSSecurityZoneRead,
		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the object",
			},
			"location_uuid": {
				Type:        schema.TypeString,
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
				Computed:    true,
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

func dataSourceGridscalePaaSSecurityZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	secZone, err := client.GetPaaSSecurityZone(emptyCtx, id)
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := secZone.Properties
	d.SetId(props.ObjectUUID)
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("error setting name: %v", err)
	}
	if err = d.Set("location_uuid", props.LocationUUID); err != nil {
		return fmt.Errorf("error setting location_uuid: %v", err)
	}
	if err = d.Set("location_country", props.LocationCountry); err != nil {
		return fmt.Errorf("error setting location_country: %v", err)
	}
	if err = d.Set("location_iata", props.LocationIata); err != nil {
		return fmt.Errorf("error setting location_iata: %v", err)
	}
	if err = d.Set("location_name", props.LocationName); err != nil {
		return fmt.Errorf("error setting location_name: %v", err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("error setting create_time: %v", err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("error setting change_time: %v", err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("error setting status: %v", err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("error setting labels: %v", err)
	}

	//Set relations
	var rels []string
	for _, val := range props.Relation.Services {
		rels = append(rels, val.ObjectUUID)
	}
	if err = d.Set("relations", rels); err != nil {
		return fmt.Errorf("error setting relations: %v", err)
	}
	return nil
}
