package stream

import (
	"io"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/proto"
	"fmt"
)

// WriteRequest writes a PluginRequest to the writer using length-prefixed framing
func WriteRequest(w io.Writer, req *pb.PluginRequest) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	return Write(w, data)
}

// ReadResponse reads a PluginResponse from the reader using length-prefixed framing
func ReadResponse(r io.Reader) (*pb.PluginResponse, error) {
	data, err := Read(r)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var resp pb.PluginResponse
	if err := proto.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// StreamResponses reads PluginResponses from the reader until EOF
func StreamResponses(r io.Reader) (<-chan *pb.PluginResponse, <-chan error) {
	resps := make(chan *pb.PluginResponse)
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

			var resp pb.PluginResponse
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
