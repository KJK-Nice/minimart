package order

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// DeliveryMethod represents how an order will be fulfilled
type DeliveryMethod int

const (
	DeliveryMethodPickup DeliveryMethod = iota
	DeliveryMethodDelivery
)

func (d DeliveryMethod) String() string {
	switch d {
	case DeliveryMethodPickup:
		return "PICKUP"
	case DeliveryMethodDelivery:
		return "DELIVERY"
	default:
		return "UNKNOWN"
	}
}

func (d DeliveryMethod) IsValid() bool {
	return d == DeliveryMethodPickup || d == DeliveryMethodDelivery
}

// Money represents a monetary value with currency
// We use Satoshis as the base unit (1 BTC = 100,000,000 Satoshis)
type Money struct {
	amount   int64  // Amount in Satoshis to avoid floating point issues
	currency string
}

// NewMoney creates a new Money value object in Satoshis
func NewMoney(amountInSatoshis int64) Money {
	return Money{
		amount:   amountInSatoshis,
		currency: "BTC",
	}
}

// NewMoneyFromBTC creates Money from BTC amount
func NewMoneyFromBTC(btc float64) Money {
	satoshis := int64(btc * 100_000_000) // 1 BTC = 100,000,000 Satoshis
	return NewMoney(satoshis)
}

// NewMoneyFromMilliBTC creates Money from mBTC amount (1 mBTC = 0.001 BTC)
func NewMoneyFromMilliBTC(mbtc float64) Money {
	satoshis := int64(mbtc * 100_000) // 1 mBTC = 100,000 Satoshis
	return NewMoney(satoshis)
}

// Amount returns the amount in Satoshis
func (m Money) Amount() int64 {
	return m.amount
}

// AmountInBTC returns the amount in BTC
func (m Money) AmountInBTC() float64 {
	return float64(m.amount) / 100_000_000
}

// AmountInMilliBTC returns the amount in mBTC
func (m Money) AmountInMilliBTC() float64 {
	return float64(m.amount) / 100_000
}

// Currency returns the currency code
func (m Money) Currency() string {
	return m.currency
}

// Add adds two money values
func (m Money) Add(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("cannot add different currencies: %s and %s", m.currency, other.currency))
	}
	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}
}

// Subtract subtracts another money value
func (m Money) Subtract(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("cannot subtract different currencies: %s and %s", m.currency, other.currency))
	}
	return Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}
}

// Multiply multiplies money by a quantity
func (m Money) Multiply(quantity int) Money {
	return Money{
		amount:   m.amount * int64(quantity),
		currency: m.currency,
	}
}

// IsPositive checks if the amount is positive
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsZero checks if the amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsNegative checks if the amount is negative
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// Equals checks if two money values are equal
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns a formatted string representation
func (m Money) String() string {
	// Display in different units based on amount size
	if m.amount >= 10_000_000 { // >= 0.1 BTC, show in BTC
		btc := float64(m.amount) / 100_000_000
		return fmt.Sprintf("%.8f BTC", btc)
	} else if m.amount >= 100_000 { // >= 1 mBTC, show in mBTC
		mbtc := float64(m.amount) / 100_000
		return fmt.Sprintf("%.3f mBTC", mbtc)
	} else { // Show in Satoshis
		return fmt.Sprintf("%d sats", m.amount)
	}
}

// Address represents a delivery address
type Address struct {
	street     string
	city       string
	state      string
	postalCode string
	country    string
	unit       string // Optional: apartment, suite, etc.
}

// NewAddress creates a new address with validation
func NewAddress(street, city, state, postalCode, country string) (*Address, error) {
	if street == "" {
		return nil, errors.New("street is required")
	}
	if city == "" {
		return nil, errors.New("city is required")
	}
	if state == "" {
		return nil, errors.New("state is required")
	}
	if postalCode == "" {
		return nil, errors.New("postal code is required")
	}
	if country == "" {
		country = "USA" // Default to USA
	}
	
	return &Address{
		street:     strings.TrimSpace(street),
		city:       strings.TrimSpace(city),
		state:      strings.TrimSpace(state),
		postalCode: strings.TrimSpace(postalCode),
		country:    strings.TrimSpace(country),
	}, nil
}

// WithUnit adds a unit number to the address
func (a *Address) WithUnit(unit string) *Address {
	a.unit = strings.TrimSpace(unit)
	return a
}

// Getters for address fields
func (a *Address) Street() string     { return a.street }
func (a *Address) City() string       { return a.city }
func (a *Address) State() string      { return a.state }
func (a *Address) PostalCode() string { return a.postalCode }
func (a *Address) Country() string    { return a.country }
func (a *Address) Unit() string       { return a.unit }

// String returns a formatted address
func (a *Address) String() string {
	lines := []string{}
	
	if a.unit != "" {
		lines = append(lines, fmt.Sprintf("%s, Unit %s", a.street, a.unit))
	} else {
		lines = append(lines, a.street)
	}
	
	lines = append(lines, fmt.Sprintf("%s, %s %s", a.city, a.state, a.postalCode))
	
	if a.country != "USA" {
		lines = append(lines, a.country)
	}
	
	return strings.Join(lines, "\n")
}

// Equals checks if two addresses are the same
func (a *Address) Equals(other *Address) bool {
	if a == nil || other == nil {
		return a == other
	}
	
	return a.street == other.street &&
		a.city == other.city &&
		a.state == other.state &&
		a.postalCode == other.postalCode &&
		a.country == other.country &&
		a.unit == other.unit
}

// TimeWindow represents an estimated time range
type TimeWindow struct {
	StartTime time.Time
	EndTime   time.Time
}

// NewTimeWindow creates a time window from now plus estimated minutes
func NewTimeWindow(from time.Time, estimatedMinutes int) *TimeWindow {
	// Add 10% buffer for uncertainty
	bufferMinutes := estimatedMinutes / 10
	if bufferMinutes < 5 {
		bufferMinutes = 5
	}
	
	startTime := from.Add(time.Duration(estimatedMinutes-bufferMinutes) * time.Minute)
	endTime := from.Add(time.Duration(estimatedMinutes+bufferMinutes) * time.Minute)
	
	return &TimeWindow{
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// IsWithinWindow checks if a given time is within the window
func (tw *TimeWindow) IsWithinWindow(t time.Time) bool {
	return t.After(tw.StartTime) && t.Before(tw.EndTime)
}

// HasPassed checks if the time window has passed
func (tw *TimeWindow) HasPassed(now time.Time) bool {
	return now.After(tw.EndTime)
}

// MinutesRemaining returns the minutes until the end of the window
func (tw *TimeWindow) MinutesRemaining(now time.Time) int {
	if tw.HasPassed(now) {
		return 0
	}
	return int(tw.EndTime.Sub(now).Minutes())
}

// DurationMinutes returns the total duration of the time window in minutes
func (tw *TimeWindow) DurationMinutes() int {
	return int(tw.EndTime.Sub(tw.StartTime).Minutes())
}

// String returns a human-readable time window
func (tw *TimeWindow) String() string {
	format := "3:04 PM"
	return fmt.Sprintf("%s - %s", tw.StartTime.Format(format), tw.EndTime.Format(format))
}
