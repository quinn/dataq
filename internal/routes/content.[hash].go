package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stoewer/go-strcase"
	"go.quinn.io/dataq/internal/middleware"
)

func Content(c echo.Context) error {
	hash := c.Param("hash")
	b := middleware.GetBoot(c)
	claims, err := b.Index.Query(c.Request().Context(), b.Index.Q.Where("content_hash = ?", hash))
	if err != nil {
		return fmt.Errorf("failed to get claims: %w", err)
	}
	if len(claims) == 0 {
		return c.Redirect(http.StatusFound, "/blob/"+hash)
	}
	claim := claims[0]

	return c.Redirect(http.StatusFound,
		"/schema/"+strcase.KebabCase(claim.SchemaKind)+"/"+claim.ContentHash)
}
