package main

import (
	"context"
	"fmt"
	"log"
	"minimart/internal/menu"
	"minimart/internal/merchant"
	"minimart/internal/order"
	"minimart/internal/shared/eventbus"
	"minimart/internal/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Event bus
	eventBus := eventbus.NewInMemoryEventBus()

	// Simple subscriber for testing
	err := eventBus.Subscribe(user.UserCreatedTopic, func(ctx context.Context, event eventbus.Event) error {
		if userEvent, ok := event.(user.UserCreatedEvent); ok {
			fmt.Printf("New user created: %+v\n", userEvent)
		}
		return nil

	})
	if err != nil {
		log.Fatal("Failed to subscribe to user created event: %v", err)
	}

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

	log.Fatal(app.Listen(":3000"))
}
