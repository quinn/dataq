package index

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/rpc"
	"go.quinn.io/dataq/schema"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Index struct {
	cas cas.Storage
	db  *sql.DB
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

type IndexableProto interface {
	protoreflect.ProtoMessage
	Indexable
}

func (i *Index) Store(ctx context.Context, data Indexable) (string, error) {
	contentHash, err := i.marshalToCAS(ctx, data)
	if err != nil {
		return "", err
	}

	claim := schema.Claim{
		Type:        "content",
		SchemaKind:  data.SchemaKind(),
		ContentHash: contentHash,
	}

	claimBytes, err := json.Marshal(claim)
	if err != nil {
		return "", err
	}

	if _, err := i.cas.Store(ctx, bytes.NewReader(claimBytes)); err != nil {
		return "", err
	}

	err = i.index(claim, data)
	return contentHash, err
}

func (i *Index) CreatePermanode(ctx context.Context, content Indexable) (string, error) {
	permanode := schema.NewPermanode(content.SchemaKind())
	permanodeHash, err := i.marshalToCAS(ctx, permanode)
	if err != nil {
		return "", err
	}

	if _, err := i.UpdatePermanode(ctx, permanodeHash, content); err != nil {
		return "", err
	}

	return permanodeHash, nil
}

func (i *Index) UpdatePermanode(ctx context.Context, permanodeHash string, content Indexable) (string, error) {
	contentHash, err := i.marshalToCAS(ctx, content)
	if err != nil {
		return "", err
	}

	permanodeVersion := schema.NewPermanodeVersion(permanodeHash, contentHash)
	permanodeVersionHash, err := i.marshalToCAS(ctx, permanodeVersion)
	if err != nil {
		return "", err
	}

	return permanodeVersionHash, nil
}

func (i *Index) GetPermanode(ctx context.Context, permanodeHash string, result Indexable) error {
	return i.Get(ctx, result, "permanode_hash = ?", permanodeHash)
}

func (i *Index) Rebuild(ctx context.Context) error {
	if _, err := i.db.ExecContext(ctx, "DROP TABLE IF EXISTS index_data"); err != nil {
		return fmt.Errorf("failed to drop table: %w", err)
	}

	// Get all hashes from CAS
	hashes, err := i.cas.Iterate(ctx)
	if err != nil {
		return fmt.Errorf("failed to get hashes: %w", err)
	}

	// Index each hash
	for hash := range hashes {
		r, err := i.cas.Retrieve(ctx, hash)
		if err != nil {
			return fmt.Errorf("failed to retrieve CAS object: %w", err)
		}
		defer r.Close()

		b, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read CAS object: %w", err)
		}

		if !strings.HasPrefix(string(b), "{\"dataq_type\":") {
			log.Println("skipping non-claim object: ", hash)
			continue
		} else {
			log.Println("rebuilding claim: ", hash)
		}

		var claim schema.Claim
		if err := json.Unmarshal(b, &claim); err != nil {
			return fmt.Errorf("failed to unmarshal typecheck: %w", err)
		}

		if claim.Type == "permanode_version" {
			if claim.PermanodeHash == "" {
				return fmt.Errorf("permanode version missing permanode hash")
			}

			var permanode schema.Claim
			if err := i.unmarshalFromCAS(ctx, claim.PermanodeHash, &permanode); err != nil {
				return fmt.Errorf("failed to unmarshal permanode: %w", err)
			}

			claim.SchemaKind = permanode.SchemaKind
		}
		if claim.SchemaKind == "" {
			return fmt.Errorf("claim missing schema kind: %w", err)
		}

		slog.Info("rebuilding claim", "hash", claim.ContentHash, "kind", claim.SchemaKind)

		var content Indexable
		switch claim.SchemaKind {
		case "ExtractRequest":
			content = &rpc.ExtractRequest{}
		case "ExtractResponse":
			content = &rpc.ExtractResponse{}
		default:
			return fmt.Errorf("unknown schema kind: %s", claim.SchemaKind)
		}

		// Get the data from CAS
		if err := i.unmarshalFromCAS(ctx, claim.ContentHash, content); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		if err := i.index(claim, content); err != nil {
			return fmt.Errorf("failed to index data: %w", err)
		}
	}

	return nil
}

