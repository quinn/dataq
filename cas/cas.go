package cas

import (
	"io"

	pb "go.quinn.io/dataq/proto"
)

type Storage interface {

	// generic
	Store(io.Reader) (hash string, err error)
	Retrieve(hash string) (data io.ReadCloser, err error)

	// data items
	StoreItem(item *pb.DataItem) (hash string, err error)
	RetrieveItem(hash string) (item *pb.DataItem, err error)

	Iterate() (hashes <-chan string, err error)
}
