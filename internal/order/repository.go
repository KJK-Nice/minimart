package order

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type OrderRepository interface {
	// Save creates or updates an order in the repository
	Save(ctx context.Context, order *Order) error
	
	// FindByID retrieves an order by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	
	// GetByID is deprecated, use FindByID
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	
	// FindByMerchantID retrieves all orders for a merchant
	FindByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Order, error)
	
	// FindPendingByMerchantID retrieves pending orders for a merchant
	FindPendingByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Order, error)
	
	// FindByCustomerID retrieves all orders for a customer
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Order, error)
}

type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: map[uuid.UUID]*Order{},
	}
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	return r.FindByID(ctx, id)
}

func (r *InMemoryOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	order, exists := r.orders[id]
	if !exists {
		return nil, nil
	}
	return order, nil
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.orders[order.ID()] = order
	return nil
}

func (r *InMemoryOrderRepository) FindByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var orders []*Order
	for _, order := range r.orders {
		if order.MerchantID() == merchantID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (r *InMemoryOrderRepository) FindPendingByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var orders []*Order
	for _, order := range r.orders {
		if order.MerchantID() == merchantID && order.Status() == OrderStatusPending {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (r *InMemoryOrderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var orders []*Order
	for _, order := range r.orders {
		if order.CustomerID() == customerID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
