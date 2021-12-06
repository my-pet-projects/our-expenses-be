package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// nolint:funlen
func TestTotalInfoAdd_CombinesTwoStructs(t *testing.T) {
	t.Parallel()
	// Arrange
	type test struct {
		t1  TotalInfo
		t2  TotalInfo
		sum TotalInfo
	}
	tests := []test{
		{
			t1: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: nil,
				ExchangeRate:   &ExchangeRates{},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: nil,
				ExchangeRate:   &ExchangeRates{},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(0)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(10)},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRates{baseCurrency: "USD"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(15)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: &Total{},
				ExchangeRate:   nil,
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{},
				ExchangeRate:   nil,
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRates{baseCurrency: "USD"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(10)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(15)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
			t2: TotalInfo{
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(0)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(50)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(1)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(7)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(6)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(57)},
				ExchangeRate:   &ExchangeRates{baseCurrency: "EUR"},
			},
		},
	}
	// Act & Assert
	for _, tc := range tests {
		res := tc.t1.Add(tc.t2)
		assert.Equal(t, tc.sum, res)
	}
}
