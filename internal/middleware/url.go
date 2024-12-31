package middleware

import (
	"net/url"

	"github.com/labstack/echo/v4"
)

func FullURL(c echo.Context) *url.URL {
	req := c.Request()
	u := req.URL.Scheme

	if u == "" {
		if req.TLS == nil {
			u = "http"
		} else {
			u = "https"
		}
	}

	// should be safe because this all comes from the request. should never be invalid
	uu, _ := url.Parse(u + "://" + req.Host + req.URL.String())
	return uu
}
