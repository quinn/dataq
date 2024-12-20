package cas

import (
	"context"
	"io"
)

type Storage interface {
	Store(context.Context, io.Reader) (hash string, err error)
	Retrieve(ctx context.Context, hash string) (data io.ReadCloser, err error)
	Iterate(context.Context) (hashes <-chan string, err error)
	Delete(ctx context.Context, hash string) error
}
