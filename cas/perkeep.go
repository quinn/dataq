package cas

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"perkeep.org/pkg/blob"
	"perkeep.org/pkg/client"
	"perkeep.org/pkg/constants"
)

// implements:
// type Storage interface {

// 	// generic
// 	Store(io.Reader) (hash string, err error)
// 	Retrieve(hash string) (data io.ReadCloser, err error)
// 	Iterate() (hashes <-chan string, err error)
// }

type Perkeep struct {
	cl *client.Client
}

func NewPerkeep() *Perkeep {
	optTransportConfig := client.OptionTransportConfig(&client.TransportConfig{})
	cl := client.NewOrFail(client.OptionInsecure(false), optTransportConfig)

	return &Perkeep{
		cl: cl,
	}
}

func (p *Perkeep) Store(ctx context.Context, r io.Reader) (hash string, err error) {
	var buf bytes.Buffer
	size, err := io.CopyN(&buf, r, constants.MaxBlobSize+1)
	if size > constants.MaxBlobSize {
		return "", fmt.Errorf("blob size cannot be bigger than %d", constants.MaxBlobSize)
	}
	if err != nil && err != io.EOF {
		return "", err
	}

	h := blob.NewHash()
	if _, err := h.Write(buf.Bytes()); err != nil {
		return "", err
	}

	put, err := p.cl.Upload(ctx, &client.UploadHandle{
		BlobRef:  blob.RefFromHash(h),
		Size:     uint32(buf.Len()),
		Contents: &buf,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload: %s", err)
	}

	return put.BlobRef.String(), nil
}

func (p *Perkeep) Retrieve(ctx context.Context, hash string) (data io.ReadCloser, err error) {
	br, ok := blob.Parse(hash)
	if !ok {
		return nil, fmt.Errorf("failed to parse argument %q as a blobref", hash)
	}

	rc, _, err := p.cl.Fetch(ctx, br)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %s", br, err)
	}

	return rc, nil
}

func (p *Perkeep) Iterate(ctx context.Context) (<-chan string, error) {
	ch := make(chan blob.SizedRef)
	hashes := make(chan string)

	go func() {
		if err := p.cl.SimpleEnumerateBlobs(ctx, ch); err != nil {
			slog.Error("failed to enumerate blobs", "error", err)
			return
		}
	}()

	go func() {
		defer close(hashes)

		for sref := range ch {
			hashes <- sref.Ref.String()
		}
	}()

	return hashes, nil
}

func (p *Perkeep) Delete(ctx context.Context, hash string) error {
	br, ok := blob.Parse(hash)
	if !ok {
		return fmt.Errorf("failed to parse argument %q as a blobref", hash)
	}

	return p.cl.RemoveBlob(ctx, br)
}
