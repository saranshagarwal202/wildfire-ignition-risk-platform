package config

import (
	"strconv"
	"time"
	"wildfire-risk-platform/shared/config"
)

type Config struct {
	GRPCPort       string
	OverpassAPIURL string
	HTTPTimeout    time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
}

func LoadConfig() (*Config, error) {

	config.LoadEnv()
	http_timeout, err := time.ParseDuration(config.GetEnv("HTTP_TIMEOUT", "60s"))
	if err != nil {
		return &Config{}, err
	}
	max_retries, err := strconv.Atoi(config.GetEnv("MAX_RETRIES", "3"))
	if err != nil {
		return &Config{}, err
	}
	retry_delay, err := time.ParseDuration(config.GetEnv("RETRY_DELAY", "2s"))
	if err != nil {
		return &Config{}, err
	}

	cfg := &Config{

		GRPCPort:       config.GetEnv("GRPC_PORT", "50052"),
		OverpassAPIURL: config.GetEnv("OVERPASS_API_URL", "https://overpass-api.de/api/interpreter"),
		HTTPTimeout:    http_timeout,
		MaxRetries:     max_retries,
		RetryDelay:     retry_delay,
	}

	return cfg, nil
}
