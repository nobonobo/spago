package main

import (
	"log"

	"github.com/nobonobo/spago/jsutil"
)

func main() {
	resp, err := jsutil.Fetch("https://jsonplaceholder.typicode.com/users/1", nil)
	if err != nil {
		log.Print(err)
		return
	}
	json, err := jsutil.Await(resp.Call("json"))
	if err != nil {
		log.Print(err)
		return
	}
	obj := jsutil.JS2Go(json)
	for k, v := range obj.(map[string]interface{}) {
		log.Printf("key: %s, val:%v", k, v)
	}
}
