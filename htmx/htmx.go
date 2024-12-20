package htmx

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Redirect(c echo.Context, path string) error {
	c.Response().Header().Set("HX-Redirect", path)
	return c.NoContent(http.StatusOK)
}
