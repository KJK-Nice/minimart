package order

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrOrderNotPending        = errors.New("order must be pending to perform this action")
	ErrEmptyOrder             = errors.New("order must have at least one item")
	ErrInvalidQuantity        = errors.New("quantity must be greater than zero")
	ErrMissingMerchant        = errors.New("merchant ID is required")
	ErrMissingCustomer        = errors.New("customer ID is required")
	ErrInvalidDeliveryMethod  = errors.New("invalid delivery method")
	ErrDeliveryAddressRequired = errors.New("delivery address required for delivery orders")
)

// Order is the aggregate root for the order domain
type Order struct {
	id              uuid.UUID
	customerID      uuid.UUID
	merchantID      uuid.UUID
	items           []OrderItem
	status          OrderStatus
	totalAmount     Money
	deliveryMethod  DeliveryMethod
	deliveryAddress *Address
	estimatedWindow *TimeWindow
	createdAt       time.Time
	updatedAt       time.Time
	statusHistory   []StatusChange
	
	// Domain events to be published
	events []DomainEvent
}

// OrderItem represents a line item in an order
type OrderItem struct {
	MenuItemID   uuid.UUID
	MenuItemName string  // Snapshot of name at order time
	Quantity     int
	PricePerItem Money   // Snapshot of price at order time
}

// OrderStatus represents the current state of an order
type OrderStatus int

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusAccepted
	OrderStatusRejected
	OrderStatusPreparing
	OrderStatusReady
	OrderStatusOutForDelivery
	OrderStatusCompleted
	OrderStatusCancelled
)

func (s OrderStatus) String() string {
	statuses := []string{
		"PENDING",
		"ACCEPTED",
		"REJECTED",
		"PREPARING",
		"READY",
		"OUT_FOR_DELIVERY",
		"COMPLETED",
		"CANCELLED",
	}
	if int(s) < len(statuses) {
		return statuses[s]
	}
	return "UNKNOWN"
}

// StatusChange records a status transition
type StatusChange struct {
	From      OrderStatus
	To        OrderStatus
	Reason    string
	ChangedAt time.Time
	ChangedBy uuid.UUID // Could be customer, merchant, or system
}

// State machine for valid transitions
var validTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:        {OrderStatusAccepted, OrderStatusRejected, OrderStatusCancelled},
	OrderStatusAccepted:       {OrderStatusPreparing, OrderStatusCancelled},
	OrderStatusPreparing:      {OrderStatusReady, OrderStatusCancelled},
	OrderStatusReady:          {OrderStatusOutForDelivery, OrderStatusCompleted, OrderStatusCancelled},
	OrderStatusOutForDelivery: {OrderStatusCompleted, OrderStatusCancelled},
	// Terminal states
	OrderStatusCompleted: {},
	OrderStatusRejected:  {},
	OrderStatusCancelled: {},
}

// NewOrder creates a new order with validation
func NewOrder(
	customerID uuid.UUID,
	merchantID uuid.UUID,
	items []OrderItem,
	deliveryMethod DeliveryMethod,
	deliveryAddress *Address,
) (*Order, error) {
	// Validate inputs
	if customerID == uuid.Nil {
		return nil, ErrMissingCustomer
	}
	if merchantID == uuid.Nil {
		return nil, ErrMissingMerchant
	}
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	if !deliveryMethod.IsValid() {
		return nil, ErrInvalidDeliveryMethod
	}
	if deliveryMethod == DeliveryMethodDelivery && deliveryAddress == nil {
		return nil, ErrDeliveryAddressRequired
	}
	
	// Validate items and calculate total
	total := NewMoney(0) // Initialize with zero value in BTC (0 Satoshis)
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		subtotal := item.CalculateSubtotal()
		total = total.Add(subtotal)
	}
	
	now := time.Now()
	order := &Order{
		id:              uuid.New(),
		customerID:      customerID,
		merchantID:      merchantID,
		items:           items,
		status:          OrderStatusPending,
		totalAmount:     total,
		deliveryMethod:  deliveryMethod,
		deliveryAddress: deliveryAddress,
		createdAt:       now,
		updatedAt:       now,
		statusHistory:   []StatusChange{},
		events:          []DomainEvent{},
	}
	
	// Create initial event
	order.events = append(order.events, OrderPlacedEvent{
		OrderID:         order.id,
		CustomerID:      customerID,
		MerchantID:      merchantID,
		TotalAmount:     total,
		DeliveryMethod:  deliveryMethod,
		PlacedAt:        now,
	})
	
	return order, nil
}

