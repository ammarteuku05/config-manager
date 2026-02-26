package usecase

import (
	"config-manager/internal/domain"
	"config-manager/internal/dto"
	"config-manager/internal/repository"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConfigUsecase interface {
	Save(req dto.ConfigRequest) error
	GetLatest() (*dto.ConfigResponse, error)
}

type configUsecase struct {
	configRepo repository.ConfigRepository
}

func NewConfigUsecase(configRepo repository.ConfigRepository) ConfigUsecase {
	return &configUsecase{configRepo: configRepo}
}

func (u *configUsecase) Save(req dto.ConfigRequest) error {
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return err
	}

	newConfig := &domain.GlobalConfig{
		Config:    string(configBytes),
		Version:   uuid.New().String(),
		CreatedAt: time.Now(),
	}

	return u.configRepo.Save(newConfig)
}

func (u *configUsecase) GetLatest() (*dto.ConfigResponse, error) {
	config, err := u.configRepo.GetLatest()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty config if none exists
			return &dto.ConfigResponse{
				Config:  map[string]interface{}{},
				Version: "0",
			}, nil
		}
		return nil, err
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(config.Config), &configMap); err != nil {
		return nil, err
	}

	return &dto.ConfigResponse{
		Config:  configMap,
		Version: config.Version,
	}, nil
}
