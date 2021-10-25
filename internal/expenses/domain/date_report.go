package domain

import "time"

type Report struct {
	From time.Time
	To   time.Time
}

type ReportByDate struct {
	Report
	CategoryByDate []*DateExpenses
	Total          Total
}

func (c *ReportByDate) CalculateTotal() Total {
	var total Total
	for _, byDate := range c.CategoryByDate {
		dateTotal := byDate.CalculateTotal()
		total = total.Add(dateTotal)
	}
	c.Total = total
	return c.Total
}
