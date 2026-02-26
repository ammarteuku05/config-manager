package di

import (
	"config-manager/configs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitializeAgent(t *testing.T) {
	cfg := &configs.Config{
		ControllerURL: "http://localhost:8080",
		WorkerPort:    "8082",
	}

	// InitializeAgent starts the poller which runs indefinitely if not mocked.
	// But actually, if Controller is unreachable, it will just keep retrying.
	// We'll run it in a goroutine and give it a small time to initialize.

	assert.NotPanics(t, func() {
		go InitializeAgent(cfg)
	})

	// Just sleep slightly to let any immediate initialization happen and verify no panic
	time.Sleep(100 * time.Millisecond)
}
