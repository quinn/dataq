package hash

import (
	"crypto/sha256"
	"encoding/base64"
)

// Generate shows an example of how to generate a proper hash.
// base64 url-safe encoded sha256
func Generate(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
