package gridscale

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"template": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Contents of the template",
				ConflictsWith: []string{"filename"},
			},
			"filename": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "file to read template from",
				// Make a "best effort" attempt to relativize the file path.
				StateFunc: func(v interface{}) string {
					if v == nil || v.(string) == "" {
						return ""
					}
					pwd, err := os.Getwd()
					if err != nil {
						return v.(string)
					}
					rel, err := filepath.Rel(pwd, v.(string))
					if err != nil {
						return v.(string)
					}
					return rel
				},
				Deprecated:    "Use the 'template' attribute instead.",
				ConflictsWith: []string{"template"},
			},
			"vars": &schema.Schema{
				Type:         schema.TypeMap,
				Optional:     true,
				Default:      make(map[string]interface{}),
				Description:  "variables to substitute",
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered template",
			},
		},
	}
}

func dataSourceFileRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}