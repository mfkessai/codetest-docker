package infra

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	instance *sql.DB
	once     sync.Once
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewConfig() *Config {
	return &Config{
		User:     "root",
		Password: "",
		Host:     "127.0.0.1",
		Port:     "3306",
		DBName:   "codetest",
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func GetConnection(cfg *Config) (*sql.DB, error) {
	var err error
	once.Do(func() {
		instance, err = sql.Open("mysql", cfg.DSN())
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return
		}

		// コネクションプールの設定
		instance.SetMaxOpenConns(25)
		instance.SetMaxIdleConns(25)

		// 接続テスト
		if err = instance.Ping(); err != nil {
			log.Printf("Failed to ping database: %v", err)
			return
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return instance, nil
}

func Close() error {
	if instance != nil {
		return instance.Close()
	}
	return nil
}
