package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var ErrDatabaseURLNotFound = errors.New("DATABASE_URL not found")

type Config struct {
	DBConfig        *DBConfig
	WebServerConfig *WebServerConfig
}

type DBConfig struct {
	DatabaseURL         string
	MaxConns            int32
	MinConns            int32
	HealthCheckInterval time.Duration
}

type WebServerConfig struct {
	Address string
	Port    string

	ShutdownTimeout time.Duration
}

func Load() (*Config, error) {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, ErrDatabaseURLNotFound
	}

	maxConnsStr, ok := os.LookupEnv("MAX_CONNS")
	if !ok {
		return nil, errors.New("MAX_CONNS not found")
	}
	minConnsStr, ok := os.LookupEnv("MIN_CONNS")
	if !ok {
		return nil, errors.New("MIN_CONNS not found")
	}
	healthCheckIntervalInSecondsStr, ok := os.LookupEnv("HEALTH_CHECK_INTERVAL_IN_SECONDS")
	if !ok {
		return nil, errors.New("HEALTH_CHECK_INTERVAL_IN_SECONDS not found")
	}

	maxConns, err := strconv.ParseInt(maxConnsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("MAX_CONNS parse error: %w", err)
	}
	minConns, err := strconv.ParseInt(minConnsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("MIN_CONNS parse error: %w", err)
	}
	healthCheckIntervalInSeconds, err := strconv.ParseInt(healthCheckIntervalInSecondsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("HEALTH_CHECK_INTERVAL_IN_SECONDS parse error: %w", err)
	}

	webServerAddress, ok := os.LookupEnv("WEB_SERVER_ADDRESS")
	if !ok {
		return nil, ErrDatabaseURLNotFound
	}

	webServerPort, ok := os.LookupEnv("WEB_SERVER_PORT")
	if !ok {
		return nil, ErrDatabaseURLNotFound
	}

	shutdownTimeoutInSecondsStr, ok := os.LookupEnv("SHUTDOWN_TIMEOUT_IN_SECONDS")
	if !ok {
		return nil, errors.New("SHUTDOWN_TIMEOUT_IN_SECONDS not found")
	}

	shutdownTimeoutInSeconds, err := strconv.ParseInt(shutdownTimeoutInSecondsStr, 10, 64)
	if err != nil {
		return nil, ErrDatabaseURLNotFound
	}

	return &Config{
		DBConfig: &DBConfig{
			DatabaseURL:         databaseURL,
			MaxConns:            int32(maxConns),
			MinConns:            int32(minConns),
			HealthCheckInterval: time.Duration(healthCheckIntervalInSeconds) * time.Second,
		},
		WebServerConfig: &WebServerConfig{
			Address:         webServerAddress,
			Port:            webServerPort,
			ShutdownTimeout: time.Duration(shutdownTimeoutInSeconds) * time.Second,
		},
	}, nil
}
