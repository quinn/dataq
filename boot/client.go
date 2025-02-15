package boot

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"go.quinn.io/dataq/cas"
	"go.quinn.io/dataq/index"
	"go.quinn.io/dataq/internal/repo"
	"go.quinn.io/dataq/rpc"
	"go.quinn.io/dataq/schema"
)

// DataQClient wraps the gRPC client with index-based request hash handling
type DataQClient struct {
	client rpc.DataQPluginClient
	index  *index.Index
	cas    cas.Storage
	repo   *repo.Repo
}

// NewDataQClient creates a new DataQClient with index-based request hash handling
func NewDataQClient(conn *grpc.ClientConn, idx *index.Index, cas cas.Storage, repo *repo.Repo) *DataQClient {
	return &DataQClient{
		client: rpc.NewDataQPluginClient(conn),
		index:  idx,
		cas:    cas,
		repo:   repo,
	}
}

func (c *DataQClient) Install(ctx context.Context, req *rpc.InstallRequest, opts ...grpc.CallOption) (*rpc.InstallResponse, error) {
	// passthru for now
	return c.client.Install(ctx, req, opts...)
}

// Extract performs an extraction with index-based request hash
func (c *DataQClient) Extract(ctx context.Context, plugin *schema.PluginInstance, req *rpc.ExtractRequest, opts ...grpc.CallOption) (*rpc.ExtractResponse, error) {
	// Store the request in the index to get a hash
	// TODO: typically, the extract is called by a request that has already been stored.
	// this may not be necessary

	hash, err := c.repo.StoreExtractRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	// always attaching here is more explicit
	req.Oauth = plugin.Oauth

	res, err := c.client.Extract(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	res.RequestHash = hash

	content := res.GetContent()
	if content == nil {
		return nil, fmt.Errorf("response content is nil")
	}

	r := bytes.NewReader(content)

	// Store the response content in the CAS
	dataHash, err := c.cas.Store(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to store response content: %w", err)
	}

	// replace the contents with a hash address to the content
	res.Data = &rpc.ExtractResponse_Hash{
		Hash: dataHash,
	}

	// For each transform in the response, create a transform request
	for _, transform := range res.GetTransforms() {
		transformReq := &rpc.TransformRequest{
			PluginId: req.PluginId,
			Data: &rpc.TransformRequest_Hash{
				Hash: dataHash,
			},
			Kind:     transform.Kind,
			Metadata: transform.Metadata,
		}

		// Store the transform request
		if _, err := c.index.Store(ctx, transformReq); err != nil {
			return nil, fmt.Errorf("failed to store transform request: %w", err)
		}
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
	// TODO: typically, the transform is called by a request that has already been stored.
	// this may not be necessary
	hash, err := c.index.Store(ctx, req)
	if err != nil {
		return nil, err
	}

	// If request contains a hash, fetch the content from CAS
	if reqHash := req.GetHash(); reqHash != "" {
		r, err := c.cas.Retrieve(ctx, reqHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get content from CAS: %w", err)
		}

		content, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("failed to read content from CAS: %w", err)
		}

		// Update request to use content instead of hash
		req.Data = &rpc.TransformRequest_Content{
			Content: content,
		}
	} else {
		return nil, fmt.Errorf("request hash is empty")
	}

	res, err := c.client.Transform(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	res.RequestHash = hash

	// Store the response in the index
	if _, err = c.index.Store(ctx, res); err != nil {
		return nil, err
	}

	// For each extract in the response, create an extract request
	for _, extract := range res.GetExtracts() {
		extractReq := &rpc.ExtractRequest{
			PluginId:   req.PluginId,
			ParentHash: hash,
			Kind:       extract.Kind,
			Metadata:   extract.Metadata,
		}

		// Store the extract request
		if _, err := c.repo.StoreExtractRequest(ctx, extractReq); err != nil {
			return nil, fmt.Errorf("failed to store extract request: %w", err)
		}
	}

	// For each permanode in the response, store the permanode version
	for _, permanode := range res.GetPermanodes() {
		// Handle the different payload types
		var content index.Indexable
		switch p := permanode.Payload.(type) {
		case *rpc.TransformResponse_Permanode_Email:
			content = p.Email
		case *rpc.TransformResponse_Permanode_FinancialTransaction:
			content = p.FinancialTransaction
		default:
			return nil, fmt.Errorf("unknown payload type: %T", p)
		}

		if _, err := c.index.CreateDataSource(ctx, req.PluginId, permanode.Key, content); err != nil {
			return nil, fmt.Errorf("failed to create data source: %w", err)
		}
	}

	return res, nil
}
