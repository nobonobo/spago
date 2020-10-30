package spago

import (
	"fmt"
	"strings"
	"syscall/js"
)

// SetTitle sets the title of the document.
func SetTitle(title string) {
	document.Set("title", title)
}

// AddMeta add meta tag
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

// LoadScript synchronous javascript loader
func LoadScript(url string) {
	ch := make(chan bool)
	script := document.Call("createElement", "script")
	script.Set("src", url)
	var fn js.Func
	fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer fn.Release()
		close(ch)
		return nil
	})
	script.Call("addeventListener", "load", fn)
	document.Get("head").Call("appendChild", script)
	<-ch
}

// LoadModule equivalent `import {'name1', 'name2', ...} from 'url'`
func LoadModule(names []string, url string) []js.Value {
	ch := make(chan js.Value, len(names))
	var sendFunc js.Func
	sendFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	})
	var closeFunc js.Func
	closeFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer sendFunc.Release()
		defer closeFunc.Release()
		close(ch)
		return nil
	})
	js.Global().Set("__spago_send__", sendFunc)
	js.Global().Set("__spago_close__", closeFunc)
	lines := []string{}
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("__spago_send__(%s);", name))
	}
	lines = append(lines, "__spago_close__();")
	script := Tag("script",
		A("type", "module"),
		T(fmt.Sprintf("import { %s } from %q;\n%s",
			strings.Join(names, ", "),
			url,
			strings.Join(lines, "\n"),
		)),
	).html(true)
	document.Get("head").Call("appendChild", script)
	res := make([]js.Value, 0, len(names))
	for v := range ch {
		res = append(res, v)
	}
	return res
}

// LoadModuleAs equivalent `import * as 'name' from 'url'`
func LoadModuleAs(name string, url string) js.Value {
	ch := make(chan js.Value, 1)
	var sendFunc js.Func
	sendFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	})
	var closeFunc js.Func
	closeFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer sendFunc.Release()
		defer closeFunc.Release()
		close(ch)
		return nil
	})
	js.Global().Set("__spago_send__", sendFunc)
	js.Global().Set("__spago_close__", closeFunc)
	lines := []string{}
	lines = append(lines, fmt.Sprintf("__spago_send__(%s);", name))
	lines = append(lines, "__spago_close__();")
	script := Tag("script",
		A("type", "module"),
		T(fmt.Sprintf("import * as %s from %q;\n%s",
			name,
			url,
			strings.Join(lines, "\n"),
		)),
	).html(true)
	document.Get("head").Call("appendChild", script)
	return <-ch
}

func expandNodes(v js.Value) []js.Value {
	if v.Get("constructor").Get("name").String() == "NodeList" {
		nodes := []js.Value{}
		for i := 0; i < v.Length(); i++ {
			nodes = append(nodes, v.Index(i))
		}
		return nodes
	}
	return []js.Value{v}
}
