package di

import (
	"config-manager/configs"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInitializeWorker(t *testing.T) {
	e := echo.New()
	cfg := &configs.Config{}

	assert.NotPanics(t, func() {
		InitializeWorker(e, cfg)
	})

	routes := e.Routes()
	assert.Greater(t, len(routes), 0)

	var hasConfig bool
	var hasHit bool
	for _, r := range routes {
		if r.Path == "/v1/config" && r.Method == "POST" {
			hasConfig = true
		}
		if r.Path == "/hit" && r.Method == "GET" {
			hasHit = true
		}
	}
	assert.True(t, hasConfig)
	assert.True(t, hasHit)
}
