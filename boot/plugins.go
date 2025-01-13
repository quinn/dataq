package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PluginManager manages plugin processes and their gRPC clients
type PluginManager struct {
	sync.RWMutex
	Clients   map[string]*DataQClient
	processes map[string]*exec.Cmd
	index     *index.Index
	cas       cas.Storage
	basePort  int
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(idx *index.Index, cas cas.Storage) *PluginManager {
	return &PluginManager{
		Clients:   make(map[string]*DataQClient),
		processes: make(map[string]*exec.Cmd),
		index:     idx,
		cas:       cas,
		basePort:  50051,
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
func (pm *PluginManager) startPlugin(cfg *config.Plugin, plugin schema.PluginInstance) (*exec.Cmd, *grpc.ClientConn, error) {
	port := fmt.Sprintf("%d", pm.basePort)

	// Create plugin state directory if it doesn't exist
	pluginStateDir := filepath.Join(config.StateDir(), cfg.ID)
	if err := os.MkdirAll(pluginStateDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("failed to create plugin state directory: %w", err)
	}

	// Start the plugin process
	cmd := exec.Command(filepath.Join(config.StateDir(), "bin", cfg.BinaryPath))
	cmd.Dir = pluginStateDir

	configJSON, err := json.Marshal(map[string]interface{}{
		"oauth_config": plugin.OauthConfig,
		"oauth_token":  plugin.OauthToken,
		"config":       plugin.Config,
	})

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%s", port),
		fmt.Sprintf("CONFIG=%s", string(configJSON)),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start plugin process: %w", err)
	}

	// Connect to the plugin
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		cmd.Process.Kill()
		return nil, nil, fmt.Errorf("failed to connect to plugin: %w", err)
	}

	pm.basePort++
	return cmd, conn, nil
}

// AddPlugin adds a new plugin process and client
func (pm *PluginManager) AddPlugin(pluginID string, cfg *config.Plugin, plugin schema.PluginInstance) error {
	pm.Lock()
	defer pm.Unlock()

	process, conn, err := pm.startPlugin(cfg, plugin)
	if err != nil {
		return fmt.Errorf("failed to start plugin %s - %s: %w", cfg.ID, pluginID, err)
	}

	pm.Clients[pluginID] = NewDataQClient(conn, pm.index, pm.cas)
	pm.processes[pluginID] = process

	return nil
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
