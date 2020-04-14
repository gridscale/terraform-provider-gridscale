package gridscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go/v2"
)

func resourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageCreate,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update: resourceGridscaleStorageUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Description:  "The capacity of a storage in GB.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Optional:    true,
				ForceNew:    true,
				Default:     "storage",
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
			"license_product_no": {
				Type:        schema.TypeInt,
				Description: "If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).",
				Computed:    true,
			},
			"last_used_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time of the last object change.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier.",
				Computed:    true,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"usage_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"template": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sshkeys": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"password_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								valid := false
								for _, passType := range passwordTypes {
									if v.(string) == passType {
										valid = true
										break
									}
								}
								if !valid {
									errors = append(errors, fmt.Errorf("%v is not a valid password type. Valid types are: %v", v.(string), strings.Join(passwordTypes, ",")))
								}
								return
							},
						},
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"template_uuid": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Update: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Delete: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
		},
	}
}

func resourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read storage (%s) resource -", d.Id())
	storage, err := client.GetStorage(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	if err = d.Set("change_time", storage.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", storage.Properties.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("status", storage.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("license_product_no", storage.Properties.LicenseProductNo); err != nil {
		return fmt.Errorf("%s error setting license_product_no: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", storage.Properties.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", storage.Properties.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("last_used_template", storage.Properties.LastUsedTemplate); err != nil {
		return fmt.Errorf("%s error setting last_used_template: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", storage.Properties.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}
	if err = d.Set("capacity", storage.Properties.Capacity); err != nil {
		return fmt.Errorf("%s error setting capacity: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", storage.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("storage_type", storage.Properties.StorageType); err != nil {
		return fmt.Errorf("%s error setting storage_type: %v", errorPrefix, err)
	}
	if err = d.Set("parent_uuid", storage.Properties.ParentUUID); err != nil {
		return fmt.Errorf("%s error setting parent_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("name", storage.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", storage.Properties.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", storage.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}

	if err = d.Set("labels", storage.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update storage (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.StorageUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate)*time.Second)
	defer cancel()
	err := client.UpdateStorage(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.StorageCreateRequest{
		Name:     d.Get("name").(string),
		Capacity: d.Get("capacity").(int),
		Labels:   convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	storageType := d.Get("storage_type").(string)
	if storageType == "storage" {
		requestBody.StorageType = gsclient.DefaultStorageType
	} else if storageType == "storage_high" {
		requestBody.StorageType = gsclient.HighStorageType
	} else if storageType == "storage_insane" {
		requestBody.StorageType = gsclient.InsaneStorageType
	}

	//since only one template can be used, we can just look at index 0
	if _, ok := d.GetOk("template"); ok {
		template := gsclient.StorageTemplate{
			Password:     d.Get("template.0.password").(string),
			Hostname:     d.Get("template.0.hostname").(string),
			TemplateUUID: d.Get("template.0.template_uuid").(string),
		}
		passType := d.Get("template.0.password_type").(string)
		if passType == "plain" {
			template.PasswordType = gsclient.PlainPasswordType
		} else if passType == "crypt" {
			template.PasswordType = gsclient.CryptPasswordType
		}

		if attr, ok := d.GetOk("template.0.sshkeys"); ok {
			for _, value := range attr.([]interface{}) {
				template.Sshkeys = append(template.Sshkeys, value.(string))
			}
		}
		requestBody.Template = &template
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate)*time.Second)
	defer cancel()
	response, err := client.CreateStorage(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for storage %s has been set to %v", requestBody.Name, response.ObjectUUID)

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete storage (%s) resource -", d.Id())
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete)*time.Second)
	defer cancel()
	storage, err := client.GetStorage(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	//Stop all server relating to this IP address if there is one
	for _, server := range storage.Properties.Relations.Servers {
		unlinkStorageAction := func(ctx context.Context) error {
			err := client.UnlinkStorage(ctx, server.ObjectUUID, d.Id())
			return err
		}
		//UnlinkStorage requires the server to be off
		err = globalServerStatusList.runActionRequireServerOff(ctx, client, server.ObjectUUID, false, unlinkStorageAction)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	err = client.DeleteStorage(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
