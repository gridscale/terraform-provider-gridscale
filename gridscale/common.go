package gridscale

import (
	"context"
	"fmt"
	"github.com/gridscale/gsclient-go"
	"log"
	"sync"
)

//serverPowerStatus represents power state of a server
//mutex is used to lock the resource when a goroutine changes/reads a server's
//power state, so it prevents other goroutines from accessing/modifying the server's power state
type serverPowerStatus struct {
	power bool
	mux   sync.Mutex
}

//listServersPowerStatus represents a list of power states of
//all servers (declared in terraform).
//mutex is used to lock when adding or removing servers.
//***NOTE: servers declared outside terraform are not included.
type listServersPowerStatus struct {
	list map[string]*serverPowerStatus
	mux  sync.Mutex
}

//actionRequireServerOff signature of a function that requires a server to be off
//in order to run
type actionRequireServerOff func(ctx context.Context) error

//addServer adds a server power state to the list
func (l *listServersPowerStatus) addServer(id string) error {
	//check if the server is already in the list
	if _, ok := l.list[id]; !ok {
		//lock the list
		l.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to add server (%v)", id)
		defer func() {
			//unlock the list
			l.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is added", id)
		}()
		l.list[id] = &serverPowerStatus{
			false,
			sync.Mutex{},
		}
		return nil
	}
	return fmt.Errorf("server (%s) ALREADY exists in current list of servers in terraform", id)
}

//removeServer removes a server power state from the list
func (l *listServersPowerStatus) removeServer(id string) error {
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//lock the list
		l.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to remove server (%v)", id)
		defer func() {
			//unlock the list
			l.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is removed", id)
		}()
		delete(l.list, id)
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//getServerPowerStatus returns power state of a server in the list (synchronously)
func (l *listServersPowerStatus) getServerPowerStatus(id string) (bool, error) {
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//lock the power state of the server
		l.list[id].mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to get server (%v) power status", id)
		defer func() {
			//unlock the power state of the server after retrieving server's power state
			l.list[id].mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Getting server (%v) power status is done", id)
		}()
		return l.list[id].power, nil
	}
	return false, fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//startServerSynchronously starts the servers synchronously. That means the server
//can only be started by one goroutine at a time.
func (l *listServersPowerStatus) startServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//lock the power state of the server
		l.list[id].mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to start server (%v)", id)
		defer func() {
			//unlock the power state of the server after server is started
			l.list[id].mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is started", id)
		}()
		err := c.StartServer(ctx, id)
		if err != nil {
			return err
		}
		l.list[id].power = true
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//shutdownServerSynchronously stop the servers synchronously. That means the server
//can only be stopped by one goroutine at a time.
func (l *listServersPowerStatus) shutdownServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//lock the power state of the server
		l.list[id].mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to stop server (%v)", id)
		defer func() {
			//unlock the power state of the server after server is stopped
			l.list[id].mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is stopped", id)
		}()
		err := c.ShutdownServer(ctx, id)
		if err != nil {
			return err
		}
		l.list[id].power = false
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//runActionRequireServerOff runs a specific action (function) after shutting down (synchronously) the server successfully
func (l *listServersPowerStatus) runActionRequireServerOff(ctx context.Context, c *gsclient.Client, id string, action actionRequireServerOff) error {
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//lock the power state of the server
		l.list[id].mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to run an action requiring server (%v) to be OFF", id)
		defer func() {
			//unlock the power state of the server after action is done
			l.list[id].mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Action requiring server (%v) is done", id)
		}()
		//shutdown the server (synchronously) before running the action
		err := c.ShutdownServer(ctx, id)
		if err != nil {
			return err
		}
		l.list[id].power = false
		//run action function
		err = action(ctx)
		return err
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//serverPowerStateList global list of all servers' power states in terraform
var serverPowerStateList = listServersPowerStatus{
	list: make(map[string]*serverPowerStatus),
}

//convSOStrings converts slice of interfaces to slice of strings
func convSOStrings(interfaceList []interface{}) []string {
	var labels []string
	for _, labelInterface := range interfaceList {
		labels = append(labels, labelInterface.(string))
	}
	return labels
}
