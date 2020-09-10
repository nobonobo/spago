package main

import (
	"log"

	"github.com/nobonobo/spago"
)

/*
spago server -p /api/=http://localhost:8000/api/
*/

func main() {
	log.SetFlags(log.Lshortfile)
	spago.RenderBody(&Top{})
	select {}
}
