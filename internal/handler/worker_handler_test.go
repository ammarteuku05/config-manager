package handler

import (
	"bytes"
	"config-manager/internal/dto"
	"config-manager/internal/logger"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfigManager is a mock for the ConfigManager interface
type MockConfigManager struct {
	mock.Mock
}

func (m *MockConfigManager) UpdateConfig(req dto.ConfigRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockConfigManager) ExecuteHit() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestWorkerHandler_ReceiveConfig(t *testing.T) {
	e := echo.New()
	log := logger.NewLogger()
	mockManager := new(MockConfigManager)

	h := &WorkerHandler{
		configManager: mockManager,
		logger:        log,
	}

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.ConfigRequest{Config: map[string]interface{}{"key": "value"}}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/config", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockManager.On("UpdateConfig", reqBody).Return(nil).Once()

		err := h.ReceiveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockManager.AssertExpectations(t)
	})

	t.Run("Bind Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/config", bytes.NewBufferString("{invalid_json}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.ReceiveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Manager Error", func(t *testing.T) {
		reqBody := dto.ConfigRequest{Config: map[string]interface{}{"key": "value"}}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/v1/config", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockManager.On("UpdateConfig", reqBody).Return(errors.New("update error")).Once()

		err := h.ReceiveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockManager.AssertExpectations(t)
	})
}

func TestWorkerHandler_HitProxy(t *testing.T) {
	e := echo.New()
	log := logger.NewLogger()
	mockManager := new(MockConfigManager)

	h := &WorkerHandler{
		configManager: mockManager,
		logger:        log,
	}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hit", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockManager.On("ExecuteHit").Return("success body", nil).Once()

		err := h.HitProxy(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockManager.AssertExpectations(t)
	})

	t.Run("Manager Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hit", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockManager.On("ExecuteHit").Return("", errors.New("hit error")).Once()

		err := h.HitProxy(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockManager.AssertExpectations(t)
	})
}
