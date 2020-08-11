package router

import (
	"fmt"
	"net/url"
	"path"
	"syscall/js"

	"github.com/nobonobo/spago"
)

var (
	document = js.Global().Get("document")
)

func parseHash(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if len(u.Fragment) == 0 {
		u.Fragment = "/"
	}
	u, err = url.Parse(u.Fragment)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetURL ...
func GetURL() *url.URL {
	u, _ := parseHash(js.Global().Get("location").Get("href").String())
	return u
}

// Router ...
type Router struct {
	current *url.URL
	f       map[string]func(key string)
	d       map[string]func(key string)
}

// Current ...
func (r *Router) Current() *url.URL {
	return r.current
}

// Navigate ...
func (r *Router) Navigate(s string) error {
	newURL, err := url.Parse(s)
	if err != nil {
		return err
	}
	key := newURL.Path
	d, f := path.Split(key)
	if len(f) > 0 {
		if fn, ok := r.f[key]; ok {
			//log.Printf("navigate to: %s", key)
			fn(key)
			return nil
		}
	}
	if fn, ok := r.d[d]; ok {
		//log.Printf("navigate to: %s", d)
		fn(key)
		return nil
	}
	return fmt.Errorf("navigate failed: unknown key %q", s)
}

func (r *Router) onHashChange(this js.Value, args []js.Value) interface{} {
	oldURL, _ := parseHash(args[0].Get("oldURL").String())
	newURL, _ := parseHash(args[0].Get("newURL").String())
	if oldURL.String() != newURL.String() {
		if err := r.Navigate(newURL.String()); err != nil {
			println(err)
			spago.RenderBody(NotFoundPage())
		}
	}
	return nil
}

// Handle ...
func (r *Router) Handle(key string, fn func(key string)) {
	d, f := path.Split(key)
	if len(f) > 0 {
		r.f[key] = fn
		return
	}
	r.d[d] = fn
}

// Start ...
func (r *Router) Start() error {
	return r.Navigate(r.current.String())
}

// New ...
func New() *Router {
	r := &Router{f: map[string]func(string){}, d: map[string]func(string){}}
	r.current = GetURL()
	js.Global().Call("addEventListener", "hashchange", js.FuncOf(r.onHashChange))
	return r
}

type defaultNotFoundPage struct {
	spago.Core
	key string
}

func (c *defaultNotFoundPage) Render() spago.HTML {
	return spago.Tag("body", spago.Tag("h1", spago.T("Not Found: "+c.key)))
}

// NotFoundPage ...
func NotFoundPage() spago.Component {
	return &defaultNotFoundPage{key: GetURL().String()}
}

// Navigate ...
func Navigate(s string) {
	js.Global().Get("location").Set("href", "#"+s)
}
