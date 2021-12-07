package domain

import "time"

// Report holds report properties.
type Report struct {
	From time.Time
	To   time.Time
}

// ReportByDate represents report by date.
type ReportByDate struct {
	Report
	CategoryByDate []*DateExpenses
	GrandTotal     GrandTotal
}

// CalculateTotal calcultares report total.
func (c *ReportByDate) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	for _, byDate := range c.CategoryByDate {
		dateTotal := byDate.CalculateTotal()
		grandTotal = grandTotal.Combine(dateTotal)
	}
	c.GrandTotal = grandTotal

	return c.GrandTotal
}
