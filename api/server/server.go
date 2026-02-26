package server

import (
	"config-manager/configs"
	"config-manager/di"
	_ "config-manager/docs" // Swagger docs

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var ControllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Start the Controller server",
	Run: func(cmd *cobra.Command, args []string) {
		startController()
	},
}

func startController() {
	cfg := configs.LoadConfig()
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Init DI
	di.InitializeControllerV1(e, cfg)

	e.Logger.Fatal(e.Start(":" + cfg.ControllerPort))
}

var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start the Agent server",
	Run: func(cmd *cobra.Command, args []string) {
		startAgent()
	},
}

func startAgent() {
	cfg := configs.LoadConfig()
	di.InitializeAgent(cfg)
}

var WorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the Worker server",
	Run: func(cmd *cobra.Command, args []string) {
		startWorker()
	},
}

func startWorker() {
	cfg := configs.LoadConfig()
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())

	di.InitializeWorker(e, cfg)

	e.Logger.Fatal(e.Start(":" + cfg.WorkerPort))
}
