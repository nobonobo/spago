package jsutil

import (
	"syscall/js"
)

var (
	global = js.Global()
	array  = global.Get("Array")
	object = global.Get("Object")
)

// JS2Bytes convert from TypedArray for JS to byte slice for Go.
func JS2Bytes(dv js.Value) []byte {
	b := make([]byte, dv.Get("byteLength").Int())
	js.CopyBytesToGo(b, global.Get("Uint8Array").New(dv.Get("buffer")))
	return b
}

// Bytes2JS convert from byte slice for Go to Uint8Array for JS.
func Bytes2JS(b []byte) js.Value {
	res := global.Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(res, b)
	return res
}

// Callback0 make auto-release callback without params.
func Callback0(fn func() interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn()
	})
	return cb
}

// Callback1 make auto-release callback with 1 param.
func Callback1(fn func(res js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args[0])
	})
	return cb
}

// CallbackN make auto-release callback with multiple params.
func CallbackN(fn func(res []js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args)
	})
	return cb
}

// RequestAnimationFrame function call for 30 or 60 fps.
// return value: cancel function
func RequestAnimationFrame(callback func(dt float64)) func() {
	var cb js.Func
	lastID := -1
	lastTick := 0
	terminate := false
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		tick := args[0].Int()
		dt := float64(tick-lastTick) / 1000.0
		lastTick = tick
		callback(dt)
		if !terminate {
			lastID = global.Call("requestAnimationFrame", cb).Int()
		}
		return js.Undefined()
	})
	cb.Invoke(js.ValueOf(0.0))
	return func() {
		terminate = true
		global.Call("cancelAnimationFrame", lastID)
		cb.Release()
	}
}

type wrappedError js.Value

func (w wrappedError) Error() string {
	return js.Value(w).Call("toString").String()
}

func (w wrappedError) JSValue() js.Value {
	return js.Value(w)
}

// Await equivalent for js await statement.
func Await(promise js.Value) (res js.Value, err error) {
	ch := make(chan bool)
	promise.Call("then",
		Callback1(func(arg js.Value) interface{} {
			res = arg
			close(ch)
			return nil
		}),
	).Call("catch",
		Callback1(func(arg js.Value) interface{} {
			err = wrappedError(arg)
			close(ch)
			return nil
		}),
	)
	<-ch
	return
}

// Releaser ...
type Releaser interface {
	Release()
}

type wrapper struct {
	fn func()
}

// Release ...
func (w *wrapper) Release() {
	w.fn()
}

// ReleaserFunc ...
func ReleaserFunc(fn func()) Releaser {
	return &wrapper{fn: fn}
}

// Bind event bind and return releaser.
func Bind(node js.Value, name string, callback func(res js.Value)) Releaser {
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback(args[0])
		return nil
	})
	node.Call("addEventListener", name, fn)
	return ReleaserFunc(func() {
		node.Call("removeEventListener", name, fn)
		fn.Release()
	})
}

// IsArray checking value is array type.
func IsArray(item js.Value) bool {
	return array.Call("isArray", item).Bool()
}

// JS2Go JS values convert to Go values.
func JS2Go(obj js.Value) interface{} {
	switch obj.Type() {
	default:
		return obj
	case js.TypeBoolean:
		return obj.Bool()
	case js.TypeNumber:
		return obj.Float()
	case js.TypeString:
		return obj.String()
	case js.TypeObject:
		if IsArray(obj) {
			res := []interface{}{}
			for i := 0; i < obj.Length(); i++ {
				res = append(res, obj.Index(i))
			}
			return res
		}
		res := map[string]interface{}{}
		entries := object.Call("entries", obj)
		for i := 0; i < entries.Length(); i++ {
			entry := entries.Index(i)
			key, value := entry.Index(0).String(), entry.Index(1)
			res[key] = JS2Go(value)
		}
		return res
	}
}

// Form2Go retrieve form values from form element.
func Form2Go(form js.Value) map[string]interface{} {
	obj := map[string]interface{}{}
	formData := global.Get("FormData").New(form)
	fp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		value, key := JS2Go(args[0]), args[1].String()
		if _, ok := obj[key]; !ok {
			obj[key] = value
			return nil
		}
		v, ok := obj[key].([]interface{})
		if !ok {
			v = []interface{}{obj[key]}
		}
		v = append(v, value)
		obj[key] = v
		return nil
	})
	defer fp.Release()
	formData.Call("forEach", fp)
	return obj
}

// Fetch wrapper fetch function.
func Fetch(url string, opt map[string]interface{}) (js.Value, error) {
	return Await(global.Call("fetch", url, opt))
}
