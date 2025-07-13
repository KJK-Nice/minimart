package order

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Items      []OrderItem
	Status     OrderStatus
	CreatedAt  time.Time
}

type OrderItem struct {
	MenuItemID uuid.UUID
	Quantity   int
}

type OrderStatus int

const (
	NEW OrderStatus = iota
	PENDING
	COMPLETED
	CANCELLED
)

func (s OrderStatus) String() string {
	return []string{"NEW", "PENDING", "COMPLETED", "CANCELLED"}[s]
}
