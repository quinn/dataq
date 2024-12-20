package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PluginExtractSendCreate(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/")
}
