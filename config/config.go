package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host               string
	Port               string
	DBURL              string
	DBToken            string
	LogLevel           slog.Level
	AppEnv             string
	UseHTTPS           bool
	GithubClientID     string
	GithubClientSecret string
}

// GetHostURL returns the host URL
func (c *Config) GetHostURL() string {
	protocol := "http"
	if c.UseHTTPS {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s:%s", protocol, c.Host, c.Port)
}

// GetDBURL returns the database URL with the auth token
func (c *Config) GetDBURL() string {
	return c.DBURL + "?authToken=" + c.DBToken
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
	useHTTPS := appEnv == "prod"
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_TOKEN")
	githubClientID := os.Getenv("OAUTH_GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("OAUTH_GITHUB_SECRET")

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

	if appEnv == "" ||
		host == "" ||
		port == "" ||
		dbURL == "" ||
		dbToken == "" ||
		githubClientID == "" ||
		githubClientSecret == "" {
		return nil, errors.New("missing one or more required environment variables")
	}

	return &Config{
		Host:               host,
		Port:               port,
		DBURL:              dbURL,
		DBToken:            dbToken,
		LogLevel:           logLevel,
		AppEnv:             appEnv,
		UseHTTPS:           useHTTPS,
		GithubClientID:     githubClientID,
		GithubClientSecret: githubClientSecret,
	}, nil
}
