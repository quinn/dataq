package tree

import "go.quinn.io/dataq/proto"

type Tree interface {
	Index() error
	Children(hash string) ([]*proto.DataItemMetadata, error)
	Node(hash string) (*proto.DataItemMetadata, error)
	Close() error
}
