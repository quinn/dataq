package boot

import (
	"context"
	"fmt"
	"log"

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/queue"
	"go.quinn.io/dataq/tree"
	"go.quinn.io/dataq/worker"
)

type Boot struct {
	Queue  queue.Queue
	Tree   tree.Tree
	Config *config.Config
	Worker *worker.Worker
}

func New() (*Boot, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Initialize queue
	q, err := queue.NewQueue("file", config.StateDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	// Create worker
	wrkr := worker.New(q, cfg.Plugins, config.DataDir())
	// w.Start(context.Background(), make(chan worker.Message, 10))

	// init tree
	t, err := tree.New(config.StateDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create tree: %w", err)
	}

	// build index
	err = t.Index()
	if err != nil {
		return nil, fmt.Errorf("failed to index tree: %w", err)
	}

	queueItems, err := q.List(context.Background(), "")
	if err != nil {
		log.Fatalf("Failed to list queue items: %v", err)
	}

	if len(queueItems) == 0 {
		for _, p := range cfg.Plugins {
			if !p.Enabled {
				continue
			}
			initialTask := queue.InitialTask(*p)
			if err := q.Push(context.Background(), initialTask); err != nil {
				log.Printf("Warning: Failed to add initial task: %v", err)
			}
		}
	}

	return &Boot{
		Config: cfg,
		Queue:  q,
		Tree:   t,
		Worker: wrkr,
	}, nil
}
