package spago

import "syscall/js"

// HTML ...
type HTML interface {
	html(bind bool) js.Value
}

// Markup ...
type Markup interface {
	apply(node *Node)
}

// Component ...
type Component interface {
	Render() HTML
	get() *Core
}

// ComponentOrHTML ...
type ComponentOrHTML interface{}

// Mounter ...
type Mounter interface {
	Mount()
}

// Unmounter ...
type Unmounter interface {
	Unmount()
}
