package main

import (
	"log"
	"runtime"
	"time"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/router"

	"github.com/nobonobo/spago/examples/spa/views"
)

func main() {
	log.SetFlags(log.Lshortfile)
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
	spago.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
	go func() {
		time.Sleep(1 * time.Second)
		runtime.GC()
	}()
	r := router.New()
	r.Handle("/login", func(key string) {
		log.Println(router.GetURL())
		spago.SetTitle("Login")
		spago.RenderBody(&views.Login{})
	})
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
