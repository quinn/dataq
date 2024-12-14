package worker

import (
	"context"
	"time"

	"go.quinn.io/dataq/queue"
)

func (w *Worker) updateTaskState(ctx context.Context, task *queue.Task, status queue.TaskStatus, errMsg string, messages chan Message) {
	task.Status = status
	task.Error = errMsg
	task.UpdatedAt = time.Now()

	if uErr := w.queue.Update(ctx, task); uErr != nil {
		sendErrorf(messages, "Failed to update task: %v", uErr)
	}
}

func (w *Worker) taskError(ctx context.Context, task *queue.Task, messages chan Message, err error) {
	sendError(messages, err)
	if task == nil {
		return
	}
	w.updateTaskState(ctx, task, queue.TaskStatusFailed, err.Error(), messages)
}
