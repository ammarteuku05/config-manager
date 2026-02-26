package usecase

import (
	"bytes"
	"config-manager/configs"
	"config-manager/internal/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AgentManager interface {
	Register() (*dto.AgentRegisterResponse, error)
	PushToWorker(req dto.ConfigRequest) error
}

type agentManager struct {
	cfg        *configs.Config
	httpClient *http.Client
}

func NewAgentManager(cfg *configs.Config) AgentManager {
	return &agentManager{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (m *agentManager) Register() (*dto.AgentRegisterResponse, error) {
	reqBody, _ := json.Marshal(dto.AgentRegisterRequest{Name: "agent-1"})

	url := fmt.Sprintf("%s/v1/register", m.cfg.ControllerURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", m.cfg.AgentAuthToken)

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to register, status: %d", resp.StatusCode)
	}

	var registerResp dto.AgentRegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return nil, err
	}

	return &registerResp, nil
}

func (m *agentManager) PushToWorker(req dto.ConfigRequest) error {
	reqBody, _ := json.Marshal(req)

	url := m.cfg.WorkerURL + "/v1/config"
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push config to worker, status: %d", resp.StatusCode)
	}
	return nil
}
