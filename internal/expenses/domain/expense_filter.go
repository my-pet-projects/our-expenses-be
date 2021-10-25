package domain

import "time"

// ExpenseFilter represents expense filter.
type ExpenseFilter struct {
	From time.Time
	To   time.Time
}
