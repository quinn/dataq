package main

import (
	"context"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
)

func (p *GmailPlugin) getMessage(ctx context.Context, preq *pb.PluginRequest, api *plugin.PluginAPI) error {
	return nil
}
