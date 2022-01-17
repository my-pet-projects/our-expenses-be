package domain

// GrandTotal holds multi total amounts.
type GrandTotal struct {
	SubTotals map[Currency]TotalInfo
	Total     Total
}

// Combine combines two grand total structs together.
func (gt GrandTotal) Combine(grandTotal GrandTotal) GrandTotal {
	if gt.SubTotals == nil {
		gt.SubTotals = make(map[Currency]TotalInfo)
	}
	for currency, currencyTotal := range grandTotal.SubTotals {
		gt.SubTotals[currency] = gt.SubTotals[currency].Add(currencyTotal)
		if currencyTotal.ConvertedTotal == nil {
			gt.Total = gt.Total.Add(&currencyTotal.OriginalTotal)
		} else {
			gt.Total = gt.Total.Add(currencyTotal.ConvertedTotal)
		}
	}

	return gt
}

// Add adds total amount to grand total based on currency.
func (gt GrandTotal) Add(t TotalInfo) GrandTotal {
	if gt.SubTotals == nil {
		gt.SubTotals = make(map[Currency]TotalInfo)
	}
	gt.SubTotals[t.OriginalTotal.Currency] = gt.SubTotals[t.OriginalTotal.Currency].Add(t)
	gt.Total = gt.Total.Add(t.ConvertedTotal)

	return gt
}
