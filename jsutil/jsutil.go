package jsutil

import "syscall/js"

var (
	global = js.Global()
)

// JS2Bytes ...
func JS2Bytes(dv js.Value) []byte {
	b := make([]byte, dv.Get("byteLength").Int())
	js.CopyBytesToGo(b, js.Global().Get("Uint8Array").New(dv.Get("buffer")))
	return b
}

// Bytes2JS ...
func Bytes2JS(b []byte) js.Value {
	res := js.Global().Get("Uint8Array").New(len(b))
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

// Releaser ...
type Releaser interface {
	Release()
}

// ReleaserFunc ...
type ReleaserFunc func()

// Release ...
func (fn ReleaserFunc) Release() {
	fn()
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
