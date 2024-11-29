package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	pb "go.quinn.io/dataq/proto"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailPlugin struct {
	credentialsPath string
	tokenPath       string
	config          map[string]string
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
	log.Printf("Gmail plugin config: %+v", config)
	if creds, ok := config["credentials_path"]; ok {
		// Convert to absolute path if relative
		if !filepath.IsAbs(creds) {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("unable to get working directory: %v", err)
			}
			log.Printf("Working directory: %s", wd)
			creds = filepath.Join(wd, creds)
		}
		log.Printf("Using credentials path: %s", creds)
		p.credentialsPath = creds
	} else {
		return fmt.Errorf("credentials_path configuration is required")
	}

	if token, ok := config["token_path"]; ok {
		// Convert to absolute path if relative
		if !filepath.IsAbs(token) {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("unable to get working directory: %v", err)
			}
			token = filepath.Join(wd, token)
		}
		log.Printf("Using token path: %s", token)
		p.tokenPath = token
	} else {
		return fmt.Errorf("token_path configuration is required")
	}

	p.config = config

	return nil
}

func (p *GmailPlugin) Extract(ctx context.Context) (<-chan *pb.DataItem, error) {
	items := make(chan *pb.DataItem)

	// Read credentials file
	b, err := os.ReadFile(p.credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file (%s): %v", p.credentialsPath, err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials (%s): %v", p.credentialsPath, err)
	}

	// Read token file
	tokenBytes, err := os.ReadFile(p.tokenPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read token file (%s): %v", p.tokenPath, err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenBytes, &token); err != nil {
		return nil, fmt.Errorf("unable to parse token (%s): %v", p.tokenPath, err)
	}

	client := config.Client(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail client: %v", err)
	}

	go func() {
		defer close(items)

		var pageToken string

		// If we have previous response data, extract the next page token
		if prevData := p.config["response_data"]; prevData != "" {
			var prevResp gmail.ListMessagesResponse
			if err := json.Unmarshal([]byte(prevData), &prevResp); err != nil {
				log.Printf("Error parsing previous response: %v", err)
				return
			}
			pageToken = prevResp.NextPageToken
		}

		// If no previous response or invalid, start from beginning
		req := srv.Users.Messages.List("me").MaxResults(100)
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		r, err := req.Do()
		if err != nil {
			log.Printf("Error listing messages: %v", err)
			return
		}

		if r.Messages == nil {
			log.Printf("No messages found")
			return
		}

		// Marshal the response to JSON
		rawJSON, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			return
		}

		// Create a DataItem for the page
		item := &pb.DataItem{
			PluginId:    p.ID(),
			SourceId:    fmt.Sprintf("page_%s", pageToken),
			ContentType: "application/json",
			RawData:     rawJSON,
		}

		select {
		case items <- item:
		case <-ctx.Done():
			return
		}
	}()

	return items, nil
}
