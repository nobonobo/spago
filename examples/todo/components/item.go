package components

import (
	"log"
	"strconv"
	"syscall/js"
	"time"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
	"github.com/nobonobo/spago/router"

	"github.com/nobonobo/spago/examples/todo/actions"
	"github.com/nobonobo/spago/examples/todo/store"
)

//go:generate spago generate -c Item -p components item.html

// Item  ...
type Item struct {
	spago.Core
	ID   int
	Data *store.Item
}

// OnCompleteClick ...
func (c *Item) OnCompleteClick(ev js.Value) {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("complete:", id)
	item := store.Get(id)
	if item != nil {
		item.Complete = time.Now()
		store.Set(item)
	}
	dispatcher.Dispatch(actions.Refresh)
}

// OnEditClick ...
func (c *Item) OnEditClick(ev js.Value) {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("edit:", id)
	item := store.Get(id)
	if item != nil {
		store.Entry = *item
		dispatcher.Dispatch(actions.Refresh)
		router.Navigate("entry-dialog")
	}
}

// OnDeleteClick ...
func (c *Item) OnDeleteClick(ev js.Value) {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return
	}
	log.Println("delete:", id)
	store.Del(id)
	dispatcher.Dispatch(actions.Refresh)
}
