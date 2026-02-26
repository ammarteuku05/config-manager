package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestStaticTokenAuth(t *testing.T) {
	e := echo.New()

	// Setting up the middleware with expected token
	headerName := "X-API-KEY"
	expectedToken := "secret123"
	middlewareFunc := StaticTokenAuth(headerName, expectedToken)

	// A dummy handler that returns 200 OK
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// The wrapped handler
	h := middlewareFunc(handler)

	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(headerName, expectedToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(headerName, "wrong_token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h(c)

		assert.NoError(t, err) // Echo middleware returns the error as JSON successfully
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Header not set
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})
}
