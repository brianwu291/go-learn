package financialservice

import (
	"testing"

	"github.com/brianwu291/go-learn/types"
	"github.com/stretchr/testify/assert"
)

func TestFinancialService_CalculateFinancial(t *testing.T) {
	tests := []struct {
		name           string
		input          types.FinancialRawInfo
		roundingDigits int
		expectedProfit float64
		expectedRatio  float64
		expectError    bool
	}{
		{
			name: "Valid calculation",
			input: types.FinancialRawInfo{
				Revenue:  1000,
				Expenses: 600,
				TaxRate:  0.2,
			},
			roundingDigits: 2,
			expectedProfit: 320.00,
			expectedRatio:  1.25,
			expectError:    false,
		},
		{
			name: "Invalid rounding digits",
			input: types.FinancialRawInfo{
				Revenue:  1000,
				Expenses: 600,
				TaxRate:  0.2,
			},
			roundingDigits: 0,
			expectError:    true,
		},
		{
			name: "Zero expenses",
			input: types.FinancialRawInfo{
				Revenue:  1000,
				Expenses: 0,
				TaxRate:  0.2,
			},
			roundingDigits: 2,
			expectedProfit: 800.00,
			expectedRatio:  1.25,
			expectError:    false,
		},
	}

	service := NewFinancialService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CalculateFinancial(tt.input, tt.roundingDigits)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedProfit, result.Profit)
			assert.Equal(t, tt.expectedRatio, result.Ratio)
		})
	}
}

func TestFinancialService_roundFloat(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		digits      int
		expected    float64
		expectError bool
	}{
		{
			name:        "Round to 2 decimal places",
			value:       3.14159,
			digits:      2,
			expected:    3.14,
			expectError: false,
		},
		{
			name:        "Round to 1 decimal place",
			value:       3.85,
			digits:      1,
			expected:    3.9,
			expectError: false,
		},
		{
			name:        "Invalid digits - zero",
			value:       3.14159,
			digits:      0,
			expectError: true,
		},
		{
			name:        "Invalid digits - negative",
			value:       3.14159,
			digits:      -1,
			expectError: true,
		},
	}

	service := &financialService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.roundFloat(tt.value, tt.digits)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}