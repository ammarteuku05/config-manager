package usecase

import (
	"config-manager/internal/repository/mocks"
	"errors"

	"config-manager/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAgentUsecase_Register(t *testing.T) {
	mockRepo := new(mocks.MockAgentRepository)
	uc := NewAgentUsecase(mockRepo, "/config", 30)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.Agent")).Return(nil).Once()

		req := dto.AgentRegisterRequest{Name: "TestAgent"}
		res, err := uc.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AgentID)
		assert.Equal(t, "/config", res.PollURL)
		assert.Equal(t, 30, res.PollIntervalSeconds)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockRepo.On("Create", mock.AnythingOfType("*domain.Agent")).Return(errors.New("db error")).Once()

		req := dto.AgentRegisterRequest{Name: "TestAgent"}
		res, err := uc.Register(req)

		assert.Error(t, err)
		assert.Nil(t, res)

		mockRepo.AssertExpectations(t)
	})
}
