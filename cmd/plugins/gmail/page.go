package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
	"google.golang.org/api/gmail/v1"
)

func (p *GmailPlugin) getPage(ctx context.Context, action *pb.Action, api *plugin.PluginAPI) error {
	srv, err := p.getClient(ctx)
	if err != nil {
		return err
	}

	req := srv.Users.Messages.List("me").MaxResults(100)

	var pageToken string
	if action.Name == "next_page" {
		if action.Config["next_page_token"] == "" {
			return fmt.Errorf("next_page action requires next_page_token")
		}

		pageToken = action.Config["next_page_token"]

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

	// at some point this will be handled separately by the worker
	// separation here allows for easier debugging and reloading
	// the queue/worker system do not have a good way to support
	// different operations on the request currently
	p.transformPage(ctx, item, api)

	return nil
}

func (p *GmailPlugin) transformPage(_ context.Context, item *pb.DataItem, api *plugin.PluginAPI) error {
	var r gmail.ListMessagesResponse
	if err := json.Unmarshal(item.RawData, &r); err != nil {
		return fmt.Errorf("error unmarshaling response: %v", err)
	}

	for _, msg := range r.Messages {
		action := &pb.Action{
			Name:       "get_message",
			Id:         msg.Id,
			ParentHash: item.Meta.Hash,
			Config: map[string]string{
				"message_id": msg.Id,
			},
		}

		api.WriteAction(action)
	}

	if r.NextPageToken != "" {
		action := &pb.Action{
			Name:       "next_page",
			Id:         r.NextPageToken,
			ParentHash: item.Meta.Hash,
			Config: map[string]string{
				"next_page_token": r.NextPageToken,
			},
		}

		api.WriteAction(action)
	}

	return nil
}
