package schema

import (
	"time"

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

type Permanode struct {
	Nonce      string `json:"nonce"`
	SchemaKind string `json:"kind"`
}

type PermanodeVersion struct {
	PermanodeHash string    `json:"permanode_hash"`
	Timestamp     time.Time `json:"timestamp"`
	ContentHash   string    `json:"content_hash"`

	// This applies to content from a plugin
	// these values will be blank if permanode is managed by dataq
	PluginID  string `json:"plugin_id,omitempty"`
	PluginKey string `json:"plugin_key,omitempty"`
}
