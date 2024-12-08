package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
)

func (p *GmailPlugin) getPage(ctx context.Context, preq *pb.PluginRequest, api *plugin.PluginAPI) error {
	srv, err := p.getClient(ctx)
	if err != nil {
		return err
	}

	req := srv.Users.Messages.List("me").MaxResults(100)

	var pageToken string
	if preq.Action.Name == "next_page" {
		if preq.Action.Config["next_page_token"] == "" {
			return fmt.Errorf("next_page action requires next_page_token")
		}

		pageToken = preq.Action.Config["next_page_token"]

		req.PageToken(pageToken)
	}

	r, err := req.Do()
	if err != nil {
		return fmt.Errorf("error listing messages: %v", err)
	}

	if r.Messages == nil {
		return fmt.Errorf("no messages found")
	}

	// Marshal the response to JSON
	rawJSON, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling response: %v", err)
	}

	// Create a DataItem for the page
	item := &pb.DataItem{
		Meta: &pb.DataItemMetadata{
			PluginId:    p.ID(),
			Id:          pageToken,
			Kind:        "page",
			ContentType: "application/json",
		},
		RawData: rawJSON,
	}

	api.WriteItem(item)

	for _, msg := range r.Messages {
		action := &pb.Action{
			Name:       "get_message",
			ParentHash: item.Meta.Hash,
			Config: map[string]string{
				"message_id": msg.Id,
			},
		}

		api.WriteAction(action)
	}

	return nil
}
