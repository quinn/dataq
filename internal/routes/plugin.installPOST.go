package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/schema"
	"golang.org/x/oauth2"
)

func PluginInstallCreate(c echo.Context) error {
	b := middleware.GetBoot(c)

	var form struct {
		AuthURL      string `form:"auth_url"`
		TokenURL     string `form:"token_url"`
		Scopes       string `form:"scopes"`
		PluginID     string `form:"plugin_id"`
		ClientID     string `form:"client_id"`
		ClientSecret string `form:"client_secret"`
	}

	if err := c.Bind(&form); err != nil {
		return fmt.Errorf("failed to bind form: %v", err)
	}

	pluginInstance := schema.PluginInstance{
		PluginID: form.PluginID,
		Label:    form.PluginID,
		OauthConfig: &oauth2.Config{
			Scopes:       strings.Split(form.Scopes, ","),
			ClientID:     form.ClientID,
			ClientSecret: form.ClientSecret,

			Endpoint: oauth2.Endpoint{
				AuthURL:  form.AuthURL,
				TokenURL: form.TokenURL,
			},
		},
	}

	permanodeHash, err := b.Index.CreatePermanode(c.Request().Context(), &pluginInstance)
	if err != nil {
		return fmt.Errorf("failed to create permanode: %v", err)
	}

	return c.Redirect(http.StatusFound, "/plugin/"+permanodeHash+"/oauth/begin")
}
