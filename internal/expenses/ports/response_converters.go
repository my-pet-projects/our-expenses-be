package ports

import "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"

func reportToResponse(domainObj domain.ReportByDate) ExpenseReport {
	dateCategoryReport := []DateCategoryReport{}
	for _, categoryByDate := range domainObj.CategoryByDate {
		categoryExpenses := make([]CategoryExpenses, 0)
		for _, category := range categoryByDate.SubCategories {
			categoryExpenses = append(categoryExpenses, categoryExpensesToResponse(*category))
		}

		dateCategoryReport = append(dateCategoryReport, DateCategoryReport{
			Date:             categoryByDate.Date,
			CategoryExpenses: categoryExpenses,
			Total:            totalToResponse(categoryByDate.Total),
		})
	}

	report := ExpenseReport{
		DateReports: dateCategoryReport,
		Total:       totalToResponse(domainObj.Total),
	}
	return report
}

func categoryExpensesToResponse(domainObj domain.CategoryExpenses) CategoryExpenses {
	expenses := []Expense{}
	if domainObj.Expenses != nil {
		for _, domainExpense := range *domainObj.Expenses {
			expenses = append(expenses, expenseToResponse(domainExpense))
		}
	}

	categoryExpenses := []CategoryExpenses{}
	for _, subCategory := range domainObj.SubCategories {
		categoryExpenses = append(categoryExpenses, categoryExpensesToResponse(*subCategory))
	}

	response := CategoryExpenses{
		Category: categoryToResponse(domainObj.Category),
		Total:    totalToResponse(domainObj.Total),
	}

	if len(expenses) != 0 {
		response.Expenses = &expenses
	}
	if len(categoryExpenses) != 0 {
		response.SubCategories = &categoryExpenses
	}

	return response
}

func categoryToResponse(domainObj domain.Category) Category {
	return Category{
		Id:    domainObj.ID(),
		Name:  domainObj.Name(),
		Icon:  domainObj.Icon(),
		Level: domainObj.Level(),
	}
}

func expenseToResponse(domainObj domain.Expense) Expense {
	return Expense{
		Id: domainObj.ID(),
		NewExpense: NewExpense{
			Comment:  domainObj.Comment(),
			Currency: domainObj.Currency(),
			Date:     domainObj.Date(),
			Price:    domainObj.Price(),
			Quantity: domainObj.Quantity(),
		},
		Category: categoryToResponse(domainObj.Category()),
	}
}

func totalToResponse(domainTotal domain.Total) Total {
	return Total{
		Debug:    domainTotal.SumDebug,
		Sum:      domainTotal.Sum.String(),
		Currency: domainTotal.Currency,
	}
}
