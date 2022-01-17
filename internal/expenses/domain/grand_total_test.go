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
		SubTotals: map[domain.Currency]domain.TotalInfo{
			domain.Currency("EUR"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(10),
				},
			},
			domain.Currency("USD"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(25),
				},
			},
		},
	}
	gt2 := domain.GrandTotal{
		SubTotals: map[domain.Currency]domain.TotalInfo{
			domain.Currency("EUR"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(30),
				},
			},
			domain.Currency("USD"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(25),
				},
			},
			domain.Currency("SEK"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(20),
				},
			},
		},
	}
	gt3 := domain.GrandTotal{
		SubTotals: map[domain.Currency]domain.TotalInfo{
			domain.Currency("RUB"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(100),
				},
			},
		},
	}

	// Act
	result := gt1.Combine(gt2).Combine(gt3)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, decimal.NewFromInt(40), result.SubTotals["EUR"].OriginalTotal.Sum)
	assert.Equal(t, decimal.NewFromInt(50), result.SubTotals["USD"].OriginalTotal.Sum)
	assert.Equal(t, decimal.NewFromInt(20), result.SubTotals["SEK"].OriginalTotal.Sum)
	assert.Equal(t, decimal.NewFromInt(100), result.SubTotals["RUB"].OriginalTotal.Sum)
}

func TestAdd_AddsTotalToGrandTotal(t *testing.T) {
	t.Parallel()
	// Arrange
	grandTotal := domain.GrandTotal{
		SubTotals: map[domain.Currency]domain.TotalInfo{
			domain.Currency("EUR"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(10),
				},
			},
			domain.Currency("USD"): {
				OriginalTotal: domain.Total{
					Sum: decimal.NewFromInt(25),
				},
			},
		},
	}
	total1 := domain.TotalInfo{
		OriginalTotal: domain.Total{
			Sum:      decimal.NewFromInt(100),
			Currency: "SEK",
		},
	}
	total2 := domain.TotalInfo{
		OriginalTotal: domain.Total{
			Sum:      decimal.NewFromInt(50),
			Currency: "USD",
		},
	}

	// Act
	result := grandTotal.Add(total1).Add(total2)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, decimal.NewFromInt(10), result.SubTotals["EUR"].OriginalTotal.Sum)
	assert.Equal(t, decimal.NewFromInt(75), result.SubTotals["USD"].OriginalTotal.Sum)
	assert.Equal(t, decimal.NewFromInt(100), result.SubTotals["SEK"].OriginalTotal.Sum)
}
