package main

import (
	"io"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	_, err := sql.Open("mysql", "root@tcp(127.0.0.1)/codetest")
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Println(err)
	}
}
