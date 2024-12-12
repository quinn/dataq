package tree

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	pb "go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/proto"
)

type SQLite struct {
	db  *sql.DB
	cas cas.Storage
}

func New(db *sql.DB, c cas.Storage) (*SQLite, error) {
	// Create table with flattened metadata fields
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS data_items (
			hash TEXT PRIMARY KEY,
			plugin_id TEXT NOT NULL,
			id TEXT NOT NULL,
			kind TEXT NOT NULL,
			timestamp INTEGER,
			content_type TEXT,
			parent_hash TEXT NOT NULL
		)
	`); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLite{db: db, cas: c}, nil
}

// Index scans the root directory and indexes all DataItems
func (s *SQLite) Index() error {
	// Walk the directory and store items
	hashes, err := s.cas.Iterate()
	if err != nil {
		return err
	}

	for hash := range hashes {
		item, err := s.cas.RetrieveItem(hash)
		if err != nil {
			return err
		}
		meta := item.GetMeta()

		if _, err := s.db.Exec(`
			INSERT OR REPLACE INTO data_items (
				hash, plugin_id, id, kind, timestamp,
				content_type, parent_hash
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
			meta.GetHash(),
			meta.GetPluginId(),
			meta.GetId(),
			meta.GetKind(),
			meta.GetTimestamp(),
			meta.GetContentType(),
			meta.GetParentHash(),
		); err != nil {
			return err
		}
	}

	return nil
}

// Children returns all DataItems that have the given hash as their parent
func (s *SQLite) Children(hash string) ([]*pb.DataItemMetadata, error) {
	rows, err := s.db.Query(`
		SELECT hash, plugin_id, id, kind, timestamp, content_type, parent_hash
		FROM data_items
		WHERE parent_hash = ?
	`, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*pb.DataItemMetadata
	for rows.Next() {
		meta := &pb.DataItemMetadata{}
		err := rows.Scan(
			&meta.Hash,
			&meta.PluginId,
			&meta.Id,
			&meta.Kind,
			&meta.Timestamp,
			&meta.ContentType,
			&meta.ParentHash,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, meta)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// Node returns the DataItem with the given hash
func (s *SQLite) Node(hash string) (*pb.DataItemMetadata, error) {
	meta := &pb.DataItemMetadata{}
	err := s.db.QueryRow(`
		SELECT hash, plugin_id, id, kind, timestamp, content_type, parent_hash
		FROM data_items
		WHERE hash = ?
	`, hash).Scan(
		&meta.Hash,
		&meta.PluginId,
		&meta.Id,
		&meta.Kind,
		&meta.Timestamp,
		&meta.ContentType,
		&meta.ParentHash,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// Close closes the database connection
func (s *SQLite) Close() error {
	return s.db.Close()
}
