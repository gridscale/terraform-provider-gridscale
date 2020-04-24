package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gridscale/gsclient-go/v2"
)

type serverStatus struct {
	//is the server is deleted
	deleted bool

	//mux for one server
	mux sync.Mutex

	//wait group is used to make sure all goroutines updating a server (requiring server to be off)
	//are done before turn a server back on
	wg sync.WaitGroup
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
			//set the shutdown timeout specifically
			shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeoutSecs*time.Second)
			defer cancel()
			err := c.ShutdownServer(shutdownCtx, id)
			//if error is returned and it is not caused by an expired context, returns error
			if err != nil && err != shutdownCtx.Err() {
				return err
			}
			// if the server cannot be shutdown gracefully, try to turn it off
			if err == shutdownCtx.Err() {
				//check if the main context is done
				select {
				//return context's error when it is done
				case <-ctx.Done():
					return ctx.Err()
				default:
				}
				//force the sever to stop
				err = c.StopServer(ctx, id)
				if err != nil {
					return err
				}
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
	if s, ok := l.list[id]; ok {
		//lock the server
		s.mux.Lock()
		log.Printf("[DEBUG] LOCK ACQUIRED to stop server (%v)", id)
		defer func() {
			//unlock the server
			s.mux.Unlock()
			log.Printf("[DEBUG] LOCK RELEASED! Shutting down server (%v) is done", id)
		}()
		if !s.deleted {
			//set the shutdown timeout specifically
			shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeoutSecs*time.Second)
			defer cancel()
			err := c.ShutdownServer(shutdownCtx, id)
			//if error is returned and it is not caused by an expired context, returns error
			if err != nil && err != shutdownCtx.Err() {
				return err
			}
			// if the server cannot be shutdown gracefully, try to turn it off
			if err == shutdownCtx.Err() {
				//check if the main context is done
				select {
				//return context's error when it is done
				case <-ctx.Done():
					return ctx.Err()
				default:
				}
				//force the sever to stop
				return c.StopServer(ctx, id)
			}
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
	//check if the server is in the list
	if s, ok := l.list[id]; ok {
		var err error
		//Add 1 to wait group of the server before running the action
		s.wg.Add(1)
		//Get the original server's state
		server, err := c.GetServer(ctx, id)
		if err != nil {
			//Tell the wait group that the action is done
			s.wg.Done()
			if reqError, ok := err.(gsclient.RequestError); ok {
				//if server is not found
				if reqError.StatusCode == http.StatusNotFound {
					//action does not need to be run,
					//if server does not present and it is not required to run the action
					if !serverRequired {
						return nil
					}
				}
			}
			return err
		}
		//if the server is on, shutdown the server (synchronously) before running the action,
		//and start the server after finishing the action.
		if server.Properties.Power {
			//shut down the server synchronously
			//If we don't turn it off synchronously, all server-update goroutines (requiring server to be off)
			//will send their shutdown requests at the same time. That causes false assumption error returned from
			//gridscale backend (as the server is being turned off by the first request).
			err = l.shutdownServerSynchronously(ctx, c, id)
			if err != nil {
				//Tell the wait group that the action is done
				s.wg.Done()
				return err
			}
			log.Printf("[DEBUG] Server (%v) is OFF to run an action", id)
			defer func() {
				//wait group of the server blocks until all actions finish
				s.wg.Wait()
				//start a server synchronously. Same explanation as why use `shutdownServerSynchronously` above
				errStartServer := l.startServerSynchronously(ctx, c, id)
				if errStartServer != nil {
					//append error from the action (if the action returns error)
					err = fmt.Errorf(
						"Error from action: %v. Error from starting server: %v",
						err,
						errStartServer,
					)
				}
			}()
		}
		err = action(ctx)
		//Tell the wait group that the action is done
		s.wg.Done()
		return err
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
		status := &serverStatus{}
		globalServerStatusList.list[uuid] = status
	}
	return nil
}

//globalServerStatusList global list of all servers' status states in terraform
var globalServerStatusList = serverStatusList{
	list: make(map[string]*serverStatus),
}
