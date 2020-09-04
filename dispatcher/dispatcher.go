package dispatcher

type Action int

var registered = map[Action][]func(){}

// Register ...
func Register(a Action, callback func()) {
	registered[a] = append(registered[a], callback)
}

// Dispatch ...
func Dispatch(a Action) {
	for _, v := range registered[a] {
		v()
	}
}
