// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package ports

import (
	"time"
)

// Error defines model for Error.
type Error struct {
	// Error code
	Code int32 `json:"code"`

	// Error message
	Message string `json:"message"`
}

// Expense defines model for Expense.
type Expense struct {
	// Embedded struct due to allOf(#/components/schemas/NewExpense)
	NewExpense `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	// Unique id of the expense
	Id string `json:"id"`
}

// ExpenseReport defines model for ExpenseReport.
type ExpenseReport struct {
	Expenses []Expense `json:"expenses"`
}

// NewExpense defines model for NewExpense.
type NewExpense struct {
	// Category ID of the expense
	CategoryId string    `json:"categoryId"`
	Comment    *string   `json:"comment,omitempty"`
	Currency   string    `json:"currency"`
	Date       time.Time `json:"date"`
	Price      float64   `json:"price"`
	Quantity   float64   `json:"quantity"`
}

// NewExpenseResponse defines model for NewExpenseResponse.
type NewExpenseResponse struct {
	// ID of the newly added expense
	Id string `json:"id"`
}

// AddExpenseJSONBody defines parameters for AddExpense.
type AddExpenseJSONBody NewExpense

// AddExpenseJSONRequestBody defines body for AddExpense for application/json ContentType.
type AddExpenseJSONRequestBody AddExpenseJSONBody
