package cas

import "io"

type Storage interface {
	Store(io.Reader) (hash string, err error)
	Retrieve(hash string) (data io.Reader, err error)
	Iterate() (hashes <-chan string, err error)
}
