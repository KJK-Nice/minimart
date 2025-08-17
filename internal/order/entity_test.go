package order

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	customerID := uuid.New()
	merchantID := uuid.New()

	t.Run("creates valid order with delivery", func(t *testing.T) {
		address, _ := NewAddress("123 Main St", "San Francisco", "CA", "94102", "USA")
		items := []OrderItem{
			{
				MenuItemID:   uuid.New(),
				MenuItemName: "Burger",
				Quantity:     2,
				PricePerItem: NewMoney(25000), // 25,000 sats (~0.25 mBTC)
			},
			{
				MenuItemID:   uuid.New(),
				MenuItemName: "Fries",
				Quantity:     1,
				PricePerItem: NewMoney(10000), // 10,000 sats (~0.10 mBTC)
			},
		}

		order, err := NewOrder(customerID, merchantID, items, DeliveryMethodDelivery, address)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.ID() == uuid.Nil {
			t.Error("expected order to have ID")
		}

		if order.Status() != OrderStatusPending {
			t.Errorf("expected status PENDING, got %s", order.Status())
		}

		expectedTotal := NewMoney(60000) // (25,000 * 2) + 10,000 = 60,000 sats
		if order.TotalAmount() != expectedTotal {
			t.Errorf("expected total %s, got %s", expectedTotal, order.TotalAmount())
		}

		if len(order.Events()) != 1 {
			t.Errorf("expected 1 event, got %d", len(order.Events()))
		}

		if _, ok := order.Events()[0].(OrderPlacedEvent); !ok {
			t.Error("expected OrderPlacedEvent")
		}
	})

	t.Run("creates valid order with pickup", func(t *testing.T) {
		items := []OrderItem{
			{
				MenuItemID:   uuid.New(),
				MenuItemName: "Pizza",
				Quantity:     1,
				PricePerItem: NewMoney(50000), // 50,000 sats (~0.50 mBTC)
			},
		}

		order, err := NewOrder(customerID, merchantID, items, DeliveryMethodPickup, nil)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.DeliveryMethod() != DeliveryMethodPickup {
			t.Error("expected pickup delivery method")
		}

		if order.DeliveryAddress() != nil {
			t.Error("expected no delivery address for pickup")
		}
	})

	t.Run("fails with missing customer", func(t *testing.T) {
		items := []OrderItem{{MenuItemID: uuid.New(), Quantity: 1, PricePerItem: NewMoney(1000)}} // 1,000 sats
		_, err := NewOrder(uuid.Nil, merchantID, items, DeliveryMethodPickup, nil)

		if err != ErrMissingCustomer {
			t.Errorf("expected ErrMissingCustomer, got %v", err)
		}
	})

	t.Run("fails with missing merchant", func(t *testing.T) {
		items := []OrderItem{{MenuItemID: uuid.New(), Quantity: 1, PricePerItem: NewMoney(1000)}} // 1,000 sats
		_, err := NewOrder(customerID, uuid.Nil, items, DeliveryMethodPickup, nil)

		if err != ErrMissingMerchant {
			t.Errorf("expected ErrMissingMerchant, got %v", err)
		}
	})

	t.Run("fails with empty order", func(t *testing.T) {
		_, err := NewOrder(customerID, merchantID, []OrderItem{}, DeliveryMethodPickup, nil)

		if err != ErrEmptyOrder {
			t.Errorf("expected ErrEmptyOrder, got %v", err)
		}
	})

	t.Run("fails with invalid quantity", func(t *testing.T) {
		items := []OrderItem{
			{MenuItemID: uuid.New(), Quantity: 0, PricePerItem: NewMoney(1000)}, // 1,000 sats
		}
		_, err := NewOrder(customerID, merchantID, items, DeliveryMethodPickup, nil)

		if err != ErrInvalidQuantity {
			t.Errorf("expected ErrInvalidQuantity, got %v", err)
		}
	})

	t.Run("fails with delivery but no address", func(t *testing.T) {
		items := []OrderItem{{MenuItemID: uuid.New(), Quantity: 1, PricePerItem: NewMoney(1000)}} // 1,000 sats
		_, err := NewOrder(customerID, merchantID, items, DeliveryMethodDelivery, nil)

		if err != ErrDeliveryAddressRequired {
			t.Errorf("expected ErrDeliveryAddressRequired, got %v", err)
		}
	})
}

