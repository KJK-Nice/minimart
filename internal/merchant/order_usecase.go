package merchant

import (
	"context"
	"errors"

	"minimart/internal/order"

	"github.com/google/uuid"
)

// MerchantOrderUsecase provides merchant-focused order management operations
type MerchantOrderUsecase interface {
	// GetPendingOrders retrieves all pending orders for a merchant
	GetPendingOrders(ctx context.Context, merchantID uuid.UUID) ([]*order.Order, error)

	// GetOrdersByStatus retrieves orders for a merchant filtered by status
	GetOrdersByStatus(ctx context.Context, merchantID uuid.UUID, status order.OrderStatus) ([]*order.Order, error)

	// GetAllOrders retrieves all orders for a merchant (across all statuses)
	GetAllOrders(ctx context.Context, merchantID uuid.UUID) ([]*order.Order, error)

	// AcceptOrderWithEstimate accepts an order with merchant's preparation time estimate
	AcceptOrderWithEstimate(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, itemCount int) error

	// AcceptOrderWithCustomTime accepts an order with a custom estimated time
	AcceptOrderWithCustomTime(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, estimatedMinutes int) error

	// RejectOrder rejects an order with a reason
	RejectOrder(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, reason string) error

	// UpdateOrderStatus updates the status of an order through the workflow
	UpdateOrderStatus(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, newStatus order.OrderStatus) error

	// GetMerchantStats retrieves order statistics for a merchant
	GetMerchantStats(ctx context.Context, merchantID uuid.UUID) (*MerchantOrderStats, error)
}

// MerchantOrderStats provides aggregated order statistics for a merchant
type MerchantOrderStats struct {
	TotalOrders     int
	PendingOrders   int
	AcceptedOrders  int
	PreparingOrders int
	CompletedOrders int
	RejectedOrders  int
	CancelledOrders int

	// Revenue statistics (in Satoshis)
	TotalRevenue     int64
	PendingRevenue   int64
	CompletedRevenue int64

	// Average preparation time
	AveragePreparationTimeMinutes float64
}

// merchantOrderUsecase implements MerchantOrderUsecase
type merchantOrderUsecase struct {
	merchantRepo MerchantRepository // Will be created later
	orderUsecase order.OrderUsecase
}

// NewMerchantOrderUsecase creates a new merchant order use case
func NewMerchantOrderUsecase(orderUsecase order.OrderUsecase) MerchantOrderUsecase {
	return &merchantOrderUsecase{
		orderUsecase: orderUsecase,
		// merchantRepo will be injected later when we create the repository
	}
}

func (u *merchantOrderUsecase) GetPendingOrders(ctx context.Context, merchantID uuid.UUID) ([]*order.Order, error) {
	// Get all orders for merchant
	orders, err := u.orderUsecase.GetOrdersByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	// Filter for pending orders
	var pendingOrders []*order.Order
	for _, ord := range orders {
		if ord.Status() == order.OrderStatusPending {
			pendingOrders = append(pendingOrders, ord)
		}
	}

	return pendingOrders, nil
}

func (u *merchantOrderUsecase) GetOrdersByStatus(ctx context.Context, merchantID uuid.UUID, status order.OrderStatus) ([]*order.Order, error) {
	// Get all orders for merchant
	orders, err := u.orderUsecase.GetOrdersByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	// Filter by status
	var filteredOrders []*order.Order
	for _, ord := range orders {
		if ord.Status() == status {
			filteredOrders = append(filteredOrders, ord)
		}
	}

	return filteredOrders, nil
}

func (u *merchantOrderUsecase) GetAllOrders(ctx context.Context, merchantID uuid.UUID) ([]*order.Order, error) {
	return u.orderUsecase.GetOrdersByMerchantID(ctx, merchantID)
}

func (u *merchantOrderUsecase) AcceptOrderWithEstimate(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, itemCount int) error {
	// For now, we'll use a simple estimate: 30 minutes base + 5 minutes per item
	// Later this could be enhanced with merchant-specific logic
	estimatedMinutes := 30 + (itemCount * 5)

	return u.orderUsecase.AcceptOrder(ctx, orderID, merchantID, estimatedMinutes)
}

func (u *merchantOrderUsecase) AcceptOrderWithCustomTime(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, estimatedMinutes int) error {
	if estimatedMinutes < 1 {
		return errors.New("estimated minutes must be at least 1")
	}
	if estimatedMinutes > 480 { // 8 hours max
		return errors.New("estimated minutes cannot exceed 480 (8 hours)")
	}

	return u.orderUsecase.AcceptOrder(ctx, orderID, merchantID, estimatedMinutes)
}

func (u *merchantOrderUsecase) RejectOrder(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, reason string) error {
	if reason == "" {
		return errors.New("rejection reason is required")
	}

	return u.orderUsecase.RejectOrder(ctx, orderID, merchantID, reason)
}

func (u *merchantOrderUsecase) UpdateOrderStatus(ctx context.Context, merchantID uuid.UUID, orderID uuid.UUID, newStatus order.OrderStatus) error {
	// Route to appropriate use case method based on status
	switch newStatus {
	case order.OrderStatusPreparing:
		return u.orderUsecase.StartPreparing(ctx, orderID, merchantID)
	case order.OrderStatusReady:
		return u.orderUsecase.MarkReady(ctx, orderID, merchantID)
	case order.OrderStatusOutForDelivery:
		return u.orderUsecase.MarkOutForDelivery(ctx, orderID, merchantID)
	case order.OrderStatusCompleted:
		return u.orderUsecase.CompleteOrder(ctx, orderID, merchantID)
	case order.OrderStatusCancelled:
		return u.orderUsecase.CancelOrder(ctx, orderID, merchantID, "Cancelled by merchant")
	default:
		return errors.New("invalid status transition")
	}
}

func (u *merchantOrderUsecase) GetMerchantStats(ctx context.Context, merchantID uuid.UUID) (*MerchantOrderStats, error) {
	// Get all orders for merchant
	orders, err := u.orderUsecase.GetOrdersByMerchantID(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	stats := &MerchantOrderStats{}

	var totalPreparationMinutes int64
	var completedOrdersWithTime int

	// Calculate statistics from orders
	for _, ord := range orders {
		stats.TotalOrders++

		// Count by status
		switch ord.Status() {
		case order.OrderStatusPending:
			stats.PendingOrders++
			stats.PendingRevenue += ord.TotalAmount().Amount()
		case order.OrderStatusAccepted:
			stats.AcceptedOrders++
		case order.OrderStatusPreparing:
			stats.PreparingOrders++
		case order.OrderStatusCompleted:
			stats.CompletedOrders++
			stats.CompletedRevenue += ord.TotalAmount().Amount()

			// Calculate average preparation time for completed orders
			if ord.EstimatedWindow() != nil {
				// This is a simple approximation - in reality you'd want actual completion times
				totalPreparationMinutes += int64(ord.EstimatedWindow().DurationMinutes())
				completedOrdersWithTime++
			}
		case order.OrderStatusRejected:
			stats.RejectedOrders++
		case order.OrderStatusCancelled:
			stats.CancelledOrders++
		}

		// Calculate total revenue (completed + pending)
		stats.TotalRevenue += ord.TotalAmount().Amount()
	}

	// Calculate average preparation time
	if completedOrdersWithTime > 0 {
		stats.AveragePreparationTimeMinutes = float64(totalPreparationMinutes) / float64(completedOrdersWithTime)
	}

	return stats, nil
}
