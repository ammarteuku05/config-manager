package usecase

import (
	"config-manager/internal/domain"
	"config-manager/internal/dto"
	"config-manager/internal/repository/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestConfigUsecase_Save(t *testing.T) {
	mockRepo := new(mocks.MockConfigRepository)
	uc := NewConfigUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		req := dto.ConfigRequest{
			Config: map[string]interface{}{"url": "http://example.com"},
		}
		mockRepo.On("Save", mock.AnythingOfType("*domain.GlobalConfig")).Return(nil).Once()

		err := uc.Save(req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("JSON Encode Error", func(t *testing.T) {
		// math.NaN is an invalid JSON value
		req := dto.ConfigRequest{
			Config: map[string]interface{}{"invalid": make(chan int)},
		}
		err := uc.Save(req)
		assert.Error(t, err)
	})

	t.Run("DB Error", func(t *testing.T) {
		req := dto.ConfigRequest{
			Config: map[string]interface{}{"url": "http://example.com"},
		}
		mockRepo.On("Save", mock.AnythingOfType("*domain.GlobalConfig")).Return(errors.New("db error")).Once()

		err := uc.Save(req)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestConfigUsecase_GetLatest(t *testing.T) {
	mockRepo := new(mocks.MockConfigRepository)
	uc := NewConfigUsecase(mockRepo)

	t.Run("Success", func(t *testing.T) {
		expectedConfig := &domain.GlobalConfig{
			Config:    `{"url":"http://example.com"}`,
			Version:   "v1.0",
			CreatedAt: time.Now(),
		}
		mockRepo.On("GetLatest").Return(expectedConfig, nil).Once()

		res, err := uc.GetLatest()
		assert.NoError(t, err)
		assert.Equal(t, "v1.0", res.Version)
		assert.Equal(t, "http://example.com", res.Config["url"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("GetLatest").Return(nil, gorm.ErrRecordNotFound).Once()

		res, err := uc.GetLatest()
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "0", res.Version)
		assert.Empty(t, res.Config)
		mockRepo.AssertExpectations(t)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockRepo.On("GetLatest").Return(nil, errors.New("db error")).Once()

		res, err := uc.GetLatest()
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		expectedConfig := &domain.GlobalConfig{
			Config:  `invalid json`,
			Version: "v1.0",
		}
		mockRepo.On("GetLatest").Return(expectedConfig, nil).Once()

		res, err := uc.GetLatest()
		assert.Error(t, err)
		assert.Nil(t, res)
		mockRepo.AssertExpectations(t)
	})
}
