package di

import (
	"config-manager/configs"
	"config-manager/internal/domain"
	"config-manager/internal/handler"
	"config-manager/internal/logger"
	"config-manager/internal/repository"
	"config-manager/internal/usecase"
	"config-manager/pkg/shared/utils"

	"github.com/labstack/echo/v4"
)

func InitializeControllerV1(e *echo.Echo, cfg *configs.Config) {
	// Init DB
	db, err := utils.InitDB(cfg.DBPath)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Migrate
	db.AutoMigrate(&domain.Agent{}, &domain.GlobalConfig{})

	// Repositories
	agentRepo := repository.NewAgentRepository(db)
	configRepo := repository.NewConfigRepository(db)

	// Usecases
	agentUsecase := usecase.NewAgentUsecase(agentRepo, cfg.PollURL, cfg.PollInterval)
	configUsecase := usecase.NewConfigUsecase(configRepo)

	// Group V1
	v1 := e.Group("/v1")

	// Logger
	log := logger.NewLogger()

	// Handlers
	handler.NewAgentHandler(v1, agentUsecase, log, cfg.AgentAuthToken)
	handler.NewConfigHandler(v1, configUsecase, log)
}
