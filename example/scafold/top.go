package main

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p main top.html

// Top  ...
type Top struct {
	spago.Core
}
