package main

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

var Now string = "unknown"

type Index struct{ spago.Core }

func (c *Index) Render() spago.HTML {
	return spago.Tag("body",
		spago.Tag("label",
			spago.T(spago.S(Now)),
		),
		spago.Tag("button",
			spago.Event("click", c.OnClick),
			spago.T("Get Now"),
		),
	)
}

func (c *Index) OnClick(ev js.Value) {
	go func() {
		resp, err := jsutil.Fetch("/api/now", nil)
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		value, err := jsutil.Await(resp.Call("json"))
		if err != nil {
			js.Global().Call("alert", err.Error())
			return
		}
		Now = js.Global().Get("JSON").Call("stringify", value).String()
		spago.Rerender(c)
	}()
}

func main() {
	spago.RenderBody(&Index{})
	select {}
}
