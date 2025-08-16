package merchant

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMerchant(t *testing.T) {
	merchant := NewMerchant("Pizza Palace", "Best pizza in town")

	assert.NotNil(t, merchant)
	assert.NotEmpty(t, merchant.ID())
	assert.Equal(t, "Pizza Palace", merchant.Name())
	assert.Equal(t, "Best pizza in town", merchant.Description())
	assert.True(t, merchant.IsActive())
	assert.Equal(t, 30, merchant.PreparationTime())
	assert.NotNil(t, merchant.OperatingHours())
	assert.False(t, merchant.CreatedAt().IsZero())
	assert.False(t, merchant.UpdatedAt().IsZero())
}

func TestOperatingHours_NewOperatingHours(t *testing.T) {
	tests := []struct {
		name        string
		openHour    int
		closeHour   int
		daysOpen    []time.Weekday
		expectError bool
	}{
		{
			name:        "valid hours",
			openHour:    9,
			closeHour:   17,
			daysOpen:    []time.Weekday{time.Monday, time.Tuesday},
			expectError: false,
		},
		{
			name:        "invalid open hour - negative",
			openHour:    -1,
			closeHour:   17,
			daysOpen:    []time.Weekday{time.Monday},
			expectError: true,
		},
		{
			name:        "invalid open hour - too high",
			openHour:    24,
			closeHour:   17,
			daysOpen:    []time.Weekday{time.Monday},
			expectError: true,
		},
		{
			name:        "invalid close hour - negative",
			openHour:    9,
			closeHour:   -1,
			daysOpen:    []time.Weekday{time.Monday},
			expectError: true,
		},
		{
			name:        "invalid close hour - too high",
			openHour:    9,
			closeHour:   24,
			daysOpen:    []time.Weekday{time.Monday},
			expectError: true,
		},
		{
			name:        "no days open",
			openHour:    9,
			closeHour:   17,
			daysOpen:    []time.Weekday{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hours, err := NewOperatingHours(tt.openHour, tt.closeHour, tt.daysOpen)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidOperatingHours, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, time.Duration(tt.openHour)*time.Hour, hours.OpenTime)
				assert.Equal(t, time.Duration(tt.closeHour)*time.Hour, hours.CloseTime)
				assert.Equal(t, tt.daysOpen, hours.DaysOpen)
			}
		})
	}
}

