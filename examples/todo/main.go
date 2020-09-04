package main

import (
	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"

	"github.com/nobonobo/spago/examples/todo/actions"
	"github.com/nobonobo/spago/examples/todo/views"
)

var top = &views.Index{}

func init() {
	dispatcher.Register(actions.Refresh, func() {
		spago.Rerender(top)
	})
}

func main() {
	spago.VerboseMode = true
	spago.AddStylesheet("assets/app.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
	spago.RenderBody(top)
	select {}
}
