package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type CategoriesByDate struct {
	SubCategories []*CategoryExpenses
	Total         Total
	Date          time.Time
}

func (c *CategoriesByDate) Sum() string {

	for _, children := range c.SubCategories {

		children.Sum()

		c.Total.Sum = c.Total.Sum.Add(children.Total.Sum)
		c.Total.SumDebug = c.Total.SumDebug + children.Total.SumDebug
	}

	return c.Total.SumDebug
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
	CategoryByDate []*CategoriesByDate
	Total          Total
}

func (c *ReportByDate) Sum() string {

	for _, byDate := range c.CategoryByDate {

		byDate.Sum()

		c.Total.Sum = c.Total.Sum.Add(byDate.Total.Sum)
		c.Total.SumDebug = c.Total.SumDebug + byDate.Total.SumDebug
	}

	return c.Total.SumDebug
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

type Total struct {
	SumDebug string
	Sum      decimal.Decimal
}

type CategoryExpenses struct {
	Expenses      *[]Expense
	SubCategories []*CategoryExpenses
	Category      Category
	Total         Total
}

func (c *CategoryExpenses) Sum() string {

	if c.Expenses != nil {
		for _, expense := range *c.Expenses {
			c.Total.Sum = c.Total.Sum.Add(expense.quantity.Mul(expense.price))
			c.Total.SumDebug = fmt.Sprintf("%s | %s: %s*%s %s", c.Total.SumDebug, expense.category.name, expense.price, expense.quantity, expense.currency)
		}
	}

	for _, children := range c.SubCategories {
		children.Sum()
		c.Total.Sum = c.Total.Sum.Add(children.Total.Sum)
		c.Total.SumDebug = c.Total.SumDebug + children.Total.SumDebug
	}

	return c.Total.SumDebug
}

// GenerateByDateReport generates report.
func (r ReportGenerator) GenerateByDateReport() ReportByDate {
	report := ReportByDate{}

	// Expenses by date map.
	dateExpenseMap := make(map[time.Time][]Expense)
	for _, expense := range r.expenses {
		expenses := dateExpenseMap[expense.date]
		if expenses == nil {
			expenses = make([]Expense, 0)
		}
		expenses = append(expenses, expense)
		dateExpenseMap[expense.date] = expenses
	}

	for date := range dateExpenseMap {
		dateExpenses := dateExpenseMap[date]

		// categoryExpenses := make([]CategoryExpenses, 0)
		categoryExpensesMap := make(map[string]*CategoryExpenses)

		categoryByDate := CategoriesByDate{
			Date: date,
		}

		rootCategoryExpense := &CategoryExpenses{
			SubCategories: []*CategoryExpenses{},
		}

		for _, expense := range dateExpenses {
			// Expense category.
			catExp := categoryExpensesMap[expense.category.id]
			if catExp == nil {
				expenseCategory := &CategoryExpenses{
					Category: expense.category,
					Expenses: &[]Expense{expense},
					// Children: nil,
				}
				catExp = expenseCategory
			} else {
				currentExpenses := *catExp.Expenses
				newExp := append(currentExpenses, expense)
				catExp.Expenses = &newExp
			}
			categoryExpensesMap[expense.category.id] = catExp

			// Parents.
			parents := expense.category.parents
			for _, parentCategory := range *parents {
				current := &CategoryExpenses{
					Category:      parentCategory,
					SubCategories: make([]*CategoryExpenses, 0),
				}
				categoryExpensesMap[parentCategory.id] = current
			}
		}

		for categoryExpenseID := range categoryExpensesMap {
			elem := categoryExpensesMap[categoryExpenseID]

			if elem.Category.parentId == nil || *elem.Category.parentId == "" {
				rootCategoryExpense.SubCategories = append(rootCategoryExpense.SubCategories, elem)
			} else {
				cur := categoryExpensesMap[*elem.Category.parentId]
				cur.SubCategories = append(cur.SubCategories, elem)
			}
		}

		categoryByDate.Date = date
		categoryByDate.Total = rootCategoryExpense.Total
		categoryByDate.SubCategories = rootCategoryExpense.SubCategories

		report.CategoryByDate = append(report.CategoryByDate, &categoryByDate)
	}

	report.Sum()

	return report
}

func getDateCategories(dateMap map[time.Time]map[Category][]Expense, date time.Time) map[Category][]Expense {
	dateCategories := dateMap[date]
	if dateCategories == nil {
		dateCategories = make(map[Category][]Expense)
	}
	return dateCategories
}

func getCategoryExpenses(categoryMap map[Category][]Expense, category Category) []Expense {
	categoryExpenses := categoryMap[category]
	if categoryExpenses == nil {
		categoryExpenses = make([]Expense, 0)
	}
	return categoryExpenses
}
