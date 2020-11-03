package gridscale

import (
	"context"
	"log"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGridscaleImportedMarketplaceApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleMarketplaceApplicationImport,
		Read:   resourceGridscaleMarketplaceApplicationRead,
		Delete: resourceGridscaleMarketplaceApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"import_unique_hash": {
				Type:        schema.TypeString,
				Description: "Hash of a specific marketplace application that you want to import",
				Required:    true,
				ForceNew:    true,
			},
			"unique_hash": {
				Type:        schema.TypeString,
				Description: "Unique hash to allow user to import the self-created marketplace application",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleMarketplaceApplicationImport(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.ImportMarketplaceApplication(ctx, gsclient.MarketplaceApplicationImportRequest{
		UniqueHash: d.Get("import_unique_hash").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)
	log.Printf("The id of the imported marketplace application has been set to %v", response.ObjectUUID)

	return resourceGridscaleMarketplaceApplicationRead(d, meta)
}
