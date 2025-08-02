package middlerware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to generate a valid JWT for testing
func generateTestToken(userID, email, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestAuthRequired(t *testing.T) {
	// Set a dummy JWT secret for testing
	viper.Set("jwt.secret", "test-secret")
	jwtSecret := viper.GetString("jwt.secret")

	// Create a new Fiber app for testing
	app := fiber.New()

	// Create a test route protected by the middlerware
	app.Get("/test", AuthRequire(), func(c *fiber.Ctx) error {
		// This handler should onlly be reached if the middleware succeeds
		userClaims := c.Locals("user").(jwt.MapClaims)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user_id": userClaims["sub"],
		})
	})

	t.Run("should return 200 OK with valid token", func(t *testing.T) {
		// Arrange
		token, err := generateTestToken("user-123", "test@example.com", jwtSecret)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Act
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return 401 Unauthorized without token", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 Unauthorized wiht invalid token", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-string")

		// Act
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return 401 Unauthorized with token signed by wrong key", func(t *testing.T) {
		// Arrange
		wrongSecret := "another-secret"
		token, err := generateTestToken("user-123", "test@example.com", wrongSecret)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		// Act
		resp, err := app.Test(req)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
