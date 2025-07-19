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
	"minimart/internal/user"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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
		logger.Error("Error reading config file", "error", err)
		os.Exit(1)
	}

	app := fiber.New()

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
	merchantRepo := merchant.NewInMemoryMerchantRepository()
	merchantUsecase := merchant.NewMerchantUsecase(merchantRepo)
	merchantHandler := merchant.NewMerchantHandler(merchantUsecase)
	merchantHandler.RegisterRoutes(app)

	// User module
	userRepo := user.NewInMemoryUserRepository()
	userUsecase := user.NewUserUsecase(userRepo, eventBus)
	userHandler := user.NewUserHandler(userUsecase)
	userHandler.RegisterRoutes(app)

	// Order module
	orderRepo := order.NewInMemoryOrderRepository()
	orderUsecase := order.NewOrderUsecase(orderRepo)
	orderHandler := order.NewOrderHandler(orderUsecase)
	orderHandler.RegisterRoutes(app)

	// Menu module
	menuRepo := menu.NewInMemoryMenuRepository()
	menuUsecase := menu.NewMenuUsecase(menuRepo)
	menuHandler := menu.NewMenuHandler(menuUsecase)
	menuHandler.RegisterRoutes(app)

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
