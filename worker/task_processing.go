package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/queue"
)

func (w *Worker) runTaskProcessingLoop(ctx context.Context, tasks <-chan *queue.Task, messages chan Message) error {
	sendInfo(messages, "Processing tasks")

	pluginReqs, taskmap, err := w.startPluginHandlers(ctx, messages)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			// Context canceled, time to stop
			return ctx.Err()
		case task, ok := <-tasks:
			if !ok {
				// No more tasks
				w.closeAllPluginRequests(pluginReqs)
				return nil
			}

			if err := w.processTask(ctx, task, pluginReqs, taskmap, messages); err != nil {
				sendError(messages, err)
			}
		}
	}
}

func (w *Worker) processTask(ctx context.Context, task *queue.Task, pluginReqs map[string]chan *pb.PluginRequest, taskmap map[string]map[string]*queue.Task, messages chan Message) error {
	sendInfo(messages, "Processing task: "+task.Key())
	taskmap[task.PluginID][task.ID] = task

	req, err := w.buildPluginRequest(ctx, task, messages)
	if err != nil {
		return err
	}

	// Send the request to the appropriate plugin
	pluginReqs[req.PluginId] <- req

	return nil
}

func (w *Worker) buildPluginRequest(ctx context.Context, task *queue.Task, messages chan Message) (*pb.PluginRequest, error) {
	request := task.Request()

	// If task has a hash, we are transforming existing data
	if task.Hash != "" {
		if task.Action != nil {
			return nil, errors.New("task has both action and data")
		}

		data, err := w.cas.Retrieve(ctx, task.Hash)
		if err != nil {
			return nil, fmt.Errorf("failed to load data: %w", err)
		}
		defer data.Close()

		bytes, err := io.ReadAll(data)
		if err != nil {
			return nil, err
		}

		var item pb.DataItem
		if err := json.Unmarshal(bytes, &item); err != nil {
			return nil, err
		}

		request.Item = &item
		request.Operation = "transform"

		sendInfo(messages, fmt.Sprintf("[input: %s] [hash: %s]", item.Meta.Id, item.Meta.Hash))
		sendInfo(messages, item.Meta.String())

	} else {
		// If task has no hash, it should have an action
		if task.Action == nil {
			return nil, errors.New("task has no action or data")
		}
		request.Operation = "extract"
		sendInfo(messages, "[input: nil]")
	}

	return request, nil
}

func (w *Worker) closeAllPluginRequests(pluginReqs map[string]chan *pb.PluginRequest) {
	for _, reqs := range pluginReqs {
		close(reqs)
	}
}
