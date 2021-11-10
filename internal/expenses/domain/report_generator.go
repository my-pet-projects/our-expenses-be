package domain

import (
	"time"
)

// ReportGenerator represents expense report generator.
type ReportGenerator struct {
	expenses []Expense
	filter   ExpenseFilter
}

// NewReportGenerator instantiates a new report.
func NewReportGenerator(expenses []Expense, filter ExpenseFilter) ReportGenerator {
	return ReportGenerator{
		expenses: expenses,
		filter:   filter,
	}
}

// GenerateByDateReport generates report.
func (r ReportGenerator) GenerateByDateReport() ReportByDate {
	dateCategoryExpenses := make([]*DateExpenses, 0)
	dateExpensesMap := r.prepareDateExpensesMap(r.expenses, r.filter.Interval())
	for date, expenses := range dateExpensesMap {
		categoryExpensesMap := r.buildCategoryFlatMap(expenses)
		rootCategoryExpense := r.buildCategoryHierarchy(categoryExpensesMap)

		dateExpense := &DateExpenses{
			Date:          date,
			SubCategories: rootCategoryExpense.SubCategories,
		}
		dateCategoryExpenses = append(dateCategoryExpenses, dateExpense)
	}

	report := ReportByDate{
		CategoryByDate: dateCategoryExpenses,
	}
	report.CalculateTotal()

	return report
}

func (r ReportGenerator) prepareDateExpensesMap(expenses []Expense, interval Interval) map[time.Time][]Expense {
	dateExpensesMap := make(map[time.Time][]Expense)
	for _, expense := range expenses {
		date := expense.date
		if interval == IntervalMonth {
			date = time.Date(expense.date.Year(), expense.date.Month(), 1, 0, 0, 0, 0, time.Local)
		} else if interval == IntervalYear {
			date = time.Date(expense.date.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		}
		dateExpenses := dateExpensesMap[date]
		if dateExpenses == nil {
			dateExpenses = make([]Expense, 0)
		}
		dateExpenses = append(dateExpenses, expense)
		dateExpensesMap[date] = dateExpenses
	}
	return dateExpensesMap
}

func (r ReportGenerator) buildCategoryFlatMap(expenses []Expense) map[string]*CategoryExpenses {
	categoryExpensesMap := make(map[string]*CategoryExpenses)
	for _, expense := range expenses {
		// Process category expenses.
		categoryExpenses := categoryExpensesMap[expense.category.id]
		if categoryExpenses == nil {
			categoryExpenses = &CategoryExpenses{
				Category: expense.category,
				Expenses: &[]Expense{expense},
			}
		} else {
			expenses := append(*categoryExpenses.Expenses, expense)
			categoryExpenses.Expenses = &expenses
		}
		categoryExpensesMap[expense.category.id] = categoryExpenses

		// Process parent categories expenses.
		if expense.category.IsRoot() {
			continue
		}
		for _, parentCategory := range *expense.category.parents {
			parentExpenses := &CategoryExpenses{
				Category:      parentCategory,
				SubCategories: make([]*CategoryExpenses, 0),
				Expenses:      &[]Expense{},
			}
			categoryExpensesMap[parentCategory.id] = parentExpenses
		}
	}
	return categoryExpensesMap
}

func (r ReportGenerator) buildCategoryHierarchy(flatCategoryExpensesMap map[string]*CategoryExpenses) CategoryExpenses {
	rootCategories := make([]*CategoryExpenses, 0)
	for _, categoryExpenses := range flatCategoryExpensesMap {
		if categoryExpenses.Category.IsRoot() {
			rootCategories = append(rootCategories, categoryExpenses)
		} else {
			category := flatCategoryExpensesMap[*categoryExpenses.Category.parentId]
			category.SubCategories = append(category.SubCategories, categoryExpenses)
		}
	}

	rootElement := CategoryExpenses{
		SubCategories: rootCategories,
	}
	return rootElement
}
