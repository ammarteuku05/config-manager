package configs

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration values.
// Environment variables can override the default values.
type Config struct {
	ControllerPort string `envconfig:"CONTROLLER_PORT" default:"8080"`
	AgentPort      string `envconfig:"AGENT_PORT" default:"8081"`
	WorkerPort     string `envconfig:"WORKER_PORT" default:"8082"`
	DBPath         string `envconfig:"DB_PATH" default:"controller.db"`
	ControllerURL  string `envconfig:"CONTROLLER_URL" default:"http://localhost:8080"`
	WorkerURL      string `envconfig:"WORKER_URL" default:"http://localhost:8082"`
	PollInterval   int    `envconfig:"POLL_INTERVAL" default:"30"`
	AgentAuthToken string `envconfig:"AGENT_AUTH_TOKEN" default:"agent-secret"`
	PollURL        string `envconfig:"POLL_URL" default:"/v1/config"`
}

// LoadConfig returns a Config populated by envconfig.
func LoadConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to process envconfig: %v", err)
	}
	return &cfg
}
