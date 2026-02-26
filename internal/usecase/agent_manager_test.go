package usecase

import (
	"config-manager/configs"
	"config-manager/internal/dto"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentManager_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/register", r.URL.Path)
			assert.Equal(t, "secret-token", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
			resp := dto.AgentRegisterResponse{
				AgentID: "agent-123",
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer ts.Close()

		cfg := &configs.Config{ControllerURL: ts.URL, AgentAuthToken: "secret-token"}
		manager := NewAgentManager(cfg)

		resp, err := manager.Register()
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "agent-123", resp.AgentID)
	})

	t.Run("ServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		cfg := &configs.Config{ControllerURL: ts.URL}
		manager := NewAgentManager(cfg)

		resp, err := manager.Register()
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		cfg := &configs.Config{ControllerURL: "http://\x00invalid"}
		manager := NewAgentManager(cfg)

		resp, err := manager.Register()
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("DecodeError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		cfg := &configs.Config{ControllerURL: ts.URL}
		manager := NewAgentManager(cfg)

		resp, err := manager.Register()
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestAgentManager_PushToWorker(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v1/config", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		cfg := &configs.Config{WorkerURL: ts.URL}
		manager := NewAgentManager(cfg)

		req := dto.ConfigRequest{Config: map[string]interface{}{"k": "v"}}
		err := manager.PushToWorker(req)
		assert.NoError(t, err)
	})

	t.Run("ServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		cfg := &configs.Config{WorkerURL: ts.URL}
		manager := NewAgentManager(cfg)

		req := dto.ConfigRequest{Config: map[string]interface{}{"k": "v"}}
		err := manager.PushToWorker(req)
		assert.Error(t, err)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		cfg := &configs.Config{WorkerURL: "http://\x00invalid"}
		manager := NewAgentManager(cfg)
		manager.(*agentManager).httpClient.Timeout = 10 * time.Millisecond

		req := dto.ConfigRequest{Config: map[string]interface{}{"k": "v"}}
		err := manager.PushToWorker(req)
		assert.Error(t, err)
	})
}
