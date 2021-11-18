package domain_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

func TestCombine_CombinesGrandTotals(t *testing.T) {
	t.Parallel()
	// Arrange
	gt1 := domain.GrandTotal{
		Totals: map[domain.Currency]domain.Total{
			domain.Currency("EUR"): {
				Sum: decimal.NewFromInt(10),
			},
			domain.Currency("USD"): {
				Sum: decimal.NewFromInt(25),
			},
		},
	}
	gt2 := domain.GrandTotal{
		Totals: map[domain.Currency]domain.Total{
			domain.Currency("EUR"): {
				Sum: decimal.NewFromInt(30),
			},
			domain.Currency("USD"): {
				Sum: decimal.NewFromInt(25),
			},
			domain.Currency("SEK"): {
				Sum: decimal.NewFromInt(20),
			},
		},
	}
	gt3 := domain.GrandTotal{
		Totals: map[domain.Currency]domain.Total{
			domain.Currency("RUB"): {
				Sum: decimal.NewFromInt(100),
			},
		},
	}

	// Act
	result := gt1.Combine(gt2).Combine(gt3)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, decimal.NewFromInt(40), result.Totals["EUR"].Sum)
	assert.Equal(t, decimal.NewFromInt(50), result.Totals["USD"].Sum)
	assert.Equal(t, decimal.NewFromInt(20), result.Totals["SEK"].Sum)
	assert.Equal(t, decimal.NewFromInt(100), result.Totals["RUB"].Sum)
}

func TestAdd_AddsTotalToGrandTotal(t *testing.T) {
	t.Parallel()
	// Arrange
	grandTotal := domain.GrandTotal{
		Totals: map[domain.Currency]domain.Total{
			domain.Currency("EUR"): {
				Sum: decimal.NewFromInt(10),
			},
			domain.Currency("USD"): {
				Sum: decimal.NewFromInt(25),
			},
		},
	}
	total1 := domain.Total{
		Sum:      decimal.NewFromInt(100),
		Currency: "SEK",
	}
	total2 := domain.Total{
		Sum:      decimal.NewFromInt(50),
		Currency: "USD",
	}

	// Act
	result := grandTotal.Add(total1).Add(total2)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, decimal.NewFromInt(10), result.Totals["EUR"].Sum)
	assert.Equal(t, decimal.NewFromInt(75), result.Totals["USD"].Sum)
	assert.Equal(t, decimal.NewFromInt(100), result.Totals["SEK"].Sum)
}
