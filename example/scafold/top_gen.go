package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("main", 
			spago.Tag("h1", 
				spago.T("Title"),
			),
			spago.Tag("p", 				
				spago.A("$", "a"),
				spago.A("$", "b"),
				spago.T("hello world!"),
			),
		),
	)
}
