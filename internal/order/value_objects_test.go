package order

import (
	"testing"
)

func TestMoneyDisplay(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		expected string
	}{
		{
			name:     "small amount in satoshis",
			amount:   500,
			expected: "500 sats",
		},
		{
			name:     "medium amount in satoshis",
			amount:   50000,
			expected: "50000 sats",
		},
		{
			name:     "amount in milliBTC",
			amount:   100000, // 1 mBTC
			expected: "1.000 mBTC",
		},
		{
			name:     "larger amount in milliBTC",
			amount:   550000, // 5.5 mBTC
			expected: "5.500 mBTC",
		},
		{
			name:     "amount in BTC",
			amount:   10000000, // 0.1 BTC
			expected: "0.10000000 BTC",
		},
		{
			name:     "full BTC",
			amount:   100000000, // 1 BTC
			expected: "1.00000000 BTC",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money := NewMoney(tt.amount)
			result := money.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestMoneyConversions(t *testing.T) {
	t.Run("from BTC", func(t *testing.T) {
		money := NewMoneyFromBTC(0.001) // 0.001 BTC
		if money.Amount() != 100000 {
			t.Errorf("expected 100000 satoshis, got %d", money.Amount())
		}
		if money.AmountInBTC() != 0.001 {
			t.Errorf("expected 0.001 BTC, got %f", money.AmountInBTC())
		}
	})
	
	t.Run("from milliBTC", func(t *testing.T) {
		money := NewMoneyFromMilliBTC(5) // 5 mBTC
		if money.Amount() != 500000 {
			t.Errorf("expected 500000 satoshis, got %d", money.Amount())
		}
		if money.AmountInMilliBTC() != 5.0 {
			t.Errorf("expected 5 mBTC, got %f", money.AmountInMilliBTC())
		}
	})
}

func TestMoneyArithmetic(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		a := NewMoney(10000) // 10,000 sats
		b := NewMoney(5000)  // 5,000 sats
		result := a.Add(b)
		
		if result.Amount() != 15000 {
			t.Errorf("expected 15000 sats, got %d", result.Amount())
		}
	})
	
	t.Run("multiplication", func(t *testing.T) {
		price := NewMoney(25000) // 25,000 sats per item
		result := price.Multiply(3)
		
		if result.Amount() != 75000 {
			t.Errorf("expected 75000 sats, got %d", result.Amount())
		}
	})
}
