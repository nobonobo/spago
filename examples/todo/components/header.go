package components

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/router"

	"todo/actions"
	"todo/store"
)

//go:generate spago generate -c Header -p components header.html

// Header  ...
type Header struct {
	spago.Core
}

// OnNewClick ...
func (c *Header) OnNewClick(ev js.Value) interface{} {
	store.Entry = *store.New("タイトル")
	dispatcher.Dispatch(actions.Refresh)
	router.Navigate("entry-dialog")
	return nil
}
