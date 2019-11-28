package gridscale

import (
	"context"
	"fmt"
	"github.com/gridscale/gsclient-go"
	"log"
	"sync"
)

type serverStatus struct {
	//power state of a server
	power bool

	//is the server is deleted
	deleted bool

	//counter of actions requiring the server to be off
	//when each action finishes, the counter is decreased by 1
	specialActionCounter int

	//mux for one server
	mux sync.Mutex
}

//serverStatusList represents a list of power states of
//all servers (declared in terraform).
//mutex is used to lock when adding or removing servers or modifying servers' power states.
//***NOTE: servers declared outside terraform are not included.
type serverStatusList struct {
	list map[string]*serverStatus
	mux  sync.Mutex
}

//actionRequireServerOff signature of a function that requires a server to be off
//in order to run
type actionRequireServerOff func(ctx context.Context) error

//addServer adds a server power state to the list
func (l *serverStatusList) addServer(id string) error {
	//lock the list
	//*Note: we don't need to lock the list anywhere else as
	//addServer always runs first
	l.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to add server (%v)", id)
	defer func() {
		//unlock the list
		l.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is added", id)
	}()
	//check if the server is already in the list
	if _, ok := l.list[id]; !ok {
		l.list[id] = &serverStatus{}
		return nil
	}
	return fmt.Errorf("server (%s) ALREADY exists in current list of servers in terraform", id)
}

//removeServerSynchronously removes a server and set `deleted` to true
//when `terraform apply` command finishes, the serverStatusList will be automatically flushed
func (l *serverStatusList) removeServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//check if the server is in the list and it is not deleted
	if s, ok := l.list[id]; ok {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to remove server (%v)", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is removed", id)
		}()
		if !s.deleted {
			err := c.ShutdownServer(ctx, id)
			if err != nil {
				return err
			}
			//Delete server
			err = c.DeleteServer(ctx, id)
			if err != nil {
				return err
			}
			s.deleted = true
			return nil
		}
		return fmt.Errorf("server (%s) is already deleted", id)
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//getServerPowerStatus returns power state of a server in the list (synchronously)
func (l *serverStatusList) getServerPowerStatus(id string) (bool, error) {
	//check if the server is in the list and it is not deleted
	if s, ok := l.list[id]; ok {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to get server (%v) power status", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Getting server (%v) power status is done", id)
		}()
		if !s.deleted {
			return s.power, nil
		}
		return false, fmt.Errorf("server (%s) is already deleted", id)
	}
	return false, fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//startServerSynchronously starts the servers synchronously. That means the server
//can only be started by one goroutine at a time.
func (l *serverStatusList) startServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//check if the server is in the list and it is not deleted
	if s, ok := l.list[id]; ok {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to start server (%v)", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Starting server (%v) is done", id)
		}()
		if !s.deleted {
			err := c.StartServer(ctx, id)
			if err != nil {
				return err
			}
			s.power = true
			return nil
		}
		return fmt.Errorf("server (%s) is already deleted", id)
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//shutdownServerSynchronously stop the servers synchronously. That means the server
//can only be stopped by one goroutine at a time.
func (l *serverStatusList) shutdownServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//check if the server is in the list and it is not deleted
	if s, ok := l.list[id]; ok && !l.list[id].deleted {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to stop server (%v)", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Shutting down server (%v) is done", id)
		}()
		if !s.deleted {
			err := c.ShutdownServer(ctx, id)
			if err != nil {
				return err
			}
			s.power = false
			return nil
		}
		return fmt.Errorf("server (%s) is already deleted", id)
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//runActionRequireServerOff runs a specific action (function) after shutting down (synchronously) the server successfully.
//Some actions are NOT necessary if the server is already deleted (such as `UnlinkXXX` methods), that means they do
//not need the server to exist strictly.
//However, the are still some others requiring the server's presence (such as some server update sequences).
func (l *serverStatusList) runActionRequireServerOff(
	ctx context.Context,
	c *gsclient.Client,
	id string,
	serverRequired bool,
	action actionRequireServerOff) error {
	var err error
	//check if the server is in the list and it is not deleted
	if s, ok := l.list[id]; ok {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to run an action requiring server (%v) to be OFF", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Action requiring server (%v) is done", id)
		}()
		if !s.deleted {
			//Get the original server's power state
			originPowerState := s.power
			//if the server is on, shutdown the server (synchronously) before running the action,
			//and start the server after finishing the action.
			if originPowerState {
				err = c.ShutdownServer(ctx, id)
				if err != nil {
					return err
				}
				log.Printf("[DEBUG] Server (%v) is OFF to run an action", id)
				//run action function
				err = action(ctx)
				if err != nil {
					return err
				}
				//Start the server (start the server after the action is done)
				err = c.StartServer(ctx, id)
				if err != nil {
					return err
				}
				log.Printf("[DEBUG] Action is done. Server (%v) is ON", id)
			}
			return nil
		}
		//if server must presents to run the action
		//return error
		if serverRequired {
			return fmt.Errorf("server (%s) is already deleted", id)
		}
		return nil
	}
	return fmt.Errorf("server (%s) does not exist in current list of servers in terraform", id)
}

//initGlobalServerStatusList fetches server list and init `globalServerStatusList`
func initGlobalServerStatusList(ctx context.Context, c *gsclient.Client) error {
	servers, err := c.GetServerList(ctx)
	if err != nil {
		return err
	}
	for _, server := range servers {
		uuid := server.Properties.ObjectUUID
		props := server.Properties
		status := &serverStatus{
			power:                props.Power,
			specialActionCounter: 0,
			mux:                  sync.Mutex{},
		}
		globalServerStatusList.list[uuid] = status
	}
	return nil
}

//globalServerStatusList global list of all servers' status states in terraform
var globalServerStatusList = serverStatusList{
	list: make(map[string]*serverStatus),
}
