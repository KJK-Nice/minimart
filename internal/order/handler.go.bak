package order

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	usecase OrderUsecase
}

func NewOrderHandler(usecase OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

func (h *OrderHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/orders", h.PlaceOrder)
}

type PlaceOrderRequest struct {
	CustomerID uuid.UUID   `json:"customer_id"`
	Items      []OrderItem `json:"items"`
}

func (h *OrderHandler) PlaceOrder(c *fiber.Ctx) error {
	var req PlaceOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.CustomerID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing customer_id",
		})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order must contain at least one item",
		})
	}

	order, err := h.usecase.PlaceOrder(c.Context(), req.CustomerID, req.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(order)
}
