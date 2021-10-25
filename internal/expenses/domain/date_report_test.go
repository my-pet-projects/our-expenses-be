package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// nolint:funlen
func TestDateReport_CalculateTotal_ReturnsTotalAmountIncludingChildren(t *testing.T) {
	t.Parallel()
	// Arrange
	cat, _ := NewCategory("id1", nil, "category 1", nil, 1, "path")
	date1 := time.Date(2021, time.July, 10, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2021, time.July, 20, 0, 0, 0, 0, time.UTC)
	expense1, _ := NewExpense(uuid.NewString(), *cat, 10, "EUR", 1, nil, nil, date1)
	expense2, _ := NewExpense(uuid.NewString(), *cat, 20, "EUR", 2, nil, nil, date1)
	expense3, _ := NewExpense(uuid.NewString(), *cat, 30, "EUR", 3, nil, nil, date1)
	expense4, _ := NewExpense(uuid.NewString(), *cat, 40, "EUR", 4, nil, nil, date1)
	expense5, _ := NewExpense(uuid.NewString(), *cat, 100, "EUR", 5, nil, nil, date1)
	expense6, _ := NewExpense(uuid.NewString(), *cat, 200, "EUR", 5, nil, nil, date2)
	expense7, _ := NewExpense(uuid.NewString(), *cat, 300, "EUR", 5, nil, nil, date2)
	total := expense1.CalculateTotal().
		Add(expense2.CalculateTotal()).
		Add(expense3.CalculateTotal()).
		Add(expense4.CalculateTotal()).
		Add(expense5.CalculateTotal()).
		Add(expense6.CalculateTotal()).
		Add(expense7.CalculateTotal())

	// SUT
	sut := ReportByDate{
		CategoryByDate: []*DateExpenses{
			{
				SubCategories: []*CategoryExpenses{
					{
						SubCategories: []*CategoryExpenses{
							{
								Expenses: &[]Expense{*expense1, *expense2},
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
			}, {
				SubCategories: []*CategoryExpenses{
					{
						SubCategories: []*CategoryExpenses{
							{
								Expenses: &[]Expense{*expense6, *expense7},
							},
						},
					},
				},
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
