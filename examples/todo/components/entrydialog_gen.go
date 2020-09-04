package components

import (
	"github.com/nobonobo/spago"
	"todo/store"
)

// Render ...
func (c *EntryDialog) Render() spago.HTML {
	return spago.Tag("div", 		
		spago.A("class", spago.S(`modal`)),
		spago.A("id", spago.S(`entry-dialog`)),
		spago.Tag("a", 			
			spago.A("href", spago.S(`#close`)),
			spago.A("class", spago.S(`modal-overlay`)),
			spago.A("aria-label", spago.S(`Close`)),
		),
		spago.Tag("form", 			
			spago.A("class", spago.S(`modal-container`)),
			spago.Event("submit", c.OnRegisterClick),
			spago.Tag("div", 				
				spago.A("class", spago.S(`modal-header`)),
				spago.Tag("a", 					
					spago.A("href", spago.S(`#/`)),
					spago.A("class", spago.S(`btn btn-clear float-right`)),
					spago.A("aria-label", spago.S(`Close`)),
				),
				spago.Tag("div", 					
					spago.A("class", spago.S(`modal-title h5`)),
					spago.T(`ToDo Item`),
				),
			),
			spago.Tag("div", 				
				spago.A("class", spago.S(`modal-body`)),
				spago.Tag("div", 					
					spago.A("class", spago.S(`content`)),
					
					spago.Tag("input", 						
						spago.A("type", spago.S(`text`)),
						spago.A("name", spago.S(`title`)),
						spago.A("class", spago.S(`form-input`)),
						spago.A("value", spago.S(``, spago.S(store.Entry.Title), ``)),
					),
				),
			),
			spago.Tag("div", 				
				spago.A("class", spago.S(`modal-footer`)),
				spago.Tag("div", 					
					spago.A("class", spago.S(`btn-group`)),
					spago.Tag("button", 						
						spago.A("class", spago.S(`btn btn-primary`)),
						spago.T(`Register`),
					),
					spago.Tag("a", 						
						spago.A("class", spago.S(`btn`)),
						spago.A("href", spago.S(`#/`)),
						spago.T(`Cancel`),
					),
				),
			),
		),
	)
}
