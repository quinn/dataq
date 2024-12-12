package cas

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.quinn.io/dataq/config"
	"go.quinn.io/dataq/dq"
	"go.quinn.io/dataq/proto"
)

type DQ struct{}

// implement the Storage interface

func (d *DQ) StoreItem(item *proto.DataItem) (string, error) {
	hash := item.Meta.Hash
	f, err := os.Create(filepath.Join(config.DataDir(), hash+".dq"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := dq.Write(f, item); err != nil {
		return "", err
	}

	return hash, nil
}

func (d *DQ) Store(r io.Reader) (hash string, err error) {
	panic("implement me")
}

func (d *DQ) RetrieveItem(hash string) (item *proto.DataItem, err error) {
	f, err := d.Retrieve(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %w", err)
	}
	if f == nil {
		return nil, nil
	}
	defer f.Close()

	return dq.Read(f)
}

func (d *DQ) Retrieve(hash string) (data io.ReadCloser, err error) {
	if hash == "" {
		return nil, nil
	}
	filename := filepath.Join(config.DataDir(), hash+".dq")

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %w", err)
	}

	return f, nil
}

func (d *DQ) Iterate() (<-chan string, error) {
	hashes := make(chan string)
	go func() {
		defer close(hashes)

		err := filepath.Walk(config.DataDir(), func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".dq" {
				return nil
			}

			hash := strings.Split(filepath.Base(path), ".")[0]
			hashes <- hash

			// // Open and process the file
			// file, err := os.Open(path)
			// if err != nil {
			// 	return err
			// }
			// defer file.Close()

			// // Try to parse as DQ file
			// item, err := dq.Read(file)
			// if err != nil && err != io.ErrUnexpectedEOF {
			// 	return err
			// }

			// // If successfully parsed, call the callback
			// if err == nil && item != nil {
			// 	if err := callback(item, path); err != nil {
			// 		return err
			// 	}
			// }

			return nil
		})

		if err != nil {
			slog.Error("failed to walk data directory", slog.Any("error", err))
		}
	}()
	return hashes, nil
}
