package config

import (
	"log"

	"wildfire-risk-platform/shared/config"
)

type Config struct {
	Port            string
	Environment     string
	OrchestratorURL string
}

func Load() *Config {
	// Load environment variables
	config.LoadEnv()

	cfg := &Config{
		Port:            config.GetEnv("PORT", "8000"),
		Environment:     config.GetEnv("ENVIRONMENT", "development"),
		OrchestratorURL: config.GetEnv("ORCHESTRATOR_URL", "orchestrator:9000"),
	}

	log.Printf("Loaded configuration: %+v", cfg)
	return cfg
}
