package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Total holds amount and currency data.
type Total struct {
	Sum      decimal.Decimal
	Currency string
	SumDebug string
}

// Add combines two total structs together.
func (t Total) Add(total Total) Total {
	if total.Sum == decimal.Zero {
		return t
	}

	t.Sum = t.Sum.Add(total.Sum)
	t.Currency = total.Currency

	// TODO: remove after stabilized.
	if t.SumDebug == "" {
		t.SumDebug = total.SumDebug
	} else {
		t.SumDebug = fmt.Sprintf("%s || %s ", t.SumDebug, total.SumDebug)
	}

	return t
}

// Equal returns whether the totals are equal.
func (t Total) Equal(t2 Total) bool {
	return t.Sum.Equal(t2.Sum) && t.Currency == t2.Currency
}
