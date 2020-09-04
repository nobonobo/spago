package main

import (
	"log"

	"github.com/nobonobo/spago"

	"github.com/nobonobo/spago/examples/pwa/views"
)

func main() {
	log.SetFlags(log.Lshortfile)
	router := spago.NewRouter()
	router.Handle("/", func(key string) {
		log.Println(spago.GetURL())
		spago.SetTitle("Top")
		spago.RenderBody(&views.Top{})
	})
	if err := router.Start(); err != nil {
		println(err)
		spago.RenderBody(spago.NotFoundPage())
	}
	select {}
}
