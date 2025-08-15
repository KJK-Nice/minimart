package menu

import (
	"testing"
	
	"github.com/google/uuid"
)

func TestNewMenuItem(t *testing.T) {
	merchantID := uuid.New()
	
	t.Run("creates valid menu item", func(t *testing.T) {
		item, err := NewMenuItem(
			merchantID,
			"Bitcoin Burger",
			"A delicious burger priced in Bitcoin",
			25000, // 25,000 sats
		)
		
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		
		if item.ID() == uuid.Nil {
			t.Error("expected item to have ID")
		}
		
		if item.Name() != "Bitcoin Burger" {
			t.Errorf("expected name 'Bitcoin Burger', got %s", item.Name())
		}
		
		if item.GetPriceInSatoshis() != 25000 {
			t.Errorf("expected price 25000 sats, got %d", item.GetPriceInSatoshis())
		}
		
		if item.StockLevel() != -1 {
			t.Errorf("expected unlimited stock (-1), got %d", item.StockLevel())
		}
		
		if !item.IsAvailable() {
			t.Error("expected item to be available by default")
		}
	})
	
	t.Run("fails with invalid merchant", func(t *testing.T) {
		_, err := NewMenuItem(uuid.Nil, "Item", "Description", 1000)
		
		if err != ErrInvalidMerchant {
			t.Errorf("expected ErrInvalidMerchant, got %v", err)
		}
	})
	
	t.Run("fails with empty name", func(t *testing.T) {
		_, err := NewMenuItem(merchantID, "", "Description", 1000)
		
		if err != ErrInvalidName {
			t.Errorf("expected ErrInvalidName, got %v", err)
		}
	})
	
	t.Run("fails with invalid price", func(t *testing.T) {
		_, err := NewMenuItem(merchantID, "Item", "Description", 0)
		
		if err != ErrInvalidPrice {
			t.Errorf("expected ErrInvalidPrice, got %v", err)
		}
		
		_, err = NewMenuItem(merchantID, "Item", "Description", -100)
		
		if err != ErrInvalidPrice {
			t.Errorf("expected ErrInvalidPrice for negative price, got %v", err)
		}
	})
}

func TestMenuItemAvailability(t *testing.T) {
	item := createTestMenuItem(t)
	
	t.Run("item is available by default", func(t *testing.T) {
		if !item.IsAvailable() {
			t.Error("expected item to be available")
		}
	})
	
	t.Run("can make item unavailable", func(t *testing.T) {
		item.MakeUnavailable()
		
		if item.IsAvailable() {
			t.Error("expected item to be unavailable")
		}
	})
	
	t.Run("can make item available again", func(t *testing.T) {
		item.MakeUnavailable()
		item.MakeAvailable()
		
		if !item.IsAvailable() {
			t.Error("expected item to be available")
		}
	})
	
	t.Run("unavailable when stock is zero", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(0)
		
		if item.IsAvailable() {
			t.Error("expected item to be unavailable when stock is 0")
		}
	})
	
	t.Run("available with positive stock", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(10)
		
		if !item.IsAvailable() {
			t.Error("expected item to be available with positive stock")
		}
	})
}

func TestMenuItemStock(t *testing.T) {
	t.Run("unlimited stock by default", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		if item.StockLevel() != -1 {
			t.Errorf("expected unlimited stock (-1), got %d", item.StockLevel())
		}
		
		// Should be able to fulfill any quantity
		if !item.CanFulfillQuantity(1000000) {
			t.Error("expected to fulfill large quantity with unlimited stock")
		}
	})
	
	t.Run("can set stock level", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.SetStockLevel(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if item.StockLevel() != 10 {
			t.Errorf("expected stock level 10, got %d", item.StockLevel())
		}
	})
	
	t.Run("cannot set negative stock (except -1)", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.SetStockLevel(-5)
		if err != ErrNegativeStockAdjustment {
			t.Errorf("expected ErrNegativeStockAdjustment, got %v", err)
		}
		
		// But -1 (unlimited) should work
		err = item.SetStockLevel(-1)
		if err != nil {
			t.Errorf("expected no error for -1 (unlimited), got %v", err)
		}
	})
	
	t.Run("can fulfill quantity with sufficient stock", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(10)
		
		if !item.CanFulfillQuantity(5) {
			t.Error("expected to fulfill quantity 5 with stock 10")
		}
		
		if !item.CanFulfillQuantity(10) {
			t.Error("expected to fulfill quantity 10 with stock 10")
		}
	})
	
	t.Run("cannot fulfill quantity with insufficient stock", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(5)
		
		if item.CanFulfillQuantity(10) {
			t.Error("should not fulfill quantity 10 with stock 5")
		}
	})
	
	t.Run("cannot fulfill invalid quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		if item.CanFulfillQuantity(0) {
			t.Error("should not fulfill quantity 0")
		}
		
		if item.CanFulfillQuantity(-1) {
			t.Error("should not fulfill negative quantity")
		}
	})
}

