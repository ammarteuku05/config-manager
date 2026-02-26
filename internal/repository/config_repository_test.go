package repository

import (
	"config-manager/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigRepository_SaveAndGetLatest(t *testing.T) {
	db := setupTestDB(t) // Reuse setup from agent_repository_test.go
	repo := NewConfigRepository(db)

	cfg1 := &domain.GlobalConfig{
		Version:   "v1",
		Config:    `{"poll_interval": 30}`,
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}

	cfg2 := &domain.GlobalConfig{
		Version:   "v2",
		Config:    `{"poll_interval": 60}`,
		CreatedAt: time.Now(),
	}

	t.Run("GetLatest Not Found", func(t *testing.T) {
		fetchedConfig, err := repo.GetLatest()
		assert.Error(t, err)
		assert.Nil(t, fetchedConfig)
	})

	t.Run("Save Success", func(t *testing.T) {
		err := repo.Save(cfg1)
		assert.NoError(t, err)

		err = repo.Save(cfg2)
		assert.NoError(t, err)
	})

	t.Run("GetLatest Returns Newest", func(t *testing.T) {
		fetchedConfig, err := repo.GetLatest()
		assert.NoError(t, err)
		assert.NotNil(t, fetchedConfig)
		assert.Equal(t, "v2", fetchedConfig.Version)
	})
}
