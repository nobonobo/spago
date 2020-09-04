package views

import (
	"log"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/examples/spa/components"
)

// Login ...
type Login struct {
	spago.Core
	UserName string
	Password string
}

// Render ...
func (c *Login) Render() spago.HTML {
	return spago.Tag("body",
		spago.A("style", "margin: 0;"),
		spago.C(&components.Header{}),
		spago.Tag("main", spago.A("class", "container"),
			spago.A("style", "padding: 1rem;"),
			spago.Tag("form",
				spago.Event("submit", c.OnSubmit),
				spago.A("class", "card col-4 col-mx-auto"),
				spago.Tag("div",
					spago.A("class", "card-header"),
					spago.Tag("div", spago.A("class", "card-title h5"), spago.T("Login")),
				),
				spago.Tag("div",
					spago.A("class", "card-body"),
					spago.Tag("div", spago.A("class", "form-group"),
						spago.Tag("label",
							spago.A("class", "form-label"),
							spago.T("ID:"),
						),
						spago.Tag("input",
							spago.A("class", "form-input"),
							spago.A("type", "text"),
							spago.A("name", "username"),
							spago.A("autocomplete", "username"),
							spago.A("value", c.UserName),
						),
					),
					spago.Tag("div", spago.A("class", "form-group"),
						spago.Tag("label",
							spago.A("class", "form-label"),
							spago.T("Password:"),
						),
						spago.Tag("input",
							spago.A("class", "form-input"),
							spago.A("type", "password"),
							spago.A("name", "password"),
							spago.A("autocomplete", "current-password"),
							spago.A("value", c.Password),
						),
					),
				),
				spago.Tag("div",
					spago.A("class", "card-footer"),
					spago.Tag("button",
						spago.A("class", "btn btn-primary"),
						spago.T("login"),
					),
				),
			),
		),
	)
}

// OnSubmit ...
func (c *Login) OnSubmit(event js.Value) interface{} {
	event.Call("preventDefault")
	log.Println("submit: login credential")
	c.UserName = event.Get("target").Get("username").Get("value").String()
	c.Password = event.Get("target").Get("password").Get("value").String()
	spago.Rerender(c)
	return nil
}

// Mount ...
func (c *Login) Mount() {
	log.Printf("mount: %T", c)
}

// Unmount ...
func (c *Login) Unmount() {
	log.Printf("unmount: %T", c)
}
