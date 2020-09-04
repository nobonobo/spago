package views

import (
	"github.com/nobonobo/spago"
	"todo/components"
)

// Render ...
func (c *Index) Render() spago.HTML {
	return spago.Tag("body", 
		spago.C(&components.Header{}),
		spago.Tag("main", 			
			spago.A("class", spago.S(`container`)),
			spago.A("style", spago.S(`padding: 1rem`)),
			spago.C(&components.ItemList{}),
		),
		spago.C(&components.EntryDialog{}),
	)
}
