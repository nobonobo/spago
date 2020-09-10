package main

import (
	"net/http"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

//go:generate spago generate -c Top -p main top.html

// Top  ...
type Top struct {
	spago.Core
	Now string
}

// OnClickButton ...
func (c *Top) OnClickButton(ev js.Value) {
	go func() {
		resp, err := jsutil.Fetch("/api/now", nil)
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		if resp.Get("status").Int() != http.StatusOK {
			js.Global().Call("alert", resp.Get("statusText"))
			return
		}
		text, err := jsutil.Await(resp.Call("text"))
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		c.Now = text.String()
		spago.Rerender(c)
	}()
}
