package main

import (
	"go.quinn.io/dataq/plugin"
)

func main() {
	p := New()
	plugin.Run(p)
}
