package boot

import (
	"context"
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
func (pm *PluginManager) Shutdown(ctx context.Context) error {
	pm.Lock()
	defer pm.Unlock()

	var wg sync.WaitGroup
	errChan := make(chan error, len(pm.processes))

	for id, proc := range pm.processes {
		if proc == nil || proc.Process == nil {
			continue
		}

		wg.Add(1)
		go func(id string, proc *exec.Cmd) {
			defer wg.Done()

			// First try SIGTERM
			if err := proc.Process.Signal(syscall.SIGTERM); err != nil {
				log.Printf("Failed to send SIGTERM to plugin %s: %v", id, err)
				if err := proc.Process.Kill(); err != nil {
					errChan <- fmt.Errorf("failed to kill plugin %s: %w", id, err)
				}
				return
			}

			// Wait for process to exit with timeout
			done := make(chan error, 1)
			go func() {
				done <- proc.Wait()
			}()

			select {
			case err := <-done:
				if err != nil {
					errChan <- fmt.Errorf("plugin %s exit error: %w", id, err)
				}
			case <-ctx.Done():
				// Context timeout - force kill
				if err := proc.Process.Kill(); err != nil {
					errChan <- fmt.Errorf("failed to kill plugin %s after timeout: %w", id, err)
				}
			}
		}(id, proc)
	}

	// Wait for all processes
	wg.Wait()
	close(errChan)

	// Collect errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("plugin shutdown errors: %v", errs)
	}
	return nil
}
