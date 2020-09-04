package main

import (
	"log"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/router"

	"github.com/nobonobo/spago/examples/pwa/views"
)

func main() {
	log.SetFlags(log.Lshortfile)
	r := router.New()
	r.Handle("/", func(key string) {
		log.Println(router.GetURL())
		spago.SetTitle("Top")
		spago.RenderBody(&views.Top{})
	})
	if err := r.Start(); err != nil {
		println(err)
		spago.RenderBody(router.NotFoundPage())
	}
	select {}
}
