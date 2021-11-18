package domain

import "time"

type DateExpenses struct {
	Date          time.Time
	SubCategories []*CategoryExpenses
	GrandTotal    GrandTotal
}

func (c *DateExpenses) CalculateTotal() GrandTotal {
	var multiTotal GrandTotal
	for _, children := range c.SubCategories {
		childTotal := children.CalculateTotal()
		multiTotal = multiTotal.Combine(childTotal)
	}
	c.GrandTotal = multiTotal
	return c.GrandTotal
}
