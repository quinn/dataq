package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ReverseContextKey struct{}

func EchoContext(e *echo.Echo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := context.WithValue(c.Request().Context(), ReverseContextKey{}, e)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func Reverse(ctx context.Context, routeName string, params ...interface{}) string {
	e := ctx.Value(ReverseContextKey{}).(*echo.Echo)
	return e.Reverse(routeName, params...)
}
