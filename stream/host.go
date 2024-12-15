package stream

import (
	"io"

	"fmt"

	"go.quinn.io/dataq/schema"
	"google.golang.org/protobuf/proto"
)

// WriteRequest writes a PluginRequest to the writer using length-prefixed framing
func WriteRequest(w io.Writer, req *schema.PluginRequest) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	return Write(w, data)
}

// ReadResponse reads a PluginResponse from the reader using length-prefixed framing
func ReadResponse(r io.Reader) (*schema.PluginResponse, error) {
	data, err := Read(r)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var resp schema.PluginResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// StreamResponses reads PluginResponses from the reader until EOF
func StreamResponses(r io.Reader) (<-chan *schema.PluginResponse, <-chan error) {
	resps := make(chan *schema.PluginResponse)
	errc := make(chan error, 1)

	go func() {
		defer close(resps)
		defer close(errc)

		for {
			data, err := Read(r)
			if err != nil {
				if err == io.EOF {
					return
				}
				errc <- fmt.Errorf("error reading message: %w", err)
				return
			}
			if data == nil {
				return
			}

			var resp schema.PluginResponse
			if err := proto.Unmarshal(data, &resp); err != nil {
				errc <- fmt.Errorf("error unmarshaling message: %w", err)
				return
			}

			if resp.GetClosed() {
				return
			}

			resps <- &resp
		}
	}()

	return resps, errc
}
