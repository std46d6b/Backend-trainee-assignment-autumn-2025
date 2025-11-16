package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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
	Port    int

	ShutdownTimeout time.Duration
}

func envOnly(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env var %s not found", key)
	}
	return value, nil
}

func envOrDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func intEnvOrDefault(key string, defaultValue int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("error parsing int env var %s: %w", key, err)
	}
	return intValue, nil
}

func int32EnvOrDefault(key string, defaultValue int32) (int32, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, nil
	}
	intValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing int32 env var %s: %w", key, err)
	}
	return int32(intValue), nil
}

func Load() (*Config, error) {
	dbCfg, err := loadDBConfig()
	if err != nil {
		return nil, err
	}

	webServerCfg, err := loadWebServerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfig:        dbCfg,
		WebServerConfig: webServerCfg,
	}, nil
}

func loadDBConfig() (*DBConfig, error) {
	databaseURL, err := envOnly("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	maxConns, err := int32EnvOrDefault("MAX_CONNS", defaultMaxConns)
	if err != nil {
		return nil, err
	}

	minConns, err := int32EnvOrDefault("MIN_CONNS", defaultMinConns)
	if err != nil {
		return nil, err
	}

	healthCheckIntervalInSeconds, err := intEnvOrDefault(
		"HEALTH_CHECK_INTERVAL_IN_SECONDS",
		defaultHealthCheckIntervalInSeconds,
	)
	if err != nil {
		return nil, err
	}

	return &DBConfig{
		DatabaseURL:         databaseURL,
		MaxConns:            maxConns,
		MinConns:            minConns,
		HealthCheckInterval: time.Duration(healthCheckIntervalInSeconds) * time.Second,
	}, nil
}

func loadWebServerConfig() (*WebServerConfig, error) {
	webServerAddress := envOrDefault("WEB_SERVER_ADDRESS", defaultAddress)

	webServerPort, err := intEnvOrDefault("WEB_SERVER_PORT", defaultPort)
	if err != nil {
		return nil, err
	}

	shutdownTimeoutInSeconds, err := intEnvOrDefault("SHUTDOWN_TIMEOUT_IN_SECONDS", defaultShutdownTimeoutInSeconds)
	if err != nil {
		return nil, err
	}

	return &WebServerConfig{
		Address:         webServerAddress,
		Port:            webServerPort,
		ShutdownTimeout: time.Duration(shutdownTimeoutInSeconds) * time.Second,
	}, nil
}
