package order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	// PlaceOrder creates a new order for a given customer with a list of items.
	PlaceOrder(ctx context.Context, customerID uuid.UUID, items []OrderItem) (*Order, error)
}

type orderUsecase struct {
	repo OrderRepository
}

func NewOrderUsecase(repo OrderRepository) OrderUsecase {
	return &orderUsecase{repo: repo}
}

func (u *orderUsecase) PlaceOrder(ctx context.Context, customerID uuid.UUID, items []OrderItem) (*Order, error) {
	order := &Order{
		ID:         uuid.New(),
		CustomerID: customerID,
		Items:      items,
		Status:     NEW,
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}