func TestOrderAccept(t *testing.T) {
	order := createTestOrder(t)
	merchantID := uuid.New()

	t.Run("accepts pending order", func(t *testing.T) {
		events, err := order.Accept(30, merchantID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.Status() != OrderStatusAccepted {
			t.Errorf("expected status ACCEPTED, got %s", order.Status())
		}

		if order.EstimatedWindow() == nil {
			t.Error("expected estimated window to be set")
		}

		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}

		acceptedEvent, ok := events[0].(OrderAcceptedEvent)
		if !ok {
			t.Error("expected OrderAcceptedEvent")
		}

		if acceptedEvent.OrderID != order.ID() {
			t.Error("event should have correct order ID")
		}

		// Check status history
		if len(order.StatusHistory()) != 1 {
			t.Errorf("expected 1 status change, got %d", len(order.StatusHistory()))
		}
	})

	t.Run("cannot accept non-pending order", func(t *testing.T) {
		order := createTestOrder(t)
		order.Accept(30, merchantID) // Accept first

		// Try to accept again
		_, err := order.Accept(30, merchantID)

		if err == nil {
			t.Error("expected error when accepting already accepted order")
		}
	})
}

func TestOrderReject(t *testing.T) {
	order := createTestOrder(t)
	merchantID := uuid.New()

	t.Run("rejects pending order", func(t *testing.T) {
		reason := "Out of stock"
		events, err := order.Reject(reason, merchantID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.Status() != OrderStatusRejected {
			t.Errorf("expected status REJECTED, got %s", order.Status())
		}

		rejectedEvent, ok := events[0].(OrderRejectedEvent)
		if !ok {
			t.Error("expected OrderRejectedEvent")
		}

		if rejectedEvent.Reason != reason {
			t.Errorf("expected reason '%s', got '%s'", reason, rejectedEvent.Reason)
		}
	})

	t.Run("cannot reject accepted order", func(t *testing.T) {
		order := createTestOrder(t)
		order.Accept(30, merchantID)

		_, err := order.Reject("Too late", merchantID)

		if err == nil {
			t.Error("expected error when rejecting accepted order")
		}
	})
}

func TestOrderWorkflow(t *testing.T) {
	t.Run("complete delivery workflow", func(t *testing.T) {
		order := createTestOrderWithDelivery(t)
		merchantID := uuid.New()
		driverID := uuid.New()

		// Accept
		events, err := order.Accept(30, merchantID)
		if err != nil {
			t.Fatalf("failed to accept: %v", err)
		}
		if len(events) != 1 {
			t.Error("expected 1 event for accept")
		}

		// Start preparing
		events, err = order.StartPreparing(merchantID)
		if err != nil {
			t.Fatalf("failed to start preparing: %v", err)
		}
		if order.Status() != OrderStatusPreparing {
			t.Errorf("expected status PREPARING, got %s", order.Status())
		}

		// Mark ready
		events, err = order.MarkReady(merchantID)
		if err != nil {
			t.Fatalf("failed to mark ready: %v", err)
		}
		if order.Status() != OrderStatusReady {
			t.Errorf("expected status READY, got %s", order.Status())
		}

		// Dispatch for delivery
		events, err = order.DispatchForDelivery(driverID)
		if err != nil {
			t.Fatalf("failed to dispatch: %v", err)
		}
		if order.Status() != OrderStatusOutForDelivery {
			t.Errorf("expected status OUT_FOR_DELIVERY, got %s", order.Status())
		}

		// Complete
		events, err = order.Complete(driverID)
		if err != nil {
			t.Fatalf("failed to complete: %v", err)
		}
		if order.Status() != OrderStatusCompleted {
			t.Errorf("expected status COMPLETED, got %s", order.Status())
		}

		// Check total events (we clear the initial placed event in createTestOrderWithDelivery)
		totalEvents := len(order.Events())
		if totalEvents < 5 { // accepted + preparing + ready + dispatched + completed
			t.Errorf("expected at least 5 events, got %d", totalEvents)
		}

		// Check status history
		if len(order.StatusHistory()) != 5 {
			t.Errorf("expected 5 status changes, got %d", len(order.StatusHistory()))
		}
	})

	t.Run("complete pickup workflow", func(t *testing.T) {
		order := createTestOrder(t) // Pickup order
		merchantID := uuid.New()
		customerID := order.CustomerID()

		// Accept -> Preparing -> Ready -> Complete
		order.Accept(20, merchantID)
		order.StartPreparing(merchantID)
		order.MarkReady(merchantID)

		// Cannot dispatch pickup order
		_, err := order.DispatchForDelivery(uuid.New())
		if err == nil {
			t.Error("expected error when dispatching pickup order")
		}

		// Complete pickup
		events, err := order.Complete(customerID)
		if err != nil {
			t.Fatalf("failed to complete pickup: %v", err)
		}

		if order.Status() != OrderStatusCompleted {
			t.Errorf("expected status COMPLETED, got %s", order.Status())
		}

		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}

		completedEvent, ok := events[0].(OrderCompletedEvent)
		if !ok {
			t.Error("expected OrderCompletedEvent")
		}
		if completedEvent.CustomerID != customerID {
			t.Error("completed event should have correct customer ID")
		}
	})
}

