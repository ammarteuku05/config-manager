package repository

import (
	"config-manager/internal/domain"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	Save(config *domain.GlobalConfig) error
	GetLatest() (*domain.GlobalConfig, error)
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) Save(config *domain.GlobalConfig) error {
	return r.db.Create(config).Error
}

func (r *configRepository) GetLatest() (*domain.GlobalConfig, error) {
	var config domain.GlobalConfig
	if err := r.db.Order("created_at desc").First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}
