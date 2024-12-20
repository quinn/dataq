package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/htmx"
	"go.quinn.io/dataq/internal/middleware"
)

func ContentDelete(c echo.Context) error {
	b := middleware.GetBoot(c)
	if err := b.CAS.Delete(c.Request().Context(), c.Param("hash")); err != nil {
		return fmt.Errorf("error deleting content: %w", err)
	}

	return htmx.Redirect(c, "/")
}
