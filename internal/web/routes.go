package web

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/routes"
)

func addRoutes(e *echo.Echo) {
	e.POST("/plugin/:id/extract/:hash/send", routes.PluginExtractSendCreate).Name = "plugin.extract.send"
	e.DELETE("/content/:hash", routes.ContentDelete).Name = "content.delete"
	e.GET("/plugin/:id/oauth/complete", routes.PluginOauthComplete).Name = "plugin.oauth.complete"
	/* insert new routes here */
}
