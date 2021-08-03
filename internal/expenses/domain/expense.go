package domain

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Expense represents a domain object.
type Expense struct {
	id         string
	categoryID string
	price      decimal.Decimal
	currency   string
	quantity   int
	comment    *string
	date       time.Time
	createdAt  time.Time
	updatedAt  *time.Time
}

// NewExpense creates a new expense domain object.
func NewExpense(
	id string,
	categoryID string,
	price string,
	currency string,
	quantity int,
	comment *string,
	date time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) (*Expense, error) {
	decPrice, decPriceErr := decimal.NewFromString(price)
	if decPriceErr != nil {
		return nil, errors.Wrap(decPriceErr, "price convert")
	}
	return &Expense{
		id:         id,
		categoryID: categoryID,
		price:      decPrice,
		currency:   currency,
		quantity:   quantity,
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

// Price returns expense price.
func (e Expense) Price() string {
	return e.price.String()
}

// Currency returns expense price.
func (e Expense) Currency() string {
	return e.currency
}

// Quantity returns expense quantity.
func (e Expense) Quantity() int {
	return e.quantity
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
