package main

import (
	"syscall/js"
	"writedown.internals/src/markdown"
	"writedown.internals/src/dobbm-markup"
)

func main() {
	js.Global().Set("fromMD", js.FuncOf(markdown.ToHTML))
	js.Global().Set("fromDML", js.FuncOf(dobbmmarkup.ToHTML))

	c := make(chan struct{})
	<-c
}
