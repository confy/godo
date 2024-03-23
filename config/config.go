package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HostIP             string
	Hostname 		   string
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

	if c.AppEnv == "prod" {
		return fmt.Sprintf("%s://%s", protocol, c.Hostname)
	}
	return fmt.Sprintf("%s://%s:%s", protocol, c.Hostname, c.Port)
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
	hostIP := os.Getenv("HOST_IP")
	hostname := os.Getenv("HOSTNAME")
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
		hostIP == "" ||
		hostname == "" ||
		port == "" ||
		dbURL == "" ||
		dbToken == "" ||
		githubClientID == "" ||
		githubClientSecret == "" {
		return nil, errors.New("missing one or more required environment variables")
	}

	return &Config{
		HostIP: 		    hostIP,
		Hostname:           hostname,
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
