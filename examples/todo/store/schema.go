package store

import "time"

// Item ...
type Item struct {
	ID       int
	Title    string
	Created  time.Time
	Complete time.Time
}
