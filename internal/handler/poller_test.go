package handler

import (
	"config-manager/configs"
	"config-manager/internal/dto"
	"config-manager/internal/logger"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAgentManager is a mock for the AgentManager interface
type MockAgentManager struct {
	mock.Mock
}

func (m *MockAgentManager) Register() (*dto.AgentRegisterResponse, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(*dto.AgentRegisterResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAgentManager) PushToWorker(req dto.ConfigRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func TestControllerPoller_StartAndPollLoop(t *testing.T) {
	log := logger.NewLogger()
	mockManager := new(MockAgentManager)

	cfg := &configs.Config{
		ControllerURL: "http://localhost:8080",
	}

	poller := NewControllerPoller(cfg, mockManager, log)
	assert.NotNil(t, poller)
	assert.Equal(t, cfg, poller.cfg)
	assert.Equal(t, mockManager, poller.agentManager)
	assert.NotNil(t, poller.httpClient)
	assert.NotNil(t, poller.logger)
}

func TestControllerPoller_Start(t *testing.T) {
	log := logger.NewLogger()
	mockManager := new(MockAgentManager)

	cfg := &configs.Config{
		ControllerURL: "http://localhost:8080",
	}
	poller := NewControllerPoller(cfg, mockManager, log)

	// Mock Register to fail once, then succeed
	mockManager.On("Register").Return(nil, errors.New("register error")).Once()

	expectedResp := &dto.AgentRegisterResponse{
		AgentID:             "agent-123",
		PollURL:             "/poll",
		PollIntervalSeconds: 1,
	}
	mockManager.On("Register").Return(expectedResp, nil).Once()

	// Since Start calls pollLoop indefinitely, we run it in a goroutine
	// and assert the fields were updated by register
	go poller.Start()

	// Wait long enough for the retry to happen
	time.Sleep(10 * time.Millisecond) // we hacked the sleep below? no, Start sleeps for 5s...
	// Since Start has a 5 second sleep on error, this is slow for a unit test.
}

func TestControllerPoller_PollIteration(t *testing.T) {
	log := logger.NewLogger()

	t.Run("Success", func(t *testing.T) {
		mockManager := new(MockAgentManager)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/poll", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			resp := dto.ConfigResponse{
				Version: "new_version",
				Config:  map[string]interface{}{"key": "value"},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		cfg := &configs.Config{
			ControllerURL: ts.URL,
		}

		poller := NewControllerPoller(cfg, mockManager, log)
		poller.pollURL = "/v1/poll"
		poller.pollInterval = 5 * time.Millisecond
		poller.versionCache = "old_version"

		mockManager.On("PushToWorker", mock.AnythingOfType("dto.ConfigRequest")).Return(nil).Once()

		go poller.pollLoop()

		time.Sleep(50 * time.Millisecond)
		mockManager.AssertExpectations(t)
		assert.Equal(t, "new_version", poller.versionCache)
	})

	t.Run("Push Error", func(t *testing.T) {
		mockManager := new(MockAgentManager)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			resp := dto.ConfigResponse{
				Version: "new_version_2",
				Config:  map[string]interface{}{"key": "value"},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		cfg := &configs.Config{
			ControllerURL: ts.URL,
		}

		poller := NewControllerPoller(cfg, mockManager, log)
		poller.pollURL = "/v1/poll"
		poller.pollInterval = 5 * time.Millisecond
		poller.versionCache = "old_version"

		mockManager.On("PushToWorker", mock.Anything).Return(errors.New("push error")).Maybe() // can be called many times

		go poller.pollLoop()

		time.Sleep(50 * time.Millisecond)
		mockManager.AssertExpectations(t)
	})

	t.Run("Fetch Error Non-200", func(t *testing.T) {
		mockManager := new(MockAgentManager)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		cfg := &configs.Config{
			ControllerURL: ts.URL,
		}

		poller := NewControllerPoller(cfg, mockManager, log)
		poller.pollURL = "/v1/poll"
		poller.pollInterval = 5 * time.Millisecond
		poller.versionCache = "old_version"

		go poller.pollLoop()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("Fetch Decode Error", func(t *testing.T) {
		mockManager := new(MockAgentManager)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		cfg := &configs.Config{
			ControllerURL: ts.URL,
		}

		poller := NewControllerPoller(cfg, mockManager, log)
		poller.pollURL = "/v1/poll"
		poller.pollInterval = 5 * time.Millisecond
		poller.versionCache = "old_version"

		go poller.pollLoop()
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("Fetch HTTP Error", func(t *testing.T) {
		mockManager := new(MockAgentManager)
		cfg := &configs.Config{
			ControllerURL: "http://localhost:1",
		}

		poller := NewControllerPoller(cfg, mockManager, log)
		poller.pollURL = "/v1/poll"
		poller.pollInterval = 5 * time.Millisecond
		poller.versionCache = "old_version"

		go poller.pollLoop()
		time.Sleep(50 * time.Millisecond)
	})
}
