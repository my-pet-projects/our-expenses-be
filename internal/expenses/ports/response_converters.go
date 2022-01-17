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
			GrandTotal:       grandTotalToResponse(categoryByDate.GrandTotal),
			ExchangeRates:    exchangeRatesToResponse(categoryByDate.ExchangeRate),
		})
	}

	report := ExpenseReport{
		DateReports: dateCategoryReport,
		GrandTotal:  grandTotalToResponse(domainObj.GrandTotal),
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
		Category:   categoryToResponse(domainObj.Category),
		GrandTotal: grandTotalToResponse(domainObj.GrandTotal),
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
			CategoryId: domainObj.Category().ID(),
			Comment:    domainObj.Comment(),
			Currency:   domainObj.Currency(),
			Date:       domainObj.Date(),
			Price:      domainObj.Price(),
			Quantity:   domainObj.Quantity(),
			TotalInfo:  totalInfoToResponse(domainObj.TotalInfo()),
		},
	}
}

func grandTotalToResponse(domainObj domain.GrandTotal) GrandTotal {
	subTotals := []TotalInfo{}
	for _, totalInfo := range domainObj.SubTotals {
		subTotals = append(subTotals, totalInfoToResponse(totalInfo))
	}
	total := &domainObj.Total
	return GrandTotal{
		SubTotals: subTotals,
		Total:     *totalToResponse(total),
	}
}

func totalToResponse(domainTotal *domain.Total) *Total {
	if domainTotal == nil {
		return nil
	}
	return &Total{
		Sum:      domainTotal.Sum.Round(2).String(),
		Currency: string(domainTotal.Currency),
	}
}

func totalInfoToResponse(domainObj domain.TotalInfo) TotalInfo {
	convertedTotal := domainObj.ConvertedTotal
	ti := TotalInfo{
		Converted: totalToResponse(convertedTotal),
		Original:  *totalToResponse(&domainObj.OriginalTotal),
	}
	if domainObj.ExchangeRate != nil {
		rate := exchangeRateToResponse(*domainObj.ExchangeRate)
		ti.Rate = &rate
	}
	return ti
}

func exchangeRatesToResponse(domainObj domain.ExchangeRates) ExchangeRates {
	rates := make([]Rate, 0)
	for currency, rate := range domainObj.Rates() {
		rates = append(rates, Rate{
			Currency: string(currency),
			Price:    rate.Round(2).String(),
		})
	}
	exchRate := ExchangeRates{
		Date:     domainObj.Date(),
		Currency: string(domainObj.BaseCurrency()),
		Rates:    rates,
	}
	return exchRate
}

func exchangeRateToResponse(domainObj domain.ExchangeRate) ExchangeRate {
	exchRate := ExchangeRate{
		Date:           domainObj.Date(),
		BaseCurrency:   string(domainObj.BaseCurrency()),
		Rate:           domainObj.Rate().Round(2).String(),
		TargetCurrency: string(domainObj.TargetCurrency()),
	}
	return exchRate
}
