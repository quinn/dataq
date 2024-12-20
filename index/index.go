package index

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"

	"go.quinn.io/dataq/cas"
)

type Index struct {
	cas cas.Storage
	db  *sql.DB
}

type Claim struct {
	SchemaKind string `json:"dataq_schema_kind"`
	PluginID   string `json:"plugin_id"`
	Hash       string `json:"hash"`
}

func NewIndex(cas cas.Storage, db *sql.DB) *Index {
	return &Index{
		cas: cas,
		db:  db,
	}
}

type Indexable interface {
	Metadata() map[string]interface{}
	SchemaKind() string
}

func (i *Index) Store(ctx context.Context, pluginID string, data Indexable) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash, err := i.cas.Store(ctx, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	claim := Claim{
		SchemaKind: data.SchemaKind(),
		PluginID:   pluginID,
		Hash:       hash,
	}

	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	return i.cas.Store(ctx, bytes.NewReader(claimBytes))
}

func (i *Index) Index(hash string, data Indexable) error {
	return nil
}
