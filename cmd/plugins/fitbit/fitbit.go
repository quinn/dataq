package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.quinn.io/dataq/rpc"
	"golang.org/x/oauth2"
)

// FitbitClient represents a client for Fitbit API
type FitbitClient struct {
	oauthConfig *oauth2.Config
	oauthToken  *oauth2.Token
}

// NewFitbitClient initializes a new FitbitClient with the given credentials
func NewFitbitClient(oauth *rpc.OAuth2) *FitbitClient {
	if oauth.Config == nil {
		log.Fatal("OauthConfig not set in CONFIG")
	}

	if oauth.Token == nil {
		log.Fatal("OauthToken not set in CONFIG")
	}

	log.Println("Using config:", oauth)

	client := &FitbitClient{
		oauthConfig: &oauth2.Config{
			ClientID:     oauth.Config.ClientId,
			ClientSecret: oauth.Config.ClientSecret,
			RedirectURL:  oauth.Config.RedirectUrl,
			Scopes:       oauth.Config.Scopes,

			Endpoint: oauth2.Endpoint{
				AuthURL:  oauth.Config.Endpoint.AuthUrl,
				TokenURL: oauth.Config.Endpoint.TokenUrl,
			},
		},
		oauthToken: &oauth2.Token{
			AccessToken:  oauth.Token.AccessToken,
			TokenType:    oauth.Token.TokenType,
			RefreshToken: oauth.Token.RefreshToken,
			ExpiresIn:    oauth.Token.ExpiresIn,
			Expiry:       time.UnixMilli(oauth.Token.Expiry),
		},
	}

	// Try to load existing token
	return client
}

// // GetAuthURL generates the URL for user authorization
// func (c *FitbitClient) GetAuthURL() string {
// 	return c.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
// }

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

	client := c.oauthConfig.Client(ctx, c.oauthToken)

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

func (c *FitbitClient) MakeRequest(ctx context.Context, path string) ([]byte, error) {
	client := c.oauthConfig.Client(ctx, c.oauthToken)

	url := fmt.Sprintf("https://api.fitbit.com%s", path)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch steps data: %v", err)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
