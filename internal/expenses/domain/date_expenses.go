package domain

import "time"

type DateExpenses struct {
	Date          time.Time
	SubCategories []*CategoryExpenses
	GrandTotal    GrandTotal
	ExchangeRate  ExchangeRates
}

func (c *DateExpenses) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	for _, children := range c.SubCategories {
		childTotal := children.CalculateTotal()
		grandTotal = grandTotal.Combine(childTotal)
	}
	c.GrandTotal = grandTotal
	return c.GrandTotal
}
