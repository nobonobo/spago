package main

import (
	"log"
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/jsutil"
)

//go:generate spago generate -c Form -p main form.html

// Form  ...
type Form struct {
	spago.Core
}

// OnSubmit ...
func (c *Form) OnSubmit(ev js.Value) {
	ev.Call("preventDefault") // cancel default behavior
	form := jsutil.Form2Go(ev.Get("target"))
	log.Printf("%#v", form)
}
