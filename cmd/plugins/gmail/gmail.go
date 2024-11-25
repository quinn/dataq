package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	pb "go.quinn.io/dataq/proto"
	"time"
)

type GmailPlugin struct {
	credentialsJSON string
	tokenJSON      string
}

func New() *GmailPlugin {
	return &GmailPlugin{}
}

func (p *GmailPlugin) ID() string {
	return "gmail"
}

func (p *GmailPlugin) Name() string {
	return "Gmail Scanner"
}

func (p *GmailPlugin) Description() string {
	return "Extracts emails and metadata from Gmail using the Gmail API"
}

func (p *GmailPlugin) Configure(config map[string]string) error {
	if creds, ok := config["credentials_json"]; ok {
		p.credentialsJSON = creds
	} else {
		return fmt.Errorf("credentials_json configuration is required")
	}

	if token, ok := config["token_json"]; ok {
		p.tokenJSON = token
	} else {
		return fmt.Errorf("token_json configuration is required")
	}

	return nil
}

func (p *GmailPlugin) Extract(ctx context.Context) (<-chan *pb.DataItem, error) {
	items := make(chan *pb.DataItem)

	config, err := google.ConfigFromJSON([]byte(p.credentialsJSON), gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal([]byte(p.tokenJSON), &token); err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}

	client := config.Client(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail client: %v", err)
	}

	go func() {
		defer close(items)

		var pageToken string
		for {
			req := srv.Users.Messages.List("me").MaxResults(100)
			if pageToken != "" {
				req.PageToken(pageToken)
			}

			r, err := req.Do()
			if err != nil {
				fmt.Printf("Error listing messages: %v\n", err)
				return
			}

			for _, msg := range r.Messages {
				message, err := srv.Users.Messages.Get("me", msg.Id).Format("raw").Do()
				if err != nil {
					fmt.Printf("Error getting message %s: %v\n", msg.Id, err)
					continue
				}

				// Create DataItem for the email
				item := &pb.DataItem{
					PluginId:     p.ID(),
					SourceId:     message.Id,
					Timestamp:    message.InternalDate / 1000, // Convert to seconds
					ContentType:  "message/rfc822",
					RawData:     []byte(message.Raw),
					Metadata:    make(map[string]string),
				}

				// Add headers to metadata
				for _, header := range message.Payload.Headers {
					item.Metadata[header.Name] = header.Value
				}

				select {
				case items <- item:
				case <-ctx.Done():
					return
				}
			}

			if r.NextPageToken == "" {
				break
			}
			pageToken = r.NextPageToken
		}
	}()

	return items, nil
}
