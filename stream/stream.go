package stream

import (
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const LengthSize = 8 // 8-byte length header

type StreamMessage interface {
	protoreflect.ProtoMessage
}

// Read reads a message from the reader using length-prefixed framing
func Read(r io.Reader, msg StreamMessage) error {
	// Read length header (8 bytes, big endian)
	var length uint64
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("failed to read length header: %w", err)
	}

	// Read data
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	// Unmarshal the PluginRequest
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PluginRequest: %w", err)
	}

	return nil
}

// Write writes a protobuf message to the writer using length-prefixed framing
func Write(w io.Writer, msg StreamMessage) error {
	// Marshal the message
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write length header (8 bytes, big endian)
	length := uint64(len(data))
	err = binary.Write(w, binary.BigEndian, length)
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

// Stream reads messages from the reader until EOF
func Stream[T StreamMessage](r io.Reader) (<-chan T, <-chan error) {
	responses := make(chan T)
	errc := make(chan error, 1)

	go func() {
		defer close(responses)
		defer close(errc)

		for {
			var msg T
			err := Read(r, msg)
			if err != nil {
				if err == io.EOF {
					return
				}
				errc <- err
				return
			}

			responses <- msg
		}
	}()

	return responses, errc
}
