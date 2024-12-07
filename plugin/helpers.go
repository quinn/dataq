package plugin

import (
	"fmt"
	"os"

	pb "go.quinn.io/dataq/proto"
	"go.quinn.io/dataq/stream"
)

type PluginAPI struct {
	plugin  Plugin
	request *pb.PluginRequest
}

func (p *PluginAPI) WriteError(err error) {
	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Error:     fmt.Sprintf("stream error: %v", err),
	})
}

func (p *PluginAPI) WriteDone() {
	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Done:      true,
	})
}

func (p *PluginAPI) WriteClosed() {
	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Closed:    true,
	})
}

func (p *PluginAPI) WriteItem(item *pb.DataItem) {
	item.Meta.Hash = GenerateHash(item.RawData)

	if p.request.Item != nil {
		item.Meta.ParentHash = p.request.Item.Meta.Hash
	}

	stream.WriteResponse(os.Stdout, &pb.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Item:      item,
	})
}
