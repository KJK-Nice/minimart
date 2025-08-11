package main

import (
	"context"
	"database/sql"
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
	_ "github.com/jackc/pgx/v5/stdlib" // Goose requires a database driver
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	RedisURL    string `mapstructure:"REDIS_URL"`
	JwtSecret   string `mapstructure:"JWT_SECRET"`
}

func main() {
	// --- Set up Structured Logger ---
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// --- Load Configuration ---
	// Load .env file (for local development)
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, continue without it")
	}

	viper.AutomaticEnv()

	// Explicitly bind environment variables to viper keys
	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("JWT_SECRET")

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info("Config file not found, relying on environment variables.")
		} else {
			logger.Error("Error reading config file", "error", err)
		}
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		logger.Error("Unable to unmarshal configuration", "error", err)
	}

	// --- Log the loaded configuration for debugging ---
	logger.Info("Configuration loaded",
		"Port", config.Port,
		"DatabaseURL", config.DatabaseURL,
		"RedisURL", config.RedisURL,
		"JwtSecret", "...", // Don't log the secret itself
	)

	// --- Run Database Migrations ---
	migrationDb, err := sql.Open("pgx", config.DatabaseURL)
	if err != nil {
		logger.Error("Failed to open database for migrations", "error", err)
		os.Exit(1)
	}
	defer migrationDb.Close()

	goose.SetBaseFS(os.DirFS("."))
	goose.SetLogger(goose.NopLogger()) // Disables goose's default logger
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("Failed to set goose dialect", "error", err)
		os.Exit(1)
	}

	logger.Info("Running database migrations...")
	if err := goose.Up(migrationDb, "migrations"); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Database migrations completed successfully.")

	// --- Database Connection Pool for Application ---
	dbpool, err := pgxpool.New(context.Background(), config.DatabaseURL)
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

	// --- Initialize Redis Client ---
	// Parse Redis URL to handle authentication
	var redisOptions *redis.Options
	if config.RedisURL != "" {
		// Try to parse as URL first (for Railway or other cloud providers)
		if opts, err := redis.ParseURL(config.RedisURL); err == nil {
			redisOptions = opts
			logger.Info("Using Redis URL configuration")
		} else {
			// Fall back to simple address format (for local development)
			redisOptions = &redis.Options{
				Addr: config.RedisURL,
			}
			logger.Info("Using Redis simple address configuration", "address", config.RedisURL)
		}
	} else {
		// Default to localhost if not specified
		redisOptions = &redis.Options{
			Addr: "localhost:6379",
		}
		logger.Info("Using default Redis configuration")
	}

	redisClient := redis.NewClient(redisOptions)

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	logger.Info("Successfully connected to Redis")

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
	userUsecase := user.NewUserUsecase(userRepo, eventBus, config.JwtSecret)
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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	addr := fmt.Sprintf(":%s", config.Port)
	logger.Info("Configuration loaded", "port", config.Port, "database_url", config.DatabaseURL, "redis_url", config.RedisURL)
	if addr == "" {
		logger.Info("No port specified, using default port 3000")
		addr = ":3000" // Default to port 3000 if not specified
	}
	logger.Info("Starting server", "address", addr)
	if err := app.Listen(addr); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
