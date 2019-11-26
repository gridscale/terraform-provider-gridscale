package gridscale

import (
	"context"
	"fmt"
	"github.com/gridscale/gsclient-go"
	"log"
	"net/http"
	"sync"
)

//serverMutex is the struct holding a mutex.
//This mutex is used to solve the terraform problem when it tries
//to deal with a server and its dependencies at the same time
type serverMutex struct {
	mux sync.Mutex
}

//actionRequireServerOff signature of a function that requires a server to be off
//in order to run
type actionRequireServerOff func(ctx context.Context) error

//removeServerSynchronously removes a server
func (m *serverMutex) removeServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock mutex
	m.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to remove server (%v)", id)
	defer func() {
		//unlock mutex
		m.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Server (%v) is removed", id)
	}()
	//Shutdown server
	err := c.ShutdownServer(ctx, id)
	if err != nil {
		return err
	}
	return c.DeleteServer(ctx, id)
}

//getServerPowerStatus returns power state of a server in the list (synchronously)
func (m *serverMutex) getServerPowerStatus(ctx context.Context, c *gsclient.Client, id string) (bool, error) {
	server, err := c.GetServer(ctx, id)
	return server.Properties.Power, err
}

//startServerSynchronously starts the servers synchronously. That means the server
//can only be started by one goroutine at a time.
func (m *serverMutex) startServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock mutex
	m.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to start server (%v)", id)
	defer func() {
		//unlock mutex
		m.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Starting server (%v) is done", id)
	}()
	return c.StartServer(ctx, id)
}

//shutdownServerSynchronously stop the servers synchronously. That means the server
//can only be stopped by one goroutine at a time.
func (m *serverMutex) shutdownServerSynchronously(ctx context.Context, c *gsclient.Client, id string) error {
	//lock mutex
	m.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to stop server (%v)", id)
	defer func() {
		//unlock mutex
		m.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Shutting down server (%v) is done", id)
	}()
	return c.ShutdownServer(ctx, id)
}

//runActionRequireServerOff runs a specific action (function) after shutting down (synchronously) the server successfully.
//Some actions are NOT necessary if the server is already deleted (such as `UnlinkXXX` methods), that means they do
//not need the server to exist strictly.
//However, the are still some others requiring the server's presence (such as some server update sequences).
func (m *serverMutex) runActionRequireServerOff(
	ctx context.Context,
	c *gsclient.Client,
	id string,
	serverRequired bool,
	action actionRequireServerOff) error {
	//lock mutex
	m.mux.Lock()
	log.Printf("[DEBUG] LOCK ACQUIRED to run an action requiring server (%v) to be OFF", id)
	defer func() {
		//unlock mutex
		m.mux.Unlock()
		log.Printf("[DEBUG] LOCK RELEASED! Action requiring server (%v) is done", id)
	}()
	var err error
	//Get the original server's power state
	originPowerState, err := m.getServerPowerStatus(ctx, c, id)
	if err != nil {
		if reqError, ok := err.(gsclient.RequestError); ok {
			//if server is not found
			if reqError.StatusCode == http.StatusNotFound {
				//if server is required to run the action
				//return error
				if serverRequired {
					return err
				} else {
					//action does not need to be run,
					//if server does not present and it is not required to run the action
					return nil
				}
			}
		}
		return err
	}
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
	err = action(ctx)
	return err
}

//globalServerMutex init a global server mutex instance
var globalServerMutex = serverMutex{}
