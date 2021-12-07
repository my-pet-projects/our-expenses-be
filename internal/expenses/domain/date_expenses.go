package domain

import "time"

// DateExpenses represents expenses for a specific date.
type DateExpenses struct {
	Date          time.Time
	SubCategories []*CategoryExpenses
	GrandTotal    GrandTotal
	ExchangeRate  ExchangeRates
}

// CalculateTotal calculates date expenses total.
func (c *DateExpenses) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	for _, children := range c.SubCategories {
		childTotal := children.CalculateTotal()
		grandTotal = grandTotal.Combine(childTotal)
	}
	c.GrandTotal = grandTotal

	return c.GrandTotal
}