func TestOrderCancel(t *testing.T) {
	t.Run("customer cancels pending order", func(t *testing.T) {
		order := createTestOrder(t)
		customerID := order.CustomerID()

		events, err := order.Cancel("Changed my mind", customerID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.Status() != OrderStatusCancelled {
			t.Errorf("expected status CANCELLED, got %s", order.Status())
		}

		cancelledEvent, ok := events[0].(OrderCancelledEvent)
		if !ok {
			t.Error("expected OrderCancelledEvent")
		}

		if cancelledEvent.CancelledBy != customerID {
			t.Error("cancelled event should have correct canceller ID")
		}
	})

	t.Run("merchant cancels accepted order", func(t *testing.T) {
		order := createTestOrder(t)
		merchantID := uuid.New()

		order.Accept(30, merchantID)

		events, err := order.Cancel("Kitchen equipment failure", merchantID)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if order.Status() != OrderStatusCancelled {
			t.Errorf("expected status CANCELLED, got %s", order.Status())
		}

		if len(events) != 1 {
			t.Errorf("expected 1 event, got %d", len(events))
		}
	})

	t.Run("cannot cancel completed order", func(t *testing.T) {
		order := createTestOrder(t)
		merchantID := uuid.New()

		// Complete the order
		order.Accept(30, merchantID)
		order.StartPreparing(merchantID)
		order.MarkReady(merchantID)
		order.Complete(merchantID)

		_, err := order.Cancel("Too late", order.CustomerID())

		if err == nil {
			t.Error("expected error when cancelling completed order")
		}
	})
}

func TestInvalidStateTransitions(t *testing.T) {
	tests := []struct {
		name          string
		setupStatus   func(*Order)
		transition    func(*Order) ([]DomainEvent, error)
		expectedError bool
	}{
		{
			name: "cannot go from pending to preparing",
			setupStatus: func(o *Order) {
				// Order starts as pending
			},
			transition: func(o *Order) ([]DomainEvent, error) {
				return o.StartPreparing(uuid.New())
			},
			expectedError: true,
		},
		{
			name: "cannot go from rejected to accepted",
			setupStatus: func(o *Order) {
				o.Reject("reason", uuid.New())
			},
			transition: func(o *Order) ([]DomainEvent, error) {
				return o.Accept(30, uuid.New())
			},
			expectedError: true,
		},
		{
			name: "cannot go from completed to any state",
			setupStatus: func(o *Order) {
				merchantID := uuid.New()
				o.Accept(30, merchantID)
				o.StartPreparing(merchantID)
				o.MarkReady(merchantID)
				o.Complete(merchantID)
			},
			transition: func(o *Order) ([]DomainEvent, error) {
				return o.Cancel("test", uuid.New())
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := createTestOrder(t)
			tt.setupStatus(order)

			_, err := tt.transition(order)

			if tt.expectedError && err == nil {
				t.Error("expected error for invalid transition")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Helper functions

func createTestOrder(t *testing.T) *Order {
	customerID := uuid.New()
	merchantID := uuid.New()
	items := []OrderItem{
		{
			MenuItemID:   uuid.New(),
			MenuItemName: "Test Item",
			Quantity:     1,
			PricePerItem: NewMoney(10000), // 10,000 sats (~0.10 mBTC)
		},
	}

	order, err := NewOrder(customerID, merchantID, items, DeliveryMethodPickup, nil)
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}

	// Clear the initial event for cleaner testing
	order.ClearEvents()

	return order
}

func createTestOrderWithDelivery(t *testing.T) *Order {
	customerID := uuid.New()
	merchantID := uuid.New()
	address, _ := NewAddress("123 Test St", "Test City", "CA", "12345", "USA")
	items := []OrderItem{
		{
			MenuItemID:   uuid.New(),
			MenuItemName: "Test Item",
			Quantity:     1,
			PricePerItem: NewMoney(10000), // 10,000 sats (~0.10 mBTC)
		},
	}

	order, err := NewOrder(customerID, merchantID, items, DeliveryMethodDelivery, address)
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}

	// Clear the initial event for cleaner testing
	order.ClearEvents()

	return order
}
