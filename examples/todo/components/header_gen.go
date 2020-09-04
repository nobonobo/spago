package components

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Header) Render() spago.HTML {
	return spago.Tag("header", 		
		spago.A("class", spago.S(`navbar`)),
		spago.Tag("section", 			
			spago.A("class", spago.S(`navber-section`)),
			spago.Tag("a", 				
				spago.A("class", spago.S(`navbar-brand mr-2`)),
				spago.T(`ToDo`),
			),
		),
		spago.Tag("section", 			
			spago.A("class", spago.S(`navber-section`)),
			spago.Tag("button", 				
				spago.Event("click", c.OnNewClick),
				spago.A("class", spago.S(`btn btn-primary btn-sm s-circle`)),
				spago.Tag("i", 					
					spago.A("class", spago.S(`icon icon-plus`)),
				),
			),
		),
	)
}
