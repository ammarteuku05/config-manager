package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("CONTROLLER_PORT", "9090")
	os.Setenv("AGENT_PORT", "9091")
	os.Setenv("WORKER_PORT", "9092")
	os.Setenv("DB_PATH", "test.db")
	os.Setenv("CONTROLLER_URL", "http://test-controller:9090")
	os.Setenv("WORKER_URL", "http://test-worker:9092")
	os.Setenv("POLL_INTERVAL", "60")

	// Cleanup
	defer func() {
		os.Unsetenv("CONTROLLER_PORT")
		os.Unsetenv("AGENT_PORT")
		os.Unsetenv("WORKER_PORT")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("CONTROLLER_URL")
		os.Unsetenv("WORKER_URL")
		os.Unsetenv("POLL_INTERVAL")
	}()

	cfg := LoadConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.ControllerPort)
	assert.Equal(t, "9091", cfg.AgentPort)
	assert.Equal(t, "9092", cfg.WorkerPort)
	assert.Equal(t, "test.db", cfg.DBPath)
	assert.Equal(t, "http://test-controller:9090", cfg.ControllerURL)
	assert.Equal(t, "http://test-worker:9092", cfg.WorkerURL)
	assert.Equal(t, 60, cfg.PollInterval)
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Ensure no environment variables are set that could conflict
	os.Clearenv()

	cfg := LoadConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.ControllerPort)
	assert.Equal(t, "8081", cfg.AgentPort)
	assert.Equal(t, "8082", cfg.WorkerPort)
	assert.Equal(t, "controller.db", cfg.DBPath)
	assert.Equal(t, "http://localhost:8080", cfg.ControllerURL)
	assert.Equal(t, "http://localhost:8082", cfg.WorkerURL)
	assert.Equal(t, 30, cfg.PollInterval)
}
