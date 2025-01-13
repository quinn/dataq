// Code generated by ccf. DO NOT EDIT.
package router

import (
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/pages"
)

// RegisterRoutes adds all page routes to the Echo instance
func RegisterRoutes(e *echo.Echo) {
	e.GET("/content/:hash", ContentHashGET)
	e.GET("/content", ContentGET)
	e.GET("/", IndexGET)
	e.GET("/plugin/:id/edit", PluginIdEditGET)
	e.POST("/plugin/:id/edit", PluginIdEditPOST)
	e.GET("/plugin/:id/extract/:hash", PluginIdExtractHashGET)
	e.GET("/plugin/:id/oauth/begin", PluginIdOauthBeginGET)
	e.POST("/plugin/:id/oauth/begin", PluginIdOauthBeginPOST)
	e.GET("/plugin/:id/send/:reqtype/:kind", PluginIdSendReqtypeKindGET)
	e.POST("/plugin/:id/send/:reqtype/:kind", PluginIdSendReqtypeKindPOST)
	e.GET("/plugin/:id", PluginIdGET)
	e.GET("/plugin/:id/transform/:hash", PluginIdTransformHashGET)
	e.GET("/plugin/install", PluginInstallGET)
	e.POST("/plugin/install", PluginInstallPOST)
}

// ContentHashGET handles GET requests to /content/:hash
func ContentHashGET(c echo.Context) error {
	result, err := pages.ContentHashGET(c, c.Param("hash"))
	if err != nil {
		return err
	}
	return pages.ContentHash(result).Render(c.Request().Context(), c.Response().Writer)
}

// ContentGET handles GET requests to /content
func ContentGET(c echo.Context) error {
	result, err := pages.ContentGET(c)
	if err != nil {
		return err
	}
	return pages.Content(result).Render(c.Request().Context(), c.Response().Writer)
}

// IndexGET handles GET requests to /
func IndexGET(c echo.Context) error {
	result, err := pages.IndexGET(c)
	if err != nil {
		return err
	}
	return pages.Index(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdEditGET handles GET requests to /plugin/:id/edit
func PluginIdEditGET(c echo.Context) error {
	result, err := pages.PluginIdEditGET(c, c.Param("id"))
	if err != nil {
		return err
	}
	return pages.PluginIdEdit(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdEditPOST handles POST requests to /plugin/:id/edit
func PluginIdEditPOST(c echo.Context) error {
	return pages.PluginIdEditPOST(c, c.Param("id"))
}

// PluginIdExtractHashGET handles GET requests to /plugin/:id/extract/:hash
func PluginIdExtractHashGET(c echo.Context) error {
	result, err := pages.PluginIdExtractHashGET(c, c.Param("id"), c.Param("hash"))
	if err != nil {
		return err
	}
	return pages.PluginIdExtractHash(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdOauthBeginGET handles GET requests to /plugin/:id/oauth/begin
func PluginIdOauthBeginGET(c echo.Context) error {
	result, err := pages.PluginIdOauthBeginGET(c, c.Param("id"))
	if err != nil {
		return err
	}
	return pages.PluginIdOauthBegin(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdOauthBeginPOST handles POST requests to /plugin/:id/oauth/begin
func PluginIdOauthBeginPOST(c echo.Context) error {
	return pages.PluginIdOauthBeginPOST(c, c.Param("id"))
}

// PluginIdSendReqtypeKindGET handles GET requests to /plugin/:id/send/:reqtype/:kind
func PluginIdSendReqtypeKindGET(c echo.Context) error {
	result, err := pages.PluginIdSendReqtypeKindGET(c, c.Param("id"), c.Param("reqtype"), c.Param("kind"))
	if err != nil {
		return err
	}
	return pages.PluginIdSendReqtypeKind(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdSendReqtypeKindPOST handles POST requests to /plugin/:id/send/:reqtype/:kind
func PluginIdSendReqtypeKindPOST(c echo.Context) error {
	return pages.PluginIdSendReqtypeKindPOST(c, c.Param("id"), c.Param("reqtype"), c.Param("kind"))
}

// PluginIdGET handles GET requests to /plugin/:id
func PluginIdGET(c echo.Context) error {
	result, err := pages.PluginIdGET(c, c.Param("id"))
	if err != nil {
		return err
	}
	return pages.PluginId(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginIdTransformHashGET handles GET requests to /plugin/:id/transform/:hash
func PluginIdTransformHashGET(c echo.Context) error {
	result, err := pages.PluginIdTransformHashGET(c, c.Param("id"), c.Param("hash"))
	if err != nil {
		return err
	}
	return pages.PluginIdTransformHash(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginInstallGET handles GET requests to /plugin/install
func PluginInstallGET(c echo.Context) error {
	result, err := pages.PluginInstallGET(c)
	if err != nil {
		return err
	}
	return pages.PluginInstall(result).Render(c.Request().Context(), c.Response().Writer)
}

// PluginInstallPOST handles POST requests to /plugin/install
func PluginInstallPOST(c echo.Context) error {
	return pages.PluginInstallPOST(c)
}
