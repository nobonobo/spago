package backend

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func now(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(io.MultiWriter(w, os.Stdout)).Encode(time.Now().Unix()); err != nil {
		log.Print(err)
		http.Error(w, "get now failed", http.StatusBadGateway)
	}
}

func init() {
	http.HandleFunc("/api/now", now)
}
