package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExpenseFilter_ValidArgs_InstantiatesFilter(t *testing.T) {
	t.Parallel()
	// Arrange
	from := time.Date(2021, 7, 2, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	type test struct {
		intervalString string
		interval       Interval
	}
	tests := []test{
		{
			intervalString: "day",
			interval:       IntervalDay,
		},
		{
			intervalString: "month",
			interval:       IntervalMonth,
		},
		{
			intervalString: "year",
			interval:       IntervalYear,
		},
	}
	// Act
	for _, tc := range tests {
		// Act
		res, resErr := NewExpenseFilter(from, to, tc.intervalString)

		// Assert
		assert.Nil(t, resErr)
		assert.NotNil(t, res)
		assert.Equal(t, from, res.From())
		assert.Equal(t, to, res.To())
		assert.Equal(t, tc.interval, res.Interval())
	}
}

func TestNewExpenseFilter_InvalidArgs_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	type test struct {
		from     time.Time
		to       time.Time
		interval string
	}
	tests := []test{
		{
			from:     time.Date(2021, 7, 2, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC),
			interval: "invalid interval",
		},
		{
			from:     time.Date(2021, 7, 2, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2020, 8, 2, 0, 0, 0, 0, time.UTC),
			interval: "day",
		},
		{
			from:     time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC),
			interval: "day",
		},
	}

	for _, tc := range tests {
		// Act
		res, resErr := NewExpenseFilter(tc.from, tc.to, tc.interval)

		// Assert
		assert.NotNil(t, resErr)
		assert.Nil(t, res)
	}
}
