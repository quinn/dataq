package index

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"go.quinn.io/dataq/cas"
)

type Index struct {
	cas cas.Storage
	db  *sql.DB
}

type Claim struct {
	SchemaKind string                 `json:"dataq_schema_kind"`
	Hash       string                 `json:"hash"`
	Metadata   map[string]interface{} `json:"-"` // Not stored in CAS
}

func NewIndex(cas cas.Storage, db *sql.DB) *Index {
	return &Index{
		cas: cas,
		db:  db,
	}
}

type Indexable interface {
	SchemaMetadata() map[string]interface{}
	SchemaKind() string
}

func (i *Index) Store(ctx context.Context, data Indexable) (string, error) {
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
		Hash:       hash,
	}

	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	if _, err := i.cas.Store(ctx, bytes.NewReader(claimBytes)); err != nil {
		return "", err
	}

	err = i.Index(hash, data)
	return hash, err
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

// Get retrieves a single object from the index and CAS store.
// The caller must provide a concrete type T that implements Indexable.
func (i *Index) Get(ctx context.Context, result Indexable, whereClause string, args ...interface{}) error {
	// Query the index to get the hash
	query := "SELECT schema_kind, hash FROM index_data"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += " LIMIT 2" // Get 2 to check for multiple matches

	rows, err := i.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var hash, schemaKind string
	var found bool

	for rows.Next() {
		if found {
			return fmt.Errorf("multiple records found for query")
		}
		if err := rows.Scan(&schemaKind, &hash); err != nil {
			return err
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no record found for query")
	}

	// Verify schema kind matches
	if schemaKind != result.SchemaKind() {
		return fmt.Errorf("schema kind mismatch: stored %s, requested %s", schemaKind, result.SchemaKind())
	}

	// Retrieve from CAS
	r, err := i.cas.Retrieve(ctx, hash)
	if err != nil {
		return err
	}
	defer r.Close()

	return json.NewDecoder(r).Decode(result)
}

// Query executes a SQL query against the index and returns matching rows
// The query should be a valid SQL WHERE clause
func (i *Index) Query(ctx context.Context, whereClause string, args ...interface{}) ([]Claim, error) {
	// Get column names first
	rows, err := i.db.Query("PRAGMA table_info(index_data)")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, type_ string
		var notnull, pk int
		var dflt_value interface{}
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}

	// Build and execute the query
	query := "SELECT " + strings.Join(columns, ", ") + " FROM index_data"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	rows, err = i.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Claim
	for rows.Next() {
		// Create a slice of interface{} to scan into
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create result object
		result := Claim{
			Metadata: make(map[string]interface{}),
		}

		// Map values to appropriate fields
		for i, col := range columns {
			val := values[i]
			switch col {
			case "schema_kind":
				if v, ok := val.(string); ok {
					result.SchemaKind = v
				}
			case "hash":
				if v, ok := val.(string); ok {
					result.Hash = v
				}
			default:
				result.Metadata[col] = val
			}
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
