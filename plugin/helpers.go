package plugin

import (
	"fmt"
	"os"

	"go.quinn.io/dataq/hash"
	"go.quinn.io/dataq/schema"
	"go.quinn.io/dataq/stream"
)

type PluginAPI struct {
	plugin  Plugin
	request *schema.PluginRequest
}

func (p *PluginAPI) WriteError(err error) {
	stream.WriteResponse(os.Stdout, &schema.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Error:     fmt.Sprintf("stream error: %v", err),
	})
}

func (p *PluginAPI) WriteDone() {
	stream.WriteResponse(os.Stdout, &schema.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Done:      true,
	})
}

func (p *PluginAPI) WriteClosed() {
	stream.WriteResponse(os.Stdout, &schema.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Closed:    true,
	})
}

func (p *PluginAPI) WriteItem(item *schema.DataItem) {
	item.Meta.Hash = hash.Generate(item.RawData)

	if p.request.Action != nil {
		item.Meta.ParentHash = p.request.Action.ParentHash
	}

	stream.WriteResponse(os.Stdout, &schema.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Item:      item,
	})
}

func (p *PluginAPI) WriteAction(action *schema.Action) {
	stream.WriteResponse(os.Stdout, &schema.PluginResponse{
		PluginId:  p.plugin.ID(),
		RequestId: p.request.Id,
		Action:    action,
	})
}
