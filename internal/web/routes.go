package web

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/routes"
)

func addRoutes(e *echo.Echo) {
	e.POST("/plugin/:id/extract/initial", routes.PluginExtractInitialCreate)
	e.POST("/plugin/:id/extract/:hash/send", routes.PluginExtractSendCreate).Name = "plugin.extract.send"
	e.DELETE("/content/:hash", routes.ContentDelete).Name = "content.delete"
	/* insert new routes here */
}
