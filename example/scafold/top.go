package main

import (
	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p main top.html

type Part struct {
	spago.Core
}

func (p *Part) Render() spago.HTML {
	return spago.Tag("button", spago.T("button"))
}

// Top  ...
type Top struct {
	spago.Core
	Title string
}
