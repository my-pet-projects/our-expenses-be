package domain_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

// nolint:funlen,gocognit
func TestGenerateByDateReport(t *testing.T) {
	t.Parallel()
	// Arrange
	// category 1
	id1 := uuid.MustParse("00000000-0000-0000-0000-000000000001").String()
	category1, _ := domain.NewCategory(id1, nil, "category 1", nil, 1,
		fmt.Sprintf("|%s", id1))
	id11 := uuid.MustParse("00000000-0000-0000-0000-000000000011").String()
	category11, _ := domain.NewCategory(id11, &id1, "category 1.1", nil, 2,
		fmt.Sprintf("|%s|%s", id1, id11))
	category11.SetParents(&[]domain.Category{*category1})

	id111 := uuid.MustParse("00000000-0000-0000-0000-000000000111").String()
	category111, _ := domain.NewCategory(id111, &id11, "category 1.1.1", nil, 3,
		fmt.Sprintf("|%s|%s|%s", id1, id11, id111))
	category111.SetParents(&[]domain.Category{*category1, *category11})

	id112 := uuid.MustParse("00000000-0000-0000-0000-000000000112").String()
	category112, _ := domain.NewCategory(id112, &id11, "category 1.1.2", nil, 3,
		fmt.Sprintf("|%s|%s|%s", id1, id11, id112))
	category112.SetParents(&[]domain.Category{*category1, *category11})

	id12 := uuid.MustParse("00000000-0000-0000-0000-000000000012").String()
	category12, _ := domain.NewCategory(id12, &id1, "category 1.2", nil, 2,
		fmt.Sprintf("|%s|%s", id1, id12))
	category12.SetParents(&[]domain.Category{*category1})

	// category 2
	id2 := uuid.MustParse("00000000-0000-0000-0000-000000000002").String()
	category2, _ := domain.NewCategory(id2, nil, "category 2", nil, 1,
		fmt.Sprintf("|%s", id2))

	id21 := uuid.MustParse("00000000-0000-0000-0000-000000000021").String()
	category21, _ := domain.NewCategory(id21, &id2, "category 2.1", nil, 2,
		fmt.Sprintf("|%s|%s", id2, id21))
	category21.SetParents(&[]domain.Category{*category2})

	// TODO: case when expense is directly in root category

	date1 := time.Date(2021, time.July, 10, 0, 0, 0, 0, time.UTC)

	expense1, _ := domain.NewExpense(uuid.NewString(), *category111, 10, "EUR", 1, nil, nil, date1)
	expense2, _ := domain.NewExpense(uuid.NewString(), *category111, 20, "EUR", 2, nil, nil, date1)
	expense3, _ := domain.NewExpense(uuid.NewString(), *category112, 30, "EUR", 3, nil, nil, date1)
	expense4, _ := domain.NewExpense(uuid.NewString(), *category12, 40, "EUR", 4, nil, nil, date1)
	expense5, _ := domain.NewExpense(uuid.NewString(), *category21, 100, "EUR", 4, nil, nil, date1)

	expenses := []domain.Expense{*expense1, *expense2, *expense3, *expense4, *expense5}

	// SUT
	sut := domain.NewReportGenerator(expenses)

	// Act
	result := sut.GenerateByDateReport()

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.CategoryByDate, 1)

	firstCategoryByDate := result.CategoryByDate[0]

	// nolint:nestif
	if firstCategoryByDate.Date == date1 {
		assert.Equal(t, date1, firstCategoryByDate.Date)
		assert.Len(t, firstCategoryByDate.SubCategories, 2)
		// assert.Equal(t, 0, firstCategoryByDate.Total)

		for _, catLevel1 := range firstCategoryByDate.SubCategories {
			if catLevel1.Category.ID() == category1.ID() {
				assert.Equal(t, category1, &catLevel1.Category)
				// assert.Equal(t, 0, catLevel1.Total)
				assert.Nil(t, catLevel1.Expenses)
				assert.Len(t, catLevel1.SubCategories, 2)

				for _, catLevel2 := range catLevel1.SubCategories {
					if catLevel2.Category.ID() == category11.ID() {
						assert.Equal(t, category11, &catLevel2.Category)
						// assert.Equal(t, 0, catLevel2.Total)
						assert.Nil(t, catLevel2.Expenses)
						assert.Len(t, catLevel2.SubCategories, 2)

						for _, catLevel3 := range catLevel2.SubCategories {
							if catLevel3.Category.ID() == category111.ID() {
								assert.Equal(t, category111, &catLevel3.Category)
								// assert.Equal(t, 0, catLevel3.Total)
								assert.Len(t, *catLevel3.Expenses, 2)
								assert.Contains(t, *catLevel3.Expenses, *expense1)
								assert.Contains(t, *catLevel3.Expenses, *expense2)
								assert.Nil(t, catLevel3.SubCategories)
							}

							if catLevel3.Category.ID() == category112.ID() {
								assert.Equal(t, category112, &catLevel3.Category)
								// assert.Equal(t, 0, catLevel3.Total)
								assert.Len(t, *catLevel3.Expenses, 1)
								assert.Contains(t, *catLevel3.Expenses, *expense3)
								assert.Nil(t, catLevel3.SubCategories)
							}
						}
					}

					if catLevel2.Category.ID() == category12.ID() {
						assert.Equal(t, category12, &catLevel2.Category)
						// assert.Equal(t, 0, catLevel2.Total)
						assert.Len(t, *catLevel2.Expenses, 1)
						assert.Contains(t, *catLevel2.Expenses, *expense4)
						assert.Nil(t, catLevel2.SubCategories)
					}
				}
			}

			if catLevel1.Category.ID() == category2.ID() {
				assert.Equal(t, category2, &catLevel1.Category)
				// assert.Equal(t, 0, catLevel1.Total)
				assert.Nil(t, catLevel1.Expenses)
				assert.Len(t, catLevel1.SubCategories, 1)

				for _, catLevel2 := range catLevel1.SubCategories {
					if catLevel2.Category.ID() == category21.ID() {
						assert.Equal(t, category21, &catLevel2.Category)
						// assert.Equal(t, 0, catLevel2.Total)
						assert.Len(t, *catLevel2.Expenses, 1)
						assert.Contains(t, *catLevel2.Expenses, *expense5)
						assert.Nil(t, catLevel2.SubCategories)
					}
				}
			}
		}
	}
}
