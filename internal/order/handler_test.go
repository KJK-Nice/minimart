package order

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("could not terminate postgres container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(nil, err)

	dbpool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// --- Run Migrations in Order ---
	runMigration(ctx, "../../migrations/001_create_users_table.sql")
	runMigration(ctx, "../../migrations/002_create_orders_tables.sql")

	// 6. Run the actual tests
	exitCode := m.Run()

	// 7. Exit with the test result code
	os.Exit(exitCode)
}

func runMigration(ctx context.Context, filePath string) {
	migrationsPath, _ := filepath.Abs(filePath)
	migrationSQL, err := os.ReadFile(migrationsPath)
	if err != nil {
		log.Fatalf("could not read migration file: %s", err)
	}
	_, err = dbpool.Exec(ctx, string(migrationSQL))
	if err != nil {
		log.Fatalf("could not run migrations: %s", err)
	}
}

func TestOrderHandler_PlaceOrder_Integration(t *testing.T) {
	// Arrange
	// 1. Setup the application using the real Postgres repository
	orderRepo := NewPostgresOrderRepository(dbpool)
	orderUsecase := NewOrderUsecase(orderRepo)
	orderHandler := NewOrderHandler(orderUsecase)

	app := fiber.New()
	orderHandler.RegisterRoutes(app)

	// 2. Seed a user in ther database to act as ther customer
	customerID := uuid.New()
	_, err := dbpool.Exec(context.Background(), "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4);", customerID, "Test Customer", "customer@example.com", "password")
	require.NoError(t, err)

	// Act
	// 3. Create the HTTP request to place an order
	reqBody := PlaceOrderRequest{
		CustomerID: customerID,
		Items: []OrderItem{
			{MenuItemID: uuid.New(), Quantity: 2},
			{MenuItemID: uuid.New(), Quantity: 1},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Assert
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdOrder Order
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &createdOrder)
	require.NoError(t, err)

	assert.Equal(t, customerID, createdOrder.CustomerID)
	assert.Len(t, createdOrder.Items, 2)
	assert.Equal(t, NEW, createdOrder.Status)
	assert.NotEmpty(t, createdOrder.ID)
}
