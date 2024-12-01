package dq

import (
	"encoding/json"
	"io"

	"go.quinn.io/dataq/proto"
)

// Delimiter used to separate metadata from raw data
var Delimiter = []byte{0x00, 0x1F, 0x00}

// DataItemWrapper is used for JSON serialization of DataItem without the raw_data field
type DataItemWrapper struct {
	PluginID    string `json:"plugin_id"`
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Timestamp   int64  `json:"timestamp"`
	ContentType string `json:"content_type"`
}

// Write writes any JSON-serializable object followed by raw data to the provided writer
func Write(w io.Writer, metadata interface{}, rawData []byte) error {
	// Marshal metadata to JSON
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Write JSON data
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	// Write delimiter
	if _, err := w.Write(Delimiter); err != nil {
		return err
	}

	// Write raw data
	_, err = w.Write(rawData)
	return err
}

// WriteDataItem writes a DataItem to the provided writer in the custom format
func WriteDataItem(w io.Writer, item *proto.DataItem) error {
	// Create wrapper without raw_data
	wrapper := DataItemWrapper{
		PluginID:    item.Meta.PluginId,
		ID:          item.Meta.Id,
		Kind:        item.Meta.Kind,
		Timestamp:   item.Meta.Timestamp,
		ContentType: item.Meta.ContentType,
	}

	return Write(w, wrapper, item.RawData)
}

// Read reads a DataItem from the provided reader in the custom format
func Read(r io.Reader) (*proto.DataItem, error) {
	// Read all data
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Find delimiter
	var delimiterIndex int = -1
	for i := 0; i <= len(data)-len(Delimiter); i++ {
		match := true
		for j := 0; j < len(Delimiter); j++ {
			if data[i+j] != Delimiter[j] {
				match = false
				break
			}
		}
		if match {
			delimiterIndex = i
			break
		}
	}

	if delimiterIndex == -1 {
		return nil, io.ErrUnexpectedEOF
	}

	// Parse JSON metadata
	var wrapper DataItemWrapper
	if err := json.Unmarshal(data[:delimiterIndex], &wrapper); err != nil {
		return nil, err
	}

	// Create DataItem
	item := &proto.DataItem{
		Meta: &proto.DataItemMetadata{
			PluginId:    wrapper.PluginID,
			Id:          wrapper.ID,
			Kind:        wrapper.Kind,
			Timestamp:   wrapper.Timestamp,
			ContentType: wrapper.ContentType,
		},
		RawData: data[delimiterIndex+len(Delimiter):],
	}

	return item, nil
}
