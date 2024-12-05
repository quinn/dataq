package tree

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.quinn.io/dataq/proto"
)

// OpenDB opens or creates a SQLite database at the given path
func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table with flattened metadata fields
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS data_items (
			hash TEXT PRIMARY KEY,
			plugin_id TEXT NOT NULL,
			id TEXT NOT NULL,
			kind TEXT NOT NULL,
			timestamp INTEGER,
			content_type TEXT,
			parent_hash TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Store scans a directory and stores all DataItems in a SQLite database
func Store(root string, db *sql.DB) error {
	// Walk the directory and store items
	return Walk(root, func(item *proto.DataItem, path string) error {
		meta := item.GetMeta()
		_, err := db.Exec(`
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
		)
		return err
	})
}

// Query retrieves metadata from the database that matches the given SQL query
func Query(db *sql.DB, sqlQuery string) ([]*proto.DataItemMetadata, error) {
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*proto.DataItemMetadata
	for rows.Next() {
		meta := &proto.DataItemMetadata{}
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

// Close closes the database connection
func Close(db *sql.DB) error {
	return db.Close()
}
