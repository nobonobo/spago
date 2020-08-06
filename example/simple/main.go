package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/nobonobo/spago"
)

var (
	window   = js.Global()
	document = window.Get("document")
)

// Sub ...
type Sub struct {
	spago.Core
}

// Render ...
func (c *Sub) Render() spago.HTML {
	return spago.Tag("span", spago.T("hello"))
}

// Top ....
type Top struct {
	spago.Core
	index int
}

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body",
		spago.Tag("header",
			spago.A("class", "navbar"),
			spago.A("style", "box-shadow: lightgrey 2px 2px 2px; padding: 1rem;"),
			spago.Tag("section",
				spago.A("class", "navbar-section"),
				spago.Tag("a",
					spago.A("class", "navbar-brand"),
					spago.A("style", "text-transform: uppercase; font-weight: bold;"),
					spago.T(fmt.Sprintf("BRAND: %d", c.index)),
				),
			),
		),
		spago.Tag("main",
			spago.A("class", "container"),
			spago.A("style", "padding: 1rem;"),
			spago.T("body: hello world"),
			spago.Tag("div",
				spago.C(&Sub{}),
			),
			spago.Tag("button",
				spago.ClassMap{"btn": true},
				spago.Event("click", c.Update),
				spago.T(fmt.Sprint(c.index)),
			),
			spago.Tag("div",
				spago.T(fmt.Sprint(c.index)),
			),
		),
	)
}

// Update ...
func (c *Top) Update(event js.Value) interface{} {
	log.Println("button click!")
	c.index++
	spago.Rerender(c)
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
	spago.Render(document.Get("body"), &Top{})
	select {}
}
