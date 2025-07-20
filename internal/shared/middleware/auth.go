package middlerware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// AuthRequired is a middleware to protect routes that require a valid JWT.
func AuthRequire() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		// 2. Check for "Bearer " prefix and extract the token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}
		tokenString := parts[1]

		// 3. Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpedted signing method")
			}

			// Return the secret key
			return []byte(viper.GetString("jwt.secret")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		}

		// 4. Extract claims and store them in context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT claims",
			})
		}

		// You can store the entire claims map or specific values
		c.Locals("user", claims)

		// 5. Continue to the next handler
		return c.Next()
	}
}
