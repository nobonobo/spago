package components

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Item) Render() spago.HTML {
	return spago.Tag("div", 		
		spago.A("class", spago.S(`tile todo-item`)),
		spago.A("id", spago.S(``, spago.S(c.Data.ID), ``)),
		spago.Tag("div", 			
			spago.A("class", spago.S(`tile-icon`)),
			spago.Tag("button", 				
				spago.A("class", spago.S(`btn s-circle`)),
				spago.A("style", spago.S(`width: 1.8rem`)),
				spago.Event("click", c.OnCompleteClick),
				spago.Tag("i", 					
					spago.A("class", spago.S(`icon icon-check`)),
				),
			),
		),
		spago.Tag("div", 			
			spago.A("class", spago.S(`tile-content`)),
			spago.Tag("div", 				
				spago.A("class", spago.S(`tile-title`)),
				spago.T(``, spago.S(c.Data.Title), ``),
			),
			spago.Tag("small", 				
				spago.A("class", spago.S(`tile-subtitle text-gray`)),
				spago.T(``, spago.S(c.Data.Created), ``),
			),
		),
		spago.Tag("div", 			
			spago.A("class", spago.S(`tile-action`)),
			spago.Tag("button", 				
				spago.A("class", spago.S(`btn btn-link`)),
				spago.Event("click", c.OnEditClick),
				spago.Tag("i", 					
					spago.A("class", spago.S(`icon icon-edit`)),
				),
			),
			spago.Tag("button", 				
				spago.A("class", spago.S(`btn btn-link`)),
				spago.Event("click", c.OnDeleteClick),
				spago.Tag("i", 					
					spago.A("class", spago.S(`icon icon-cross`)),
				),
			),
		),
	)
}
