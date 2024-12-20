package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/htmx"
	"go.quinn.io/dataq/rpc"
)

func PluginTransformSendCreate(c echo.Context) error {
	b := c.Get("boot").(*boot.Boot)

	id := c.Param("id")
	plugin, ok := b.Plugins.Clients[id]
	if !ok {
		return c.String(http.StatusNotFound, "plugin not found")
	}

	var req rpc.TransformRequest
	if err := b.Index.Get(c.Request().Context(), &req, "hash = ?", c.Param("hash")); err != nil {
		return fmt.Errorf("error getting request from index: %w", err)
	}

	if _, err := plugin.Transform(c.Request().Context(), &req); err != nil {
		return fmt.Errorf("error transforming: %w", err)
	}

	return htmx.Refresh(c)
}
