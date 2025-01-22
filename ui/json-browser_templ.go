// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.819
package ui

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/google/uuid"
)

func JsonBrowser(data any) templ.Component {
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

		var encoded bool
		if _, ok := data.([]byte); ok {
			encoded = true
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div><div role=\"data-container\" data-is-encoded=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%t", encoded))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `ui/json-browser.templ`, Line: 16, Col: 73}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if encoded {
			templ_7745c5c3_Err = templ.JSONScript(uuid.NewString(), string(data.([]byte))).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templ.JSONScript(uuid.NewString(), data).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div role=\"json-viewer\" class=\"p-1 rounded bg-slate-100 dark:bg-slate-950 border border-slate-200 dark:border-slate-800\"></div><script>\n\t\t\t{\n\t\t\t\tconst jsonViewer = me('[role=\"json-viewer\"]', me())\n\t\t\t\tconst dataContainer = me('[role=\"data-container\"]', me())\n\t\t\t\tconst script = dataContainer.querySelector('script')\n\t\t\t\tconst jsonData = dataContainer.dataset.isEncoded === 'false'\n\t\t\t\t\t? script.textContent\n\t\t\t\t\t: JSON.parse(script.textContent)\n\n\t\t\t\tconst data = JSON.parse(jsonData)\n\n\t\t\t\tif (!window.renderJson) {\n\t\t\t\t\tdocument.addEventListener('renderJsonReady', () => {\n\t\t\t\t\t\tjsonViewer.appendChild(window.renderJson(data))\n\t\t\t\t\t})\n\t\t\t\t} else {\n\t\t\t\t\tjsonViewer.appendChild(window.renderJson(data))\n\t\t\t\t}\n\t\t\t\t// new JsonViewer({ \n\t\t\t\t// \tvalue: data, \n\t\t\t\t// \ttheme: 'auto',\n\t\t\t\t// \tvalueTypes: [\n\t\t\t\t// \t\tJsonViewer.Utils.defineDataType({\n\t\t\t\t// \t\t\tis: function(value, path) {\n\t\t\t\t// \t\t\t\t// check for string\n\t\t\t\t// \t\t\t\tif (typeof value === 'string') {\n\t\t\t\t// \t\t\t\t\tif (value.startsWith('sha224-')) {\n\t\t\t\t// \t\t\t\t\t\treturn true;\n\t\t\t\t// \t\t\t\t\t}\n\t\t\t\t// \t\t\t\t}\n\t\t\t\t// \t\t\t}\n\t\t\t\t// \t\t\t// Component:\n\t\t\t\t// \t\t})\n\t\t\t\t// \t]\n\t\t\t\t// }).render(jsonViewer)\n\t\t\t}\n\t\t</script></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
