package spago

import (
	"fmt"
	"syscall/js"
)

var (
	global   = js.Global()
	document = global.Get("document")
	console  = js.Global().Get("console")
	mounts   []Mounter

	// VerboseMode ...
	VerboseMode bool
)

// Tag markup
func Tag(tag string, markups ...Markup) *Node {
	t := &Node{tag: tag, classmap: ClassMap{}}
	for _, m := range markups {
		m.apply(t)
	}
	return t
}

// TagNS markup
func TagNS(namespace, tag string, markups ...Markup) *Node {
	t := &Node{namespace: namespace, tag: tag, classmap: ClassMap{}}
	for _, m := range markups {
		m.apply(t)
	}
	return t
}

// Render ...
func render(old js.Value, c Component) {
	core := c.get()
	if v, ok := c.(Unmounter); !core.target.IsUndefined() && ok {
		v.Unmount()
	}
	core.target = c.Render().html(true)
	for _, v := range expandNodes(core.target) {
		old.Get("parentNode").Call("replaceChild", v, old)
	}
	if v := old.Get("release"); !v.IsUndefined() {
		v.Invoke()
	}
	//old.Call("remove")
}

func mount() {
	for _, v := range mounts {
		v.Mount()
	}
	mounts = nil
}

var lastComponent = map[string]Component{}

// Render ...
func Render(q string, c Component) {
	e := document.Call("querySelector", q)
	if e.IsUndefined() {
		panic(fmt.Sprintf("not found element: %q", q))
	}
	if v, ok := lastComponent[q].(Unmounter); ok {
		v.Unmount()
	}
	render(e, c)
	if v, ok := c.(Mounter); ok {
		mounts = append(mounts, v)
	}
	lastComponent[q] = c
	mount()
}

var lastBody Component

// RenderBody ...
func RenderBody(c Component) {
	if v, ok := lastBody.(Unmounter); ok {
		v.Unmount()
	}
	render(document.Get("body"), c)
	if v, ok := c.(Mounter); ok {
		mounts = append(mounts, v)
	}
	lastBody = c
	mount()
}

// Rerender ...
func Rerender(c Component) {
	core := c.get()
	old := core.target
	if old.IsNull() {
		panic(fmt.Sprint("this component rendering not yet:", c))
	}
	if v, ok := c.(Unmounter); ok {
		v.Unmount()
	}
	next := c.Render().html(true)
	var patch Patches
	for _, v := range expandNodes(next) {
		patch = append(patch, Diff(old, v)...)
		v.Call("release")
	}
	patch.Do()
	if v, ok := c.(Mounter); ok {
		mounts = append(mounts, v)
	}
	mount()
}
