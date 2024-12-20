package boot

import (
	"bytes"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/rpc"
)

// DataQClient wraps the gRPC client with index-based request hash handling
type DataQClient struct {
	client rpc.DataQPluginClient
	index  *index.Index
	cas    cas.Storage
}

// NewDataQClient creates a new DataQClient with index-based request hash handling
func NewDataQClient(conn *grpc.ClientConn, idx *index.Index, cas cas.Storage) *DataQClient {
	return &DataQClient{
		client: rpc.NewDataQPluginClient(conn),
		index:  idx,
		cas:    cas,
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

	content := res.GetContent()
	if content == nil {
		return nil, fmt.Errorf("response content is nil")
	}

	r := bytes.NewReader(content)

	dataHash, err := c.cas.Store(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to store response content: %w", err)
	}

	res.Data = &rpc.ExtractResponse_Hash{
		Hash: dataHash,
	}

	// Store the response in the index
	if _, err = c.index.Store(ctx, res); err != nil {
		return nil, err
	}

	return res, nil
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
