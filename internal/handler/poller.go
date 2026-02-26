package handler

import (
	"config-manager/configs"
	"config-manager/internal/usecase"

	"config-manager/internal/dto"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"
)

type ControllerPoller struct {
	cfg          *configs.Config
	agentManager usecase.AgentManager
	httpClient   *http.Client
	agentID      string
	pollURL      string
	versionCache string
	pollInterval time.Duration
	logger       *slog.Logger
}

func NewControllerPoller(cfg *configs.Config, agentManager usecase.AgentManager, logger *slog.Logger) *ControllerPoller {
	return &ControllerPoller{
		cfg:          cfg,
		agentManager: agentManager,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
		logger:       logger,
	}
}

func (p *ControllerPoller) Start() {
	p.logger.Info("Starting Agent Controller Poller...")

	// 1. Register with controller (with exp backoff if needed, kept simple here to match test)
	var regResp *dto.AgentRegisterResponse
	var err error
	for {
		regResp, err = p.agentManager.Register()
		if err == nil {
			break
		}
		p.logger.Error("Failed to register agent", "error", err)
		time.Sleep(5 * time.Second)
	}

	p.agentID = regResp.AgentID
	p.pollURL = regResp.PollURL
	p.pollInterval = time.Duration(regResp.PollIntervalSeconds) * time.Second
	p.logger.Info("Successfully registered agent", "agent_id", p.agentID, "poll_interval", p.pollInterval)

	p.pollLoop()
}

func (p *ControllerPoller) pollLoop() {
	backoffRetries := 0

	for {
		time.Sleep(p.pollInterval)

		url := fmt.Sprintf("%s%s", p.cfg.ControllerURL, p.pollURL)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", p.cfg.AgentAuthToken) // agent credentials

		resp, err := p.httpClient.Do(req)
		if err != nil {
			backoffRetries++
			backoffTime := time.Duration(math.Pow(2, float64(backoffRetries))) * time.Second
			p.logger.Error("Failed to poll controller", "error", err, "backoff_time", backoffTime)
			time.Sleep(backoffTime)
			continue
		}

		// Reset backoff on success
		backoffRetries = 0

		if resp.StatusCode != http.StatusOK {
			p.logger.Error("Unexpected status code from controller", "status_code", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		var configResp dto.ConfigResponse
		if err := json.NewDecoder(resp.Body).Decode(&configResp); err != nil {
			p.logger.Error("Failed to parse config from controller", "error", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Detect config changes
		if configResp.Version != p.versionCache {
			p.logger.Info("Configuration change detected!", "new_version", configResp.Version)
			p.versionCache = configResp.Version

			// Push to worker
			if err := p.agentManager.PushToWorker(dto.ConfigRequest{Config: configResp.Config}); err != nil {
				p.logger.Error("Failed to push config to worker", "error", err)
			} else {
				p.logger.Info("Successfully pushed config to worker")
			}
		}
	}
}
