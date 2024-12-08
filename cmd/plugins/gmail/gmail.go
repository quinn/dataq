package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.quinn.io/dataq/plugin"
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

func (p *GmailPlugin) Extract(ctx context.Context, preq *pb.PluginRequest, api *plugin.PluginAPI) error {
	switch preq.Action.Name {
	case "initial", "next_page":
		return p.getPage(ctx, preq, api)
	case "get_message":
		return p.getMessage(ctx, preq, api)
	default:
		return fmt.Errorf("unknown action: %s", preq.Action.Name)
	}
}

func (p *GmailPlugin) getClient(ctx context.Context) (*gmail.Service, error) {
	// Read credentials file
	b, err := os.ReadFile(p.credentialsPath)
	if err != nil {
		log.Printf("Error reading credentials: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Printf("Error parsing credentials (%s): %v", p.credentialsPath, err)
		return nil, err
	}

	// Read token file
	tokenBytes, err := os.ReadFile(p.tokenPath)
	if err != nil {
		return nil, fmt.Errorf("error reading token file (%s): %v", p.tokenPath, err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(tokenBytes, &token); err != nil {
		return nil, fmt.Errorf("error parsing token (%s): %v", p.tokenPath, err)
	}

	client := config.Client(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Error creating Gmail client: %v", err)
		return nil, err
	}
	return srv, nil
}
