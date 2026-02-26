package usecase

import (
	"config-manager/internal/domain"
	"config-manager/internal/dto"
	"config-manager/internal/repository"
	"time"

	"github.com/google/uuid"
)

type AgentUsecase interface {
	Register(req dto.AgentRegisterRequest) (*dto.AgentRegisterResponse, error)
}

type agentUsecase struct {
	agentRepo    repository.AgentRepository
	pollURL      string
	pollInterval int
}

func NewAgentUsecase(agentRepo repository.AgentRepository, pollURL string, pollInterval int) AgentUsecase {
	return &agentUsecase{
		agentRepo:    agentRepo,
		pollURL:      pollURL,
		pollInterval: pollInterval,
	}
}

func (u *agentUsecase) Register(req dto.AgentRegisterRequest) (*dto.AgentRegisterResponse, error) {
	agent := &domain.Agent{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := u.agentRepo.Create(agent); err != nil {
		return nil, err
	}

	return &dto.AgentRegisterResponse{
		AgentID:             agent.ID,
		PollURL:             u.pollURL,
		PollIntervalSeconds: u.pollInterval,
	}, nil
}
