package index

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	Q   sq.SelectBuilder
}

func NewIndex(cas cas.Storage, db *sql.DB) *Index {
	return &Index{
		cas: cas,
		db:  db,
		Q:   sq.Select("*").From("index_data"),
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

	claim := schema.NewContent(data.SchemaKind(), contentHash)
	if _, err := i.marshalToCAS(ctx, claim); err != nil {
		return "", err
	}

	err = i.index(ctx, *claim, data)
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

	if err := i.index(ctx, *permanodeVersion, content); err != nil {
		return "", err
	}

	return permanodeVersionHash, nil
}

func (i *Index) Delete(ctx context.Context, hash string) error {
	del := schema.Delete(hash)
	if _, err := i.marshalToCAS(ctx, del); err != nil {
		return err
	}

	if err := i.index(ctx, *del, nil); err != nil {
		return err
	}

	return nil
}

func (i *Index) GetPermanode(ctx context.Context, permanodeHash string, result Indexable) error {
	sel := i.Q.
		Where("permanode_hash = ?", permanodeHash).
		OrderBy("timestamp DESC").
		Limit(1)
	if err := i.Get(ctx, result, sel); err != nil {
		return fmt.Errorf("failed to get permanode: %w", err)
	}

	return nil
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

	log.Println("rebuilding index")
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

		if claim.Type == "permanode" {
			slog.Info("skipping permanode", "hash", hash)
			continue
		}

		if claim.Type == "delete" {
			if claim.DeleteHash == "" {
				return fmt.Errorf("delete claim missing delete hash")
			}
			slog.Info("processing delete claim", "delete_hash", claim.DeleteHash)

			if err := i.index(ctx, claim, nil); err != nil {
				return fmt.Errorf("failed to index data: %w", err)
			}

			continue
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
			return fmt.Errorf("claim missing schema kind: (%v) %w", claim, err)
		}

		slog.Info("rebuilding claim", "hash", claim.ContentHash, "kind", claim.SchemaKind)

		var content Indexable
		switch claim.SchemaKind {
		case "ExtractRequest":
			content = &rpc.ExtractRequest{}
		case "ExtractResponse":
			content = &rpc.ExtractResponse{}
		case "PluginInstance":
			content = &schema.PluginInstance{}
		case "TransformRequest":
			content = &rpc.TransformRequest{}
		case "TransformResponse":
			content = &rpc.TransformResponse{}
		default:
			return fmt.Errorf("unknown schema kind: %s", claim.SchemaKind)
		}

		// Get the content from CAS
		if err := i.unmarshalFromCAS(ctx, claim.ContentHash, content); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		if err := i.index(ctx, claim, content); err != nil {
			return fmt.Errorf("failed to index data: %w", err)
		}
	}

	return nil
}

