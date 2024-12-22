package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

type Config struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectURL"`
	AuthURL      string `json:"authURL"`
	RefreshURL   string `json:"refreshURL"`

	// These come from the authed Oauth2 token
	Token string `json:"token"`
}

// FitbitClient represents a client for Fitbit API
type FitbitClient struct {
	oauthConfig *oauth2.Config
	config      Config
}

// NewFitbitClient initializes a new FitbitClient with the given credentials
func NewFitbitClient() *FitbitClient {
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	homeDir = "."
	// }
	// tokenPath := filepath.Join(homeDir, ".dataq", "fitbit_token.json")

	// Load config from environment
	configJSON := os.Getenv("CONFIG")
	if configJSON == "" {
		log.Fatal("CONFIG environment variable not set")
	}

	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		log.Fatalf("Error parsing CONFIG: %v", err)
	}

	log.Println("Using config:", config)

	client := &FitbitClient{
		oauthConfig: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes: []string{
				"activity",
				"heartrate",
				"profile",
				"sleep",
				"weight",
			},
			Endpoint: fitbit.Endpoint,
		},
		config: config,
	}

	// Try to load existing token
	return client
}

// GetAuthURL generates the URL for user authorization
func (c *FitbitClient) GetAuthURL() string {
	return c.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

// RefreshToken refreshes the OAuth token if it's expired
// func (c *FitbitClient) RefreshToken(ctx context.Context) error {
// 	if c.config.AccessToken == nil {
// 		return fmt.Errorf("no token available to refresh")
// 	}

// 	token, err := c.oauthConfig.TokenSource(ctx, c.config.AccessToken).Token()
// 	if err != nil {
// 		return fmt.Errorf("failed to refresh token: %v", err)
// 	}

// 	c.config.AccessToken = token
// 	if err := c.saveToken(); err != nil {
// 		return fmt.Errorf("failed to save refreshed token: %v", err)
// 	}

// 	return nil
// }

// loadToken attempts to load the token from the filesystem
// func (c *FitbitClient) loadToken() error {
// 	data, err := os.ReadFile(c.tokenPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read token file: %v", err)
// 	}

// 	token := &oauth2.Token{}
// 	if err := json.Unmarshal(data, token); err != nil {
// 		return fmt.Errorf("failed to parse token: %v", err)
// 	}

// 	c.token = token
// 	return nil
// }

// saveToken saves the current token to the filesystem
// func (c *FitbitClient) saveToken() error {
// 	if c.token == nil {
// 		return fmt.Errorf("no token to save")
// 	}

// 	// Ensure directory exists
// 	dir := filepath.Dir(c.tokenPath)
// 	if err := os.MkdirAll(dir, 0700); err != nil {
// 		return fmt.Errorf("failed to create token directory: %v", err)
// 	}

// 	data, err := json.Marshal(c.token)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal token: %v", err)
// 	}

// 	if err := os.WriteFile(c.tokenPath, data, 0600); err != nil {
// 		return fmt.Errorf("failed to save token: %v", err)
// 	}

// 	return nil
// }

// // Token returns the current OAuth2 token
// func (c *FitbitClient) Token() *oauth2.Token {
// 	return c.token
// }

// // HasValidToken checks if the client has a valid token
// func (c *FitbitClient) HasValidToken() bool {
// 	return c.token != nil && c.token.Valid()
// }

// StepData represents the step count data from Fitbit API
type StepData struct {
	Activities []struct {
		Steps int `json:"steps"`
	} `json:"activities-steps"`
}

// GetTodaySteps fetches the step count for the current day
func (c *FitbitClient) GetTodaySteps(ctx context.Context) ([]byte, error) {
	// if !c.HasValidToken() {
	// 	return nil, fmt.Errorf("no valid token available")
	// }

	var oauth2Token oauth2.Token
	if err := json.Unmarshal([]byte(c.config.Token), &oauth2Token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %v. (%s)", err, c.config.ClientSecret)
	}

	client := c.oauthConfig.Client(ctx, &oauth2Token)

	// Format today's date in the required format (YYYY-MM-DD)
	today := time.Now().Format("2006-01-02")

	// Construct the Fitbit API URL for daily activity summary
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/steps/date/%s/30d.json", today)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch steps data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
