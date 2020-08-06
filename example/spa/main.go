package main

import (
	"log"
	"runtime"
	"time"

	"github.com/nobonobo/spago"

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
	router := spago.NewRouter()
	router.Handle("/login", func(key string) {
		log.Println(spago.GetURL())
		spago.SetTitle("Login")
		spago.RenderBody(&views.Login{})
	})
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
