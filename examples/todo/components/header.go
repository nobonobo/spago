package components

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/router"

	"github.com/nobonobo/spago/examples/todo/actions"
	"github.com/nobonobo/spago/examples/todo/store"
)

//go:generate spago generate -c Header -p components header.html

// Header  ...
type Header struct {
	spago.Core
}

// OnNewClick ...
func (c *Header) OnNewClick(ev js.Value) {
	store.Entry = *store.New("タイトル")
	dispatcher.Dispatch(actions.Refresh)
	router.Navigate("entry-dialog")
}
