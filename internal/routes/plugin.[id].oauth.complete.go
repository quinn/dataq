package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/schema"
)

func PluginOauthComplete(c echo.Context) error {
	b := middleware.GetBoot(c)
	id := c.Param("id")
	code := c.QueryParam("code")

	var plugin schema.PluginInstance
	if err := b.Index.GetPermanode(c.Request().Context(), id, &plugin); err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	token, err := plugin.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web: (%v) %w", plugin.OauthConfig, err)
	}

	// Save the token
	plugin.OauthToken = token
	if _, err := b.Index.UpdatePermanode(c.Request().Context(), id, &plugin); err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}

	return c.Redirect(http.StatusFound, "/plugin/"+id+"/edit")
}
