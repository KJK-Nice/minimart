package main

import (
	"fmt"
	"log"
	"minimart/pages"
	"minimart/types"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		Network:      "tcp",
		ServerHeader: "Fiber",
		AppName:      "Minimart Test App v0.0.1",
	})

	// --- Setup Middleware ---
	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// --- Static File Serving ---
	app.Static("/static", "./static")

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Home page - test without logged-in user
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		var user *types.User = nil
		return pages.Home(user).Render(c.Context(), c.Response().BodyWriter())
	})

	// Home page - test with logged-in user
	app.Get("/user", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		user := &types.User{
			ID:       "user123",
			Username: "satoshi",
			Role:     "customer",
		}
		return pages.Home(user).Render(c.Context(), c.Response().BodyWriter())
	})

	// Home page - test with merchant user
	app.Get("/merchant", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		user := &types.User{
			ID:       "merchant456",
			Username: "pizza_place",
			Role:     "merchant",
		}
		return pages.Home(user).Render(c.Context(), c.Response().BodyWriter())
	})

	port := ":8080"
	fmt.Printf("Starting test server on http://localhost%s\n", port)
	fmt.Println("Routes available:")
	fmt.Println("  / - Home page (no user)")
	fmt.Println("  /user - Home page (customer)")
	fmt.Println("  /merchant - Home page (merchant)")
	fmt.Println("  /health - Health check")
	
	log.Fatal(app.Listen(port))
}
