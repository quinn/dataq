package boot

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"syscall"

	rpc "go.quinn.io/dataq/rpc"
)

// PluginManager manages plugin processes and their gRPC clients
type PluginManager struct {
	sync.RWMutex
	Clients   map[string]rpc.DataQPluginClient
	processes map[string]*exec.Cmd
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		Clients:   make(map[string]rpc.DataQPluginClient),
		processes: make(map[string]*exec.Cmd),
	}
}

// GetClient returns the gRPC client for a plugin
func (pm *PluginManager) GetClient(pluginID string) (rpc.DataQPluginClient, error) {
	pm.RLock()
	defer pm.RUnlock()

	client, exists := pm.Clients[pluginID]
	if !exists {
		return nil, fmt.Errorf("no client found for plugin %s", pluginID)
	}
	return client, nil
}

// AddPlugin adds a new plugin process and client
func (pm *PluginManager) AddPlugin(pluginID string, client rpc.DataQPluginClient, process *exec.Cmd) {
	pm.Lock()
	defer pm.Unlock()

	pm.Clients[pluginID] = client
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
