package hash

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

// Generate shows an example of how to generate a proper hash.
// base64 url-safe encoded sha256
func Generate(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return Encode(h.Sum(nil))
}

func Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Decode(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

func UID() string {
	id := [16]byte(uuid.New())
	return Encode(id[:])
}
