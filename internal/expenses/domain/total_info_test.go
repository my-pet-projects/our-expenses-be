package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

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
				ExchangeRate:   &ExchangeRate{},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: nil,
				ExchangeRate:   &ExchangeRate{},
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
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRate{baseCurrency: "USD"},
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
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{},
				ConvertedTotal: &Total{},
				ExchangeRate:   &ExchangeRate{baseCurrency: "USD"},
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
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(10)},
				ConvertedTotal: nil,
				ExchangeRate:   nil,
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(15)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
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
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(0)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(20)},
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
		},
		{
			t1: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(5)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(50)},
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			t2: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(1)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(7)},
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
			sum: TotalInfo{
				OriginalTotal:  Total{Sum: decimal.NewFromInt(6)},
				ConvertedTotal: &Total{Sum: decimal.NewFromInt(57)},
				ExchangeRate:   &ExchangeRate{baseCurrency: "EUR"},
			},
		},
	}
	// Act & Assert
	for _, tc := range tests {
		res := tc.t1.Add(tc.t2)
		assert.Equal(t, tc.sum, res)
	}
}
