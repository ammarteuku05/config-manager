package di

import (
	"config-manager/configs"
	"config-manager/internal/handler"
	"config-manager/internal/logger"
	"config-manager/internal/usecase"

	"github.com/labstack/echo/v4"
)

func InitializeWorker(e *echo.Echo, cfg *configs.Config) {
	configManager := usecase.NewConfigManager()
	log := logger.NewLogger()

	handler.NewWorkerHandler(e, configManager, log)
}
