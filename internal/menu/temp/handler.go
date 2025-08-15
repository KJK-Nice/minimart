package menu

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MenuHandler is reposible for handling HTTP requests from menu items.
type MenuHandler struct {
	usecase MenuUsecase
}

// NewMenuHandler creates a new instance of MenuHandler.
func NewMenuHandler(usecase MenuUsecase) *MenuHandler {
	return &MenuHandler{
		usecase: usecase,
	}
}

// RegisterRoutes adds the menu routes to the Fiber app.
func (h *MenuHandler) RegisterRoutes(app *fiber.App) {
	// Group routes for a specific merchant's menu
	menuRoutes := app.Group("/merchants/:merchantID/menu")
	menuRoutes.Post("/", h.CreateMenuItem)
	menuRoutes.Get("/", h.GetMenuForMerchant)
}

// CreateMenuITemRequest defines the JSON request body for creating a menu item.
type CreateMenuItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

// CreateMenuItem handles the creation of a new menu item.
func (h *MenuHandler) CreateMenuItem(c *fiber.Ctx) error {
	merchantID, err := uuid.Parse(c.Params("merchantID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid merchant ID"})
	}

	var req CreateMenuItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	item, err := h.usecase.CreateMenuItem(c.Context(), merchantID, req.Name, req.Description, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// GetMenuForMerchant handles fetching the menu for a specific merchant.
func (h *MenuHandler) GetMenuForMerchant(c *fiber.Ctx) error {
	merchantID, err := uuid.Parse(c.Params("merchantID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid merchant ID"})
	}

	items, err := h.usecase.GetMenuForMerchant(c.Context(), merchantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(items)
}
