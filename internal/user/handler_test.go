package user

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var dbpool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Define the PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		log.Fatalf("could not start Postgres container: %s", err)
	}

	// 2. Set up a teardown function to be called when tests are done
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("could not terminate postgres container: %s", err)
		}
	}()

	// 3. Get the connection URL for the PostgreSQL container
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(nil, err)

	// 4. Create a connection pool to the PostgreSQL database
	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// 5. Run the database migrations
	migrationsPath, _ := filepath.Abs("../../migrations/001_create_users_table.sql")
	migrationSQL, err := os.ReadFile(migrationsPath)
	if err != nil {
		log.Fatalf("could not read migration file: %s", err)
	}
	_, err = dbpool.Exec(ctx, string(migrationSQL))
	if err != nil {
		log.Fatalf("could not run migrations: %s", err)
	}

	// 6. Run the actual tests
	exitCode := m.Run()

	// 7. Exit with the test result code
	os.Exit(exitCode)
}

func TestUserHandler_RegisterUser_Integration(t *testing.T) {
	// 1. Arrange: Set up our application and dependencies
	// userRepo := NewInMemoryUserRepository()
	userRepo := NewPostgresUserRepository(dbpool)
	userUsecase := NewUserUsecase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	// Create a new Fiber app for testing
	app := fiber.New()
	userHandler.RegisterRoutes(app)

	// 2. Act: Create the HTTP request
	// Create the request body
	reqBody := map[string]string{
		"name":  "Test User",
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
