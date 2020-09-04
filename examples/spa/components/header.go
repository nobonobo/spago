package components

import (
	"log"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/router"
)

// Header ...
type Header struct {
	spago.Core
}

// Render ...
func (c *Header) Render() spago.HTML {
	return spago.Tag("header",
		spago.A("class", "navbar"),
		spago.A("style", "box-shadow: lightgrey 2px 2px 2px; padding: 1rem;"),
		spago.Tag("section",
			spago.A("class", "navbar-section"),
			spago.Tag("a",
				spago.A("class", "navbar-brand mr-2"),
				spago.A("style", "text-transform: uppercase; font-weight: bold;"),
				spago.T("BRAND"),
			),
			spago.Tag("a",
				spago.A("class", "btn btn-link"),
				spago.ClassMap{"disabled": router.GetURL().String() == "/"},
				spago.A("href", "#/"),
				spago.T("Top"),
			),
			spago.Tag("a",
				spago.A("class", "btn btn-link"),
				spago.ClassMap{"disabled": router.GetURL().String() == "/login"},
				spago.A("href", "#/login"),
				spago.T("Login"),
			),
		),
	)
}

// Mount ...
func (c *Header) Mount() {
	log.Printf("mount: %T", c)
}

// Unmount ...
func (c *Header) Unmount() {
	log.Printf("unmount: %T", c)
}
