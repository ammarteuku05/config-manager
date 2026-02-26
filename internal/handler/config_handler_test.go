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

// MockConfigUsecase is a mock for the ConfigUsecase interface
type MockConfigUsecase struct {
	mock.Mock
}

func (m *MockConfigUsecase) Save(req dto.ConfigRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockConfigUsecase) GetLatest() (*dto.ConfigResponse, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*dto.ConfigResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestConfigHandler_SaveConfig(t *testing.T) {
	e := echo.New()
	log := logger.NewLogger()
	mockUsecase := new(MockConfigUsecase)
	h := &ConfigHandler{configUsecase: mockUsecase, logger: log}

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.ConfigRequest{Config: map[string]interface{}{"key": "value"}}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("Save", reqBody).Return(nil).Once()

		err := h.SaveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Bind Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString("{invalid_json}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.SaveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Usecase Error", func(t *testing.T) {
		reqBody := dto.ConfigRequest{Config: map[string]interface{}{"key": "value"}}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("Save", reqBody).Return(errors.New("db error")).Once()

		err := h.SaveConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestConfigHandler_GetConfig(t *testing.T) {
	e := echo.New()
	log := logger.NewLogger()
	mockUsecase := new(MockConfigUsecase)
	h := &ConfigHandler{configUsecase: mockUsecase, logger: log}

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedResp := &dto.ConfigResponse{
			Config:  map[string]interface{}{"key": "value"},
			Version: "123",
		}
		mockUsecase.On("GetLatest").Return(expectedResp, nil).Once()

		err := h.GetConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "123", rec.Header().Get("ETag"))
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Usecase Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("GetLatest").Return(nil, errors.New("db error")).Once()

		err := h.GetConfig(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}
