package main

import (
	"go.quinn.io/dataq/pkg/pluginutil"
)

func main() {
	plugin := New()
	pluginutil.HandlePlugin(plugin)
}
