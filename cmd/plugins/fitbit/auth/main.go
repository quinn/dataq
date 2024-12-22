package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

var (
	tokenFile = flag.String("token", "token.json", "Path to save OAuth2 token")
)

func main() {
	flag.Parse()

	// TODO: Maybe should read these from the config file?
	config := &oauth2.Config{
		ClientID:     os.Getenv("FitBitClientID"),
		ClientSecret: os.Getenv("FitBitClientSecret"),
		RedirectURL:  os.Getenv("FitBitRedirectURL"),
		Scopes: []string{
			"activity",
			"heartrate",
			"profile",
			"sleep",
			"weight",
		},
		Endpoint: fitbit.Endpoint,
	}

	token, err := getTokenFromWeb(config)
	if err != nil {
		log.Fatalf("Unable to get token: %v", err)
	}

	// Save the token
	f, err := os.OpenFile(*tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	// Print the credentials and token in the format needed for the plugin
	fmt.Println("\nUse these values in your plugin configuration:")

	fmt.Println("\ntoken_json:")
	tokenBytes, _ := json.Marshal(token)
	fmt.Println(string(tokenBytes))
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Use localhost callback with a random port
	config.RedirectURL = "http://localhost:8085"

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n\n", authURL)

	// Start local server to receive the code
	codeCh := make(chan string)
	server := &http.Server{Addr: ":8085"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			codeCh <- code
			fmt.Fprintf(w, "Authorization code received! You can close this window.")
		}
	})

	go server.ListenAndServe()
	defer server.Close()

	code := <-codeCh

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return token, nil
}
