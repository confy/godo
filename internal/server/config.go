package server

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	DB_URL   string
	DB_token string
	LogLevel slog.Level
}

func (c *Config) GetDbURL() (string, error) {
	url := c.DB_URL
	if len(c.DB_token) > 0 {
		url = url + "?authToken=" + c.DB_token
	}

	if len(url) == 0 {
		return "", errors.New("No DB url given.")
	}
	return url, nil
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stdout, ".env not found: %s\n", err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	db_url := os.Getenv("DB_URL")
	db_token := os.Getenv("DB_TOKEN")

	if host == "" || port == "" || db_url == "" || db_token == "" {
		return nil, errors.New("Missing required environment variable")
	}

	return &Config{
		Host:     host,
		Port:     port,
		DB_URL:   db_url,
		DB_token: db_token,
	}, nil
}