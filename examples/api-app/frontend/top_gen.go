package main

import (
	"github.com/nobonobo/spago"
)

// Render ...
func (c *Top) Render() spago.HTML {
	return spago.Tag("body", 
		spago.Tag("button", 			
			spago.A("class", spago.S(`btn`)),
			spago.Event("click", c.OnClickButton),
			spago.T(`getNow()`),
		),
		spago.Tag("label", 
			spago.T(``, spago.S(c.Now), ``),
		),
	)
}
