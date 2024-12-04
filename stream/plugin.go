package stream

import (
	"io"

	pb "go.quinn.io/dataq/proto"
)

// ReadRequest reads a PluginRequest from the reader using length-prefixed framing
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func ReadRequest(r io.Reader) (*pb.PluginRequest, error) {
	var req pb.PluginRequest

	if err := Read(r, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

// StreamRequests reads PluginRequests from the reader until EOF
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func StreamRequests(r io.Reader) (<-chan *pb.PluginRequest, <-chan error) {
	return Stream[*pb.PluginRequest](r)
}

// WriteResponse writes a PluginResponse to the writer using length-prefixed framing
// currently, this is ONLY USED by the plugin harness. Which means this code is not
// necessary for the implementation of dataq, if no plugins are written in Go
func WriteResponse(w io.Writer, resp *pb.PluginResponse) error {
	return Write(w, resp)
}
