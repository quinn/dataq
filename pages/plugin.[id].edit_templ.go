// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/schema"
	"go.quinn.io/dataq/ui"
	"net/http"
)

type PluginIdEditData struct {
	plugin schema.PluginInstance
	id     string
}

func PluginIdEditGET(c echo.Context, id string) (PluginIdEditData, error) {
	b := middleware.GetBoot(c)

	var data PluginIdEditData
	data.id = id

	if err := b.Index.GetPermanode(c.Request().Context(), id, &data.plugin); err != nil {
		return data, fmt.Errorf("failed to get plugin: %w", err)
	}

	return data, nil
}

func PluginIdEditPOST(c echo.Context, id string) error {
	b := middleware.GetBoot(c)

	if c.FormValue("delete") == "true" {
		if err := b.Index.Delete(c.Request().Context(), id); err != nil {
			return fmt.Errorf("failed to delete plugin: %w", err)
		}

		return c.Redirect(http.StatusFound, "/")
	}

	var plugin schema.PluginInstance
	if err := b.Index.GetPermanode(c.Request().Context(), id, &plugin); err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	var form struct {
		ClientID     string `form:"client_id"`
		ClientSecret string `form:"client_secret"`
	}

	if err := c.Bind(&form); err != nil {
		return fmt.Errorf("failed to bind form: %w", err)
	}

	plugin.OauthConfig.ClientID = form.ClientID
	plugin.OauthConfig.ClientSecret = form.ClientSecret

	if _, err := b.Index.UpdatePermanode(c.Request().Context(), id, &plugin); err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}

	return c.Redirect(http.StatusFound, c.Request().RequestURI)
}

func PluginIdEdit(data PluginIdEditData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		plugin := data.plugin
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2>Edit Plugin</h2><div class=\"space-y-3\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = ui.JsonBrowser(plugin).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form method=\"post\" class=\"space-y-3\"><div><label for=\"client_id\">Client ID</label><br><input class=\"input\" type=\"text\" name=\"client_id\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(plugin.OauthConfig.ClientID)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `pages/plugin.[id].edit.templ`, Line: 74, Col: 90}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div><label for=\"client_secret\">Client Secret</label><br><input class=\"input\" type=\"text\" name=\"client_secret\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(plugin.OauthConfig.ClientSecret)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `pages/plugin.[id].edit.templ`, Line: 78, Col: 98}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><button class=\"underline\" type=\"submit\">Save</button></form><form method=\"post\"><input type=\"hidden\" name=\"delete\" value=\"true\"> <button class=\"underline text-red-700\" type=\"submit\">Delete</button></form><a class=\"underline block\" href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 templ.SafeURL = templ.URL("/plugin/" + data.id + "/oauth/begin")
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var5)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Connect Oauth</a></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = ui.Layout().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
