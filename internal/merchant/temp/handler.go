package merchant

import (
	"github.com/gofiber/fiber/v2"
)

type MerchantHandler struct {
	usecase MerchantUsecase
}

func NewMerchantHandler(usecase MerchantUsecase) *MerchantHandler {
	return &MerchantHandler{
		usecase: usecase,
	}
}

func (h *MerchantHandler) RegisterRoutes(app *fiber.App) {
	app.Post("merchants/register", h.CreateMerchant)
}

func (h *MerchantHandler) CreateMerchant(c *fiber.Ctx) error {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	user, err := h.usecase.CreateMerchant(c.Context(), req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}
