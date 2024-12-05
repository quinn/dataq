package stream

import (
	"io"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/proto"
	"fmt"
)

// ReadRequest reads a PluginRequest from the reader using length-prefixed framing
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func ReadRequest(r io.Reader) (*pb.PluginRequest, error) {
	data, err := Read(r)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var req pb.PluginRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

// StreamRequests reads PluginRequests from the reader until EOF
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func StreamRequests(r io.Reader) (<-chan *pb.PluginRequest, <-chan error) {
	reqs := make(chan *pb.PluginRequest)
	errc := make(chan error, 1)

	go func() {
		defer close(reqs)
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

			var req pb.PluginRequest
			if err := proto.Unmarshal(data, &req); err != nil {
				errc <- fmt.Errorf("error unmarshaling message: %w", err)
				return
			}

			if req.GetClosed() {
				return
			}

			reqs <- &req
		}
	}()

	return reqs, errc
}

// WriteResponse writes a PluginResponse to the writer using length-prefixed framing
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func WriteResponse(w io.Writer, resp *pb.PluginResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	return Write(w, data)
}
