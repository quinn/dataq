package stream

import (
	"encoding/binary"
	"fmt"
	"io"
)

const LengthSize = 8 // 8-byte length header

// Read reads a message from the reader using length-prefixed framing
func Read(r io.Reader) ([]byte, error) {
	// Read length header (8 bytes, big endian)
	var length uint64
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read length header: %w", err)
	}

	// Read data
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return data, nil
}

// Write writes a protobuf message to the writer using length-prefixed framing
func Write(w io.Writer, data []byte) error {
	// Write length header (8 bytes, big endian)
	err := binary.Write(w, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to write length header: %w", err)
	}

	// Write data
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}
