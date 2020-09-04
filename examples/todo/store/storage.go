package store

import (
	"sort"
	"time"
)

var (
	id    = 0
	items = map[int]*Item{}
	// Entry ...
	Entry Item
)

// New ...
func New(title string) *Item {
	return &Item{ID: id + 1, Title: title, Created: time.Now()}
}

// Items ...
func Items() []*Item {
	res := []*Item{}
	for _, v := range items {
		if v.Complete.IsZero() {
			res = append(res, v)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Created.After(res[j].Created)
	})
	return res
}

// Get ...
func Get(id int) *Item {
	return items[id]
}

// Set ...
func Set(item *Item) {
	items[item.ID] = item
	if id < item.ID {
		id = item.ID
	}
}

// Del ...
func Del(id int) {
	delete(items, id)
}

func init() {
	Set(New("サンプル１"))
	Set(New("サンプル２"))
	Set(New("サンプル３"))
	Set(New("サンプル４"))
	Set(New("サンプル５"))
}