func TestMenuItemReserveStock(t *testing.T) {
	t.Run("reserve stock with limited quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(10)
		
		err := item.ReserveStock(3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if item.StockLevel() != 7 {
			t.Errorf("expected stock level 7 after reserving 3, got %d", item.StockLevel())
		}
	})
	
	t.Run("reserve stock with unlimited quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		// Default is unlimited (-1)
		
		err := item.ReserveStock(1000)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Stock should remain unlimited
		if item.StockLevel() != -1 {
			t.Errorf("expected stock to remain unlimited, got %d", item.StockLevel())
		}
	})
	
	t.Run("cannot reserve more than available", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(5)
		
		err := item.ReserveStock(10)
		if err != ErrInsufficientStock {
			t.Errorf("expected ErrInsufficientStock, got %v", err)
		}
		
		// Stock should remain unchanged
		if item.StockLevel() != 5 {
			t.Errorf("expected stock to remain 5, got %d", item.StockLevel())
		}
	})
	
	t.Run("cannot reserve when unavailable", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.MakeUnavailable()
		
		err := item.ReserveStock(1)
		if err != ErrItemNotAvailable {
			t.Errorf("expected ErrItemNotAvailable, got %v", err)
		}
	})
	
	t.Run("cannot reserve invalid quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.ReserveStock(0)
		if err != ErrInvalidQuantity {
			t.Errorf("expected ErrInvalidQuantity for 0, got %v", err)
		}
		
		err = item.ReserveStock(-1)
		if err != ErrInvalidQuantity {
			t.Errorf("expected ErrInvalidQuantity for negative, got %v", err)
		}
	})
}

func TestMenuItemReleaseStock(t *testing.T) {
	t.Run("release stock increases quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(5)
		
		err := item.ReleaseStock(3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if item.StockLevel() != 8 {
			t.Errorf("expected stock level 8 after releasing 3, got %d", item.StockLevel())
		}
	})
	
	t.Run("release stock with unlimited does nothing", func(t *testing.T) {
		item := createTestMenuItem(t)
		// Default is unlimited (-1)
		
		err := item.ReleaseStock(10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Should remain unlimited
		if item.StockLevel() != -1 {
			t.Errorf("expected stock to remain unlimited, got %d", item.StockLevel())
		}
	})
	
	t.Run("cannot release invalid quantity", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.ReleaseStock(0)
		if err != ErrInvalidQuantity {
			t.Errorf("expected ErrInvalidQuantity for 0, got %v", err)
		}
		
		err = item.ReleaseStock(-1)
		if err != ErrInvalidQuantity {
			t.Errorf("expected ErrInvalidQuantity for negative, got %v", err)
		}
	})
}

func TestMenuItemPriceUpdate(t *testing.T) {
	t.Run("can update price", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.UpdatePrice(50000) // 50,000 sats
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if item.GetPriceInSatoshis() != 50000 {
			t.Errorf("expected price 50000 sats, got %d", item.GetPriceInSatoshis())
		}
	})
	
	t.Run("cannot set invalid price", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.UpdatePrice(0)
		if err != ErrInvalidPrice {
			t.Errorf("expected ErrInvalidPrice for 0, got %v", err)
		}
		
		err = item.UpdatePrice(-100)
		if err != ErrInvalidPrice {
			t.Errorf("expected ErrInvalidPrice for negative, got %v", err)
		}
	})
}

