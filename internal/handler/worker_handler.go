package handler

import (
	"config-manager/internal/dto"
	"config-manager/internal/usecase"

	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WorkerHandler struct {
	configManager usecase.ConfigManager
	logger        *slog.Logger
}

func NewWorkerHandler(e *echo.Echo, configManager usecase.ConfigManager, logger *slog.Logger) {
	handler := &WorkerHandler{
		configManager: configManager,
		logger:        logger,
	}

	v1 := e.Group("/v1")
	v1.POST("/config", handler.ReceiveConfig)

	e.GET("/hit", handler.HitProxy)
}

func (h *WorkerHandler) ReceiveConfig(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	var req dto.ConfigRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusBadRequest,
			"request_id": reqID,
		})
	}

	if err := h.configManager.UpdateConfig(req); err != nil {
		h.logger.Error("failed to update config", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusInternalServerError,
			"request_id": reqID,
		})
	}

	h.logger.Info("Worker received new config", "config", req.Config, "request_id", reqID)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "config updated",
		"code":       http.StatusOK,
		"request_id": reqID,
	})
}

func (h *WorkerHandler) HitProxy(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	result, err := h.configManager.ExecuteHit()
	if err != nil {
		h.logger.Error("failed to execute hit proxy", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":      "Worker Error: " + err.Error(),
			"code":       http.StatusInternalServerError,
			"request_id": reqID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"result":     result,
		"code":       http.StatusOK,
		"request_id": reqID,
	})
}
