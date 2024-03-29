package gridscale

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
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
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			storageVariant := d.Get("storage_variant").(string)
			if storageVariant == "local" {
				if d.HasChange("storage_type") {
					return errors.New("storage_type cannot be set when storage_variant is set to \"local\"")
				}
			}
			return nil
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Description:  "The capacity of a storage in GB.",
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "The location this object is placed.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Optional:    true,
				Computed:    true,
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
			"storage_variant": {
				Type:        schema.TypeString,
				Description: "Storage variant (one of local or distributed).",
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, sVariant := range storageVariants {
						if v.(string) == sVariant {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid storage variant. Valid variants are: %v", v.(string), strings.Join(storageTypes, ",")))
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
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Two digit country code (ISO 3166-2) of the location where this object is placed.",
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
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
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
			"rollback_from_backup_uuid": {
				Type:        schema.TypeString,
				Description: "Rollback the storage from a specific storage backup.",
				Optional:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
		Name:     d.Get("name").(string),
		Capacity: d.Get("capacity").(int),
		Labels:   &labels,
	}

	// Only distributed storage variant allows
	// to set storage type.
	storageVariant, _ := d.Get("storage_variant").(string)
	if storageVariant == "" || storageVariant == "distributed" {
		storageType := d.Get("storage_type").(string)
		if storageType == "storage" || storageType == "" {
			requestBody.StorageType = gsclient.DefaultStorageType
		} else if storageType == "storage_high" {
			requestBody.StorageType = gsclient.HighStorageType
		} else if storageType == "storage_insane" {
			requestBody.StorageType = gsclient.InsaneStorageType
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateStorage(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	//If rollback_from_backup_uuid is updated and set, start rolling back process
	if hasChanged := d.HasChange("rollback_from_backup_uuid"); hasChanged {
		if attr, ok := d.GetOk("rollback_from_backup_uuid"); ok {
			log.Printf("Start rolling back storage %s with backup %s", d.Id(), attr.(string))
			err = client.RollbackStorageBackup(
				ctx,
				d.Id(),
				attr.(string),
				gsclient.StorageRollbackRequest{
					Rollback: true,
				},
			)
			if err != nil {
				return err
			}
		}
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

	// Only distributed storage variant allows
	// to set storage type.
	storageVariant, _ := d.Get("storage_variant").(string)
	if storageVariant == "" || storageVariant == "distributed" {
		storageType := d.Get("storage_type").(string)
		if storageType == "storage" || storageType == "" {
			requestBody.StorageType = gsclient.DefaultStorageType
		} else if storageType == "storage_high" {
			requestBody.StorageType = gsclient.HighStorageType
		} else if storageType == "storage_insane" {
			requestBody.StorageType = gsclient.InsaneStorageType
		}
	} else if storageVariant == "local" {
		requestBody.StorageVariant = gsclient.LocalStorageVariant
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

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateStorage(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)
	log.Printf("The id for storage %s has been set to %v", requestBody.Name, response.ObjectUUID)

	//If rollback_from_backup_uuid is set, start rolling back process
	if attr, ok := d.GetOk("rollback_from_backup_uuid"); ok {
		log.Printf("Start rolling back storage %s with backup %s", response.ObjectUUID, attr.(string))
		err = client.RollbackStorageBackup(
			ctx,
			response.ObjectUUID,
			attr.(string),
			gsclient.StorageRollbackRequest{
				Rollback: true,
			},
		)
		if err != nil {
			return err
		}
	}
	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete storage (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	storage, err := client.GetStorage(ctx, d.Id())
	//In case of 404, don't catch the error
	if errHandler.SuppressHTTPErrorCodes(err, http.StatusNotFound) != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	//Stop all server relating to this IP address if there is one
	for _, server := range storage.Properties.Relations.Servers {
		unlinkStorageAction := func(ctx context.Context) error {
			//No need to unlink when server returns 409 or 404
			err := errHandler.SuppressHTTPErrorCodes(
				client.UnlinkStorage(ctx, server.ObjectUUID, d.Id()),
				http.StatusConflict,
				http.StatusNotFound,
			)
			return err
		}
		//UnlinkStorage requires the server to be off
		err = globalServerStatusList.runActionRequireServerOff(ctx, client, server.ObjectUUID, false, unlinkStorageAction)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	err = errHandler.SuppressHTTPErrorCodes(
		client.DeleteStorage(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
