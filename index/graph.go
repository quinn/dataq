package index

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

type Rel struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
}

func (i *Index) GetRels(ctx context.Context, hash string) ([]Rel, error) {
	var or sq.Or
	var fields []string
	for name, err := range i.IterateFields(ctx) {
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(name, "_hash") {
			fields = append(fields, name)

			if name == "content_hash" {
				continue
			}

			or = append(or, sq.Eq{name: hash})
		}
	}

	rows, err := sq.
		Select(fields...).
		From("index_data").
		Where(or).
		OrderBy("timestamp DESC").
		RunWith(i.db).
		Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var results []Rel
	for rows.Next() {
		// Prepare a slice for the column values and a slice of pointers to each item.
		// The pointers slice (`valuePtrs`) is passed to rows.Scan.
		values := make([]interface{}, len(fields))
		valuePtrs := make([]interface{}, len(fields))
		for i := range fields {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the pointers slice.
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var entry Rel
		for i, field := range fields {
			if field == "content_hash" {
				sval, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("unexpected value for content_hash %q", values[i])
				}
				entry.Hash = sval
				continue
			}
			ival := values[i]
			if ival == nil {
				continue
			}

			sval, ok := ival.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected type for field %q", field)
			}

			if sval == "" {
				continue
			}

			if sval != hash {
				return nil, fmt.Errorf("unexpected value for field %q: %q", field, sval)
			}

			typeName := strings.TrimSuffix(field, "_hash")
			entry.Type = typeName
		}

		results = append(results, entry)
	}

	return results, nil
}
