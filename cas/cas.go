package cas

import (
	"io"

	"go.quinn.io/dataq/proto"
)

type Storage interface {

	// generic
	Store(io.Reader) (hash string, err error)
	Retrieve(hash string) (data io.ReadCloser, err error)

	// data items
	StoreItem(item *proto.DataItem) (hash string, err error)
	RetrieveItem(hash string) (item *proto.DataItem, err error)

	Iterate() (hashes <-chan string, err error)
}
