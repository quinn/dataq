package boot

import (
	"fmt"

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
	q, err := queue.NewQueue("file", config.DataDir())
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}
	defer q.Close()

	// Initialize tree
	// t, err := tree.NewTree(q)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create tree: %w", err)
	// }

	// Create worker
	wrkr := worker.New(q, cfg.Plugins, config.DataDir())
	// w.Start(context.Background(), make(chan worker.Message, 10))

	return &Boot{
		Config: cfg,
		Queue:  q,
		Tree:   nil,
		Worker: wrkr,
	}, nil
}
