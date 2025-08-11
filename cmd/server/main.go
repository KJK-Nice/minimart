package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"minimart/internal/menu"
	"minimart/internal/merchant"
	"minimart/internal/notifications"
	"minimart/internal/order"
	"minimart/internal/shared/eventbus"
	middlerware "minimart/internal/shared/middleware"
	"minimart/internal/user"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	// --- Set up Structured Logger ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// --- Load Configuration ---
	// Load .env file (for local development)
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, continue without it")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Info("No config file found, using environment variables")
	}

	// --- Database Connection ---
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		dbUser := viper.GetString("database.user")
		dbName := viper.GetString("database.dbname")
		dbHost := viper.GetString("database.host")
		dbPort := viper.GetString("database.port")
		dbPassword := os.Getenv("DB_PASSWORD")
		connStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	}

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	app := fiber.New(fiber.Config{
		Network:      "tcp",
		ServerHeader: "Fiber",
		AppName:      "Minimart App v0.0.1",
	})

	redisAdrr := viper.GetString("redis.address")
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAdrr,
	})

	// Event bus
	eventBus := eventbus.NewRedisEventBus(redisClient)

	userSubscriber := notifications.NewUserSubscriber(logger)

	go func() {
		pubsub := redisClient.Subscribe(context.Background(), user.UserCreatedTopic)
		defer pubsub.Close()

		ch := pubsub.Channel()
		logger.Info("Subscribed to Redis topic", "topic", user.UserCreatedTopic)

		for msg := range ch {
			// When a message comes in, we handle it
			var event user.UserCreatedEvent
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				logger.Info("Error unmarshaling event", "error", err, "payload", msg.Payload)
				continue
			}

			// We call the same handler as before
			_ = userSubscriber.HandleUserCreatedEvent(context.Background(), event)
		}
	}()

	// Merchant module
	merchantRepo := merchant.NewPostgresMerchantRepository(dbpool)
	merchantUsecase := merchant.NewMerchantUsecase(merchantRepo)
	merchantHandler := merchant.NewMerchantHandler(merchantUsecase)
	merchantHandler.RegisterRoutes(app)

	// User module
	userRepo := user.NewPostgresUserRepository(dbpool)
	userUsecase := user.NewUserUsecase(userRepo, eventBus)
	userHandler := user.NewUserHandler(userUsecase)
	userHandler.RegisterRoutes(app)

	// Order module
	orderRepo := order.NewPostgresOrderRepository(dbpool)
	orderUsecase := order.NewOrderUsecase(orderRepo)
	orderHandler := order.NewOrderHandler(orderUsecase)
	orderHandler.RegisterRoutes(app)

	// Menu module
	menuRepo := menu.NewPostgresMenuRepository(dbpool)
	menuUsecase := menu.NewMenuUsecase(menuRepo)
	menuHandler := menu.NewMenuHandler(menuUsecase)
	menuHandler.RegisterRoutes(app)

	api := app.Group("/api", middlerware.AuthRequire())

	api.Get("/profile", func(c *fiber.Ctx) error {
		// The middlerware has already validated the token and stored the user claims.
		// We can safely access it from c.Locals.
		userClaims := c.Locals("user").(jwt.MapClaims)

		// You can now use the claims, for example, to fetch user details from the DB.
		// For this example, we'll just return the claims.
		return c.JSON(fiber.Map{
			"message": "Welcome to your profile!",
			"user_id": userClaims["sub"],
			"email":   userClaims["email"],
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	serverPort := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", serverPort)
	logger.Info("Starting server", "address", addr)
	if err := app.Listen(addr); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
