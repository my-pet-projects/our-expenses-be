package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryExpenses_CalculateTotal_ReturnsTotalAmountIncludingChildren(t *testing.T) {
	t.Parallel()
	// Arrange
	cat, _ := NewCategory("id1", nil, "category 1", nil, 1, "path")
	date := time.Date(2021, time.July, 10, 0, 0, 0, 0, time.UTC)
	expense1, _ := NewExpense(uuid.NewString(), *cat, 10, "EUR", 1, nil, nil, date)
	expense2, _ := NewExpense(uuid.NewString(), *cat, 20, "EUR", 2, nil, nil, date)
	expense3, _ := NewExpense(uuid.NewString(), *cat, 30, "EUR", 3, nil, nil, date)
	expense4, _ := NewExpense(uuid.NewString(), *cat, 40, "EUR", 4, nil, nil, date)
	expense5, _ := NewExpense(uuid.NewString(), *cat, 100, "EUR", 5, nil, nil, date)
	total := expense1.CalculateTotal().
		Add(expense2.CalculateTotal()).
		Add(expense3.CalculateTotal()).
		Add(expense4.CalculateTotal()).
		Add(expense5.CalculateTotal())

	// SUT
	sut := CategoryExpenses{
		Expenses: &[]Expense{*expense1, *expense2},
		SubCategories: []*CategoryExpenses{
			{
				SubCategories: []*CategoryExpenses{
					{
						SubCategories: []*CategoryExpenses{
							{
								Expenses: &[]Expense{*expense3, *expense4},
							},
						},
					},
				},
			}, {
				Expenses: &[]Expense{*expense5},
			},
		},
	}

	// Act
	result := sut.CalculateTotal()

	// Assert
	assert.NotNil(t, result)
	assert.True(t, total.Equal(result))
	assert.True(t, total.Equal(sut.Total))
}
