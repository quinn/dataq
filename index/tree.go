package index

import (
	"context"

	"go.quinn.io/dataq/schema"
)

type Tree interface {
	Index(context.Context) error
	Children(hash string) ([]*schema.DataItemMetadata, error)
	Node(hash string) (*schema.DataItemMetadata, error)
	Close() error
}
