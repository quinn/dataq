package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func PluginInstallCreate(c echo.Context) error {
	var form struct {
		AuthURL      string `form:"auth_url"`
		TokenURL     string `form:"token_url"`
		Scopes       string `form:"scopes"`
		PluginID     string `form:"plugin_id"`
		ClientID     string `form:"client_id"`
		ClientSecret string `form:"client_secret"`
	}

	return fmt.Errorf("form data: %v", form)
	// return c.Redirect(http.StatusFound, "/")
}