// Accept accepts the order with an estimated preparation time
func (o *Order) Accept(estimatedMinutes int, acceptedBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusAccepted) {
		return nil, fmt.Errorf("%w: cannot transition from %s to ACCEPTED", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusAccepted
	o.estimatedWindow = NewTimeWindow(time.Now(), estimatedMinutes)
	o.recordStatusChange(OrderStatusAccepted, "Order accepted by merchant", acceptedBy)
	o.updatedAt = time.Now()
	
	event := OrderAcceptedEvent{
		OrderID:       o.id,
		MerchantID:    o.merchantID,
		CustomerID:    o.customerID,
		EstimatedTime: o.estimatedWindow.EndTime,
		AcceptedAt:    time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// Reject rejects the order with a reason
func (o *Order) Reject(reason string, rejectedBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusRejected) {
		return nil, fmt.Errorf("%w: cannot transition from %s to REJECTED", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusRejected
	o.recordStatusChange(OrderStatusRejected, reason, rejectedBy)
	o.updatedAt = time.Now()
	
	event := OrderRejectedEvent{
		OrderID:    o.id,
		MerchantID: o.merchantID,
		CustomerID: o.customerID,
		Reason:     reason,
		RejectedAt: time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// StartPreparing marks the order as being prepared
func (o *Order) StartPreparing(preparedBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusPreparing) {
		return nil, fmt.Errorf("%w: cannot transition from %s to PREPARING", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusPreparing
	o.recordStatusChange(OrderStatusPreparing, "Order preparation started", preparedBy)
	o.updatedAt = time.Now()
	
	event := OrderPreparingEvent{
		OrderID:     o.id,
		MerchantID:  o.merchantID,
		CustomerID:  o.customerID,
		StartedAt:   time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// MarkReady marks the order as ready for pickup or delivery
func (o *Order) MarkReady(markedBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusReady) {
		return nil, fmt.Errorf("%w: cannot transition from %s to READY", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusReady
	o.recordStatusChange(OrderStatusReady, "Order is ready", markedBy)
	o.updatedAt = time.Now()
	
	event := OrderReadyEvent{
		OrderID:        o.id,
		MerchantID:     o.merchantID,
		CustomerID:     o.customerID,
		DeliveryMethod: o.deliveryMethod,
		ReadyAt:        time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// DispatchForDelivery marks the order as out for delivery
func (o *Order) DispatchForDelivery(driverID uuid.UUID) ([]DomainEvent, error) {
	if o.deliveryMethod != DeliveryMethodDelivery {
		return nil, errors.New("can only dispatch delivery orders")
	}
	
	if !o.canTransitionTo(OrderStatusOutForDelivery) {
		return nil, fmt.Errorf("%w: cannot transition from %s to OUT_FOR_DELIVERY", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusOutForDelivery
	o.recordStatusChange(OrderStatusOutForDelivery, "Order out for delivery", driverID)
	o.updatedAt = time.Now()
	
	event := OrderOutForDeliveryEvent{
		OrderID:    o.id,
		CustomerID: o.customerID,
		DriverID:   driverID,
		Address:    o.deliveryAddress,
		DispatchedAt: time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// Complete marks the order as completed
func (o *Order) Complete(completedBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusCompleted) {
		return nil, fmt.Errorf("%w: cannot transition from %s to COMPLETED", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusCompleted
	o.recordStatusChange(OrderStatusCompleted, "Order completed", completedBy)
	o.updatedAt = time.Now()
	
	event := OrderCompletedEvent{
		OrderID:     o.id,
		MerchantID:  o.merchantID,
		CustomerID:  o.customerID,
		CompletedAt: time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// Cancel cancels the order with a reason
func (o *Order) Cancel(reason string, cancelledBy uuid.UUID) ([]DomainEvent, error) {
	if !o.canTransitionTo(OrderStatusCancelled) {
		return nil, fmt.Errorf("%w: cannot transition from %s to CANCELLED", 
			ErrInvalidStateTransition, o.status.String())
	}
	
	o.status = OrderStatusCancelled
	o.recordStatusChange(OrderStatusCancelled, reason, cancelledBy)
	o.updatedAt = time.Now()
	
	event := OrderCancelledEvent{
		OrderID:     o.id,
		MerchantID:  o.merchantID,
		CustomerID:  o.customerID,
		Reason:      reason,
		CancelledBy: cancelledBy,
		CancelledAt: time.Now(),
	}
	o.events = append(o.events, event)
	
	return []DomainEvent{event}, nil
}

// Internal methods

func (o *Order) canTransitionTo(newStatus OrderStatus) bool {
	validStates, exists := validTransitions[o.status]
	if !exists {
		return false
	}
	for _, valid := range validStates {
		if valid == newStatus {
			return true
		}
	}
	return false
}

func (o *Order) recordStatusChange(newStatus OrderStatus, reason string, changedBy uuid.UUID) {
	change := StatusChange{
		From:      o.status,
		To:        newStatus,
		Reason:    reason,
		ChangedAt: time.Now(),
		ChangedBy: changedBy,
	}
	o.statusHistory = append(o.statusHistory, change)
}

// CalculateSubtotal calculates the subtotal for an order item
func (oi OrderItem) CalculateSubtotal() Money {
	return oi.PricePerItem.Multiply(oi.Quantity)
}

// Getters for accessing private fields

func (o *Order) ID() uuid.UUID              { return o.id }
func (o *Order) CustomerID() uuid.UUID      { return o.customerID }
func (o *Order) MerchantID() uuid.UUID      { return o.merchantID }
func (o *Order) Items() []OrderItem         { return o.items }
func (o *Order) Status() OrderStatus        { return o.status }
func (o *Order) TotalAmount() Money         { return o.totalAmount }
func (o *Order) DeliveryMethod() DeliveryMethod { return o.deliveryMethod }
func (o *Order) DeliveryAddress() *Address  { return o.deliveryAddress }
func (o *Order) EstimatedWindow() *TimeWindow { return o.estimatedWindow }
func (o *Order) CreatedAt() time.Time       { return o.createdAt }
func (o *Order) UpdatedAt() time.Time       { return o.updatedAt }
func (o *Order) StatusHistory() []StatusChange { return o.statusHistory }
func (o *Order) Events() []DomainEvent      { return o.events }

// ClearEvents clears the events after they've been published
func (o *Order) ClearEvents() {
	o.events = []DomainEvent{}
}
