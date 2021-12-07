package domain

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// ExpenseFilter represents expense filter.
type ExpenseFilter struct {
	from     time.Time
	to       time.Time
	interval Interval
}

// NewExpenseFilter instantiates expense filter.
func NewExpenseFilter(from time.Time, to time.Time, intervalString string) (*ExpenseFilter, error) {
	if from.After(to) {
		return nil, errors.New("'from' date could not be after 'to' date")
	}
	if from.Equal(to) {
		return nil, errors.New("'from' is equal to 'to' date")
	}

	var interval Interval
	switch intervalString {
	case "day":
		interval = IntervalDay
	case "month":
		interval = IntervalMonth
	case "year":
		interval = IntervalYear
	default:
		return nil, fmt.Errorf("unknown interval %s", intervalString)
	}

	filter := &ExpenseFilter{
		from:     from,
		to:       to,
		interval: interval,
	}

	return filter, nil
}

// Interval returns expense filter interval.
func (f ExpenseFilter) Interval() Interval {
	return f.interval
}

// From returns expense filter from date.
func (f ExpenseFilter) From() time.Time {
	return f.from
}

// To returns expense filter to date.
func (f ExpenseFilter) To() time.Time {
	return f.to
}
