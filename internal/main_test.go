package internal

import (
	"net/http"
	"testing"
)

/*
 内部の処理をTDDで書いていく用のテストファイル
*/

func TestPOSTTransactions(t *testing.T) {
	t.Run("returns HTTP 200 status code", func(t *testing.T) {
		api_key := "secure-api-key-1"
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/transactions", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("apikey", api_key)

		want := http.StatusOK

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != want {
			t.Errorf("got %d, want %d", resp.StatusCode, want)
		}
	})

	t.Run("returns HTTP 401 when API Key is not in the header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8888/transactions", nil)

		want := http.StatusUnauthorized

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != want {
			t.Errorf("got %d, want %d", resp.StatusCode, want)
		}
	})

	t.Run("returns HTTP 401 when API Key is invalid", func(t *testing.T) {
		api_key := "invalidKey"
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/transactions", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("apikey", api_key)

		want := http.StatusUnauthorized

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != want {
			t.Errorf("got %d, want %d", resp.StatusCode, want)
		}
	})
}
