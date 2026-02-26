package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// StaticTokenAuth checks a specific header for a specific token
func StaticTokenAuth(headerName, token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get(headerName) != token {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}
			return next(c)
		}
	}
}
