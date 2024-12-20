package boot

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"syscall"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/index"
	"google.golang.org/grpc"
)

// PluginManager manages plugin processes and their gRPC clients
type PluginManager struct {
	sync.RWMutex
	Clients   map[string]*DataQClient
	processes map[string]*exec.Cmd
	index     *index.Index
	cas       cas.Storage
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(idx *index.Index, cas cas.Storage) *PluginManager {
	return &PluginManager{
		Clients:   make(map[string]*DataQClient),
		processes: make(map[string]*exec.Cmd),
		index:     idx,
		cas:       cas,
	}
}

// GetClient returns the gRPC client for a plugin
func (pm *PluginManager) GetClient(pluginID string) (*DataQClient, error) {
	pm.RLock()
	defer pm.RUnlock()

	client, exists := pm.Clients[pluginID]
	if !exists {
		return nil, fmt.Errorf("no client found for plugin %s", pluginID)
	}
	return client, nil
}

// AddPlugin adds a new plugin process and client
func (pm *PluginManager) AddPlugin(pluginID string, conn *grpc.ClientConn, process *exec.Cmd) {
	pm.Lock()
	defer pm.Unlock()

	pm.Clients[pluginID] = NewDataQClient(conn, pm.index, pm.cas)
	pm.processes[pluginID] = process
}

// Shutdown gracefully stops all plugin processes
func (pm *PluginManager) Shutdown() {
	pm.Lock()
	defer pm.Unlock()

	for id, proc := range pm.processes {
		if proc != nil && proc.Process != nil {
			if err := proc.Process.Signal(syscall.SIGTERM); err != nil {
				log.Printf("Failed to send SIGTERM to plugin %s: %v", id, err)
				proc.Process.Kill()
			}
		}
	}
}
