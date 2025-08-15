package menu

import (
	"errors"
	"github.com/google/uuid"
	"minimart/internal/order"
)

// Domain errors
var (
	ErrItemOutOfStock        = errors.New("item is out of stock")
	ErrInvalidPrice          = errors.New("price must be positive")
	ErrInvalidQuantity       = errors.New("quantity must be positive")
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrItemNotAvailable      = errors.New("item is not available")
	ErrInvalidName           = errors.New("item name is required")
	ErrInvalidMerchant       = errors.New("merchant ID is required")
	ErrNegativeStockAdjustment = errors.New("stock cannot be negative")
)

// MenuItem represents a product or service that can be ordered
// This is now a rich domain entity with business logic
type MenuItem struct {
	id          uuid.UUID
	merchantID  uuid.UUID
	name        string
	description string
	price       order.Money // Using Bitcoin/Satoshis
	stockLevel  int         // -1 means unlimited stock
	isAvailable bool        // Can be made unavailable even with stock
	category    string
	imageURL    string
}

// NewMenuItem creates a new menu item with validation
func NewMenuItem(
	merchantID uuid.UUID,
	name string,
	description string,
	priceInSatoshis int64,
) (*MenuItem, error) {
	if merchantID == uuid.Nil {
		return nil, ErrInvalidMerchant
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if priceInSatoshis <= 0 {
		return nil, ErrInvalidPrice
	}
	
	return &MenuItem{
		id:          uuid.New(),
		merchantID:  merchantID,
		name:        name,
		description: description,
	price:       order.NewMoney(priceInSatoshis),
		stockLevel:  -1, // Unlimited by default
		isAvailable: true,
	}, nil
}

// Business logic methods

// IsAvailable checks if the item can be ordered
func (m *MenuItem) IsAvailable() bool {
	if !m.isAvailable {
		return false
	}
	// If stock tracking is enabled (stockLevel >= 0), check stock
	if m.stockLevel == 0 {
		return false
	}
	return true
}

// CanFulfillQuantity checks if the requested quantity can be fulfilled
func (m *MenuItem) CanFulfillQuantity(quantity int) bool {
	if quantity <= 0 {
		return false
	}
	if !m.IsAvailable() {
		return false
	}
	// Unlimited stock
	if m.stockLevel < 0 {
		return true
	}
	// Check if we have enough stock
	return m.stockLevel >= quantity
}

// ReserveStock reserves stock for an order (decreases available stock)
func (m *MenuItem) ReserveStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if !m.IsAvailable() {
		return ErrItemNotAvailable
	}
	
	// Unlimited stock - no need to reserve
	if m.stockLevel < 0 {
		return nil
	}
	
	// Check sufficient stock
	if m.stockLevel < quantity {
		return ErrInsufficientStock
	}
	
	m.stockLevel -= quantity
	return nil
}

// ReleaseStock releases previously reserved stock (increases available stock)
func (m *MenuItem) ReleaseStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	
	// Unlimited stock - no need to release
	if m.stockLevel < 0 {
		return nil
	}
	
	m.stockLevel += quantity
	return nil
}

// SetStockLevel sets the stock level (-1 for unlimited)
func (m *MenuItem) SetStockLevel(level int) error {
	if level < -1 {
		return ErrNegativeStockAdjustment
	}
	m.stockLevel = level
	return nil
}

// MakeAvailable makes the item available for ordering
func (m *MenuItem) MakeAvailable() {
	m.isAvailable = true
}

// MakeUnavailable makes the item unavailable for ordering
func (m *MenuItem) MakeUnavailable() {
	m.isAvailable = false
}

// UpdatePrice updates the item's price
func (m *MenuItem) UpdatePrice(priceInSatoshis int64) error {
	if priceInSatoshis <= 0 {
		return ErrInvalidPrice
	}
	m.price = order.NewMoney(priceInSatoshis)
	return nil
}

// UpdateDetails updates name and description
func (m *MenuItem) UpdateDetails(name, description string) error {
	if name == "" {
		return ErrInvalidName
	}
	m.name = name
	m.description = description
	return nil
}

// SetCategory sets the item's category
func (m *MenuItem) SetCategory(category string) {
	m.category = category
}

// SetImageURL sets the item's image URL
func (m *MenuItem) SetImageURL(url string) {
	m.imageURL = url
}

// Getters for accessing private fields

func (m *MenuItem) ID() uuid.UUID          { return m.id }
func (m *MenuItem) MerchantID() uuid.UUID  { return m.merchantID }
func (m *MenuItem) Name() string           { return m.name }
func (m *MenuItem) Description() string    { return m.description }
func (m *MenuItem) Price() order.Money     { return m.price }
func (m *MenuItem) StockLevel() int        { return m.stockLevel }
func (m *MenuItem) Category() string       { return m.category }
func (m *MenuItem) ImageURL() string       { return m.imageURL }

// GetPriceInSatoshis returns the price in Satoshis for compatibility
func (m *MenuItem) GetPriceInSatoshis() int64 {
	return m.price.Amount()
}

// CreateOrderItem creates an OrderItem from this MenuItem
// This is used when adding items to an order
func (m *MenuItem) CreateOrderItem(quantity int) (*order.OrderItem, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	
	if !m.IsAvailable() {
		return nil, ErrItemNotAvailable
	}
	
	if !m.CanFulfillQuantity(quantity) {
		return nil, ErrInsufficientStock
	}
	
	// Reserve stock (this handles unlimited stock internally)
	if err := m.ReserveStock(quantity); err != nil {
		return nil, err
	}
	
	return &order.OrderItem{
		MenuItemID:   m.id,
		MenuItemName: m.name,
		Quantity:     quantity,
		PricePerItem: m.price,
	}, nil
}
