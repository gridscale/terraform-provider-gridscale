package gridscale

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	version string
	commit  string
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_UUID", nil),
				Description: "User-UUID for the gridscale API.",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_TOKEN", nil),
				Description: "API-token for the gridscale API.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_URL", nil),
				Description: "the url for the gridscale API.",
			},
			"http_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRIDSCALE_TF_HEADERS", nil),
				Description: "Custom HTTP headers",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gridscale_server":                   dataSourceGridscaleServer(),
			"gridscale_storage":                  dataSourceGridscaleStorage(),
			"gridscale_network":                  dataSourceGridscaleNetwork(),
			"gridscale_public_network":           dataSourceGridscalePublicNetwork(),
			"gridscale_ipv4":                     dataSourceGridscaleIpv4(),
			"gridscale_ipv6":                     dataSourceGridscaleIpv6(),
			"gridscale_sshkey":                   dataSourceGridscaleSshkey(),
			"gridscale_template":                 dataSourceGridscaleTemplate(),
			"gridscale_loadbalancer":             dataSourceGridscaleLoadBalancer(),
			"gridscale_snapshot":                 dataSourceGridscaleStorageSnapshot(),
			"gridscale_backup_list":              dataSourceGridscaleStorageBackupList(),
			"gridscale_snapshotschedule":         dataSourceGridscaleStorageSnapshotSchedule(),
			"gridscale_backupschedule":           dataSourceGridscaleStorageBackupSchedule(),
			"gridscale_paas":                     dataSourceGridscalePaaS(),
			"gridscale_paas_securityzone":        dataSourceGridscalePaaSSecurityZone(),
			"gridscale_object_storage_accesskey": dataSourceGridscaleObjectStorage(),
			"gridscale_isoimage":                 dataSourceGridscaleISOImage(),
			"gridscale_firewall":                 dataSourceGridscaleFirewall(),
			"gridscale_marketplace_application":  dataSourceGridscaleMarketplaceApplication(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gridscale_server":                         resourceGridscaleServer(),
			"gridscale_storage":                        resourceGridscaleStorage(),
			"gridscale_storage_clone":                  resourceGridscaleStorageClone(),
			"gridscale_storage_import":                 resourceGridscaleStorageImport(),
			"gridscale_network":                        resourceGridscaleNetwork(),
			"gridscale_ipv4":                           resourceGridscaleIpv4(),
			"gridscale_ipv6":                           resourceGridscaleIpv6(),
			"gridscale_sshkey":                         resourceGridscaleSshkey(),
			"gridscale_loadbalancer":                   resourceGridscaleLoadBalancer(),
			"gridscale_snapshot":                       resourceGridscaleStorageSnapshot(),
			"gridscale_snapshotschedule":               resourceGridscaleStorageSnapshotSchedule(),
			"gridscale_backupschedule":                 resourceGridscaleStorageBackupSchedule(),
			"gridscale_paas":                           resourceGridscalePaaS(),
			"gridscale_k8s":                            resourceGridscaleK8s(),
			"gridscale_paas_securityzone":              resourceGridscalePaaSSecurityZone(),
			"gridscale_object_storage_accesskey":       resourceGridscaleObjectStorage(),
			"gridscale_template":                       resourceGridscaleTemplate(),
			"gridscale_isoimage":                       resourceGridscaleISOImage(),
			"gridscale_firewall":                       resourceGridscaleFirewall(),
			"gridscale_marketplace_application":        resourceGridscaleMarketplaceApplication(),
			"gridscale_marketplace_application_import": resourceGridscaleImportedMarketplaceApplication(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	headers := convertStrToHeaderMap(d.Get("http_headers").(string))
	headers["User-Agent"] = fmt.Sprintf("terraform-provider-gridscale/%s-%s-%s", version, commit, runtime.GOOS)
	config := Config{
		UserUUID:    d.Get("uuid").(string),
		APIToken:    d.Get("token").(string),
		APIUrl:      d.Get("api_url").(string),
		HTTPHeaders: headers,
	}

	return config.Client()
}

// getHeaderMapFromStr converts string (format: "key1:val1,key2:val2")
// to a HTTP header map
func convertStrToHeaderMap(str string) map[string]string {
	result := make(map[string]string)
	// split string into comma separated headers
	headers := strings.Split(str, ",")
	for _, header := range headers {
		if header != "" {
			// split each header into a key and a value
			kv := strings.Split(header, ":")
			if len(kv) == 2 {
				result[kv[0]] = kv[1]
			}
		}
	}
	return result
}
