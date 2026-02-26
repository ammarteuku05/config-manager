package usecase

import (
	"config-manager/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigManager_UpdateAndHit(t *testing.T) {
	// Create mock target server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mocked response"))
	}))
	defer ts.Close()

	cm := NewConfigManager()

	t.Run("Hit Without Config", func(t *testing.T) {
		_, err := cm.ExecuteHit()
		assert.Error(t, err)
		assert.Equal(t, "URL not configured", err.Error())
	})

	t.Run("Hit With Invalid Config Type", func(t *testing.T) {
		req := dto.ConfigRequest{
			Config: map[string]interface{}{
				"url": 123, // not a string
			},
		}
		err := cm.UpdateConfig(req)
		assert.NoError(t, err)

		_, err = cm.ExecuteHit()
		assert.Error(t, err)
		assert.Equal(t, "configured URL is not a string", err.Error())
	})

	t.Run("Hit Failure Target Unavailable", func(t *testing.T) {
		req := dto.ConfigRequest{
			Config: map[string]interface{}{
				"url": "http://localhost:1", // unavailable port
			},
		}
		err := cm.UpdateConfig(req)
		assert.NoError(t, err)

		cm.(*configManager).httpClient.Timeout = 50 * time.Millisecond // fast fail
		_, err = cm.ExecuteHit()
		assert.Error(t, err)
	})

	t.Run("Hit Success", func(t *testing.T) {
		req := dto.ConfigRequest{
			Config: map[string]interface{}{
				"url": ts.URL,
			},
		}
		err := cm.UpdateConfig(req)
		assert.NoError(t, err)

		res, err := cm.ExecuteHit()
		assert.NoError(t, err)
		assert.Equal(t, "mocked response", res)
	})
}
