package dq

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"go.quinn.io/dataq/proto"
)

func TestWriteReadWithDelimiterInMetadata(t *testing.T) {
	// Create a DataItem with the delimiter sequence in its metadata
	item := &proto.DataItem{
		PluginId:    "test-plugin",
		Id:          "test-source",
		Kind:        "test",
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
	err := WriteDataItem(&buf, item)
	if err != nil {
		t.Fatalf("WriteDataItem failed: %v", err)
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
	if item.Id != readItem.Id {
		t.Errorf("Id mismatch: got %v, want %v", readItem.Id, item.Id)
	}
	if item.Kind != readItem.Kind {
		t.Errorf("Kind mismatch: got %v, want %v", readItem.Kind, item.Kind)
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
		Id:          "test-source",
		Kind:        "test",
		Timestamp:   123456789,
		ContentType: "application/octet-stream",
		Metadata: map[string]string{
			"test": "value",
		},
		RawData: rawData,
	}

	// Write the item to a buffer
	var buf bytes.Buffer
	err := WriteDataItem(&buf, item)
	if err != nil {
		t.Fatalf("WriteDataItem failed: %v", err)
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
		Id:          "test-source",
		Kind:        "test",
		Timestamp:   123456789,
		ContentType: "text/plain",
		Metadata:    map[string]string{},
		RawData:     []byte{},
	}

	var buf bytes.Buffer
	err := WriteDataItem(&buf, item)
	if err != nil {
		t.Fatalf("WriteDataItem failed: %v", err)
	}

	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if item.PluginId != readItem.PluginId {
		t.Errorf("PluginId mismatch: got %v, want %v", readItem.PluginId, item.PluginId)
	}
	if item.Id != readItem.Id {
		t.Errorf("Id mismatch: got %v, want %v", readItem.Id, item.Id)
	}
	if item.Kind != readItem.Kind {
		t.Errorf("Kind mismatch: got %v, want %v", readItem.Kind, item.Kind)
	}
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
	}
	if len(readItem.Metadata) != 0 {
		t.Errorf("Expected empty metadata, got %v", readItem.Metadata)
	}
}

func TestGenericWrite(t *testing.T) {
	type CustomMetadata struct {
		Name        string            `json:"name"`
		Tags        []string          `json:"tags"`
		Properties  map[string]string `json:"properties"`
	}

	metadata := CustomMetadata{
		Name: "test-item",
		Tags: []string{"tag1", "tag2"},
		Properties: map[string]string{
			"prop1": "value1",
			"prop2": "value2",
		},
	}
	rawData := []byte("custom raw data")

	var buf bytes.Buffer
	err := Write(&buf, metadata, rawData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read the buffer and split at delimiter
	data := buf.Bytes()
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
		t.Fatal("Delimiter not found in output")
	}

	// Verify metadata
	var readMetadata CustomMetadata
	if err := json.Unmarshal(data[:delimiterIndex], &readMetadata); err != nil {
		t.Fatalf("Failed to unmarshal metadata: %v", err)
	}

	if !reflect.DeepEqual(metadata, readMetadata) {
		t.Errorf("Metadata mismatch: got %v, want %v", readMetadata, metadata)
	}

	// Verify raw data
	readRawData := data[delimiterIndex+len(Delimiter):]
	if !bytes.Equal(rawData, readRawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readRawData, rawData)
	}
}
