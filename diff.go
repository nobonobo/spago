package spago

import (
	"log"
	"syscall/js"
)

func getNodeIndex(node js.Value) int {
	var i int
	for i = 0; !node.IsNull(); i++ {
		node = node.Get("previousSibling")
	}
	return i
}

func getChildByIndex(parent js.Value, index int) js.Value {
	return parent.Get("childNodes").Index(index)
}

func getAttributes(node js.Value) map[string]js.Value {
	attributes := map[string]js.Value{}
	attrs := node.Get("attributes")
	if !attrs.IsUndefined() {
		for i := 0; i < attrs.Get("length").Int(); i++ {
			a := attrs.Index(i)
			attributes[a.Get("name").String()] = a.Get("value")
		}
	}
	return attributes
}

func getChildNodes(node js.Value) []js.Value {
	childNodes := node.Get("childNodes")
	children := []js.Value{}
	for i := 0; i < childNodes.Length(); i++ {
		children = append(children, childNodes.Index(i))
	}
	return children
}

func without(s []js.Value, value js.Value) []js.Value {
	result := make([]js.Value, 0, len(s))
	for _, v := range s {
		if v.Equal(value) {
			continue
		}
		result = append(result, v)
	}
	return result
}

func remove(args ...js.Value) {
	if VerboseMode {
		console.Call("log", "remove node:", args[0])
	}
	parent := args[0].Get("parentNode")
	if !parent.IsNull() {
		parent.Call("removeChild", args[0])
	}
}

func replace(args ...js.Value) {
	from, to := args[0], args[1]
	if VerboseMode {
		console.Call("log", "replace node:", from, to)
	}
	parent := from.Get("parentNode")
	if parent.IsNull() {
		parent.Call("replaceChild", to, from)
	}
}

func insert(args ...js.Value) {
	parent, node, index := args[0], args[1], args[2]
	if VerboseMode {
		console.Call("log", "insert node:", parent, node, index)
	}
	if !parent.IsNull() {
		if !index.Equal(js.ValueOf(-1)) {
			parent.Call("removeChild", node)
			parent.Call("insertBefore", node, getChildByIndex(parent, index.Int()))
			return
		}
		parent.Call("appendChild", node)
	}
}

func changeAttribute(args ...js.Value) {
	node, attr, value := args[0], args[1], args[2]
	if VerboseMode {
		console.Call("log", "change attr:", node, attr, value)
	}
	node.Call("setAttribute", attr, value)
}

func removeAttribute(args ...js.Value) {
	node, attr := args[0], args[1]
	if VerboseMode {
		console.Call("log", "remove attr:", node, attr)
	}
	node.Call("removeAttribute", attr)
}

func changeProperty(args ...js.Value) {
	node, prop, value := args[0], args[1], args[2]
	if VerboseMode {
		console.Call("log", "change prop:", node, prop, value)
	}
	node.Set(prop.String(), value)
}

// PROPS ...
var PROPS = []string{"value", "selected", "checked", "data"}

type patch struct {
	Func func(args ...js.Value)
	Args []js.Value
}

// Patches ...
type Patches []patch

// Do ...
func (p Patches) Do() {
	for _, patch := range p {
		patch.Func(patch.Args...)
	}
}

func diffAttributes(a, b js.Value) Patches {
	patches := Patches{}
	attrsA, attrsB := getAttributes(a), getAttributes(b)
	for k, va := range attrsA {
		vb, ok := attrsB[k]
		if ok {
			if !vb.Equal(va) {
				patches = append(patches, patch{changeAttribute,
					[]js.Value{a, js.ValueOf(k), vb},
				})
			}
			delete(attrsB, k)
		} else {
			patches = append(patches, patch{removeAttribute,
				[]js.Value{a, js.ValueOf(k)},
			})
		}
	}
	for k, vb := range attrsB {
		patches = append(patches, patch{changeAttribute, []js.Value{a, js.ValueOf(k), vb}})
	}
	return patches
}

func diffProperties(a, b js.Value) Patches {
	patches := Patches{}
	for _, p := range PROPS {
		propA := a.Get(p)
		propB := b.Get(p)
		if !propA.Equal(propB) {
			patches = append(patches, patch{changeProperty, []js.Value{a, js.ValueOf(p), propB}})
		}
	}
	return patches
}

func diffChildren(a, b js.Value) Patches {
	childNodesA := getChildNodes(a)
	childNodesB := getChildNodes(b)
	if VerboseMode && len(childNodesA) != len(childNodesB) {
		log.Println("nodes-A:", childNodesA)
		log.Println("nodes-B:", childNodesB)
	}
	remainA := childNodesA[:]
	remainB := childNodesB[:]
	patches := Patches{}
	remain := append([]js.Value{}, remainB...)
	for ia := range remainA {
		if ia >= len(remainB) {
			patches = append(patches, patch{remove, []js.Value{remainA[ia]}})
			continue
		}
		cb := remainB[ia]
		patches = append(patches, Diff(remainA[ia], cb)...)
		remain = without(remain, cb)
	}
	for i := range remain {
		patches = append(patches, patch{insert, []js.Value{a, remain[i], js.ValueOf(-1)}})
	}
	return patches
}

// Diff ...
func Diff(a, b js.Value) Patches {
	if a.Call("isEqualNode", b).Bool() {
		return append(diffProperties(a, b), diffChildren(a, b)...)
	}
	if !a.Get("nodeType").Equal(b.Get("nodeType")) || !a.Get("tagName").Equal(b.Get("tagName")) {
		return Patches{{replace, []js.Value{a, b}}}
	}
	patches := Patches{}
	patches = append(patches, diffAttributes(a, b)...)
	patches = append(patches, diffProperties(a, b)...)
	patches = append(patches, diffChildren(a, b)...)
	return patches
}
