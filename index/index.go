package index

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"strings"

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
	metadata := data.Metadata()
	schemaKind := data.SchemaKind()

	// Build create table statement with dynamic columns
	var columns []string
	columns = append(columns, "schema_kind TEXT", "hash TEXT")
	
	// Build insert statement parts
	insertColumns := []string{"schema_kind", "hash"}
	placeholders := []string{"?", "?"}
	values := []interface{}{schemaKind, hash}

	for key, value := range metadata {
		var colType string
		switch v := value.(type) {
		case int, int32, int64:
			colType = "INTEGER"
		case float32, float64:
			colType = "REAL"
		case bool:
			colType = "BOOLEAN"
		case string:
			colType = "TEXT"
		default:
			// For complex types, store as JSON
			colType = "TEXT"
			b, err := json.Marshal(v)
			if err != nil {
				return err
			}
			value = string(b)
		}
		columns = append(columns, key+" "+colType)
		insertColumns = append(insertColumns, key)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	// Create table if not exists
	createTableSQL := "CREATE TABLE IF NOT EXISTS index_data (" + 
		strings.Join(columns, ", ") + ")"
	
	_, err := i.db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Insert the data
	insertSQL := "INSERT INTO index_data (" + 
		strings.Join(insertColumns, ", ") + 
		") VALUES (" + strings.Join(placeholders, ", ") + ")"
	
	_, err = i.db.Exec(insertSQL, values...)
	return err
}
