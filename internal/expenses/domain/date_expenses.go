package domain

import "time"

type DateExpenses struct {
	Date          time.Time
	SubCategories []*CategoryExpenses
	Total         Total
}

func (c *DateExpenses) CalculateTotal() Total {
	var total Total
	for _, children := range c.SubCategories {
		childTotal := children.CalculateTotal()
		total = total.Add(childTotal)
	}
	c.Total = total
	return c.Total
}
