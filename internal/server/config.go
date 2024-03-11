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
	DbURL    string
	DbToken  string
	LogLevel slog.Level
}

func (c *Config) GetDbURL() (string, error) {
	url := c.DbURL
	if len(c.DbToken) > 0 {
		url = url + "?authToken=" + c.DbToken
	}

	if len(url) == 0 {
		return "", errors.New("no db url given")
	}
	return url, nil
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			fmt.Fprintf(os.Stdout, ".env not found: %s\n", err)
			return nil, err
		}
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	DbURL := os.Getenv("DB_URL")
	DbToken := os.Getenv("DB_TOKEN")

	if host == "" || port == "" || DbURL == "" || DbToken == "" {
		return nil, errors.New("missing required environment variable")
	}

	return &Config{
		Host:     host,
		Port:     port,
		DbURL:    DbURL,
		DbToken:  DbToken,
		LogLevel: slog.LevelDebug,
	}, nil
}
