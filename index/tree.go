package index

import (
	"context"

	pb "go.quinn.io/dataq/proto"
)

type Tree interface {
	Index(context.Context) error
	Children(hash string) ([]*pb.DataItemMetadata, error)
	Node(hash string) (*pb.DataItemMetadata, error)
	Close() error
}
