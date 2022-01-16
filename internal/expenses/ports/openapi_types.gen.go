// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package ports

import (
	"time"
)

// Defines values for Interval.
const (
	IntervalDay Interval = "day"

	IntervalMonth Interval = "month"

	IntervalYear Interval = "year"
)

// Category defines model for Category.
type Category struct {
	Icon *string `json:"icon,omitempty"`

	// Unique id of the category
	Id    string `json:"id"`
	Level int    `json:"level"`

	// Name of the category
	Name    string      `json:"name"`
	Parents *[]Category `json:"parents,omitempty"`
}

// CategoryExpenses defines model for CategoryExpenses.
type CategoryExpenses struct {
	Category      Category            `json:"category"`
	Expenses      *[]Expense          `json:"expenses,omitempty"`
	GrandTotal    GrandTotal          `json:"grandTotal"`
	SubCategories *[]CategoryExpenses `json:"subCategories,omitempty"`
}

// DateCategoryReport defines model for DateCategoryReport.
type DateCategoryReport struct {
	CategoryExpenses []CategoryExpenses `json:"categoryExpenses"`
	Date             time.Time          `json:"date"`
	ExchangeRates    ExchangeRates      `json:"exchangeRates"`
	GrandTotal       GrandTotal         `json:"grandTotal"`
}

// Error defines model for Error.
type Error struct {
	// Error code
	Code int32 `json:"code"`

	// Error message
	Message string `json:"message"`
}

// ExchangeRate defines model for ExchangeRate.
type ExchangeRate struct {
	BaseCurrency   string    `json:"baseCurrency"`
	Date           time.Time `json:"date"`
	Rate           string    `json:"rate"`
	TargetCurrency string    `json:"targetCurrency"`
}

// ExchangeRates defines model for ExchangeRates.
type ExchangeRates struct {
	Currency string    `json:"currency"`
	Date     time.Time `json:"date"`
	Rates    []Rate    `json:"rates"`
}

// Expense defines model for Expense.
type Expense struct {
	// Embedded struct due to allOf(#/components/schemas/NewExpense)
	NewExpense `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Category *Category `json:"category,omitempty"`

	// Unique id of the expense
	Id string `json:"id"`
}

// ExpenseReport defines model for ExpenseReport.
type ExpenseReport struct {
	DateReports []DateCategoryReport `json:"dateReports"`
	GrandTotal  GrandTotal           `json:"grandTotal"`
}

// GrandTotal defines model for GrandTotal.
type GrandTotal struct {
	Converted Total       `json:"converted"`
	Totals    []TotalInfo `json:"totals"`
}

// Interval defines model for Interval.
type Interval string

// NewExpense defines model for NewExpense.
type NewExpense struct {
	// Category ID of the expense
	CategoryId string    `json:"categoryId"`
	Comment    *string   `json:"comment,omitempty"`
	Currency   string    `json:"currency"`
	Date       time.Time `json:"date"`
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
	TotalInfo  TotalInfo `json:"totalInfo"`
	Trip       *string   `json:"trip,omitempty"`
}

// NewExpenseResponse defines model for NewExpenseResponse.
type NewExpenseResponse struct {
	// ID of the newly added expense
	Id string `json:"id"`
}

// Rate defines model for Rate.
type Rate struct {
	Currency string `json:"currency"`
	Price    string `json:"price"`
}

// Total defines model for Total.
type Total struct {
	// Total currency
	Currency string `json:"currency"`

	// Total sum amount
	Sum string `json:"sum"`
}

// TotalInfo defines model for TotalInfo.
type TotalInfo struct {
	Converted *Total        `json:"converted,omitempty"`
	Original  Total         `json:"original"`
	Rate      *ExchangeRate `json:"rate,omitempty"`
}

// AddExpenseJSONBody defines parameters for AddExpense.
type AddExpenseJSONBody NewExpense

// GenerateReportParams defines parameters for GenerateReport.
type GenerateReportParams struct {
	// from date to filter by
	From time.Time `json:"from"`

	// to date to filter by
	To time.Time `json:"to"`

	// results interval
	Interval Interval `json:"interval"`
}

// AddExpenseJSONRequestBody defines body for AddExpense for application/json ContentType.
type AddExpenseJSONRequestBody AddExpenseJSONBody
