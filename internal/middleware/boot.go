package middleware

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/boot"
)

const BootContextKey = "boot"

func BootContext(b *boot.Boot) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(BootContextKey, b)
			return next(c)
		}
	}
}

func GetBoot(c echo.Context) *boot.Boot {
	return c.Get(BootContextKey).(*boot.Boot)
}
