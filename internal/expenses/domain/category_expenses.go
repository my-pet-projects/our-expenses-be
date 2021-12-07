package domain

// CategoryExpenses holds category expenses.
type CategoryExpenses struct {
	Category      Category
	Expenses      *[]Expense
	SubCategories []*CategoryExpenses
	GrandTotal    GrandTotal
}

// CalculateTotal calculates category expenses total.
func (c *CategoryExpenses) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	if c.Expenses != nil {
		for _, expense := range *c.Expenses {
			grandTotal = grandTotal.Add(expense.totalInfo)
		}
	}

	for _, subCategory := range c.SubCategories {
		subCategoryTotal := subCategory.CalculateTotal()
		grandTotal = grandTotal.Combine(subCategoryTotal)
	}

	c.GrandTotal = grandTotal
	return c.GrandTotal
}
