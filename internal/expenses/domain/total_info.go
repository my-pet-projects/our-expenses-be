package domain

// TotalInfo represents total info.
type TotalInfo struct {
	OriginalTotal  Total
	ConvertedTotal *Total
	ExchangeRate   *ExchangeRates
}

// Add combines two total info structs together.
func (t TotalInfo) Add(t2 TotalInfo) TotalInfo {
	origTotal := t.OriginalTotal.Add(&t2.OriginalTotal)

	if !t.canCombineWith(t2) {
		return TotalInfo{
			OriginalTotal: origTotal,
		}
	}

	convTotal := Total{}
	exchRate := &ExchangeRates{}
	if t.ConvertedTotal != nil {
		convTotal = convTotal.Add(t.ConvertedTotal)
		exchRate = t.ExchangeRate
	}

	if t2.ConvertedTotal != nil {
		convTotal = convTotal.Add(t2.ConvertedTotal)
		exchRate = t2.ExchangeRate
	}

	ti := TotalInfo{
		OriginalTotal:  origTotal,
		ConvertedTotal: &convTotal,
		ExchangeRate:   exchRate,
	}

	return ti
}

func (t TotalInfo) canCombineWith(t2 TotalInfo) bool {
	if t.ConvertedTotal == nil && t2.ConvertedTotal == nil {
		return false
	}
	if t.ExchangeRate == nil && t2.ExchangeRate == nil {
		return false
	}
	if (t.ExchangeRate == nil && t2.ExchangeRate != nil) ||
		(t.ExchangeRate != nil && t2.ExchangeRate == nil) {
		return true
	}
	if t.ExchangeRate.baseCurrency != t2.ExchangeRate.baseCurrency {
		return false
	}
	return true
}
