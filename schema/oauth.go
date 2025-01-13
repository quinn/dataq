package schema

import (
	"time"

	"go.quinn.io/dataq/rpc"
	"golang.org/x/oauth2"
)

func NewOauthConfig(inp *rpc.OAuth2) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     inp.Config.ClientId,
		ClientSecret: inp.Config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  inp.Config.Endpoint.AuthUrl,
			TokenURL: inp.Config.Endpoint.TokenUrl,
		},
		RedirectURL: inp.Config.RedirectUrl,
		Scopes:      inp.Config.Scopes,
	}
}

func NewOauthToken(inp *rpc.OAuth2) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  inp.Token.AccessToken,
		TokenType:    inp.Token.TokenType,
		RefreshToken: inp.Token.RefreshToken,
		ExpiresIn:    inp.Token.ExpiresIn,
		Expiry:       time.UnixMilli(inp.Token.Expiry),
	}
}

func NewRPCOauthToken(inp *oauth2.Token) *rpc.OAuth2_Token {
	return &rpc.OAuth2_Token{
		AccessToken:  inp.AccessToken,
		TokenType:    inp.TokenType,
		RefreshToken: inp.RefreshToken,
		ExpiresIn:    inp.ExpiresIn,
		Expiry:       inp.Expiry.UnixMilli(),
	}
}
