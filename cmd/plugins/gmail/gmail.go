package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.quinn.io/dataq/rpc"
	"go.quinn.io/dataq/schema"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// type Config struct {
// 	CredentialsPath string `json:"credentials_path"`
// 	TokenPath       string `json:"token_path"`
// }

type GmailPlugin struct {
	// credentialsPath string
	// tokenPath       string
}

func New() *GmailPlugin {
	return &GmailPlugin{
		// credentialsPath: config.CredentialsPath,
		// tokenPath:       config.TokenPath,
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

func (p *GmailPlugin) getClient(ctx context.Context, rpcoauth *rpc.OAuth2) (*gmail.Service, error) {
	if rpcoauth == nil {
		return nil, fmt.Errorf("no OAuth2 credentials provided")
	}
	// // Read credentials file
	// b, err := os.ReadFile(p.credentialsPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("error reading credentials (%s): %v", p.credentialsPath, err)
	// }

	type cred struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
	}
	var j struct {
		Installed *cred `json:"installed"`
	}

	j.Installed = &cred{
		ClientID:     rpcoauth.Config.ClientId,
		ClientSecret: rpcoauth.Config.ClientSecret,
		RedirectURIs: []string{rpcoauth.Config.RedirectUrl},
		AuthURI:      rpcoauth.Config.Endpoint.AuthUrl,
		TokenURI:     rpcoauth.Config.Endpoint.TokenUrl,
	}

	b, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("error marshaling credentials: %w", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing credentials: %w", err)
	}

	token := schema.NewOauthToken(rpcoauth)

	client := config.Client(ctx, token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating Gmail client: %v", err)
	}

	return srv, nil
}
