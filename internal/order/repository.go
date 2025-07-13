package order

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	// Save creates or updates an order in the repositroy
	Save(ctx context.Context, order *Order) error

	// GetByID retrieves an order by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
}

type InMemoryOrderRepository struct {
	orders map[uuid.UUID]*Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: map[uuid.UUID]*Order{},
	}
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, nil
	}
	return order, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *Order) error {
	r.orders[order.ID] = order
	return nil
}
