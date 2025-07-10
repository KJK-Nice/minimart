package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_RegisterUser(t *testing.T) {
	// 1. Arrange: Set up our application and dependencies
	userRepo := NewInMemoryUserRepository()
	userUsecase := NewUserUsecase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	// Create a new Fiber app for testing
	app := fiber.New()
	userHandler.RegisterRoutes(app)

	// 2. Act: Create the HTTP request
	// Create the request body
	reqBody := map[string]string{
		"name": "Test User",
		"email": "test@example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Create the POST request
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// 3. Assert: Perform the request and check the response
	// The app. Test function sends the request to the app and returns the response
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Check the status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Check the response body
	repsBody, _ := io.ReadAll(resp.Body)
	var createdUser User
	err = json.Unmarshal(repsBody, &createdUser)
	require.NoError(t, err)

	assert.NotEmpty(t, createdUser.ID, "Expected user ID to be generated")
	assert.Equal(t, "Test User", createdUser.Name, "Expected user name to match")
	assert.Equal(t, "test@example.com", createdUser.Email, "Expected user email to match")
}
