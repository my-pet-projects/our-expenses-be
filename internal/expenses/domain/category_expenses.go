package domain

type CategoryExpenses struct {
	Category      Category
	Expenses      *[]Expense
	SubCategories []*CategoryExpenses
	Total         Total
}

func (c *CategoryExpenses) CalculateTotal() Total {
	var total Total
	if c.Expenses != nil {
		for _, expense := range *c.Expenses {
			expenseTotal := expense.CalculateTotal()
			total = total.Add(expenseTotal)
		}
	}

	for _, subCategory := range c.SubCategories {
		subCategoryTotal := subCategory.CalculateTotal()
		total = total.Add(subCategoryTotal)
	}
	c.Total = total
	return c.Total
}
