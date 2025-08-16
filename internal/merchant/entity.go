package merchant

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Domain errors
var (
	ErrMerchantNotActive     = errors.New("merchant is not active")
	ErrOutsideOperatingHours = errors.New("merchant is outside operating hours")
	ErrInvalidOperatingHours = errors.New("invalid operating hours")
)

// OperatingHours represents the business hours for a merchant
type OperatingHours struct {
	OpenTime  time.Duration // Duration since midnight (e.g., 9*time.Hour for 9:00 AM)
	CloseTime time.Duration // Duration since midnight (e.g., 17*time.Hour for 5:00 PM)
	DaysOpen  []time.Weekday
}

// NewOperatingHours creates operating hours with validation
func NewOperatingHours(openHour, closeHour int, daysOpen []time.Weekday) (OperatingHours, error) {
	if openHour < 0 || openHour > 23 || closeHour < 0 || closeHour > 23 {
		return OperatingHours{}, ErrInvalidOperatingHours
	}
	if len(daysOpen) == 0 {
		return OperatingHours{}, ErrInvalidOperatingHours
	}

	return OperatingHours{
		OpenTime:  time.Duration(openHour) * time.Hour,
		CloseTime: time.Duration(closeHour) * time.Hour,
		DaysOpen:  daysOpen,
	}, nil
}

// IsOpenAt checks if the merchant is open at a given time
func (oh OperatingHours) IsOpenAt(t time.Time) bool {
	// Check if the day is in operating days
	dayOpen := false
	for _, day := range oh.DaysOpen {
		if t.Weekday() == day {
			dayOpen = true
			break
		}
	}
	if !dayOpen {
		return false
	}

	// Convert time to duration since midnight
	timeOfDay := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute

	// Handle same-day hours
	if oh.CloseTime > oh.OpenTime {
		return timeOfDay >= oh.OpenTime && timeOfDay <= oh.CloseTime
	}

	// Handle overnight hours (e.g., 22:00 to 06:00)
	return timeOfDay >= oh.OpenTime || timeOfDay <= oh.CloseTime
}

type Merchant struct {
	id              uuid.UUID
	name            string
	description     string
	isActive        bool
	operatingHours  *OperatingHours
	preparationTime int // Default preparation time in minutes
	createdAt       time.Time
	updatedAt       time.Time
}

// NewMerchant creates a new merchant with validation
func NewMerchant(name, description string) *Merchant {
	// Default operating hours: 9 AM to 9 PM, Monday to Sunday
	defaultHours, _ := NewOperatingHours(9, 21, []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	})

	now := time.Now()
	return &Merchant{
		id:              uuid.New(),
		name:            name,
		description:     description,
		isActive:        true,
		operatingHours:  &defaultHours,
		preparationTime: 30, // Default 30 minutes
		createdAt:       now,
		updatedAt:       now,
	}
}

// Getters
func (m *Merchant) ID() uuid.UUID {
	return m.id
}

func (m *Merchant) Name() string {
	return m.name
}

func (m *Merchant) Description() string {
	return m.description
}

func (m *Merchant) IsActive() bool {
	return m.isActive
}

func (m *Merchant) OperatingHours() *OperatingHours {
	return m.operatingHours
}

func (m *Merchant) PreparationTime() int {
	return m.preparationTime
}

func (m *Merchant) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Merchant) UpdatedAt() time.Time {
	return m.updatedAt
}

// Business logic methods

// CanAcceptOrders checks if the merchant can currently accept orders
func (m *Merchant) CanAcceptOrders() error {
	if !m.isActive {
		return ErrMerchantNotActive
	}

	if m.operatingHours != nil && !m.operatingHours.IsOpenAt(time.Now()) {
		return ErrOutsideOperatingHours
	}

	return nil
}

// CanAcceptOrdersAt checks if the merchant can accept orders at a specific time
func (m *Merchant) CanAcceptOrdersAt(t time.Time) error {
	if !m.isActive {
		return ErrMerchantNotActive
	}

	if m.operatingHours != nil && !m.operatingHours.IsOpenAt(t) {
		return ErrOutsideOperatingHours
	}

	return nil
}

// EstimatePreparationTime returns the merchant's standard preparation time
// This can be enhanced later with item-specific logic
func (m *Merchant) EstimatePreparationTime(itemCount int) int {
	// Simple heuristic: base preparation time + 2 minutes per additional item
	baseTime := m.preparationTime
	if itemCount > 1 {
		baseTime += (itemCount - 1) * 2
	}
	return baseTime
}

// UpdateOperatingHours updates the merchant's operating hours
func (m *Merchant) UpdateOperatingHours(hours OperatingHours) {
	m.operatingHours = &hours
	m.updatedAt = time.Now()
}

// UpdatePreparationTime updates the default preparation time
func (m *Merchant) UpdatePreparationTime(minutes int) error {
	if minutes < 1 || minutes > 240 { // 1 minute to 4 hours
		return errors.New("preparation time must be between 1 and 240 minutes")
	}
	m.preparationTime = minutes
	m.updatedAt = time.Now()
	return nil
}

// Deactivate marks the merchant as inactive
func (m *Merchant) Deactivate() {
	m.isActive = false
	m.updatedAt = time.Now()
}

// Activate marks the merchant as active
func (m *Merchant) Activate() {
	m.isActive = true
	m.updatedAt = time.Now()
}
