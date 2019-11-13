package gridscale

import (
	"context"
	"fmt"
	"github.com/gridscale/gsclient-go"
	"log"
	"sync"
)

//listServersPowerStatus represents a list of power states of
//all servers (declared in terraform).
//mutex is used to lock when adding or removing servers or modifying servers' power states.
//***NOTE: servers declared outside terraform are not included.
type listServersPowerStatus struct {
	list map[string]bool
	mux  sync.Mutex
}

//actionRequireServerOff signature of a function that requires a server to be off
//in order to run
type actionRequireServerOff func(ctx context.Context) error

//addServer adds a server power state to the list
func (l *listServersPowerStatus) addServer(id string) error {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to add server (%v)", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is added", id)
	}()
	//check if the server is already in the list
	if _, ok := l.list[id]; !ok {
		l.list[id] = false
		return nil
	}
	return fmt.Errorf("server (%s) ALREADY exists in current list of servers in terraform", id)
}

//removeServerSynchronously removes a server
func (l *listServersPowerStatus) removeServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to remove server (%v)", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is removed", id)
	}()
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		//Shutdown server
		err := c.ShutdownServer(ctx, id)
		if err != nil {
			return err
		}
		//Delete server
		err = c.DeleteServer(ctx, id)
		if err != nil {
			return err
		}
		delete(l.list, id)
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform 1", id)
}

//getServerPowerStatus returns power state of a server in the list (synchronously)
func (l *listServersPowerStatus) getServerPowerStatus(id string) (bool, error) {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to get server (%v) power status", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Getting server (%v) power status is done", id)
	}()
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		return l.list[id], nil
	}
	return false, fmt.Errorf("server (%s) does not exist in current list of servers in terraform 2", id)
}

//startServerSynchronously starts the servers synchronously. That means the server
//can only be started by one goroutine at a time.
func (l *listServersPowerStatus) startServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to start server (%v)", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Starting server (%v) is done", id)
	}()
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		err := c.StartServer(ctx, id)
		if err != nil {
			return err
		}
		l.list[id] = true
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform 3", id)
}

//shutdownServerSynchronously stop the servers synchronously. That means the server
//can only be stopped by one goroutine at a time.
func (l *listServersPowerStatus) shutdownServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to stop server (%v)", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Shutting down server (%v) is done", id)
	}()
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		err := c.ShutdownServer(ctx, id)
		if err != nil {
			return err
		}
		l.list[id] = false
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform 4", id)
}

//runActionRequireServerOff runs a specific action (function) after shutting down (synchronously) the server successfully.
//Some actions are NOT necessary if the server is already deleted (such as `UnlinkXXX` methods), that means they do
//not need the server to exist strictly.
//However, the are still some others requiring the server's presence (such as some server update sequences).
func (l *listServersPowerStatus) runActionRequireServerOff(
	ctx context.Context,
	c *gsclient.Client,
	id string,
	serverRequired bool,
	action actionRequireServerOff) error {
	//lock the list
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to run an action requiring server (%v) to be OFF", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Action requiring server (%v) is done", id)
	}()
	//check if the server is in the list
	if _, ok := l.list[id]; ok {
		var err error
		//Get the original server's power state
		originPowerState := l.list[id]
		//if the server is on, shutdown the server (synchronously) before running the action,
		//and start the server after finishing the action.
		if originPowerState {
			err = c.ShutdownServer(ctx, id)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Server (%v) is OFF to run an action", id)
			//DEFER Start the server (start the server after the action is done)
			defer func() {
				errStartServer := c.StartServer(ctx, id)
				if errStartServer != nil {
					//append error from the action (if the action returns error)
					err = fmt.Errorf(
						"Error from action: %v. Error from starting server: %v",
						err,
						errStartServer,
					)
				}
				log.Printf("[DEBUG] Action is done. Server (%v) is ON", id)
			}()
		}
		//run action function
		err = action(ctx)
		return err
	}
	if serverRequired {
		err := fmt.Errorf("server (%s) does not exist in current list of servers in terraform 5", id)
		return err
	}
	return nil
}

//initServerPowerStateList fetches server list and init `serverPowerStateList`
func initServerPowerStateList(ctx context.Context, c *gsclient.Client) error {
	servers, err := c.GetServerList(ctx)
	if err != nil {
		return err
	}
	for _, server := range servers {
		uuid := server.Properties.ObjectUUID
		powerState := server.Properties.Power
		serverPowerStateList.list[uuid] = powerState
	}
	return nil
}

//serverPowerStateList global list of all servers' power states in terraform
var serverPowerStateList = listServersPowerStatus{
	list: make(map[string]bool),
}
