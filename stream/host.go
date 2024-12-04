package stream

import (
	"io"

	pb "go.quinn.io/dataq/proto"
)

// WriteRequest writes a PluginRequest to the writer using length-prefixed framing
func WriteRequest(w io.Writer, req *pb.PluginRequest) error {
	return Write(w, req)
}

// ReadResponse reads a PluginResponse from the reader using length-prefixed framing
func ReadResponse(r io.Reader) (*pb.PluginResponse, error) {
	var resp pb.PluginResponse

	if err := Read(r, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// StreamResponses reads PluginResponses from the reader until EOF
func StreamResponses(r io.Reader) (<-chan *pb.PluginResponse, <-chan error) {
	return Stream[*pb.PluginResponse](r)
}
