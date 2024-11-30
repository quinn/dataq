package dq

import (
	"bytes"
	"reflect"
	"testing"

	"go.quinn.io/dataq/proto"
)

func TestWriteReadWithDelimiterInMetadata(t *testing.T) {
	// Create a DataItem with the delimiter sequence in its metadata
	item := &proto.DataItem{
		PluginId:    "test-plugin",
		SourceId:    "test-source",
		Timestamp:   123456789,
		ContentType: "text/plain",
		Metadata: map[string]string{
			"key_with_delimiter": string([]byte{0x00, 0x1F, 0x00}),
			"normal_key":        "normal_value",
		},
		RawData: []byte("Hello, World!"),
	}

	// Write the item to a buffer
	var buf bytes.Buffer
	err := Write(&buf, item)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read it back
	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Verify all fields match
	if item.PluginId != readItem.PluginId {
		t.Errorf("PluginId mismatch: got %v, want %v", readItem.PluginId, item.PluginId)
	}
	if item.SourceId != readItem.SourceId {
		t.Errorf("SourceId mismatch: got %v, want %v", readItem.SourceId, item.SourceId)
	}
	if item.Timestamp != readItem.Timestamp {
		t.Errorf("Timestamp mismatch: got %v, want %v", readItem.Timestamp, item.Timestamp)
	}
	if item.ContentType != readItem.ContentType {
		t.Errorf("ContentType mismatch: got %v, want %v", readItem.ContentType, item.ContentType)
	}
	if !reflect.DeepEqual(item.Metadata, readItem.Metadata) {
		t.Errorf("Metadata mismatch: got %v, want %v", readItem.Metadata, item.Metadata)
	}
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
	}

	// Specifically verify the metadata with delimiter is preserved
	if item.Metadata["key_with_delimiter"] != readItem.Metadata["key_with_delimiter"] {
		t.Errorf("Delimiter in metadata not preserved: got %v, want %v",
			readItem.Metadata["key_with_delimiter"],
			item.Metadata["key_with_delimiter"])
	}
}

func TestWriteReadWithDelimiterInRawData(t *testing.T) {
	// Create raw data that contains the delimiter sequence multiple times
	rawData := bytes.Join([][]byte{
		[]byte("start"),
		Delimiter,
		[]byte("middle"),
		Delimiter,
		[]byte("end"),
	}, nil)

	item := &proto.DataItem{
		PluginId:    "test-plugin",
		SourceId:    "test-source",
		Timestamp:   123456789,
		ContentType: "application/octet-stream",
		Metadata: map[string]string{
			"test": "value",
		},
		RawData: rawData,
	}

	// Write the item to a buffer
	var buf bytes.Buffer
	err := Write(&buf, item)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read it back
	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Verify the raw data is preserved exactly
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
	}
}

func TestWriteReadEmpty(t *testing.T) {
	// Test with empty raw data and metadata
	item := &proto.DataItem{
		PluginId:    "test-plugin",
		SourceId:    "test-source",
		Timestamp:   123456789,
		ContentType: "text/plain",
		Metadata:    map[string]string{},
		RawData:     []byte{},
	}

	var buf bytes.Buffer
	err := Write(&buf, item)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if item.PluginId != readItem.PluginId {
		t.Errorf("PluginId mismatch: got %v, want %v", readItem.PluginId, item.PluginId)
	}
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
	}
	if len(readItem.Metadata) != 0 {
		t.Errorf("Expected empty metadata, got %v", readItem.Metadata)
	}
}
