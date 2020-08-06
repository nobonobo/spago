package spago

import "syscall/js"

// Core ...
type Core struct {
	target js.Value
}

func (c *Core) get() *Core { return c }
