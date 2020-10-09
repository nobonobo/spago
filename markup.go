package spago

import (
	"fmt"
	"syscall/js"
)

type none struct{}

func (m none) apply(n *Node) {}

// If ...
func If(cond bool, m Markup) Markup {
	if cond {
		return m
	}
	return none{}
}

// attribute attribute
type attribute struct {
	Key   string
	Value interface{}
}

func (a attribute) apply(n *Node) {
	n.attributes = append(n.attributes, a)
}

// A ...
func A(k string, v interface{}) Markup {
	return attribute{Key: k, Value: v}
}

// AttrMap ...
type AttrMap map[string]interface{}

func (a AttrMap) apply(n *Node) {
	for k, v := range a {
		attribute{k, v}.apply(n)
	}
}

// listener ...
type listener struct {
	Name string
	Func func(ev js.Value)
}

func (e listener) apply(n *Node) {
	n.listeners = append(n.listeners, e)
}

// Event ...
func Event(name string, fn func(ev js.Value)) Markup {
	return &listener{name, fn}
}

// Markups ...
type Markups []Markup

func (c Markups) apply(n *Node) {
	for _, v := range c {
		v.apply(n)
	}
}

type components []Component

func (c components) apply(n *Node) {
	for _, v := range c {
		n.children = append(n.children, v)
	}
}

// C ....
func C(c ...Component) Markup {
	return components(c)
}

// S ...
func S(s ...interface{}) string {
	return fmt.Sprint(s...)
}
