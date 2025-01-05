package main

import (
	"context"

	"go.quinn.io/dataq/rpc"
	"golang.org/x/oauth2/fitbit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	rpc.UnimplementedDataQPluginServer
	client *FitbitClient
}

func NewServer(client *FitbitClient) *server {
	return &server{client: client}
}

func (s *server) Install(ctx context.Context, req *rpc.InstallRequest) (*rpc.InstallResponse, error) {
	return &rpc.InstallResponse{
		PluginId: "fitbit",
		OauthConfig: &rpc.OauthConfig{
			AuthUrl:  fitbit.Endpoint.AuthURL,
			TokenUrl: fitbit.Endpoint.TokenURL,
			Scopes: []string{
				"activity",
				"heartrate",
				"profile",
				"sleep",
				"weight",
			},
		},
		Extracts: []*rpc.InstallResponse_Extract{
			{
				Kind:        "make_request",
				Label:       "Make a Request",
				Description: "Make a request to the Fitbit API by specifing the request path",
				Configs: []*rpc.PluginConfig{
					{
						Key:   "request_path",
						Label: "Request Path",
					},
				},
			},
		},
	}, nil
}

func (s *server) Extract(ctx context.Context, req *rpc.ExtractRequest) (*rpc.ExtractResponse, error) {
	reqHash, err := getReqHash(ctx)
	if err != nil {
		return nil, err
	}

	data, err := s.client.GetTodaySteps(ctx)
	if err != nil {
		return nil, err
	}

	return &rpc.ExtractResponse{
		Kind: "steps",
		Data: &rpc.ExtractResponse_Content{
			Content: data,
		},
		RequestHash: reqHash,
	}, nil
}

func (*server) Transform(context.Context, *rpc.TransformRequest) (*rpc.TransformResponse, error) {
	return nil, nil
}

func getReqHash(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "failed to get metadata")
	}

	hashes := md.Get("request-hash")
	if len(hashes) != 1 {
		return "", status.Errorf(codes.InvalidArgument, "request-hash metadata has multiple or zero values")
	}
	return hashes[0], nil
}
