package tree

import (
	"io"
	"os"
	"path/filepath"

	"go.quinn.io/dataq/dq"
	"go.quinn.io/dataq/proto"
)

// Walk walks through the directory tree and processes each file
func Walk(root string, callback func(*proto.DataItem, string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Open and process the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Try to parse as DQ file
		item, err := dq.Read(file)
		if err != nil && err != io.ErrUnexpectedEOF {
			return err
		}

		// If successfully parsed, call the callback
		if err == nil && item != nil {
			if err := callback(item, path); err != nil {
				return err
			}
		}

		return nil
	})
}
