package infra

import "testing"

func TestUserRepository(t *testing.T) {
	testCfg := NewConfig()

	t.Run("FindByApiKey returns a user when the API key exists", func(t *testing.T) {
		conn, err := GetConnection(testCfg)

		repo := NewUserRepository(conn)
		apiKey := "secure-api-key-1"

		user, err := repo.FindByApiKey(apiKey)
		if err != nil {
			t.Fatal(err)
		}

		if user == nil {
			t.Errorf("got %v, want not nil", user)
		}

		if user.ID != 1 {
			t.Errorf("got %d, want 1", user.ID)
		}
	})

	t.Run("FindByApiKey returns nil when the API key does not exist", func(t *testing.T) {
		conn, err := GetConnection(testCfg)

		repo := NewUserRepository(conn)
		apiKey := "invalidKey"

		user, err := repo.FindByApiKey(apiKey)
		if err != nil {
			t.Fatal(err)
		}

		if user != nil {
			t.Errorf("got %v, want nil", user)
		}

	})
}
