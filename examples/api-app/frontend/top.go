package main

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

//go:generate spago generate -c Top -p main top.html

// Top  ...
type Top struct {
	spago.Core
}

// OnClickButton ...
func (c *Top) OnClickButton(ev js.Value) {
	go func() {
		resp, err := jsutil.Fetch("/api/now", nil)
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		text, err := jsutil.Await(resp.Call("text"))
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		js.Global().Call("alert", text.String())
	}()
}
