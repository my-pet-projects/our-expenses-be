package domain

import "time"

type Report struct {
	From time.Time
	To   time.Time
}

type ReportByDate struct {
	Report
	CategoryByDate []*DateExpenses
	GrandTotal     GrandTotal
}

func (c *ReportByDate) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	for _, byDate := range c.CategoryByDate {
		dateTotal := byDate.CalculateTotal()
		grandTotal = grandTotal.Combine(dateTotal)
	}
	c.GrandTotal = grandTotal
	return c.GrandTotal
}
