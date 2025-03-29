package main

import (
	"database/sql"
	"log"
	"net/http"
	"fmt"
	"sync"
	"encoding/json"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Transaction struct {
	UserID      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

const amountLimit = 1000 // 登録可能な取引上限金額
var mu sync.Mutex

// すでにDBに登録されているユーザーが確認する
func isValidUser(conn *sql.DB, key string) (bool, error) {
	stmt, err := conn.Prepare("SELECT COUNT(*) FROM users WHERE api_key = ?;")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int
	if err = stmt.QueryRow(key).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

// 上限金額を超えていない確認する
func reachesAmountLimit(conn *sql.DB, t Transaction) (bool, error) {
	stmt, err := conn.Prepare("SELECT SUM(amount) FROM transactions WHERE user_id = ?;")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var sum sql.NullInt64
	if err = stmt.QueryRow(t.UserID).Scan(&sum); err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return int(sum.Int64) + t.Amount > amountLimit, nil
}

func insertTranaction(conn *sql.DB, t Transaction) error {
	stmt, err := conn.Prepare("INSERT INTO transactions(user_id, amount, description) VALUES(?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(t.UserID, t.Amount, t.Description); err != nil {
		return err
	}
	return nil
}

func handleResponse(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(""))
		return
	}
	
	if strings.Contains(err.Error(), "is not allowd") {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	} else if strings.Contains(err.Error(), "is reaching the limit") {
		w.WriteHeader(http.StatusPaymentRequired)
		w.Write([]byte(err.Error()))
	} else {
		log.Fatal(err)
	}
	return
}

func main() {
	log.Println("api start")
	conn, err := sql.Open("mysql", "root@tcp(db)/codetest")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		defer handleResponse(w, err)

		log.Println("transaction start")

		key := r.Header.Get("apikey")
		isValid, err := isValidUser(conn, key)
		if err != nil {
			return
		}
		if !isValid {
			err = fmt.Errorf("key %s is not allowed", key)
			return
		}

		// 並列で挿入するとレース競合が起きるのでロック
		mu.Lock()
		defer mu.Unlock()

		var t Transaction
		if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
			return
		}

		reachesAmountLimit, err := reachesAmountLimit(conn, t)
		if err != nil {
			return
		}
		if reachesAmountLimit {
			err = fmt.Errorf("%d is reaching the limit", t.UserID)
			return
		}

		err = insertTranaction(conn, t)
		if err != nil {
			return
		}
		log.Printf("done inserting transaction %+v \n", t)
		return
	})
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
