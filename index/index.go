package index

import (
	"bytes"
	"context"
	"encoding/json"

	"go.quinn.io/dataq/cas"
)

type Index struct {
	CAS cas.Storage
}

type Claim struct {
	SchemaKind string `json:"dataq_schema_kind"`
	PluginID   string `json:"plugin_id"`
	Hash       string `json:"hash"`
}

func NewIndex(cas cas.Storage) *Index {
	return &Index{
		CAS: cas,
	}
}

type Indexable interface {
	QueryMetadata() map[string]string
}

func (i *Index) Store(ctx context.Context, schemaKind string, pluginID string, data Indexable) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash, err := i.CAS.Store(ctx, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	claim := Claim{
		SchemaKind: schemaKind,
		PluginID:   pluginID,
		Hash:       hash,
	}

	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	return i.CAS.Store(ctx, bytes.NewReader(claimBytes))
}