// Get retrieves a single object from the index and CAS store.
// The caller must provide a concrete type T that implements Indexable.
func (i *Index) Get(ctx context.Context, result Indexable, query sq.SelectBuilder) error {
	query = query.RemoveColumns().Columns("schema_kind", "content_hash")
	rows, err := query.RunWith(i.db).Query()
	if err != nil {
		return fmt.Errorf("failed to query index: %w", err)
	}
	defer rows.Close()

	var contentHash, schemaKind string
	var found bool

	for rows.Next() {
		if found {
			return fmt.Errorf("multiple records found for query")
		}
		if err := rows.Scan(&schemaKind, &contentHash); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		found = true
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate rows: %w", err)
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

// IterateFields returns a sequence of field names for the index_data table.
func (i *Index) IterateFields(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		rows, err := i.db.QueryContext(ctx, "PRAGMA table_info(index_data)")
		if err != nil {
			// For demonstration, we panic on error.
			// In real code, you might want a different mechanism.
			yield("", fmt.Errorf("failed to query table info: %w", err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var cid int
			var name, type_ string
			var notnull, pk int
			var dflt_value interface{}

			if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
				yield("", err)
				return
			}

			if !yield(name, nil) {
				return
			}
		}

		// Check for iteration error
		if rowsErr := rows.Err(); rowsErr != nil {
			yield("", rowsErr)
			return
		}
	}
}

// Query executes a SQL query against the index and returns matching rows
// The query should be a valid SQL WHERE clause
func (i *Index) Query(ctx context.Context, query sq.SelectBuilder) ([]schema.Claim, error) {
	var columns []string
	for name, err := range i.IterateFields(ctx) {
		if err != nil {
			return nil, err
		}

		columns = append(columns, name)
	}

	// If columns is zero, it means that the table does not exist yet.
	if len(columns) == 0 {
		return nil, nil
	}
	rows, err := query.RunWith(i.db).Query()
	if err != nil {
		if strings.Contains(err.Error(), "no such column") {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to execute query (%v): %w", query, err)
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
			case "content_hash":
				if v, ok := val.(string); ok {
					result.ContentHash = v
				}
			case "permanode_hash":
				if v, ok := val.(string); ok {
					result.PermanodeHash = v
				}
			case "timestamp":
				if v, ok := val.(int64); ok {
					result.Timestamp = time.Unix(v, 0)
				}
			case "delete_hash":
				if v, ok := val.(string); ok {
					result.DeleteHash = v
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

func (i *Index) index(ctx context.Context, claim schema.Claim, data Indexable) error {
	var metadata map[string]interface{}
	var schemaKind string

	if data != nil {
		metadata = data.SchemaMetadata()
		schemaKind = data.SchemaKind()
	} else if claim.Type == "delete" && claim.DeleteHash != "" {
		// For delete claims, we don't need the data
		metadata = make(map[string]interface{})
		schemaKind = "delete"
	} else {
		return fmt.Errorf("data cannot be nil for non-delete claims")
	}

	// Create base table if not exists
	// it is possible to index and store content that is not part of a permanode.
	// This content cannot be edited
	createTableSQL := `CREATE TABLE IF NOT EXISTS index_data (
		schema_kind TEXT NOT NULL,
		permanode_hash TEXT,
		timestamp INTEGER,
		content_hash TEXT,
		delete_hash TEXT
	)`

	if _, err := i.db.Exec(createTableSQL); err != nil {
		return err
	}

	// Get existing columns
	existingColumns := make(map[string]bool)
	for name, err := range i.IterateFields(ctx) {
		if err != nil {
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

	var values []interface{}

	// Build insert statement using Squirrel
	insertBuilder := sq.Insert("index_data").
		Columns("schema_kind", "permanode_hash", "content_hash")
	values = append(values, schemaKind, claim.PermanodeHash, claim.ContentHash)

	if !claim.Timestamp.IsZero() {
		insertBuilder = insertBuilder.Columns("timestamp")
		values = append(values, claim.Timestamp.UnixMilli())
	}

	if claim.DeleteHash != "" {
		insertBuilder = insertBuilder.Columns("delete_hash")
		values = append(values, claim.DeleteHash)
	}

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
		insertBuilder = insertBuilder.Columns(key)
		values = append(values, value)
	}

	// Add all values at once
	insertBuilder = insertBuilder.Values(values...)

	if claim.Type == "delete" {
		if claim.DeleteHash == "" {
			return fmt.Errorf("delete claim must have a delete_hash")
		}

		// First delete existing entries if the delete claim is newer
		_, err := sq.Delete("index_data").
			Where(sq.Or{
				sq.Eq{"content_hash": claim.DeleteHash},
				sq.Eq{"permanode_hash": claim.DeleteHash},
			}).
			Where(sq.Or{
				sq.Eq{"timestamp": nil},
				sq.Lt{"timestamp": claim.Timestamp.UnixMilli()},
			}).
			RunWith(i.db).
			Exec()
		if err != nil {
			return fmt.Errorf("failed to delete entries: %w", err)
		}
	} else {
		// Check if content_hash already exists
		var contentHash string
		err := sq.Select("content_hash").
			From("index_data").
			Where(sq.Eq{"content_hash": claim.ContentHash}).
			RunWith(i.db).
			QueryRow().
			Scan(&contentHash)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to check for existing content: %w", err)
		}
		if err != sql.ErrNoRows {
			slog.Warn("content hash already exists in index", "hash", claim.ContentHash)
			return nil
		}

		// Check if there's a newer delete claim for this content or permanode
		var deleteTimestamp int64
		err = sq.Select("timestamp").
			From("index_data").
			Where(sq.And{
				sq.Or{
					sq.Eq{"delete_hash": claim.ContentHash},
					sq.Eq{"delete_hash": claim.PermanodeHash},
				},
			}).
			OrderBy("timestamp DESC").
			Limit(1).
			RunWith(i.db).
			QueryRow().
			Scan(&deleteTimestamp)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to check for delete claims: %w", err)
		}
		if err != sql.ErrNoRows && deleteTimestamp > claim.Timestamp.UnixMilli() {
			// Skip this claim as there's a newer delete
			return nil
		}
	}

	_, err := insertBuilder.RunWith(i.db).Exec()
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
