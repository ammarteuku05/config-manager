package main

import (
	"config-manager/api/server"

	"github.com/spf13/cobra"
)

// @title Distributed Config Manager API
// @version 1.0
// @description This is a distributed configuration management system server.
// @host localhost:8080
// @BasePath /v1
func main() {
	var rootCmd = &cobra.Command{
		Use:   "config-manager",
		Short: "Distributed Configuration Management",
		Long:  "A distributed configuration management system with Controller, Agent, and Worker",
	}

	rootCmd.AddCommand(server.ControllerCmd)
	rootCmd.AddCommand(server.AgentCmd)
	rootCmd.AddCommand(server.WorkerCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
