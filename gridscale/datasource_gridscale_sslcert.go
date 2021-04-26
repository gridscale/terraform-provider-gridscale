package gridscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go/v3"
)

func dataSourceGridscaleSSLCert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleSSLCertRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a SSL certificate resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Computed:    true,
			},
			"common_name": {
				Type:        schema.TypeString,
				Description: "The common domain name of the SSL certificate.",
				Computed:    true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change.",
				Computed:    true,
			},
			"not_valid_after": {
				Type:        schema.TypeString,
				Description: "Defines the date after which the certificate is not valid.",
				Computed:    true,
			},
			"fingerprints": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Description: "Defines a list of unique identifiers generated from the MD5, SHA-1, and SHA-256 fingerprints of the certificate.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"md5": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MD5 fingerprint of the certificate.",
						},
						"sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SHA256 fingerprint of the certificate.",
						},
						"sha1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SHA1 fingerprint of the certificate.",
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func dataSourceGridscaleSSLCertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read SSL certificate (%s) datasource -", id)

	cert, err := client.GetSSLCertificate(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(cert.Properties.ObjectUUID)

	if err = d.Set("name", cert.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("common_name", cert.Properties.CommonName); err != nil {
		return fmt.Errorf("%s error setting common_name: %v", errorPrefix, err)
	}
	if err = d.Set("status", cert.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", cert.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", cert.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("not_valid_after", cert.Properties.NotValidAfter.String()); err != nil {
		return fmt.Errorf("%s error setting not_valid_after: %v", errorPrefix, err)
	}
	if err = d.Set("labels", cert.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	fingerprints := []interface{}{
		map[string]interface{}{
			"md5":    cert.Properties.Fingerprints.MD5,
			"sha256": cert.Properties.Fingerprints.SHA256,
			"sha1":   cert.Properties.Fingerprints.SHA1,
		},
	}
	if err = d.Set("fingerprints", fingerprints); err != nil {
		return fmt.Errorf("%s error setting fingerprints: %v", errorPrefix, err)
	}
	return nil
}
