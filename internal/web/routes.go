package web

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/routes"
)

func addRoutes(e *echo.Echo) {
	e.POST("/plugin/:id/extract/initial", routes.PluginExtractInitialCreate)
	e.POST("/plugin/:id/extract/:hash/send", routes.PluginExtractSendCreate).Name = "plugin.extract.send"
	e.DELETE("/content/:hash", routes.ContentDelete).Name = "content.delete"
	e.POST("/plugin/:id/transform/:hash/send", routes.PluginTransformSendCreate).Name = "plugin.transform.send"
	e.POST("/plugin/install", routes.PluginInstallCreate).Name = "plugin.install"
	e.GET("/plugin/:id/oauth/complete", routes.PluginOauthComplete).Name = "plugin.oauth.complete"
	/* insert new routes here */
}
