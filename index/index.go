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

	// Create base table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS index_data (
		schema_kind TEXT NOT NULL,
		hash TEXT PRIMARY KEY
	)`
	
	if _, err := i.db.Exec(createTableSQL); err != nil {
		return err
	}

	// Get existing columns
	rows, err := i.db.Query("PRAGMA table_info(index_data)")
	if err != nil {
		return err
	}
	defer rows.Close()

	existingColumns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, type_ string
		var notnull, pk int
		var dflt_value interface{}
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return err
		}
		existingColumns[name] = true
	}

	// Add new columns as needed
	for key, value := range metadata {
		if existingColumns[key] {
			continue
		}

		var colType string
		switch value.(type) {
		case int, int32, int64:
			colType = "INTEGER"
		case float32, float64:
			colType = "REAL"
		case bool:
			colType = "BOOLEAN"
		case string:
			colType = "TEXT"
		default:
			colType = "TEXT"
		}

		alterSQL := "ALTER TABLE index_data ADD COLUMN " + key + " " + colType
		if _, err := i.db.Exec(alterSQL); err != nil {
			return err
		}
	}

	// Build insert statement
	insertColumns := []string{"schema_kind", "hash"}
	placeholders := []string{"?", "?"}
	values := []interface{}{schemaKind, hash}

	for key, value := range metadata {
		switch v := value.(type) {
		case int, int32, int64, float32, float64, bool, string:
			// Use value as is
		default:
			// Convert complex types to JSON
			b, err := json.Marshal(v)
			if err != nil {
				return err
			}
			value = string(b)
		}
		insertColumns = append(insertColumns, key)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	// Insert the data
	insertSQL := "INSERT INTO index_data (" + 
		strings.Join(insertColumns, ", ") + 
		") VALUES (" + strings.Join(placeholders, ", ") + ")"
	
	_, err = i.db.Exec(insertSQL, values...)
	return err
}
