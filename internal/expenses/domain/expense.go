package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// Expense represents a domain object.
type Expense struct {
	id         string
	categoryID string
	category   Category
	price      decimal.Decimal
	currency   string
	quantity   decimal.Decimal
	comment    *string
	date       time.Time
	createdAt  time.Time
	updatedAt  *time.Time
}

// NewExpense creates a new expense domain object.
func NewExpense(
	id string,
	categoryID string,
	price float64,
	currency string,
	quantity float64,
	comment *string,
	date time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) (*Expense, error) {
	if categoryID == "" {
		return nil, errors.New("empty categoryID")
	}

	decPrice := decimal.NewFromFloat(price)
	decQuantity := decimal.NewFromFloat(quantity)

	return &Expense{
		id:         id,
		categoryID: categoryID,
		price:      decPrice,
		currency:   currency,
		quantity:   decQuantity,
		comment:    comment,
		date:       date,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}, nil
}

// ID returns expense id.
func (e Expense) ID() string {
	return e.id
}

// CategoryID returns expense id.
func (e Expense) CategoryID() string {
	return e.categoryID
}

// Category returns expense category.
func (e Expense) Category() Category {
	return e.category
}

// SetCategory sets expense category.
func (e *Expense) SetCategory(category Category) {
	e.category = category
}

// Price returns expense price.
func (e Expense) Price() float64 {
	price, _ := e.price.Float64()
	return price
}

// Currency returns expense price.
func (e Expense) Currency() string {
	return e.currency
}

// Quantity returns expense quantity.
func (e Expense) Quantity() float64 {
	quantity, _ := e.quantity.Float64()
	return quantity
}

// Comment returns expense comment.
func (e Expense) Comment() *string {
	return e.comment
}

// Date returns expense date.
func (e Expense) Date() time.Time {
	return e.date
}

// CreatedAt returns expense creation date.
func (e Expense) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns expense update date.
func (e Expense) UpdatedAt() *time.Time {
	return e.updatedAt
}
