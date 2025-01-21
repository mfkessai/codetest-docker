package domain

type User struct {
	ID     int
	Name   string
	ApiKey string
}

type UserRepository interface {
	FindByApiKey(apiKey string) (*User, error)
}
