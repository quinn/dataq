package cas

import (
	"context"
	"io"
)

type Storage interface {

	// generic
	Store(context.Context, io.Reader) (hash string, err error)
	Retrieve(ctx context.Context, hash string) (data io.ReadCloser, err error)

	// claims
	// StoreItem(item *pb.DataItem) (hash string, err error)
	// RetrieveItem(hash string) (item *pb.DataItem, err error)
	// StoreClaim()

	Iterate(context.Context) (hashes <-chan string, err error)
}
