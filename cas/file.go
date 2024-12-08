package cas

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
)

// FileStorage implements Storage interface using the filesystem
type FileStorage struct {
	root string
}

// NewFileStorage creates a new FileStorage with the given root directory
func NewFileStorage(root string) (*FileStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &FileStorage{root: root}, nil
}

// Store saves the content from the reader and returns its hash
func (f *FileStorage) Store(r io.Reader) (string, error) {
	// Create a temporary file to store the content
	tmpFile, err := os.CreateTemp(f.root, "tmp-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Calculate hash while copying to temp file
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(tmpFile, h), r); err != nil {
		return "", err
	}

	// Get the hash and create the final path
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	finalPath := filepath.Join(f.root, hash)

	// Rewind temp file
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return "", err
	}

	// Create the final file
	finalFile, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer finalFile.Close()

	// Copy content to final location
	if _, err := io.Copy(finalFile, tmpFile); err != nil {
		os.Remove(finalPath)
		return "", err
	}

	return hash, nil
}

// Retrieve returns a reader for the content with the given hash
func (f *FileStorage) Retrieve(hash string) (io.Reader, error) {
	path := filepath.Join(f.root, hash)
	return os.Open(path)
}

// Iterate returns a channel that will receive all stored hashes
func (f *FileStorage) Iterate() (<-chan string, error) {
	hashes := make(chan string)

	go func() {
		defer close(hashes)

		entries, err := os.ReadDir(f.root)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				hashes <- entry.Name()
			}
		}
	}()

	return hashes, nil
}
