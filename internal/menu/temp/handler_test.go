package menu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minimart/internal/merchant"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	runMigration(ctx, "../../migrations/004_create_merchants_table.sql")
	runMigration(ctx, "../../migrations/003_create_menu_items_table.sql")

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

func TestMenuHandler_Integration(t *testing.T) {
	// Arrange: Set up a full Fiber app with both Merchant and Menu handlers
	app := fiber.New()

	// Merchant dependencies
	merchantRepo := merchant.NewPostgresMerchantRepository(dbpool)
	merchantUsecase := merchant.NewMerchantUsecase(merchantRepo)
	merchantHandler := merchant.NewMerchantHandler(merchantUsecase)
	merchantHandler.RegisterRoutes(app)

	// Menu dependencies
	menuRepo := NewPostgresMenuRepository(dbpool)
	menuUsecase := NewMenuUsecase(menuRepo)
	menuHandler := NewMenuHandler(menuUsecase)
	menuHandler.RegisterRoutes(app)

	// --- Seed a merchant to be the owner of the menu items ---
	seededMerchant := merchant.NewMerchant("The Berger Joint", "Best burgers in town")
	err := merchantRepo.Save(context.Background(), seededMerchant)
	require.NoError(t, err)

	// --- Test Case 1: Create a new menu item ---
	t.Run("should create a new menu item", func(t *testing.T) {
		// Act
		reqBody := CreateMenuItemRequest{
			Name:        "Classic Burger",
			Description: "A delicious classic burger",
			Price:       1000,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		// Note the URL includes the seeded merchant ID
		url := fmt.Sprintf("/merchants/%s/menu", seededMerchant.ID)
		req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Assert
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var createdItem MenuItem
		respBody, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(respBody, &createdItem)
		require.NoError(t, err)

		assert.Equal(t, "Classic Burger", createdItem.Name)
		assert.Equal(t, 1000, createdItem.Price)
		assert.Equal(t, seededMerchant.ID, createdItem.MerchantID)
	})

	// --- Test Case 2: Get all menu items for the merchant ---
	t.Run("should get all menu items for a merchant", func(t *testing.T) {
		// Act
		url := fmt.Sprintf("/merchants/%s/menu", seededMerchant.ID)
		req := httptest.NewRequest(http.MethodGet, url, nil)

		// Assert
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var items []*MenuItem
		respBody, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(respBody, &items)
		require.NoError(t, err)

		require.Len(t, items, 1)
		assert.Equal(t, "Classic Burger", items[0].Name)
	})

}
