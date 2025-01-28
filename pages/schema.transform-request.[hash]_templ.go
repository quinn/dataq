// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"go.quinn.io/dataq/boot"
	"go.quinn.io/dataq/htmx"
	"go.quinn.io/dataq/internal/middleware"
	"go.quinn.io/dataq/rpc"
	"go.quinn.io/dataq/ui"
)

type SchemaTransformRequestHashData struct {
	req *rpc.TransformRequest
	res []struct {
		res *rpc.TransformResponse
	}
	hash string
}

func SchemaTransformRequestHashGET(c echo.Context, hash string) (SchemaTransformRequestHashData, error) {
	var data SchemaTransformRequestHashData
	b := middleware.GetBoot(c)

	var req rpc.TransformRequest
	sel := b.Index.Q.Where("content_hash = ?", hash)
	if err := b.Index.Get(c.Request().Context(), &req, sel); err != nil {
		return data, fmt.Errorf("transform request not found: %w", err)
	}

	sel = b.Index.Q.Where("request_hash = ?", hash)
	claims, err := b.Index.Query(c.Request().Context(), sel)
	if err != nil {
		return data, fmt.Errorf("failed to query response claims: %w", err)
	}

	for _, claim := range claims {
		var res rpc.TransformResponse
		sel := b.Index.Q.Where("content_hash = ?", claim.ContentHash)
		if err := b.Index.Get(c.Request().Context(), &res, sel); err != nil {
			return data, fmt.Errorf("transform response not found (%s): %w", claim.ContentHash, err)
		}

		data.res = append(data.res, struct {
			res *rpc.TransformResponse
		}{res: &res})
	}

	data.req = &req
	data.hash = hash
	return data, nil
}

func SchemaTransformRequestHashPOST(c echo.Context, hash string) error {
	b := c.Get("boot").(*boot.Boot)

	var req rpc.TransformRequest
	sel := b.Index.Q.Where(squirrel.Eq{"content_hash": hash})
	if err := b.Index.Get(c.Request().Context(), &req, sel); err != nil {
		return fmt.Errorf("error getting request from index: %w", err)
	}

	plugin, err := b.Repo.GetPluginInstance(c.Request().Context(), req.PluginId)
	if err != nil {
		return fmt.Errorf("error getting plugin instance: %w", err)
	}

	client, ok := b.Plugins.Clients[plugin.PluginID]
	if !ok {
		var keys []string
		for k := range b.Plugins.Clients {
			keys = append(keys, k)
		}
		return fmt.Errorf("plugin not found: %s. Plugin IDs are: %v", plugin.PluginID, keys)
	}

	if _, err := client.Transform(c.Request().Context(), &req); err != nil {
		return fmt.Errorf("error transforming: %w", err)
	}

	return htmx.Refresh(c)
}

func SchemaTransformRequestHash(data SchemaTransformRequestHashData) templ.Component {
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
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"space-y-3\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = ui.JsonBrowser(data.req).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<hr><div class=\"font-bold\">Responses</div><ul class=\"list-disc list-inside\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, res := range data.res {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<li class=\"list-item\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = ui.JsonBrowser(res.res).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</ul><hr><div class=\"font-bold\">Actions</div><ul class=\"list-disc list-inside\"><li class=\"list-item\"><button hx-post class=\"underline\">Send</button></li></ul></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = ui.Layout().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
