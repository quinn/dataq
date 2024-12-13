package cas

import (
	"context"
	"fmt"
	"io"
	"log"

	"perkeep.org/pkg/blob"
	"perkeep.org/pkg/cacher"
	"perkeep.org/pkg/client"
)

// implements:
// type Storage interface {

// 	// generic
// 	Store(io.Reader) (hash string, err error)
// 	Retrieve(hash string) (data io.ReadCloser, err error)

// 	// data items
// 	StoreItem(item *pb.DataItem) (hash string, err error)
// 	RetrieveItem(hash string) (item *pb.DataItem, err error)

// 	Iterate() (hashes <-chan string, err error)
// }

type Perkeep struct{}

func NewPerkeep() *Perkeep {
	return &Perkeep{}
}

func (p *Perkeep) Store(r io.Reader) (hash string, err error) {
	return "", nil
}

func (p *Perkeep) Retrieve(hash string) (data io.ReadCloser, err error) {
	var cl *client.Client
	optTransportConfig := client.OptionTransportConfig(&client.TransportConfig{})

	ctx := context.Background()

	cl = client.NewOrFail(client.OptionInsecure(false), optTransportConfig)
	br, ok := blob.Parse(hash)
	if !ok {
		return nil, fmt.Errorf("Failed to parse argument %q as a blobref.", hash)
	}

	src, err := cacher.NewDiskCache(cl)
	if err != nil {
		log.Fatalf("Error setting up local disk cache: %v", err)
	}
	defer src.Clean()

	rc, _, err := src.Fetch(ctx, br)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch %s: %s", br, err)
	}

	return rc, nil
}

func fetch(ctx context.Context, src blob.Fetcher, br blob.Ref) (rc io.ReadCloser, err error) {
	rc, _, err = src.Fetch(ctx, br)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch %s: %s", br, err)
	}
	return rc, err
}

// A little less than the sniffer will take, so we don't truncate.
const sniffSize = 900 * 1024
