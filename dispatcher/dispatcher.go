package dispatcher

// Action enum-type for action.
type Action int

// Callback action callback type.
type Callback func(args ...interface{})

var registered = map[Action][]Callback{}

// Register action with callback.
func Register(a Action, callback Callback) {
	registered[a] = append(registered[a], callback)
}

// Dispatch function call for action.
func Dispatch(a Action, args ...interface{}) {
	for _, v := range registered[a] {
		v(args...)
	}
}
