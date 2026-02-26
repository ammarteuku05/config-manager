package di

import (
	"config-manager/configs"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInitializeControllerV1(t *testing.T) {
	e := echo.New()
	cfg := &configs.Config{
		DBPath:       "file::memory:?cache=shared",
		PollInterval: 30,
	}

	// Given a valid config and echo instance, InitializeControllerV1 should not panic
	assert.NotPanics(t, func() {
		InitializeControllerV1(e, cfg)
	})

	// Check if routes are registered
	routes := e.Routes()
	assert.Greater(t, len(routes), 0)

	var hasRegister bool
	var hasConfig bool
	for _, r := range routes {
		if r.Path == "/v1/register" && r.Method == "POST" {
			hasRegister = true
		}
		if r.Path == "/v1/config" {
			hasConfig = true
		}
	}
	assert.True(t, hasRegister)
	assert.True(t, hasConfig)
}
