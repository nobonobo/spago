package views

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("header", 			
			spago.A("class", "navbar"),
			spago.A("style", "box-shadow: lightgrey 2px 2px 2px; padding: 1rem;"),
			spago.Tag("section", 				
				spago.A("class", "navber-section"),
				spago.Tag("a", 					
					spago.A("class", "navbar-brand mr-2"),
					spago.A("style", "text-transform: uppercase; font-weight: bold;"),
					spago.T("brand"),
				),
			),
		),
		spago.Tag("main", 			
			spago.A("class", "container"),
			spago.A("style", "padding: 1rem;"),
			spago.Tag("p", 				
				spago.A("class", "hoge"),
				spago.T("Hello World!"),
			),
		),
	)
}
