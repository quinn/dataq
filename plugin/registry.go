package plugin

import (
    "fmt"
    "sync"
)

// Registry manages the collection of available plugins
type Registry struct {
    plugins map[string]Plugin
    mu      sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
    return &Registry{
        plugins: make(map[string]Plugin),
    }
}

// Register adds a plugin to the registry
func (r *Registry) Register(p Plugin) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.plugins[p.ID()]; exists {
        return fmt.Errorf("plugin with ID %s already registered", p.ID())
    }

    r.plugins[p.ID()] = p
    return nil
}

// Get retrieves a plugin by its ID
func (r *Registry) Get(id string) (Plugin, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    p, exists := r.plugins[id]
    if !exists {
        return nil, fmt.Errorf("plugin with ID %s not found", id)
    }

    return p, nil
}

// List returns all registered plugins
func (r *Registry) List() []Plugin {
    r.mu.RLock()
    defer r.mu.RUnlock()

    plugins := make([]Plugin, 0, len(r.plugins))
    for _, p := range r.plugins {
        plugins = append(plugins, p)
    }

    return plugins
}
