package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/schema"
)

func PluginOauthComplete(c echo.Context) error {
	b := middleware.GetBoot(c)
	hash := c.Param("hash")
	code := c.QueryParam("code")

	var plugin schema.PluginInstance
	if err := b.Index.GetPermanode(c.Request().Context(), hash, &plugin); err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	oauthConfig := schema.NewOauthConfig(plugin.Oauth)
	token, err := oauthConfig.Exchange(c.Request().Context(), code)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web: (%v) %w", plugin.Oauth.Config, err)
	}

	// Save the token
	plugin.Oauth.Token = schema.NewRPCOauthToken(token)
	if _, err := b.Index.UpdatePermanode(c.Request().Context(), hash, &plugin); err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}

	return c.Redirect(http.StatusFound, "/plugin/"+hash+"/edit")
}
