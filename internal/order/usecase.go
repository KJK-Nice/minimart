package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	// PlaceOrder creates a new order for a given customer with a list of items.
	// This is now a thin orchestration layer - business logic is in the entity
	PlaceOrder(ctx context.Context, customerID uuid.UUID, merchantID uuid.UUID, items []OrderItem, deliveryMethod DeliveryMethod, deliveryAddress *Address) (*Order, error)
	
	// GetOrderByID retrieves an order by its ID
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*Order, error)
	
	// AcceptOrder accepts an order (merchant action)
	AcceptOrder(ctx context.Context, orderID uuid.UUID, merchantID uuid.UUID, estimatedMinutes int) error
	
	// RejectOrder rejects an order (merchant action)
	RejectOrder(ctx context.Context, orderID uuid.UUID, merchantID uuid.UUID, reason string) error
}

type orderUsecase struct {
	repo OrderRepository
	// eventPublisher will be added later
}

func NewOrderUsecase(repo OrderRepository) OrderUsecase {
	return &orderUsecase{repo: repo}
}

// PlaceOrder is now a thin orchestration layer
func (u *orderUsecase) PlaceOrder(
	ctx context.Context,
	customerID uuid.UUID,
	merchantID uuid.UUID,
	items []OrderItem,
	deliveryMethod DeliveryMethod,
	deliveryAddress *Address,
) (*Order, error) {
	// Create order using the factory (business logic is in the entity)
	order, err := NewOrder(customerID, merchantID, items, deliveryMethod, deliveryAddress)
	if err != nil {
		return nil, err
	}
	
	// Save to repository
	if err := u.repo.Save(ctx, order); err != nil {
		return nil, err
	}
	
	// TODO: Publish events from order.Events()
	
	return order, nil
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*Order, error) {
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (u *orderUsecase) AcceptOrder(ctx context.Context, orderID uuid.UUID, merchantID uuid.UUID, estimatedMinutes int) error {
	// Load the order
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	
	// Verify merchant owns this order
	if order.MerchantID() != merchantID {
		return errors.New("unauthorized: merchant does not own this order")
	}
	
	// Call domain method (business logic is in the entity)
	events, err := order.Accept(estimatedMinutes, merchantID)
	if err != nil {
		return err
	}
	
	// Save the updated order
	if err := u.repo.Save(ctx, order); err != nil {
		return err
	}
	
	// TODO: Publish events
	_ = events
	
	return nil
}

func (u *orderUsecase) RejectOrder(ctx context.Context, orderID uuid.UUID, merchantID uuid.UUID, reason string) error {
	// Load the order
	order, err := u.repo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}
	
	// Verify merchant owns this order
	if order.MerchantID() != merchantID {
		return errors.New("unauthorized: merchant does not own this order")
	}
	
	// Call domain method (business logic is in the entity)
	events, err := order.Reject(reason, merchantID)
	if err != nil {
		return err
	}
	
	// Save the updated order
	if err := u.repo.Save(ctx, order); err != nil {
		return err
	}
	
	// TODO: Publish events
	_ = events
	
	return nil
}
