package index

import (
	"context"

	"go.quinn.io/dataq/schema"
)

func (ci *ClaimsIndexer) ListTasks(ctx context.Context) ([]schema.Task, error) {
	var tasks []schema.Task
	err := ci.db.Find(&tasks).Error
	return tasks, err
}
