package usecase

import (
	"config-manager/internal/dto"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type ConfigManager interface {
	UpdateConfig(req dto.ConfigRequest) error
	ExecuteHit() (string, error)
}

type configManager struct {
	mu         sync.RWMutex
	config     map[string]interface{}
	httpClient *http.Client
}

func NewConfigManager() ConfigManager {
	return &configManager{
		config:     make(map[string]interface{}),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *configManager) UpdateConfig(req dto.ConfigRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = req.Config
	return nil
}

func (m *configManager) ExecuteHit() (string, error) {
	m.mu.RLock()
	urlInter, ok := m.config["url"]
	m.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("URL not configured")
	}

	urlStr, ok := urlInter.(string)
	if !ok {
		return "", fmt.Errorf("configured URL is not a string")
	}

	resp, err := m.httpClient.Get(urlStr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
