package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTotalAdd_CombinesTotalsTogether(t *testing.T) {
	t.Parallel()
	// Arrange
	total1 := &Total{
		Currency: "EUR",
		Sum:      decimal.NewFromFloat(15.58),
	}
	total2 := &Total{
		Currency: "EUR",
		Sum:      decimal.NewFromFloat(15.12),
	}
	total3 := &Total{
		Currency: "EUR",
		Sum:      decimal.NewFromFloat(5.05),
	}
	total4 := &Total{
		Currency: "EUR",
		Sum:      decimal.Zero,
	}

	// Act
	result := total1.Add(total2).Add(total3).Add(total4).Add(nil)

	// Assert
	assert.Equal(t, "35.75", result.Sum.String())
	assert.Equal(t, total1.Currency, result.Currency)
}

func TestTotalEqual_ChecksEquality(t *testing.T) {
	t.Parallel()
	// Arrange
	tests := []struct {
		total1 Total
		total2 Total
		equal  bool
	}{
		{
			total1: Total{Currency: "EUR", Sum: decimal.NewFromFloat(15.58)},
			total2: Total{Currency: "EUR", Sum: decimal.NewFromFloat(15.58)},
			equal:  true,
		},
		{
			total1: Total{Currency: "EUR", Sum: decimal.NewFromFloat(15.58)},
			total2: Total{Currency: "USD", Sum: decimal.NewFromFloat(15.58)},
			equal:  false,
		},
		{
			total1: Total{Currency: "EUR", Sum: decimal.NewFromFloat(15.58)},
			total2: Total{Currency: "EUR", Sum: decimal.NewFromFloat(15)},
			equal:  false,
		},
	}

	for _, tc := range tests {
		// Act
		result := tc.total1.Equal(tc.total2)

		// Assert
		assert.Equal(t, tc.equal, result)
	}
}
