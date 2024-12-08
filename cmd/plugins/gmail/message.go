package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.quinn.io/dataq/plugin"
	pb "go.quinn.io/dataq/proto"
)

func (p *GmailPlugin) getMessage(ctx context.Context, action *pb.Action, api *plugin.PluginAPI) error {
	srv, err := p.getClient(ctx)
	if err != nil {
		return err
	}

	messageID := action.Config["message_id"]
	if messageID == "" {
		return fmt.Errorf("get_message action requires message_id")
	}

	msg, err := srv.Users.Messages.Get("me", messageID).Do()
	if err != nil {
		return fmt.Errorf("error getting message %s: %v", messageID, err)
	}

	// Marshal the message to JSON
	rawJSON, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	// Create a DataItem for the message
	item := &pb.DataItem{
		Meta: &pb.DataItemMetadata{
			PluginId:    p.ID(),
			Id:          messageID,
			Kind:        "message",
			ContentType: "application/json",
			Timestamp:   msg.InternalDate, // Use the message's internal date
		},
		RawData: rawJSON,
	}

	api.WriteItem(item)
	return nil
}
