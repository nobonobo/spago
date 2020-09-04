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
func (c *Item) OnCompleteClick(ev js.Value) interface{} {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return nil
	}
	log.Println("complete:", id)
	item := store.Get(id)
	if item != nil {
		item.Complete = time.Now()
		store.Set(item)
	}
	dispatcher.Dispatch(actions.Refresh)
	return nil
}

// OnEditClick ...
func (c *Item) OnEditClick(ev js.Value) interface{} {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return nil
	}
	log.Println("edit:", id)
	item := store.Get(id)
	if item != nil {
		store.Entry = *item
		dispatcher.Dispatch(actions.Refresh)
		router.Navigate("entry-dialog")
	}
	return nil
}

// OnDeleteClick ...
func (c *Item) OnDeleteClick(ev js.Value) interface{} {
	s := ev.Get("target").Call("closest", ".todo-item").Get("id").String()
	id, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
		return nil
	}
	log.Println("delete:", id)
	store.Del(id)
	dispatcher.Dispatch(actions.Refresh)
	return nil
}
