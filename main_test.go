package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMainExecution(t *testing.T) {
	// Let's run the main command but pass help so it doesn't block
	// or we can test if we can execute the root command.

	// Temporarily redirect args to simulate running with "help"
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"config-manager", "help"}

	assert.NotPanics(t, func() {
		// Just run it in a goroutine in case it blocks
		go main()
		time.Sleep(100 * time.Millisecond)
	})
}
