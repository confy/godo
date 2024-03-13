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
	AppEnv   string
	UseHttps bool
}

// GetDbURL returns the database URL with the auth token
func (c *Config) GetDbURL() string {
	return c.DbURL + "?authToken=" + c.DbToken
}

// LoadConfig loads configuration from environment variables or an .env file if not in production
func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			fmt.Fprintf(os.Stdout, ".env not found: %s\n", err)
			return nil, err
		}
	}
	appEnv := os.Getenv("APP_ENV")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")
	useHttps := appEnv == "prod"

	var logLevel slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	default:
		logLevel = slog.LevelDebug
	}

	if appEnv == "" || host == "" || port == "" || dbURL == "" || dbToken == "" {
		return nil, errors.New("missing required environment variable")
	}

	return &Config{
		Host:     host,
		Port:     port,
		DbURL:    dbURL,
		DbToken:  dbToken,
		LogLevel: logLevel,
		AppEnv:   appEnv,
		UseHttps: useHttps,
	}, nil
}
