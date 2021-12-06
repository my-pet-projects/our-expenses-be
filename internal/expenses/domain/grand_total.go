package domain

// GrandTotal holds multi total amounts.
type GrandTotal struct {
	TotalInfos     map[Currency]TotalInfo
	ConvertedTotal Total
}

// Combine combines two grand total structs together.
func (gt GrandTotal) Combine(grandTotal GrandTotal) GrandTotal {
	if gt.TotalInfos == nil {
		gt.TotalInfos = make(map[Currency]TotalInfo, 0)
	}
	for currency, currencyTotal := range grandTotal.TotalInfos {
		gt.TotalInfos[currency] = gt.TotalInfos[currency].Add(currencyTotal)
		gt.ConvertedTotal = gt.ConvertedTotal.Add(currencyTotal.ConvertedTotal)
	}
	return gt
}

// Add adds total amount to grand total based on currency.
func (gt GrandTotal) Add(t TotalInfo) GrandTotal {
	if gt.TotalInfos == nil {
		gt.TotalInfos = make(map[Currency]TotalInfo, 0)
	}
	gt.TotalInfos[t.OriginalTotal.Currency] = gt.TotalInfos[t.OriginalTotal.Currency].Add(t)
	gt.ConvertedTotal = gt.ConvertedTotal.Add(t.ConvertedTotal)
	return gt
}
