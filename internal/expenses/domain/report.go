package domain

import (
	"fmt"
	"time"
)

type CategoriesByDate struct {
	ExpensesByCategory []ExpensesByCategory
	Date               time.Time
}

type ExpensesByCategory struct {
	Expenses []Expense
	Category Category
}

type Report struct {
	From time.Time
	To   time.Time
}
type ReportByDate struct {
	Report
	CategoryByDate []CategoriesByDate
}

// ReportGenerator represents expense report generator.
type ReportGenerator struct {
	expenses []Expense
}

// NewReportGenerator instantiates a new report.
func NewReportGenerator(expenses []Expense) ReportGenerator {
	return ReportGenerator{
		expenses: expenses,
	}
}

// GenerateReport generates report.
func (r ReportGenerator) GenerateReport() ReportByDate {
	dateCategoryMap := make(map[time.Time]map[Category][]Expense)

	for _, expense := range r.expenses {

		if expense.categoryID == "610ab575bb6f9675995afcea" {
			fmt.Print("")
		}

		category := expense.category
		date := expense.date

		dateCategories := dateCategoryMap[date]
		if dateCategories == nil {
			dateCategories = make(map[Category][]Expense)
		}

		var parentCategory Category
		for _, c := range *category.parents {
			if c.level == category.level-1 {
				parentCategory = c
			}
		}
		// parentCategory.SetParents(category.parents)

		categoryExpenses := dateCategories[parentCategory]
		if categoryExpenses == nil {
			categoryExpenses = make([]Expense, 0)
		}
		categoryExpenses = append(categoryExpenses, expense)

		dateCategories[parentCategory] = categoryExpenses
		dateCategoryMap[date] = dateCategories
	}

	categoriesByDate := make([]CategoriesByDate, 0)
	for date, dateCategories := range dateCategoryMap {
		expensesByCategory := make([]ExpensesByCategory, 0)
		for category, expenses := range dateCategories {
			expensesByCategory = append(expensesByCategory, ExpensesByCategory{
				Category: category,
				Expenses: expenses,
			})
		}

		categoriesByDate = append(categoriesByDate, CategoriesByDate{
			Date:               date,
			ExpensesByCategory: expensesByCategory,
		})
	}

	return ReportByDate{
		Report: Report{
			To: time.Now(),
		},
		CategoryByDate: categoriesByDate,
	}
}
