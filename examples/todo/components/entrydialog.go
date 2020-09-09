package components

import (
	"syscall/js"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/router"

	"github.com/nobonobo/spago/examples/todo/actions"
	"github.com/nobonobo/spago/examples/todo/store"
)

//go:generate spago generate -c EntryDialog -p components entrydialog.html

// EntryDialog  ...
type EntryDialog struct {
	spago.Core
}

// OnRegisterClick ...
func (c *EntryDialog) OnRegisterClick(ev js.Value) {
	ev.Call("preventDefault")
	title := ev.Get("target").Get("title").Get("value").String()
	js.Global().Get("console").Call("log", title)
	store.Entry.Title = title
	dup := store.Entry
	store.Set(&dup)
	router.Navigate("")
	dispatcher.Dispatch(actions.Refresh)
}
