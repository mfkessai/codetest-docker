package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/killinsun/codetest-docker/app/infra"
)

func main() {
	config := infra.NewConfig()
	conn, err := infra.GetConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🚀 Starting server...")
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if !isValidRequest(r, conn) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		if _, err := io.Copy(w, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	go func() {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("👉 Server is running on :8888")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🚪 Shutting down server...")
}

func isValidRequest(req *http.Request, conn *sql.DB) bool {
	apiKey := req.Header.Get("apikey")
	if apiKey == "" {
		return false
	}

	return isUserExist(apiKey, conn)
}

func isUserExist(apiKey string, conn *sql.DB) bool {
	userRepo := infra.NewUserRepository(conn)
	user, err := userRepo.FindByApiKey(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	if user == nil {
		return false
	}

	return true
}
