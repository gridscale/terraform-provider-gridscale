package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleSSLCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleSSLCertCreate,
		Read:   resourceGridscaleSSLCertRead,
		Delete: resourceGridscaleSSLCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
				ForceNew:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "The PEM-formatted private-key of the SSL certificate.",
				Required:    true,
				ForceNew:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace((val.(string)))
				},
			},
			"leaf_certificate": {
				Type:        schema.TypeString,
				Description: "The PEM-formatted public SSL of the SSL certificate.",
				Required:    true,
				ForceNew:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace((val.(string)))
				},
			},
			"certificate_chain": {
				Type:        schema.TypeString,
				Description: "The PEM-formatted full-chain between the certificate authority and the domain's SSL certificate.",
				Optional:    true,
				ForceNew:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace((val.(string)))
				},
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
							Optional:    true,
							Description: "MD5 fingerprint of the certificate.",
						},
						"sha256": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SHA256 fingerprint of the certificate.",
						},
						"sha1": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SHA1 fingerprint of the certificate.",
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				ForceNew:    true,
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

func resourceGridscaleSSLCertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read SSL Certificate (%s) resource -", d.Id())
	cert, err := client.GetSSLCertificate(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

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

func resourceGridscaleSSLCertCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	privKey := d.Get("private_key").(string)
	leafCert := d.Get("leaf_certificate").(string)
	certChain := d.Get("certificate_chain").(string)
	requestBody := gsclient.SSLCertificateCreateRequest{
		Name:             d.Get("name").(string),
		PrivateKey:       strings.TrimSpace(privKey),
		LeafCertificate:  strings.TrimSpace(leafCert),
		CertificateChain: strings.TrimSpace(certChain),
		Labels:           convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateSSLCertificate(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new SSL Certificate has been set to %v", response.ObjectUUID)

	return resourceGridscaleSSLCertRead(d, meta)
}

func resourceGridscaleSSLCertDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete SSL Certificate (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeleteSSLCertificate(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
