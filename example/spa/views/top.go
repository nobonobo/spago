package views

import (
	"log"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/examples/spa/components"
)

// Top ...
type Top struct {
	spago.Core
}

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body",
		spago.A("style", "margin: 0;"),
		spago.C(&components.Header{}),
		spago.Tag("main", spago.A("class", "container"),
			spago.A("style", "padding: 1rem;"),
			spago.Tag("h1", spago.T("Top")),
			spago.Tag("p", spago.T("Hello World!")),
		),
	)
}

// Mount ...
func (c *Top) Mount() {
	log.Printf("mount: %T", c)
}

// Unmount ...
func (c *Top) Unmount() {
	log.Printf("unmount: %T", c)
}
