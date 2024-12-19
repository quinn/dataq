package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/rpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func PluginExtractInitialCreate(c echo.Context) error {
	b := middleware.GetBoot(c)
	plugin := c.Param("id")

	// Create a new extract
	req := rpc.ExtractRequest{
		PluginId: plugin,
		Kind:     "initial",
	}

	jsn, err := protojson.Marshal(&req)
	if err != nil {
		return fmt.Errorf("failed to marshal extract request: %w", err)
	}

	r := bytes.NewReader(jsn)
	hash, err := b.CAS.Store(c.Request().Context(), r)
	if err != nil {
		return fmt.Errorf("failed to store extract request: %w", err)
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf("/plugin/%s/extract/%s", plugin, hash))
}
