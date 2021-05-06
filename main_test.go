package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	baseURL     = "http://localhost:8888" // テスト対象サーバー
	amountLimit = 1000                    // 1日あたりの登録可能な取引金額上限
)

type Transaction struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

// TestCreate 取引登録処理が仕様を満たしているかテストする。
func TestCreate(t *testing.T) {
	conn, err := sql.Open("mysql", "root@tcp(127.0.0.1)/codetest")
	if err != nil {
		t.Fatal(err)
	}
	// クリーンアップ
	if _, err := conn.Exec("delete from transactions"); err != nil {
		t.Fatal(err)
	}

	const uID = 1 // テスト対象のユーザーID

	// 並列で取引登録リクエストをPOSTする
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 3; j++ {
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

				// 想定外のレスポンスステータスが返ってきていないかをテスト
				if resp.StatusCode != http.StatusPaymentRequired && resp.StatusCode != http.StatusCreated {
					t.Errorf("POST /transactions status %d", resp.StatusCode)
				}

				body, err := ioutil.ReadAll(resp.Body)
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

	// 1日あたりの登録可能な取引金額上限を超えて登録されていないかをテスト
	var amount int
	if err := conn.QueryRow(
		"select sum(amount) from transactions where user_id=?",
		uID,
	).Scan(&amount); err != nil {
		t.Fatal(err)
	}
	if amount > amountLimit {
		t.Errorf("amount %d over the amountLimit %d", amount, amountLimit)
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
