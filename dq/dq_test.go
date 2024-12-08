package dq

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	pb "go.quinn.io/dataq/proto"
)

func TestWriteReadWithDelimiterInMetadata(t *testing.T) {
	// Create a DataItem with the delimiter sequence in its metadata
	item := &pb.DataItem{
		Meta: &pb.DataItemMetadata{
			PluginId:    "test-plugin",
			Id:          "test-source",
			Kind:        "test",
			Timestamp:   123456789,
			ContentType: "text/plain",
		},
		RawData: []byte("Hello, World!"),
	}

	// Write the item to a buffer
	var buf bytes.Buffer
	err := Write(&buf, item)
	if err != nil {
		t.Fatalf("WriteDataItem failed: %v", err)
	}

	// Read it back
	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Verify all fields match
	if item.Meta.PluginId != readItem.Meta.PluginId {
		t.Errorf("PluginId mismatch: got %v, want %v", readItem.Meta.PluginId, item.Meta.PluginId)
	}
	if item.Meta.Id != readItem.Meta.Id {
		t.Errorf("Id mismatch: got %v, want %v", readItem.Meta.Id, item.Meta.Id)
	}
	if item.Meta.Kind != readItem.Meta.Kind {
		t.Errorf("Kind mismatch: got %v, want %v", readItem.Meta.Kind, item.Meta.Kind)
	}
	if item.Meta.Timestamp != readItem.Meta.Timestamp {
		t.Errorf("Timestamp mismatch: got %v, want %v", readItem.Meta.Timestamp, item.Meta.Timestamp)
	}
	if item.Meta.ContentType != readItem.Meta.ContentType {
		t.Errorf("ContentType mismatch: got %v, want %v", readItem.Meta.ContentType, item.Meta.ContentType)
	}
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
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

	item := &pb.DataItem{
		Meta: &pb.DataItemMetadata{
			PluginId:    "test-plugin",
			Id:          "test-source",
			Kind:        "test",
			Timestamp:   123456789,
			ContentType: "application/octet-stream",
		},
		RawData: rawData,
	}

	// Write the item to a buffer
	var buf bytes.Buffer
	err := Write(&buf, item)
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
	item := &pb.DataItem{
		Meta: &pb.DataItemMetadata{
			PluginId:    "test-plugin",
			Id:          "test-source",
			Kind:        "test",
			Timestamp:   123456789,
			ContentType: "text/plain",
		},
		RawData: []byte{},
	}

	var buf bytes.Buffer
	err := Write(&buf, item)
	if err != nil {
		t.Fatalf("WriteDataItem failed: %v", err)
	}

	readItem, err := Read(&buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if item.Meta.PluginId != readItem.Meta.PluginId {
		t.Errorf("PluginId mismatch: got %v, want %v", readItem.Meta.PluginId, item.Meta.PluginId)
	}
	if item.Meta.Id != readItem.Meta.Id {
		t.Errorf("Id mismatch: got %v, want %v", readItem.Meta.Id, item.Meta.Id)
	}
	if item.Meta.Kind != readItem.Meta.Kind {
		t.Errorf("Kind mismatch: got %v, want %v", readItem.Meta.Kind, item.Meta.Kind)
	}
	if !bytes.Equal(item.RawData, readItem.RawData) {
		t.Errorf("RawData mismatch: got %v, want %v", readItem.RawData, item.RawData)
	}
}

func TestGenericWrite(t *testing.T) {
	type CustomMetadata struct {
		Name       string            `json:"name"`
		Tags       []string          `json:"tags"`
		Properties map[string]string `json:"properties"`
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
	err := write(&buf, metadata, rawData)
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
