package spago

import "syscall/js"

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
	Func func(ev js.Value) interface{}
}

func (e listener) apply(n *Node) {
	n.listeners = append(n.listeners, e)
}

// Event ...
func Event(name string, fn func(ev js.Value) interface{}) Markup {
	return &listener{name, fn}
}

type component struct {
	Component
}

func (c *component) apply(n *Node) {
	n.children = append(n.children, c)
}

// C ....
func C(c Component) Markup {
	return &component{Component: c}
}