func TestOperatingHours_IsOpenAt(t *testing.T) {
	// Create operating hours: 9 AM to 5 PM, Monday to Friday
	hours, err := NewOperatingHours(9, 17, []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
	})
	require.NoError(t, err)

	tests := []struct {
		name     string
		testTime time.Time
		expected bool
	}{
		{
			name:     "open during business hours - Monday 10 AM",
			testTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
		{
			name:     "closed before opening - Monday 8 AM",
			testTime: time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC), // Monday
			expected: false,
		},
		{
			name:     "closed after closing - Monday 6 PM",
			testTime: time.Date(2024, 1, 1, 18, 0, 0, 0, time.UTC), // Monday
			expected: false,
		},
		{
			name:     "closed on weekend - Saturday 10 AM",
			testTime: time.Date(2024, 1, 6, 10, 0, 0, 0, time.UTC), // Saturday
			expected: false,
		},
		{
			name:     "closed on Sunday",
			testTime: time.Date(2024, 1, 7, 10, 0, 0, 0, time.UTC), // Sunday
			expected: false,
		},
		{
			name:     "open at opening time - Monday 9 AM",
			testTime: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
		{
			name:     "open at closing time - Monday 5 PM",
			testTime: time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC), // Monday
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hours.IsOpenAt(tt.testTime)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOperatingHours_IsOpenAt_OvernightHours(t *testing.T) {
	// Create overnight operating hours: 10 PM to 6 AM, Friday to Saturday
	hours, err := NewOperatingHours(22, 6, []time.Weekday{time.Friday, time.Saturday})
	require.NoError(t, err)

	tests := []struct {
		name     string
		testTime time.Time
		expected bool
	}{
		{
			name:     "open late night - Friday 11 PM",
			testTime: time.Date(2024, 1, 5, 23, 0, 0, 0, time.UTC), // Friday
			expected: true,
		},
		{
			name:     "open early morning - Saturday 3 AM",
			testTime: time.Date(2024, 1, 6, 3, 0, 0, 0, time.UTC), // Saturday
			expected: true,
		},
		{
			name:     "closed during day - Friday 2 PM",
			testTime: time.Date(2024, 1, 5, 14, 0, 0, 0, time.UTC), // Friday
			expected: false,
		},
		{
			name:     "closed on non-operating day - Monday 11 PM",
			testTime: time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC), // Monday
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hours.IsOpenAt(tt.testTime)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMerchant_CanAcceptOrders(t *testing.T) {
	merchant := NewMerchant("Test Merchant", "Description")

	t.Run("active merchant during operating hours", func(t *testing.T) {
		// Set operating hours to include current time (using a mock time would be better in real code)
		hours, err := NewOperatingHours(0, 23, []time.Weekday{
			time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
			time.Friday, time.Saturday, time.Sunday,
		})
		require.NoError(t, err)
		merchant.UpdateOperatingHours(hours)

		err = merchant.CanAcceptOrders()
		assert.NoError(t, err)
	})

	t.Run("inactive merchant", func(t *testing.T) {
		merchant.Deactivate()
		err := merchant.CanAcceptOrders()
		assert.Error(t, err)
		assert.Equal(t, ErrMerchantNotActive, err)
	})

	t.Run("active merchant outside operating hours", func(t *testing.T) {
		merchant.Activate()
		// Set very restrictive hours that definitely don't include current time
		hours, err := NewOperatingHours(1, 2, []time.Weekday{time.Monday})
		require.NoError(t, err)
		merchant.UpdateOperatingHours(hours)

		// Test on a Tuesday (outside operating days)
		testTime := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC) // Tuesday
		err = merchant.CanAcceptOrdersAt(testTime)
		assert.Error(t, err)
		assert.Equal(t, ErrOutsideOperatingHours, err)
	})
}

func TestMerchant_EstimatePreparationTime(t *testing.T) {
	merchant := NewMerchant("Test Merchant", "Description")
	assert.Equal(t, 30, merchant.PreparationTime()) // Default

	tests := []struct {
		name      string
		itemCount int
		expected  int
	}{
		{
			name:      "single item",
			itemCount: 1,
			expected:  30,
		},
		{
			name:      "two items",
			itemCount: 2,
			expected:  32, // 30 + 2
		},
		{
			name:      "five items",
			itemCount: 5,
			expected:  38, // 30 + 8 (4 additional items * 2 minutes each)
		},
		{
			name:      "zero items (edge case)",
			itemCount: 0,
			expected:  30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := merchant.EstimatePreparationTime(tt.itemCount)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMerchant_UpdatePreparationTime(t *testing.T) {
	merchant := NewMerchant("Test Merchant", "Description")

	t.Run("valid preparation time", func(t *testing.T) {
		err := merchant.UpdatePreparationTime(45)
		assert.NoError(t, err)
		assert.Equal(t, 45, merchant.PreparationTime())
	})

	t.Run("invalid preparation time - too low", func(t *testing.T) {
		err := merchant.UpdatePreparationTime(0)
		assert.Error(t, err)
		// Should remain unchanged
		assert.Equal(t, 45, merchant.PreparationTime())
	})

	t.Run("invalid preparation time - too high", func(t *testing.T) {
		err := merchant.UpdatePreparationTime(300) // 5 hours
		assert.Error(t, err)
		// Should remain unchanged
		assert.Equal(t, 45, merchant.PreparationTime())
	})

	t.Run("boundary values", func(t *testing.T) {
		// Minimum valid value
		err := merchant.UpdatePreparationTime(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, merchant.PreparationTime())

		// Maximum valid value
		err = merchant.UpdatePreparationTime(240)
		assert.NoError(t, err)
		assert.Equal(t, 240, merchant.PreparationTime())
	})
}

func TestMerchant_ActivateDeactivate(t *testing.T) {
	merchant := NewMerchant("Test Merchant", "Description")
	initialTime := merchant.UpdatedAt()

	// Initially active
	assert.True(t, merchant.IsActive())

	// Wait a moment to ensure timestamp changes
	time.Sleep(time.Millisecond)

	t.Run("deactivate merchant", func(t *testing.T) {
		merchant.Deactivate()
		assert.False(t, merchant.IsActive())
		assert.True(t, merchant.UpdatedAt().After(initialTime))
	})

	t.Run("reactivate merchant", func(t *testing.T) {
		time.Sleep(time.Millisecond)
		deactivationTime := merchant.UpdatedAt()
		
		merchant.Activate()
		assert.True(t, merchant.IsActive())
		assert.True(t, merchant.UpdatedAt().After(deactivationTime))
	})
}

func TestMerchant_UpdateOperatingHours(t *testing.T) {
	merchant := NewMerchant("Test Merchant", "Description")
	initialTime := merchant.UpdatedAt()

	// Wait a moment to ensure timestamp changes
	time.Sleep(time.Millisecond)

	newHours, err := NewOperatingHours(8, 20, []time.Weekday{time.Monday, time.Tuesday})
	require.NoError(t, err)

	merchant.UpdateOperatingHours(newHours)

	assert.Equal(t, &newHours, merchant.OperatingHours())
	assert.True(t, merchant.UpdatedAt().After(initialTime))
}
