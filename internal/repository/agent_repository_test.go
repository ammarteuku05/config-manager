package repository

import (
	"config-manager/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&domain.Agent{}, &domain.GlobalConfig{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	return db
}

func TestAgentRepository_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAgentRepository(db)

	agent := &domain.Agent{
		ID:        "agent-123",
		Name:      "test-agent",
		CreatedAt: time.Now(),
	}

	t.Run("Create Success", func(t *testing.T) {
		err := repo.Create(agent)
		assert.NoError(t, err)
	})

	t.Run("GetByID Success", func(t *testing.T) {
		fetchedAgent, err := repo.GetByID(agent.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetchedAgent)
		assert.Equal(t, agent.ID, fetchedAgent.ID)
		assert.Equal(t, agent.Name, fetchedAgent.Name)
	})

	t.Run("GetByID Not Found", func(t *testing.T) {
		fetchedAgent, err := repo.GetByID("invalid-id")
		assert.Error(t, err)
		assert.Nil(t, fetchedAgent)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
