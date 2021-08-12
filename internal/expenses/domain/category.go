package domain

import (
	"strings"

	"github.com/pkg/errors"
)

// Category represents a domain object.
type Category struct {
	id       string
	name     string
	icon     *string
	level    int
	parents  *[]Category
	expenses *[]Expense
}

// NewCategory creates a new category domain object.
func NewCategory(id string, name string, icon *string, level int) (*Category, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	if icon != nil {
		trimmedIcon := strings.TrimSpace(*icon)
		icon = &trimmedIcon
	}

	return &Category{
		id:    id,
		name:  strings.TrimSpace(name),
		icon:  icon,
		level: level,
	}, nil
}

// ID returns category id.
func (c Category) ID() string {
	return c.id
}

// Name returns category name.
func (c Category) Name() string {
	return c.name
}

// Icon returns category icon.
func (c Category) Icon() *string {
	return c.icon
}

// Level returns category level.
func (c Category) Level() int {
	return c.level
}

// Parents returns category parents.
func (c Category) Parents() *[]Category {
	return c.parents
}

// SetParents sets category parents.
func (c *Category) SetParents(parents *[]Category) {
	c.parents = parents
}

// Expenses returns category expenses.
func (c Category) Expenses() *[]Expense {
	return c.expenses
}

// SetExpenses sets category expenses.
func (c *Category) SetExpenses(expenses *[]Expense) {
	c.expenses = expenses
}
