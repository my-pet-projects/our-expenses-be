package domain

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// Expense represents a domain object.
type Expense struct {
	id        string
	category  Category
	price     decimal.Decimal
	currency  string
	quantity  decimal.Decimal
	comment   *string
	trip      *string
	date      time.Time
	createdAt time.Time
	createdBy string
	updatedAt *time.Time
	updatedBy *string
}

// SetCreateMetadata sets expense create metadata.
func SetCreateMetadata(createdBy string, createdAt time.Time) func(*Expense) {
	return func(e *Expense) {
		e.createdBy = createdBy
		e.createdAt = createdAt
	}
}

// SetUpdateMetadata sets expense update metadata.
func SetUpdateMetadata(updatedBy string, updatedAt time.Time) func(*Expense) {
	return func(e *Expense) {
		e.updatedBy = &updatedBy
		e.updatedAt = &updatedAt
	}
}

// NewExpense creates a new expense domain object.
func NewExpense(
	id string,
	category Category,
	price float64,
	currency string,
	quantity float64,
	comment *string,
	trip *string,
	date time.Time,
	opts ...func(*Expense),
) (*Expense, error) {
	if price == 0 {
		return nil, errors.New("price could not be empty")
	}

	// TODO: add more business checks.

	decPrice := decimal.NewFromFloat(price)
	decQuantity := decimal.NewFromFloat(quantity)

	expense := &Expense{
		id:       id,
		category: category,
		price:    decPrice,
		currency: currency,
		quantity: decQuantity,
		comment:  comment,
		trip:     trip,
		date:     date,
	}

	for _, opt := range opts {
		opt(expense)
	}

	return expense, nil
}

// ID returns expense id.
func (e Expense) ID() string {
	return e.id
}

// Category returns expense category.
func (e Expense) Category() Category {
	return e.category
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

// Trip returns expense trip.
func (e Expense) Trip() *string {
	return e.trip
}

// Date returns expense date.
func (e Expense) Date() time.Time {
	return e.date
}

// CreatedAt returns expense creation date.
func (e Expense) CreatedAt() time.Time {
	return e.createdAt
}

// CreatedBy returns expense creator.
func (e Expense) CreatedBy() string {
	return e.createdBy
}

// UpdatedAt returns expense update date.
func (e Expense) UpdatedAt() *time.Time {
	return e.updatedAt
}

// UpdatedBy returns expense updater.
func (e Expense) UpdatedBy() *string {
	return e.updatedBy
}
