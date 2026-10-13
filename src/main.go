package main

import (
	"regexp"
	"syscall/js"
)

var (
	reBold   = regexp.MustCompile(`\*\*([^*]+)\*\*|__([^*]+)__`)
	reItalic = regexp.MustCompile(`\*([^*]+)\*|_([^*]+)_`)
)

func toHTML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return "Error: Missing text argument"
	}

	md := args[0].String()

	md = reBold.ReplaceAllString(md, "<b>$1$2</b>")
	md = reItalic.ReplaceAllString(md, "<i>$1$2</i>")

	return md
}

func main() {
	js.Global().Set("toHTML", js.FuncOf(toHTML))

	c := make(chan struct{})
	<-c
}
