package views

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p views top.html

// Top  ...
type Top struct {
	spago.Core
}