func TestMenuItemDetails(t *testing.T) {
	t.Run("can update details", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.UpdateDetails("New Name", "New Description")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if item.Name() != "New Name" {
			t.Errorf("expected name 'New Name', got %s", item.Name())
		}
		
		if item.Description() != "New Description" {
			t.Errorf("expected description 'New Description', got %s", item.Description())
		}
	})
	
	t.Run("cannot set empty name", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		err := item.UpdateDetails("", "Description")
		if err != ErrInvalidName {
			t.Errorf("expected ErrInvalidName, got %v", err)
		}
	})
	
	t.Run("can set category and image", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		item.SetCategory("Burgers")
		if item.Category() != "Burgers" {
			t.Errorf("expected category 'Burgers', got %s", item.Category())
		}
		
		item.SetImageURL("https://example.com/burger.jpg")
		if item.ImageURL() != "https://example.com/burger.jpg" {
			t.Errorf("expected correct image URL, got %s", item.ImageURL())
		}
	})
}

func TestCreateOrderItem(t *testing.T) {
	t.Run("creates order item successfully", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(10)
		
		orderItem, err := item.CreateOrderItem(2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if orderItem.MenuItemID != item.ID() {
			t.Error("order item should have correct menu item ID")
		}
		
		if orderItem.MenuItemName != item.Name() {
			t.Errorf("expected name %s, got %s", item.Name(), orderItem.MenuItemName)
		}
		
		if orderItem.Quantity != 2 {
			t.Errorf("expected quantity 2, got %d", orderItem.Quantity)
		}
		
		// Price should match
		if orderItem.PricePerItem != item.Price() {
			t.Error("order item should have correct price")
		}
	})
	
	t.Run("fails when item unavailable", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.MakeUnavailable()
		
		_, err := item.CreateOrderItem(1)
		if err != ErrItemNotAvailable {
			t.Errorf("expected ErrItemNotAvailable, got %v", err)
		}
	})
	
	t.Run("fails with insufficient stock", func(t *testing.T) {
		item := createTestMenuItem(t)
		item.SetStockLevel(2)
		
		_, err := item.CreateOrderItem(5)
		if err != ErrInsufficientStock {
			t.Errorf("expected ErrInsufficientStock, got %v", err)
		}
	})
	
	t.Run("works with unlimited stock", func(t *testing.T) {
		item := createTestMenuItem(t)
		// Default is unlimited
		
		orderItem, err := item.CreateOrderItem(1000)
		if err != nil {
			t.Fatalf("unexpected error with unlimited stock: %v", err)
		}
		
		if orderItem.Quantity != 1000 {
			t.Errorf("expected quantity 1000, got %d", orderItem.Quantity)
		}
	})
}

func TestStockWorkflow(t *testing.T) {
	t.Run("complete stock management workflow", func(t *testing.T) {
		item := createTestMenuItem(t)
		
		// Start with limited stock
		item.SetStockLevel(20)
		
		// Reserve some stock for orders
		err := item.ReserveStock(5)
		if err != nil {
			t.Fatalf("failed to reserve stock: %v", err)
		}
		if item.StockLevel() != 15 {
			t.Errorf("expected stock 15, got %d", item.StockLevel())
		}
		
		// Reserve more
		err = item.ReserveStock(10)
		if err != nil {
			t.Fatalf("failed to reserve more stock: %v", err)
		}
		if item.StockLevel() != 5 {
			t.Errorf("expected stock 5, got %d", item.StockLevel())
		}
		
		// Try to reserve more than available
		err = item.ReserveStock(10)
		if err != ErrInsufficientStock {
			t.Error("expected insufficient stock error")
		}
		
		// Release some stock (cancelled order)
		err = item.ReleaseStock(5)
		if err != nil {
			t.Fatalf("failed to release stock: %v", err)
		}
		if item.StockLevel() != 10 {
			t.Errorf("expected stock 10, got %d", item.StockLevel())
		}
		
		// Now we can reserve again
		err = item.ReserveStock(8)
		if err != nil {
			t.Fatalf("failed to reserve after release: %v", err)
		}
		if item.StockLevel() != 2 {
			t.Errorf("expected stock 2, got %d", item.StockLevel())
		}
	})
}

// Helper functions

func createTestMenuItem(t *testing.T) *MenuItem {
	merchantID := uuid.New()
	item, err := NewMenuItem(
		merchantID,
		"Test Item",
		"Test Description",
		10000, // 10,000 sats
	)
	if err != nil {
		t.Fatalf("failed to create test menu item: %v", err)
	}
	return item
}

func createMenuItemWithStock(t *testing.T, stock int) *MenuItem {
	item := createTestMenuItem(t)
	err := item.SetStockLevel(stock)
	if err != nil {
		t.Fatalf("failed to set stock level: %v", err)
	}
	return item
}
