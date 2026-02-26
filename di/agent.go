package di

import (
	"config-manager/configs"
	"config-manager/internal/handler"
	"config-manager/internal/logger"
	"config-manager/internal/usecase"
)

func InitializeAgent(cfg *configs.Config) {
	agentManager := usecase.NewAgentManager(cfg)
	log := logger.NewLogger()
	poller := handler.NewControllerPoller(cfg, agentManager, log)

	// Block and run
	poller.Start()
}
