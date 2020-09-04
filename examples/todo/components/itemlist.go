package components

import (
	"log"

	"github.com/nobonobo/spago/examples/todo/store"

	"github.com/nobonobo/spago"
)

// ItemList  ...
type ItemList struct {
	spago.Core
}

// Items ...
func (c *ItemList) Items() []spago.Markup {
	items := []*Item{}
	for _, v := range store.Items() {
		if v.Complete.IsZero() {
			items = append(items, &Item{Data: v})
		}
	}
	log.Println(items)
	res := make([]spago.Markup, 0, len(items))
	for _, v := range items {
		res = append(res, spago.C(v))
	}
	return res
}

// Render ...
func (c *ItemList) Render() spago.HTML {
	return spago.Tag("div",
		spago.A("class", spago.S(`card p-centered`)),
		spago.A("style", spago.S(`max-width: 640px`)),
		spago.Tag("div", append(
			[]spago.Markup{spago.A("class", spago.S(`card-body`))},
			c.Items()...,
		)...),
	)
}
