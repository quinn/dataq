package tree

import pb "go.quinn.io/dataq/proto"

type Tree interface {
	Index() error
	Children(hash string) ([]*pb.DataItemMetadata, error)
	Node(hash string) (*pb.DataItemMetadata, error)
	Close() error
}
