package jsutil

import (
	"syscall/js"
)

var (
	global = js.Global()
	array  = global.Get("Array")
	object = global.Get("Object")
)

// JS2Bytes ...
func JS2Bytes(dv js.Value) []byte {
	b := make([]byte, dv.Get("byteLength").Int())
	js.CopyBytesToGo(b, global.Get("Uint8Array").New(dv.Get("buffer")))
	return b
}

// Bytes2JS ...
func Bytes2JS(b []byte) js.Value {
	res := global.Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(res, b)
	return res
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

// Bind ...
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

// IsArray ...
func IsArray(item js.Value) bool {
	return array.Call("isArray", item).Bool()
}

// JS2Go ...
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
		return obj
	}
}

// Form2Go ...
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

// Fetch ...
func Fetch(url string, opt map[string]interface{}) (js.Value, error) {
	return Await(global.Call("fetch", url, opt))
}
