package boot

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/index"
)

type Boot struct {
	// Queue  queue.Queue
	// Tree   tree.Tree
	Config *config.Config
	// Worker *worker.Worker
	CAS cas.Storage
	// Claim   *claims.ClaimsService
	Index   *index.Index
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

	return &Boot{
		Config: cfg,
		Index:  idx,
		// Queue:  q,
		// Tree:   t,
		// Worker:  wrkr,
		CAS:     pk,
		Plugins: NewPluginManager(idx),
	}, nil
}

// StartPlugins initializes and starts all enabled plugins
func (b *Boot) StartPlugins() error {
	basePort := 50051
	for _, plugin := range b.Config.Plugins {
		if !plugin.Enabled {
			continue
		}

		port := fmt.Sprintf("%d", basePort)
		if err := b.startPlugin(plugin, port); err != nil {
			return fmt.Errorf("failed to start plugin %s: %w", plugin.ID, err)
		}
		basePort++
	}
	return nil
}

func (b *Boot) startPlugin(plugin *config.Plugin, port string) error {
	// Start the plugin process
	cmd := exec.Command(filepath.Join(config.StateDir(), "bin", plugin.BinaryPath))
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%s", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin process: %w", err)
	}

	// Connect to the plugin
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to connect to plugin: %w", err)
	}

	b.Plugins.AddPlugin(plugin.ID, conn, cmd)
	return nil
}

// Shutdown gracefully stops all components
func (b *Boot) Shutdown() {
	if b.Plugins != nil {
		b.Plugins.Shutdown()
	}
}
