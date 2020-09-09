package main

import "github.com/nobonobo/spago"

func main() {
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	spago.RenderBody(&Form{})
	select {}
}
