package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Println(err)
	}
}
