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

func NewExpenseFilter(from time.Time, to time.Time, intervalString string) (*ExpenseFilter, error) {
	// TODO: from/to validation

	var interval Interval
	switch intervalString {
	case "day":
		interval = IntervalDay
	case "month":
		interval = IntervalMonth
	case "year":
		interval = IntervalYear
	default:
		return nil, errors.New(fmt.Sprintf("unknown interval %s", intervalString))
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
