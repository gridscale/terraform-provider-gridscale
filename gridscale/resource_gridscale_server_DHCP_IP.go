package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
)

const serverDHCPIPUUIDDelimeter = "&"

func resourceGridscaleServerDHCPIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleServerDHCPIPCreate,
		Read:   resourceGridscaleServerDHCPIPRead,
		Delete: resourceGridscaleServerDHCPIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"server_uuid": {
				Type:        schema.TypeString,
				Description: "The server UUID which will be assigned a DHCP IP.",
				Required:    true,
				ForceNew:    true,
			},
			"network_uuid": {
				Type:        schema.TypeString,
				Description: "The DHCP-enabled network UUID which the server resides in.",
				Required:    true,
				ForceNew:    true,
			},
			"ip": {
				Type:        schema.TypeString,
				Description: "The IP address which will be assigne to the server.",
				Required:    true,
				ForceNew:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleServerDHCPIPRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	idList := strings.Split(id, serverDHCPIPUUIDDelimeter)
	networkUUID := idList[0]
	serverUUID := idList[1]

	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read Server-DHCP IP assignment (%s) resource -", d.Id())
	pinnedServerList, err := client.GetPinnedServerList(context.Background(), networkUUID)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	var found bool
	var foundObject gsclient.ServerWithIP
	for _, pinnedServer := range pinnedServerList.List {
		if pinnedServer.ServerUUID == serverUUID {
			foundObject = pinnedServer
			found = true
		}
	}
	if !found {
		return fmt.Errorf("%s error: %s", errorPrefix, "server-DHCP IP assignment is NOT found")
	}
	if err = d.Set("server_uuid", foundObject.ServerUUID); err != nil {
		return fmt.Errorf("%s error setting server_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("network_uuid", networkUUID); err != nil {
		return fmt.Errorf("%s error setting network_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("ip", foundObject.IP); err != nil {
		return fmt.Errorf("%s error setting ip: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleServerDHCPIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	serverUUID := d.Get("server_uuid").(string)
	networkUUID := d.Get("network_uuid").(string)
	ip := d.Get("ip").(string)
	errorPrefix := fmt.Sprintf("assign IP %s to server %s in network %s -", ip, serverUUID, networkUUID)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	err := client.UpdateNetworkPinnedServer(ctx, networkUUID, serverUUID, gsclient.PinServerRequest{
		IP: ip,
	})
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", networkUUID, serverDHCPIPUUIDDelimeter, serverUUID, serverDHCPIPUUIDDelimeter, ip))
	log.Printf("Successfully assign IP %s to server %s in network %s", ip, serverUUID, networkUUID)
	return resourceGridscaleServerDHCPIPRead(d, meta)
}

func resourceGridscaleServerDHCPIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	serverUUID := d.Get("server_uuid").(string)
	networkUUID := d.Get("network_uuid").(string)
	ip := d.Get("ip").(string)
	errorPrefix := fmt.Sprintf("remove IP %s from server %s in network %s -", ip, serverUUID, networkUUID)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeleteNetworkPinnedServer(ctx, networkUUID, serverUUID),
		http.StatusConflict,
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	log.Printf("Successfully remove IP %s from server %s in network %s", ip, serverUUID, networkUUID)
	return nil
}
