package main

import (
	"log"
	"minimart/internal/merchant"
	"minimart/internal/order"
	"minimart/internal/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Merchant module
	merchantRepo := merchant.NewInMemoryMerchantRepository()
	merchantUsecase := merchant.NewMerchantUsecase(merchantRepo)
	merchantHandler := merchant.NewMerchantHandler(merchantUsecase)
	merchantHandler.RegisterRoutes(app)

	// User module
	userRepo := user.NewInMemoryUserRepository()
	userUsecase := user.NewUserUsecase(userRepo)
	userHandler := user.NewUserHandler(userUsecase)
	userHandler.RegisterRoutes(app)

	// Order module
	orderRepo := order.NewInMemoryOrderRepository()
	orderUsecase := order.NewOrderUsecase(orderRepo)
	orderHandler := order.NewOrderHandler(orderUsecase)
	orderHandler.RegisterRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":3000"))
}
