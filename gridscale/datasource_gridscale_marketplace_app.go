package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleMarketplaceApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleMarketplaceApplicationRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Computed:    true,
			},
			"category": {
				Type:        schema.TypeString,
				Description: "Category of marketplace application. Accepted values: \"CMS\", \"project management\", \"Adminpanel\", \"Collaboration\", \"Cloud Storage\", \"Archiving\"",
				Computed:    true,
			},
			"object_storage_path": {
				Type:        schema.TypeString,
				Description: "Path to the images for the application, must be in .gz format and started with s3//",
				Computed:    true,
			},
			"setup_cores": {
				Type:        schema.TypeInt,
				Description: "Number of server's cores",
				Computed:    true,
			},
			"setup_memory": {
				Type:        schema.TypeInt,
				Description: "The capacity of server's memory in GB",
				Computed:    true,
			},
			"setup_storage_capacity": {
				Type:        schema.TypeInt,
				Description: "The capacity of server's storage in GB",
				Computed:    true,
			},
			"meta_license": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_components": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"meta_overview": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_hints": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_terms_of_use": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_icon": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_features": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_author": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta_advices": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_hash": {
				Type:        schema.TypeString,
				Description: "Unique hash to allow user to import the self-created marketplace application",
				Computed:    true,
			},
			"is_application_owner": {
				Type:        schema.TypeBool,
				Description: "Whether the you are the owner of application or not",
				Computed:    true,
			},
			"is_published": {
				Type:        schema.TypeBool,
				Description: "Whether the template is published by the partner to their tenant",
				Computed:    true,
			},
			"published_date": {
				Type:        schema.TypeString,
				Description: "The date when the template is published into other tenant in the same partner",
				Computed:    true,
			},
			"is_publish_requested": {
				Type:        schema.TypeBool,
				Description: "Whether the tenants want their template to be published or not",
				Computed:    true,
			},
			"publish_requested_date": {
				Type:        schema.TypeString,
				Description: "The date when the tenant requested their template to be published",
				Computed:    true,
			},
			"is_publish_global_requested": {
				Type:        schema.TypeBool,
				Description: "Whether a partner wants their tenant template published to other partners",
				Computed:    true,
			},
			"publish_global_requested_date": {
				Type:        schema.TypeString,
				Description: "The date when a partner requested their tenants template to be published",
				Computed:    true,
			},
			"is_publish_global": {
				Type:        schema.TypeBool,
				Description: "Whether a template is published to other partner or not",
				Computed:    true,
			},
			"published_global_date": {
				Type:        schema.TypeString,
				Description: "The date when a template is published to other partner",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of template",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time the object was initially created",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time of the last object change",
				Computed:    true,
			},
		},
	}
}

func dataSourceGridscaleMarketplaceApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read marketplace application (%s) datasource -", id)
	marketApp, err := client.GetMarketplaceApplication(context.Background(), id)
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	d.SetId(marketApp.Properties.ObjectUUID)

	if err = d.Set("name", marketApp.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("category", marketApp.Properties.Category); err != nil {
		return fmt.Errorf("%s error setting category: %v", errorPrefix, err)
	}
	if err = d.Set("object_storage_path", marketApp.Properties.ObjectStoragePath); err != nil {
		return fmt.Errorf("%s error setting object_storage_path: %v", errorPrefix, err)
	}
	if err = d.Set("setup_cores", marketApp.Properties.Setup.Cores); err != nil {
		return fmt.Errorf("%s error setting setup_cores: %v", errorPrefix, err)
	}
	if err = d.Set("setup_memory", marketApp.Properties.Setup.Memory); err != nil {
		return fmt.Errorf("%s error setting setup_memory: %v", errorPrefix, err)
	}
	if err = d.Set("setup_storage_capacity", marketApp.Properties.Setup.Capacity); err != nil {
		return fmt.Errorf("%s error setting setup_storage_capacity: %v", errorPrefix, err)
	}
	if err = d.Set("meta_license", marketApp.Properties.Metadata.License); err != nil {
		return fmt.Errorf("%s error setting meta_license: %v", errorPrefix, err)
	}
	if err = d.Set("meta_os", marketApp.Properties.Metadata.OS); err != nil {
		return fmt.Errorf("%s error setting meta_os: %v", errorPrefix, err)
	}
	if err = d.Set("meta_components", marketApp.Properties.Metadata.Components); err != nil {
		return fmt.Errorf("%s error setting meta_components: %v", errorPrefix, err)
	}
	if err = d.Set("meta_overview", marketApp.Properties.Metadata.Overview); err != nil {
		return fmt.Errorf("%s error setting meta_overview: %v", errorPrefix, err)
	}
	if err = d.Set("meta_hints", marketApp.Properties.Metadata.Hints); err != nil {
		return fmt.Errorf("%s error setting meta_hints: %v", errorPrefix, err)
	}
	if err = d.Set("meta_terms_of_use", marketApp.Properties.Metadata.TermsOfUse); err != nil {
		return fmt.Errorf("%s error setting meta_terms_of_use: %v", errorPrefix, err)
	}
	if err = d.Set("meta_icon", marketApp.Properties.Metadata.Icon); err != nil {
		return fmt.Errorf("%s error setting meta_icon: %v", errorPrefix, err)
	}
	if err = d.Set("meta_features", marketApp.Properties.Metadata.Features); err != nil {
		return fmt.Errorf("%s error setting meta_features: %v", errorPrefix, err)
	}
	if err = d.Set("meta_author", marketApp.Properties.Metadata.Authors); err != nil {
		return fmt.Errorf("%s error setting meta_author: %v", errorPrefix, err)
	}
	if err = d.Set("meta_advices", marketApp.Properties.Metadata.Advices); err != nil {
		return fmt.Errorf("%s error setting meta_advices: %v", errorPrefix, err)
	}
	if err = d.Set("unique_hash", marketApp.Properties.UniqueHash); err != nil {
		return fmt.Errorf("%s error setting unique_hash: %v", errorPrefix, err)
	}
	if err = d.Set("is_application_owner", marketApp.Properties.IsApplicationOwner); err != nil {
		return fmt.Errorf("%s error setting is_application_owner: %v", errorPrefix, err)
	}
	if err = d.Set("is_published", marketApp.Properties.Published); err != nil {
		return fmt.Errorf("%s error setting is_published: %v", errorPrefix, err)
	}
	if (marketApp.Properties.PublishedDate != gsclient.GSTime{}) {
		if err = d.Set("published_date", marketApp.Properties.PublishedDate); err != nil {
			return fmt.Errorf("%s error setting published_date: %v", errorPrefix, err)
		}
	}
	if err = d.Set("is_publish_requested", marketApp.Properties.PublishRequested); err != nil {
		return fmt.Errorf("%s error setting is_publish_requested: %v", errorPrefix, err)
	}
	if (marketApp.Properties.PublishRequestedDate != gsclient.GSTime{}) {
		if err = d.Set("publish_requested_date", marketApp.Properties.PublishRequestedDate); err != nil {
			return fmt.Errorf("%s error setting publish_requested_date: %v", errorPrefix, err)
		}
	}
	if err = d.Set("is_publish_global_requested", marketApp.Properties.PublishGlobalRequested); err != nil {
		return fmt.Errorf("%s error setting is_publish_global_requested: %v", errorPrefix, err)
	}
	if (marketApp.Properties.PublishGlobalRequestedDate != gsclient.GSTime{}) {
		if err = d.Set("publish_global_requested_date", marketApp.Properties.PublishGlobalRequestedDate); err != nil {
			return fmt.Errorf("%s error setting publish_global_requested_date: %v", errorPrefix, err)
		}
	}
	if err = d.Set("is_publish_global", marketApp.Properties.PublishedGlobal); err != nil {
		return fmt.Errorf("%s error setting is_publish_global: %v", errorPrefix, err)
	}
	if (marketApp.Properties.PublishedGlobalDate != gsclient.GSTime{}) {
		if err = d.Set("published_global_date", marketApp.Properties.PublishedGlobalDate); err != nil {
			return fmt.Errorf("%s error setting published_global_date: %v", errorPrefix, err)
		}
	}
	if err = d.Set("type", marketApp.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting type: %v", errorPrefix, err)
	}
	if err = d.Set("status", marketApp.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", marketApp.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", marketApp.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	return nil
}
