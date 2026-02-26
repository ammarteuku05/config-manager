package handler

import (
	"config-manager/internal/dto"
	"config-manager/internal/usecase"
	"log/slog"
	"net/http"

	"config-manager/pkg/shared/middleware"

	"github.com/labstack/echo/v4"
)

type AgentHandler struct {
	agentUsecase usecase.AgentUsecase
	logger       *slog.Logger
}

func NewAgentHandler(e *echo.Group, agentUsecase usecase.AgentUsecase, logger *slog.Logger, authToken string) {
	handler := &AgentHandler{
		agentUsecase: agentUsecase,
		logger:       logger,
	}

	e.POST("/register", handler.Register, middleware.StaticTokenAuth("Authorization", authToken))
}

// Register godoc
// @Summary Register a new agent
// @Description Register a new agent and get polling details
// @Tags Agent
// @Accept json
// @Produce json
// @Param req body dto.AgentRegisterRequest true "Agent Registration"
// @Success 200 {object} dto.AgentRegisterResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AgentHandler) Register(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	var req dto.AgentRegisterRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusBadRequest,
			"request_id": reqID,
		})
	}

	// Assuming basic static validation or handled in middleware
	res, err := h.agentUsecase.Register(req)
	if err != nil {
		h.logger.Error("failed to register agent", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusInternalServerError,
			"request_id": reqID,
		})
	}

	res.Code = http.StatusOK
	res.RequestID = reqID
	return c.JSON(http.StatusOK, res)
}
