package domain

// GrandTotal holds multi total amounts.
type GrandTotal struct {
	Totals map[Currency]Total
}

// Combine combines two grand total structs together.
func (gt GrandTotal) Combine(grandTotal GrandTotal) GrandTotal {
	if gt.Totals == nil {
		gt.Totals = make(map[Currency]Total, 0)
	}
	for currency, currencyTotal := range grandTotal.Totals {
		gt.Totals[currency] = gt.Totals[currency].Add(currencyTotal)
	}
	return gt
}

// Add adds total amount to grand total based on currency.
func (gt GrandTotal) Add(t Total) GrandTotal {
	if gt.Totals == nil {
		gt.Totals = make(map[Currency]Total, 0)
	}
	gt.Totals[t.Currency] = gt.Totals[t.Currency].Add(t)
	return gt
}
