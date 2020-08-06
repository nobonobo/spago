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
				spago.ClassMap{},
				spago.T(c.Title),
			),
			spago.C(&Part{}),
		),
	)
}
