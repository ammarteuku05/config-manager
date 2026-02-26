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

// MockAgentUsecase is a mock for the AgentUsecase interface
type MockAgentUsecase struct {
	mock.Mock
}

func (m *MockAgentUsecase) Register(req dto.AgentRegisterRequest) (*dto.AgentRegisterResponse, error) {
	args := m.Called(req)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.AgentRegisterResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAgentHandler_Register(t *testing.T) {
	e := echo.New()
	log := logger.NewLogger()
	mockUsecase := new(MockAgentUsecase)

	h := &AgentHandler{
		agentUsecase: mockUsecase,
		logger:       log,
	}

	t.Run("Success", func(t *testing.T) {
		reqBody := dto.AgentRegisterRequest{Name: "test-agent"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expectedResp := &dto.AgentRegisterResponse{
			AgentID:             "123",
			PollURL:             "/poll",
			PollIntervalSeconds: 30,
		}
		mockUsecase.On("Register", reqBody).Return(expectedResp, nil).Once()

		err := h.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resResp dto.AgentRegisterResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resResp)
		assert.NoError(t, err)
		assert.Equal(t, expectedResp.AgentID, resResp.AgentID)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Bind Error", func(t *testing.T) {
		// Invalid json
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("{invalid_json}"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Usecase Error", func(t *testing.T) {
		reqBody := dto.AgentRegisterRequest{Name: "test-agent"}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase.On("Register", reqBody).Return(nil, errors.New("db error")).Once()

		err := h.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockUsecase.AssertExpectations(t)
	})
}
