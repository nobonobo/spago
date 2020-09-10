package main

import "github.com/nobonobo/spago"

func main() {
	spago.RenderBody(&Top{})
	select {}
}
