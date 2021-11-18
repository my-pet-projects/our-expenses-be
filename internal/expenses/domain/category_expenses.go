package domain

type CategoryExpenses struct {
	Category      Category
	Expenses      *[]Expense
	SubCategories []*CategoryExpenses
	GrandTotal    GrandTotal
}

func (c *CategoryExpenses) CalculateTotal() GrandTotal {
	var grandTotal GrandTotal
	if c.Expenses != nil {
		for _, expense := range *c.Expenses {
			expenseTotal := expense.CalculateTotal()
			grandTotal = grandTotal.Add(expenseTotal)
		}
	}

	for _, subCategory := range c.SubCategories {
		subCategoryTotal := subCategory.CalculateTotal()
		grandTotal = grandTotal.Combine(subCategoryTotal)
	}
	c.GrandTotal = grandTotal
	return c.GrandTotal
}