// Get retrieves a single object from the index and CAS store.
// The caller must provide a concrete type T that implements Indexable.
func (i *Index) Get(ctx context.Context, result Indexable, whereClause string, args ...interface{}) error {
	// Query the index to get the hash
	query := "SELECT schema_kind, content_hash FROM index_data"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}
	query += " LIMIT 2" // Get 2 to check for multiple matches

	rows, err := i.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var contentHash, schemaKind string
	var found bool

	for rows.Next() {
		if found {
			return fmt.Errorf("multiple records found for query")
		}
		if err := rows.Scan(&schemaKind, &contentHash); err != nil {
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
	return i.unmarshalFromCAS(ctx, contentHash, result)
}

// Query executes a SQL query against the index and returns matching rows
// The query should be a valid SQL WHERE clause
func (i *Index) Query(ctx context.Context, whereClause string, args ...interface{}) ([]schema.Claim, error) {
	// Get column names first
	rows, err := i.db.Query("PRAGMA table_info(index_data)")
	if err != nil {
		return nil, fmt.Errorf("failed to get table info: %w", err)
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

	if len(columns) == 0 {
		// If columns is zero, it means that the table does not exist yet.
		return nil, nil
	}

	// Build and execute the query
	query := "SELECT " + strings.Join(columns, ", ") + " FROM index_data"
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	rows, err = i.db.Query(query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "no such column") {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to execute query (%s): %w", query, err)
	}
	defer rows.Close()

	var results []schema.Claim
	for rows.Next() {
		// Create a slice of interface{} to scan into
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create result object
		result := schema.Claim{
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
					result.ContentHash = v
				}
			default:
				result.Metadata[col] = val
			}
		}

		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err() would like to speak with you: %w", err)
	}

	return results, nil
}

// unmarshalFromCAS reads and unmarshals data from CAS storage into the provided object
func (i *Index) unmarshalFromCAS(ctx context.Context, contentHash string, result any) error {
	r, err := i.cas.Retrieve(ctx, contentHash)
	if err != nil {
		return fmt.Errorf("failed to retrieve CAS object: %w", err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read CAS object: %w", err)
	}

	if presult, ok := result.(IndexableProto); ok {
		err = protojson.Unmarshal(b, presult)
	} else {
		err = json.Unmarshal(b, result)
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return nil
}

func (i *Index) index(claim schema.Claim, data Indexable) error {
	metadata := data.SchemaMetadata()
	schemaKind := data.SchemaKind()

	// Create base table if not exists
	// it is possible to index and store content that is not part of a permanode.
	// This content cannot be edited
	createTableSQL := `CREATE TABLE IF NOT EXISTS index_data (
		schema_kind TEXT NOT NULL,
		permanode_hash TEXT,
		timestamp INTEGER,
		content_hash TEXT PRIMARY KEY
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
	insertColumns := []string{"schema_kind", "permanode_hash", "content_hash", "timestamp"}
	placeholders := []string{"?", "?", "?", "?"}
	values := []interface{}{schemaKind, claim.PermanodeHash, claim.ContentHash, claim.Timestamp.UnixMilli()}

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
	insertSQL := "INSERT OR IGNORE INTO index_data (" +
		strings.Join(insertColumns, ", ") +
		") VALUES (" + strings.Join(placeholders, ", ") + ")"

	_, err = i.db.Exec(insertSQL, values...)
	return err
}

// marshalToCAS marshals the provided object and stores it in CAS storage
func (i *Index) marshalToCAS(ctx context.Context, data any) (string, error) {
	var b []byte
	var err error

	if pdata, ok := data.(IndexableProto); ok {
		// Trying to make this as deterministic as possible
		b, err = protojson.MarshalOptions{
			Multiline:     true,
			Indent:        "\t",
			UseProtoNames: true,
		}.Marshal(pdata)
	} else {
		b, err = json.Marshal(data)
	}
	if err != nil {
		return "", err
	}

	hash, err := i.cas.Store(ctx, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return hash, nil
}
