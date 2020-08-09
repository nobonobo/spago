package spago

import (
	"fmt"
	"strings"
	"syscall/js"
)

type wrapper func()

func (fn wrapper) Release() {
	fn()
}

// Bind ...
func Bind(node js.Value, name string, callback func(res js.Value)) Releaser {
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(args[0])
		return nil
	})
	node.Call("addEventListener", name, fn)
	return wrapper(func() {
		node.Call("removeEventListener", name, fn)
		fn.Release()
	})
}

// SetTitle sets the title of the document.
func SetTitle(title string) {
	document.Set("title", title)
}

// AddMeta ...
func AddMeta(name, content string) {
	meta := document.Call("createElement", "meta")
	meta.Set("name", name)
	meta.Set("content", content)
	document.Get("head").Call("appendChild", meta)
}

// AddStylesheet adds an external stylesheet to the document.
func AddStylesheet(url string) {
	link := document.Call("createElement", "link")
	link.Set("rel", "stylesheet")
	link.Set("href", url)
	document.Get("head").Call("appendChild", link)
}

// LoadScript ...
func LoadScript(url string) {
	ch := make(chan bool)
	script := Tag("script",
		A("src", url),
		Event("load", func(js.Value) interface{} {
			println("load!")
			close(ch)
			return nil
		}),
	).html(true)
	document.Get("head").Call("appendChild", script)
	<-ch
}

// LoadModule ...
func LoadModule(names []string, url string) <-chan js.Value {
	ch := make(chan js.Value)
	js.Global().Set("__wecty_send__", Callback1(func(obj js.Value) interface{} {
		ch <- obj
		return nil
	}))
	js.Global().Set("__wecty_close__", Callback0(func() interface{} {
		close(ch)
		return nil
	}))
	lines := []string{}
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("__wecty_send__(%s);", name))
	}
	lines = append(lines, "__wecty_close__();")
	script := Tag("script",
		A("type", "module"),
		T(fmt.Sprintf("import { %s } from %q;\n%s",
			strings.Join(names, ", "),
			url,
			strings.Join(lines, "\n"),
		)),
	).html(true)
	document.Get("head").Call("appendChild", script)
	return ch
}

// RequestAnimationFrame ...
func RequestAnimationFrame(callback func(float64)) int {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cb.Release()
		if callback != nil {
			callback(args[0].Float())
		}
		return js.Undefined()
	})
	return global.Call("requestAnimationFrame", cb).Int()
}

// Callback0 ...
func Callback0(fn func() interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn()
	})
	return cb
}

// Callback1 ...
func Callback1(fn func(res js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args[0])
	})
	return cb
}

// CallbackN ...
func CallbackN(fn func(res []js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args)
	})
	return cb
}

// S ...
func S(s ...interface{}) string {
	return fmt.Sprint(s...)
}

type wrappedError js.Value

func (w wrappedError) Error() string {
	return js.Value(w).Call("toString").String()
}

func (w wrappedError) JSValue() js.Value {
	return js.Value(w)
}

// Await ...
func Await(promise js.Value) (res js.Value, err error) {
	ch := make(chan bool)
	promise.Call("then",
		Callback1(func(r js.Value) interface{} {
			res = r
			close(ch)
			return nil
		}),
	).Call("catch",
		Callback1(func(res js.Value) interface{} {
			err = wrappedError(res)
			close(ch)
			return nil
		}),
	)
	<-ch
	return
}
