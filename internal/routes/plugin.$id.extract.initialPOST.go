package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/rpc"
)

func PluginExtractInitialCreate(c echo.Context) error {
	b := middleware.GetBoot(c)
	plugin := c.Param("id")

	// Create a new extract
	req := rpc.ExtractRequest{
		PluginId: plugin,
		Kind:     "initial",
	}

	hash, err := b.Index.Store(c.Request().Context(), &req)
	if err != nil {
		return fmt.Errorf("failed to store extract: %w", err)
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/plugin/%s/extract/%s", plugin, hash))
	return c.NoContent(http.StatusCreated)
}
