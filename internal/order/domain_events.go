package order

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is a marker interface for all domain events
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// OrderPlacedEvent is emitted when a new order is placed
type OrderPlacedEvent struct {
	OrderID        uuid.UUID
	CustomerID     uuid.UUID
	MerchantID     uuid.UUID
	TotalAmount    Money
	DeliveryMethod DeliveryMethod
	PlacedAt       time.Time
}

func (e OrderPlacedEvent) EventName() string    { return "order.placed" }
func (e OrderPlacedEvent) OccurredAt() time.Time { return e.PlacedAt }

// OrderAcceptedEvent is emitted when an order is accepted by the merchant
type OrderAcceptedEvent struct {
	OrderID       uuid.UUID
	MerchantID    uuid.UUID
	CustomerID    uuid.UUID
	EstimatedTime time.Time
	AcceptedAt    time.Time
}

func (e OrderAcceptedEvent) EventName() string    { return "order.accepted" }
func (e OrderAcceptedEvent) OccurredAt() time.Time { return e.AcceptedAt }

// OrderRejectedEvent is emitted when an order is rejected by the merchant
type OrderRejectedEvent struct {
	OrderID    uuid.UUID
	MerchantID uuid.UUID
	CustomerID uuid.UUID
	Reason     string
	RejectedAt time.Time
}

func (e OrderRejectedEvent) EventName() string    { return "order.rejected" }
func (e OrderRejectedEvent) OccurredAt() time.Time { return e.RejectedAt }

// OrderPreparingEvent is emitted when order preparation starts
type OrderPreparingEvent struct {
	OrderID    uuid.UUID
	MerchantID uuid.UUID
	CustomerID uuid.UUID
	StartedAt  time.Time
}

func (e OrderPreparingEvent) EventName() string    { return "order.preparing" }
func (e OrderPreparingEvent) OccurredAt() time.Time { return e.StartedAt }

// OrderReadyEvent is emitted when an order is ready for pickup/delivery
type OrderReadyEvent struct {
	OrderID        uuid.UUID
	MerchantID     uuid.UUID
	CustomerID     uuid.UUID
	DeliveryMethod DeliveryMethod
	ReadyAt        time.Time
}

func (e OrderReadyEvent) EventName() string    { return "order.ready" }
func (e OrderReadyEvent) OccurredAt() time.Time { return e.ReadyAt }

// OrderOutForDeliveryEvent is emitted when an order is out for delivery
type OrderOutForDeliveryEvent struct {
	OrderID      uuid.UUID
	CustomerID   uuid.UUID
	DriverID     uuid.UUID
	Address      *Address
	DispatchedAt time.Time
}

func (e OrderOutForDeliveryEvent) EventName() string    { return "order.out_for_delivery" }
func (e OrderOutForDeliveryEvent) OccurredAt() time.Time { return e.DispatchedAt }

// OrderCompletedEvent is emitted when an order is completed
type OrderCompletedEvent struct {
	OrderID     uuid.UUID
	MerchantID  uuid.UUID
	CustomerID  uuid.UUID
	CompletedAt time.Time
}

func (e OrderCompletedEvent) EventName() string    { return "order.completed" }
func (e OrderCompletedEvent) OccurredAt() time.Time { return e.CompletedAt }

// OrderCancelledEvent is emitted when an order is cancelled
type OrderCancelledEvent struct {
	OrderID     uuid.UUID
	MerchantID  uuid.UUID
	CustomerID  uuid.UUID
	Reason      string
	CancelledBy uuid.UUID
	CancelledAt time.Time
}

func (e OrderCancelledEvent) EventName() string    { return "order.cancelled" }
func (e OrderCancelledEvent) OccurredAt() time.Time { return e.CancelledAt }
