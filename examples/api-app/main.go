package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"go.pyspa.org/brbundle/brhttp"
)

//go:generate sh -c "cd frontend && go generate ./..."
//go:generate sh -c "cd frontend && spago deploy ../www"
//go:generate brbundle embedded -p main www

func getNow(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now())
}

func main() {
	http.Handle("/api/now", http.HandlerFunc(getNow))
	http.Handle("/", brhttp.Mount())
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	if err := http.Serve(l, nil); err != nil {
		log.Fatal(err)
	}
}
