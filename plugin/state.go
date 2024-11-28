package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// State represents the execution state of a plugin
type State struct {
	PluginID string            `json:"plugin_id"`
	LastRun  int64             `json:"last_run"`
	Metadata map[string]string `json:"metadata"`
}

// StateManager manages plugin execution states
type StateManager struct {
	statePath string
	mu        sync.RWMutex
	states    map[string]*State
}

// NewStateManager creates a new state manager
func NewStateManager(statePath string) (*StateManager, error) {
	sm := &StateManager{
		statePath: statePath,
		states:    make(map[string]*State),
	}

	// Create state directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %v", err)
	}

	// Load existing state if it exists
	if _, err := os.Stat(statePath); err == nil {
		data, err := os.ReadFile(statePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read state file: %v", err)
		}

		if err := json.Unmarshal(data, &sm.states); err != nil {
			return nil, fmt.Errorf("failed to parse state file: %v", err)
		}
	}

	return sm, nil
}

// GetState retrieves the state for a plugin
func (sm *StateManager) GetState(pluginID string) *State {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	state, exists := sm.states[pluginID]
	if !exists {
		state = &State{
			PluginID: pluginID,
			Metadata: make(map[string]string),
		}
		sm.states[pluginID] = state
	}
	return state
}

// UpdateState updates the state for a plugin
func (sm *StateManager) UpdateState(state *State) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.states[state.PluginID] = state

	// Save state to disk
	data, err := json.MarshalIndent(sm.states, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if err := os.WriteFile(sm.statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}
