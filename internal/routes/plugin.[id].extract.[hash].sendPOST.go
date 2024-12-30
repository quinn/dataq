package routes

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/htmx"
	"go.quinn.io/dataq/rpc"
)

func PluginExtractSendCreate(c echo.Context) error {
	b := c.Get("boot").(*boot.Boot)

	squirrel.Select("hash").From("requests").Where(squirrel.Eq{"hash": c.Param("hash")})

	id := c.Param("id")
	plugin, ok := b.Plugins.Clients[id]
	if !ok {
		return c.String(http.StatusNotFound, "plugin not found")
	}

	var req rpc.ExtractRequest
	sel := b.Index.Q.Where(squirrel.Eq{"hash": c.Param("hash")})
	if err := b.Index.Get(c.Request().Context(), &req, sel); err != nil {
		return fmt.Errorf("error getting request from index: %w", err)
	}

	if _, err := plugin.Extract(c.Request().Context(), &req); err != nil {
		return fmt.Errorf("error extracting: %w", err)
	}

	return htmx.Refresh(c)
}
