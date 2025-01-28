package web

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/routes"
)

func addRoutes(e *echo.Echo) {
	e.GET("/plugin/:hash/oauth/complete", routes.PluginOauthComplete).Name = "plugin.oauth.complete"
	e.GET("/content/:hash", routes.Content).Name = "content"
	/* insert new routes here */
}
