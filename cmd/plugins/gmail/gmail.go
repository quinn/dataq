package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Config struct {
	CredentialsPath string `json:"credentials_path"`
	TokenPath       string `json:"token_path"`
}

type GmailPlugin struct {
	credentialsPath string
	tokenPath       string
}

func New() *GmailPlugin {
	// Load config from environment
	configJSON := os.Getenv("CONFIG")
	if configJSON == "" {
		log.Fatal("CONFIG environment variable not set")
	}

	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		log.Fatalf("Error parsing CONFIG: %v", err)
	}

	// Convert paths to absolute if they're relative
	if !filepath.IsAbs(config.CredentialsPath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Unable to get working directory: %v", err)
		}
		config.CredentialsPath = filepath.Join(wd, config.CredentialsPath)
	}

	if !filepath.IsAbs(config.TokenPath) {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Unable to get working directory: %v", err)
		}
		config.TokenPath = filepath.Join(wd, config.TokenPath)
	}

	log.Printf("Using credentials path: %s", config.CredentialsPath)
	log.Printf("Using token path: %s", config.TokenPath)

	return &GmailPlugin{
		credentialsPath: config.CredentialsPath,
		tokenPath:       config.TokenPath,
	}
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

func (p *GmailPlugin) getClient(ctx context.Context) (*gmail.Service, error) {
	// Read credentials file
	b, err := os.ReadFile(p.credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading credentials (%s): %v", p.credentialsPath, err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing credentials (%s): %v", p.credentialsPath, err)
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
		return nil, fmt.Errorf("error creating Gmail client: %v", err)
	}

	return srv, nil
}
