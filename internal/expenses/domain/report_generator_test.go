package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

func TestGenerateReport(t *testing.T) {
	// Arrange
	categoryID1 := uuid.NewString()
	categoryID2 := uuid.NewString()
	categoryID3 := uuid.NewString()
	categoryID4 := uuid.NewString()
	date1 := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2021, time.July, 15, 0, 0, 0, 0, time.UTC)
	expense1, _ := domain.NewExpense(uuid.NewString(), categoryID1, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense2, _ := domain.NewExpense(uuid.NewString(), categoryID1, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense3, _ := domain.NewExpense(uuid.NewString(), categoryID3, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense4, _ := domain.NewExpense(uuid.NewString(), categoryID3, 0, "EUR", 0, nil, date2, time.Now(), nil)
	expense5, _ := domain.NewExpense(uuid.NewString(), categoryID3, 0, "EUR", 0, nil, date3, time.Now(), nil)
	expense6, _ := domain.NewExpense(uuid.NewString(), categoryID4, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense7, _ := domain.NewExpense(uuid.NewString(), categoryID1, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense8, _ := domain.NewExpense(uuid.NewString(), categoryID2, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense9, _ := domain.NewExpense(uuid.NewString(), categoryID1, 0, "EUR", 0, nil, date1, time.Now(), nil)
	expense10, _ := domain.NewExpense(uuid.NewString(), categoryID1, 0, "EUR", 0, nil, date1, time.Now(), nil)

	expenses := []domain.Expense{*expense1, *expense2, *expense3, *expense4, *expense5,
		*expense6, *expense7, *expense8, *expense9, *expense10}

	// SUT
	sut := domain.NewReportGenerator(expenses)

	// Act
	result := sut.GenerateReport()

	// Assert
	assert.NotNil(t, result)

}
