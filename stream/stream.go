package stream

import (
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

const LengthSize = 8 // 8-byte length header

type StreamMessage interface {
	proto.Message
	GetClosed() bool
}

// Read reads a message from the reader using length-prefixed framing
func Read[T StreamMessage](r io.Reader, msg *T) error {
	if msg == nil {
		panic("msg must not be nil")
	}

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

	// Unmarshal the message
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return nil
}

// Write writes a protobuf message to the writer using length-prefixed framing
func Write[T StreamMessage](w io.Writer, msg T) error {
	// Marshal the message
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Write length header (8 bytes, big endian)
	err = binary.Write(w, binary.BigEndian, uint64(len(data)))
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
	msgs := make(chan T)
	errc := make(chan error, 1)

	go func() {
		defer close(msgs)
		defer close(errc)

		for {
			var msg T
			err := Read(r, &msg)
			if err != nil {
				if err == io.EOF {
					return
				}
				errc <- fmt.Errorf("error reading message: %w", err)
				return
			}

			// closing, without sending final msg
			if msg.GetClosed() {
				return
			}

			msgs <- msg
		}
	}()

	return msgs, errc
}
