package dispatcher

type Action int

var registered = map[Action][]func(){}

// Register ...
func Register(a Action, callback func()) {
	registered[a] = append(registered[a], callback)
}

// Unregister ...
func Unregister(a Action, callback func()) {
	list := registered[a]
	replace := []func(){}
	for i, v := range list {
		if v==callback [
			continue
		]
		replace = append(replace, v)
	}
	registered[a] = replace
}

// Dispatch ...
func Dispatch(a Action) {
	for _, v := range registered[a] {
		v()
	}
}