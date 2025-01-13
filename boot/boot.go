package boot

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/internal/repo"
	"go.quinn.io/dataq/schema"
)

type Boot struct {
	// Queue  queue.Queue
	// Tree   tree.Tree
	Config *config.Config
	// Worker *worker.Worker
	CAS cas.Storage
	// Claim   *claims.ClaimsService
	Index   *index.Index
	Repo    *repo.Repo
	Plugins *PluginManager
}

func New() (*Boot, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	config.StateDir()
	db, err := sql.Open("sqlite3", config.StateDir()+"/state.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// // Initialize queue
	// q, err := queue.NewSQLiteQueue(db)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create queue: %w", err)
	// }

	// casDQ := &cas.DQ{}

	pk := cas.NewPerkeep()

	// // Create worker
	// wrkr := worker.New(q, cfg.Plugins, config.DataDir(), pk)

	// // init tree
	// t, err := tree.New(db, pk)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create tree: %w", err)
	// }

	// Initialize index
	idx := index.NewIndex(pk, db)

	// queueItems, err := q.List(context.Background(), "")
	// if err != nil {
	// 	log.Fatalf("Failed to list queue items: %v", err)
	// }

	// if len(queueItems) == 0 {
	// 	for _, p := range cfg.Plugins {
	// 		if !p.Enabled {
	// 			continue
	// 		}
	// 		initialTask := queue.InitialTask(*p)
	// 		if err := q.Push(context.Background(), initialTask); err != nil {
	// 			log.Printf("Warning: Failed to add initial task: %v", err)
	// 		}
	// 	}
	// }

	repo := repo.NewRepo(idx)

	return &Boot{
		Config: cfg,
		Index:  idx,
		// Queue:  q,
		// Tree:   t,
		// Worker:  wrkr,
		CAS:     pk,
		Plugins: NewPluginManager(idx, pk),
		Repo:    repo,
	}, nil
}

// StartPlugins initializes and starts all enabled plugins
func (b *Boot) StartPlugins(ctx context.Context) error {
	claims, err := b.Repo.PluginClaims(ctx)
	if err != nil {
		return fmt.Errorf("failed to get plugins: %w", err)
	}

	for _, claim := range claims {
		var cfg *config.Plugin
		pluginID := claim.Metadata["plugin_id"]
		permanodeHash := claim.PermanodeHash

		for _, c := range b.Config.Plugins {
			if c.ID == pluginID {
				cfg = c
				break
			}
		}

		if cfg == nil {
			return fmt.Errorf("plugin %s not found in config", pluginID)
		}

		if !cfg.Enabled {
			continue
		}

		var plugin schema.PluginInstance
		if err := b.Index.GetPermanode(ctx, claim.PermanodeHash, &plugin); err != nil {
			return err
		}

		if err := b.Plugins.AddPlugin(permanodeHash, cfg, plugin); err != nil {
			return fmt.Errorf("failed to start plugin %s: %w", pluginID, err)
		}
	}

	return nil
}

// Shutdown gracefully stops all components
func (b *Boot) Shutdown(ctx context.Context) error {
	var shutdownErr error
	if b.Plugins != nil {
		if err := b.Plugins.Shutdown(ctx); err != nil {
			shutdownErr = fmt.Errorf("plugin shutdown error: %w", err)
		}
	}
	return shutdownErr
}
