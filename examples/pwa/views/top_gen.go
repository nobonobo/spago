package views

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("header", 			
			spago.A("class", spago.S(`navbar`)),
			spago.A("style", spago.S(`box-shadow: lightgrey 2px 2px 2px; padding: 1rem;`)),
			spago.Tag("section", 				
				spago.A("class", spago.S(`navber-section`)),
				spago.Tag("a", 					
					spago.A("class", spago.S(`navbar-brand mr-2`)),
					spago.A("style", spago.S(`text-transform: uppercase; font-weight: bold;`)),
					spago.T(`brand`),
				),
			),
		),
		spago.Tag("main", 			
			spago.A("class", spago.S(`container`)),
			spago.A("style", spago.S(`padding: 1rem;`)),
			spago.Tag("p", 				
				spago.A("class", spago.S(`hoge`)),
				spago.T(`Hello World!`),
			),
		),
	)
}
