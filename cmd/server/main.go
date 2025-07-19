package main

import (
	"log"
	"minimart/internal/menu"
	"minimart/internal/merchant"
	"minimart/internal/notifications"
	"minimart/internal/order"
	"minimart/internal/shared/eventbus"
	"minimart/internal/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Event bus
	eventBus := eventbus.NewInMemoryEventBus()

	userSubscriber := notifications.NewUserSubscriber()
	err := eventBus.Subscribe(user.UserCreatedTopic, userSubscriber.HandleUserCreatedEvent)
	if err != nil {
		log.Fatalf("Failed to subscribe to user created event: %v", err)
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
