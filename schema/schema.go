package schema

import (
	"time"

	"go.quinn.io/dataq/hash"
	"golang.org/x/oauth2"
)

// contains schema not in protobuf format

// PluginInstance
type PluginInstance struct {
	PluginID    string         `json:"plugin_id"`
	Label       string         `json:"label"`
	OauthConfig *oauth2.Config `json:"oauth_config,omitempty"`
	OauthToken  *oauth2.Token  `json:"oauth_token,omitempty"`
}

func (p *PluginInstance) SchemaMetadata() map[string]interface{} {
	return map[string]interface{}{
		"plugin_id": p.PluginID,
		"label":     p.Label,
	}
}

func (p *PluginInstance) SchemaKind() string {
	return "PluginInstance"
}

// Combined object for all types of claims
type Claim struct {
	Type          string    `json:"dataq_type"`
	SchemaKind    string    `json:"schema_kind,omitempty"`
	ContentHash   string    `json:"content_hash,omitempty"`
	Nonce         string    `json:"nonce,omitempty"`
	PermanodeHash string    `json:"permanode_hash,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitzero"`

	// Used for deleting
	DeleteHash string `json:"delete_hash,omitempty"`

	// This applies to content from a plugin
	// these values will be blank if permanode is managed by dataq
	PluginID              string `json:"plugin_id,omitempty"`
	PluginKey             string `json:"plugin_key,omitempty"`
	TransformResponseHash string `json:"transform_response_hash,omitempty"`

	// Not stored in CAS
	// Used for returning search results
	Metadata map[string]interface{} `json:"-"`
}

// Below are the types of claims. May not use any of these structs, for now. Instead use Claim struct above.

// Content is immutable content. Exclusive to permanode.
type Content struct {
	DataQType   string `json:"dataq_type"` // "content"
	SchemaKind  string `json:"schema_kind"`
	ContentHash string `json:"content_hash"`
}

func NewContent(schemaKind, contentHash string) *Claim {
	return &Claim{
		Type:        "content",
		SchemaKind:  schemaKind,
		ContentHash: contentHash,
	}
}

// Permanode is mutable content. Exclusive to Content.
type Permanode struct {
	DataQType  string `json:"dataq_type"` // "permanode"
	SchemaKind string `json:"schema_kind"`
	Nonce      string `json:"nonce"`
}

func NewPermanode(schemaKind string) *Claim {
	return &Claim{
		Type:       "permanode",
		SchemaKind: schemaKind,
		Nonce:      hash.UID(),
	}
}

// PermanodeVersion is a version of a permanode.
type PermanodeVersion struct {
	DataQType     string    `json:"dataq_type"` // "permanode_version"
	PermanodeHash string    `json:"permanode_hash"`
	Timestamp     time.Time `json:"timestamp"`
	ContentHash   string    `json:"content_hash"`

	// This applies to content from a plugin
	// these values will be blank if permanode is managed by dataq
	PluginID              string `json:"plugin_id,omitempty"`
	PluginKey             string `json:"plugin_key,omitempty"`
	TransformResponseHash string `json:"transform_response_hash,omitempty"`
}

func NewPermanodeVersion(permanodeHash, contentHash string) *Claim {
	return &Claim{
		Type:          "permanode_version",
		PermanodeHash: permanodeHash,
		ContentHash:   contentHash,
		Timestamp:     time.Now(),
	}
}

func Delete(hash string) *Claim {
	return &Claim{
		Type:       "delete",
		DeleteHash: hash,
	}
}
