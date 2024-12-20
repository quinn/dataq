package boot

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/rpc"
)

// DataQClient wraps the gRPC client with index-based request hash handling
type DataQClient struct {
	client rpc.DataQPluginClient
	index  *index.Index
}

// NewDataQClient creates a new DataQClient with index-based request hash handling
func NewDataQClient(conn *grpc.ClientConn, idx *index.Index) *DataQClient {
	return &DataQClient{
		client: rpc.NewDataQPluginClient(conn),
		index:  idx,
	}
}

// Extract performs an extraction with index-based request hash
func (c *DataQClient) Extract(ctx context.Context, req *rpc.ExtractRequest, opts ...grpc.CallOption) (*rpc.ExtractResponse, error) {
	// Store the request in the index to get a hash
	hash, err := c.index.Store(ctx, req)
	if err != nil {
		return nil, err
	}

	// Add the hash to the context metadata
	md := metadata.Pairs("request-hash", hash)
	newCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := c.client.Extract(newCtx, req, opts...)
	if err != nil {
		return nil, err
	}

	// Store the response in the index
	hash, err = c.index.Store(ctx, res)

	return res, err
}

// Transform performs a transformation with index-based request hash
func (c *DataQClient) Transform(ctx context.Context, req *rpc.TransformRequest, opts ...grpc.CallOption) (*rpc.TransformResponse, error) {
	// Store the request in the index to get a hash
	hash, err := c.index.Store(ctx, req)
	if err != nil {
		return nil, err
	}

	// Add the hash to the context metadata
	md := metadata.Pairs("request-hash", hash)
	newCtx := metadata.NewOutgoingContext(ctx, md)

	return c.client.Transform(newCtx, req, opts...)
}
