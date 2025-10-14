package main_test

import (
	"bytes"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	baseURL     = "http://localhost:8888" // Test server URL
	amountLimit = 1000                    // Maximum total transaction amount per user
)

type Transaction struct {
	UserID      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

//go:embed db/init.sql
var initSQL string

// TestCreate Test whether the transaction registration implementation passes the test.
func TestCreate(t *testing.T) {
	conn, err := sql.Open("mysql", "root@tcp(127.0.0.1)/codetest")
	if err != nil {
		t.Fatal(err)
	}
	// Cleanup
	for _, q := range strings.Split(initSQL, ";") {
		if q == "" {
			continue
		}
		if _, err := conn.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := conn.Exec("delete from codetest.transactions"); err != nil {
		t.Fatal(err)
	}

	// POST transaction registration requests in parallel
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 6; j++ {
				uID := (i+j)%2 + 1 // User ID to be tested. Either 1 or 2.
				req, err := request(uID)
				if err != nil {
					t.Error(err)
					return
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				// Test whether an unexpected response status was returned
				if resp.StatusCode != http.StatusPaymentRequired && resp.StatusCode != http.StatusCreated {
					t.Errorf("POST /transactions status %d", resp.StatusCode)
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				t.Log(string(body))

				if err := resp.Body.Close(); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()

	// Check if there is a user that has an amount greater than the per-user limit registered
	for _, uID := range []int{1, 2} {
		var got int
		if err := conn.QueryRow("select sum(amount) from codetest.transactions where user_id=?", uID).
			Scan(&got); err != nil {
			t.Fatal(err)
		}
		want := amountLimit
		if got != want {
			t.Errorf("sum(amount) of user:%d = %d, want %d", uID, got, want)
		}
	}
}

func request(uID int) (*http.Request, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, 128))
	if err := json.NewEncoder(buffer).Encode(Transaction{
		UserID:      uID,
		Amount:      100,
		Description: fmt.Sprintf("商品%d", uID),
	}); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/transactions",
		buffer,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", fmt.Sprintf("secure-api-key-%d", uID))
	return req, nil
}
