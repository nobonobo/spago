package main

import (
	"log"

	"github.com/nobonobo/spago"
)

func main() {
	log.SetFlags(log.Lshortfile)
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
	spago.RenderBody(&Top{})
	select {}
}
