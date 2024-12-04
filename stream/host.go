package stream

import (
	"io"

	pb "go.quinn.io/dataq/proto"
	"google.golang.org/protobuf/proto"
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
	msgs := make(chan StreamMessage)
	resps := make(chan *pb.PluginResponse)
	
	unmarshal := func(data []byte) (StreamMessage, error) {
		var resp pb.PluginResponse
		if err := proto.Unmarshal(data, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}
	
	errc := Stream(r, msgs, unmarshal)
	
	go func() {
		defer close(resps)
		for msg := range msgs {
			resps <- msg.(*pb.PluginResponse)
		}
	}()
	
	return resps, errc
}
