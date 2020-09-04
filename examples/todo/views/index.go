package views

import (
	"github.com/nobonobo/spago"

	"todo/store"
)

//go:generate spago generate -c Index -p views index.html

// Index  ...
type Index struct {
	spago.Core
	Entry store.Item
}
