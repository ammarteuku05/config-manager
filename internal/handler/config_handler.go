package handler

import (
	"config-manager/internal/dto"
	"config-manager/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	configUsecase usecase.ConfigUsecase
	logger        *slog.Logger
}

func NewConfigHandler(e *echo.Group, configUsecase usecase.ConfigUsecase, logger *slog.Logger) {
	handler := &ConfigHandler{
		configUsecase: configUsecase,
		logger:        logger,
	}

	e.POST("/config", handler.SaveConfig) // Requires admin auth
	e.GET("/config", handler.GetConfig)   // Requires agent auth
}

// SaveConfig godoc
// @Summary Save global config
// @Description Update the global configuration for all workers
// @Tags Config
// @Accept json
// @Produce json
// @Param req body dto.ConfigRequest true "New Configuration"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /config [post]
func (h *ConfigHandler) SaveConfig(c echo.Context) error {
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

	if err := h.configUsecase.Save(req); err != nil {
		h.logger.Error("failed to save config", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusInternalServerError,
			"request_id": reqID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "success",
		"code":       http.StatusOK,
		"request_id": reqID,
	})
}

// GetConfig godoc
// @Summary Get global config
// @Description Get the global configuration for workers
// @Tags Config
// @Produce json
// @Success 200 {object} dto.ConfigResponse
// @Failure 500 {object} map[string]string
// @Router /config [get]
func (h *ConfigHandler) GetConfig(c echo.Context) error {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	res, err := h.configUsecase.GetLatest()
	if err != nil {
		h.logger.Error("failed to get config", "error", err.Error(), "request_id", reqID)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":      err.Error(),
			"code":       http.StatusInternalServerError,
			"request_id": reqID,
		})
	}

	c.Response().Header().Set("ETag", res.Version)
	res.Code = http.StatusOK
	res.RequestID = reqID
	return c.JSON(http.StatusOK, res)
}
