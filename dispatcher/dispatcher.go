package dispatcher

// Action ...
type Action int

// Callback ...
type Callback func(args ...interface{})

var registered = map[Action][]Callback{}

// Register ...
func Register(a Action, callback Callback) {
	registered[a] = append(registered[a], callback)
}

// Dispatch ...
func Dispatch(a Action, args ...interface{}) {
	for _, v := range registered[a] {
		v(args...)
	}
}
