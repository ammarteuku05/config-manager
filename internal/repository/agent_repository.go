package repository

import (
	"config-manager/internal/domain"

	"gorm.io/gorm"
)

type AgentRepository interface {
	Create(agent *domain.Agent) error
	GetByID(id string) (*domain.Agent, error)
}

type agentRepository struct {
	db *gorm.DB
}

func NewAgentRepository(db *gorm.DB) AgentRepository {
	return &agentRepository{db: db}
}

func (r *agentRepository) Create(agent *domain.Agent) error {
	return r.db.Create(agent).Error
}

func (r *agentRepository) GetByID(id string) (*domain.Agent, error) {
	var agent domain.Agent
	if err := r.db.First(&agent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}
