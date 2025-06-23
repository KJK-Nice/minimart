package user

import (
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase UserUsecase
}

func NewUserHandler(usecase UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Post("users/register", h.RegisterUser)
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	user, err := h.usecase.RegisterUser(c.Context(), req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}
