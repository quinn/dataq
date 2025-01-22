package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
)

type BreadcrumbContextKey struct{}

type Breadcrumb struct {
	Name string
	URL  string
}

type BreadcrumbBuilder struct {
	c      echo.Context
	crumbs []Breadcrumb
}

func SetBreadcrumb(c echo.Context) *BreadcrumbBuilder {
	return &BreadcrumbBuilder{
		c: c,
	}
}

func (b *BreadcrumbBuilder) Add(name string, url string) *BreadcrumbBuilder {
	b.crumbs = append(b.crumbs, Breadcrumb{
		Name: name,
		URL:  url,
	})

	ctx := context.WithValue(b.c.Request().Context(), BreadcrumbContextKey{}, b.crumbs)
	b.c.SetRequest(b.c.Request().WithContext(ctx))

	return b
}

func Breadcrumbs(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var crumbs []Breadcrumb

		if strings.HasPrefix(c.Path(), "/plugin") {
			crumbs = append(crumbs, Breadcrumb{
				Name: "Plugins",
				URL:  "/plugin/install",
			})

			if id := c.Param("id"); id != "" {
				crumbs = append(crumbs, Breadcrumb{
					Name: id,
					URL:  "/plugin/" + id,
				})
			}
		}

		ctx := context.WithValue(c.Request().Context(), BreadcrumbContextKey{}, crumbs)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func GetBreadcrumbs(ctx context.Context) []Breadcrumb {
	return ctx.Value(BreadcrumbContextKey{}).([]Breadcrumb)
}
