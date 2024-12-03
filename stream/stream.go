package stream

import (
	"encoding/binary"
	"fmt"
	"io"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/proto"
)

const LengthSize = 8 // 8-byte length header

// WriteResponse writes a PluginResponse to the writer using length-prefixed framing
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func WriteResponse(w io.Writer, resp *pb.PluginResponse) error {
	// Marshal the PluginResponse
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal PluginResponse: %w", err)
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

// ReadResponse reads a PluginResponse from the reader using length-prefixed framing
func ReadResponse(r io.Reader) (*pb.PluginResponse, error) {
	// Read length header (8 bytes, big endian)
	var length uint64
	err := binary.Read(r, binary.BigEndian, &length)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read length header: %w", err)
	}

	// Read data
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// Unmarshal the PluginResponse
	var resp pb.PluginResponse
	err = proto.Unmarshal(data, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PluginResponse: %w", err)
	}

	return &resp, nil
}

// StreamResponses reads PluginResponses from the reader until EOF
func StreamResponses(r io.Reader) (<-chan *pb.PluginResponse, <-chan error) {
	responses := make(chan *pb.PluginResponse)
	errc := make(chan error, 1)

	go func() {
		defer close(responses)
		defer close(errc)

		for {
			resp, err := ReadResponse(r)
			if err != nil {
				if err == io.EOF {
					return
				}
				errc <- err
				return
			}

			// Check for error in response
			if resp.Error != "" {
				errc <- fmt.Errorf("plugin error: %s", resp.Error)
				return
			}

			responses <- resp
		}
	}()

	return responses, errc
}
