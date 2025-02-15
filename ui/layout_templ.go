// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package ui

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"go.quinn.io/ccf/assets"
	"go.quinn.io/dataq/internal/middleware"
)

func Layout() templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\" class=\"h-full\"><head><meta charset=\"UTF-8\"><title>DataQ</title><script src=\"https://unpkg.com/htmx.org@1.9.11\" integrity=\"sha384-0gxUXCCR8yv9FM2b+U3FDbsKthCI66oH5IA9fHppQq9DDMHuMauqq1ZHBpJxQ0J0\" crossorigin=\"anonymous\"></script><script src=\"https://cdn.jsdelivr.net/gh/gnat/surreal@main/surreal.js\"></script><script type=\"module\" id=\"json-browser-loader\">\n\t\t\t\timport { renderJson } from 'https://esm.sh/jsr/@quinn/json-browser@0.1.3'\n\t\t\t\trenderJson.setMaxStringLength(80)\n\t\t\t\trenderJson.setShowToLevel(9)\n\n\t\t\t\trenderJson.setReplacer((key, value) => {\n\t\t\t\t\tif (typeof value === 'string' && value.startsWith('sha224-')) {\n\t\t\t\t\t\tconst anchor = document.createElement('a');\n\t\t\t\t\t\tanchor.href = `/content/${value}`;\n\t\t\t\t\t\tanchor.textContent = value;\n\t\t\t\t\t\tanchor.className = 'underline text-blue-600 hover:text-blue-800';\n\t\t\t\t\t\treturn anchor;\n\t\t\t\t\t}\n\t\t\t\t\treturn value\n\t\t\t\t})\n\n\t\t\t\t// tailwind hack\n\t\t\t\tlet classes = \"renderjson disclosure syntax string number boolean key keyword object array\"\n\n\t\t\t\t// for surreal\n\t\t\t\twindow.renderJson = renderJson\n\t\t\t\tdocument.dispatchEvent(new Event('renderJsonReady'))\n\t\t\t</script><script src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(assets.Path("toast.js"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/layout.templ`, Line: 40, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" defer></script><link rel=\"stylesheet\" href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(assets.Path("styles.css"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/layout.templ`, Line: 41, Col: 58}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\"></head><body class=\"font-mono dark:bg-black dark:text-white min-h-full\"><div class=\"bg-slate-400 p-3 flex justify-between\"><a href=\"/\">dataq</a><nav><a href=\"/plugin/install\" class=\"underline\">install plugin</a></nav></div><div class=\"p-3\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i, crumb := range middleware.GetBreadcrumbs(ctx) {
			if i != 0 {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "&nbsp;<span class=\"text-slate-500\">></span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, " <a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 templ.SafeURL = templ.SafeURL(crumb.URL)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(crumb.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/layout.templ`, Line: 55, Col: 54}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</a>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</div><hr class=\"mx-3\"><div class=\"p-3\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</div><div id=\"toast-container\" style=\"position: fixed; top: 1rem; right: 1rem; z-index: 9999;\"></div><script>\n\t\t\t\tdocument.body.addEventListener('showToast', function(e) {\n\t\t\t\t\tlet { message, type } = e.detail;\n\t\t\t\t\tshowToast(message, type);\n\t\t\t\t});\n\n\t\t\t\tdocument.body.addEventListener('htmx:sendError', function(e) {\n\t\t\t\t\tshowToast(e.detail.error, 'error');\n\t\t\t\t});\n\t\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
