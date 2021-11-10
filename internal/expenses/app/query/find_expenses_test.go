package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewFindExpensesHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.ReportRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := query.NewFindExpensesHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindExpensesHandle_FilterError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.ReportRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	from := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, time.August, 3, 0, 0, 0, 0, time.UTC)
	interval := "unknown"
	findQuery := query.FindExpensesQuery{
		From:     from,
		To:       to,
		Interval: interval,
	}

	// SUT
	sut := query.NewFindExpensesHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindExpensesHandle_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.ReportRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	from := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, time.August, 3, 0, 0, 0, 0, time.UTC)
	interval := "month"
	findQuery := query.FindExpensesQuery{
		From:     from,
		To:       to,
		Interval: interval,
	}

	matchFilterFn := func(filter domain.ExpenseFilter) bool {
		return filter.To() == to &&
			filter.From() == from &&
			string(filter.Interval()) == interval
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindExpensesHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindExpensesHandle_RepoSuccess_ReturnsExpenses(t *testing.T) {
	// Arrange
	repo := new(mocks.ReportRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	from := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, time.August, 3, 0, 0, 0, 0, time.UTC)
	interval := "month"
	findQuery := query.FindExpensesQuery{
		From:     from,
		To:       to,
		Interval: interval,
	}
	expenses := []domain.Expense{}

	matchFilterFn := func(filter domain.ExpenseFilter) bool {
		return filter.To() == to &&
			filter.From() == from &&
			string(filter.Interval()) == interval
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(expenses, nil)

	// SUT
	sut := query.NewFindExpensesHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, expenses, query, "Should return expenses.")
	assert.Nil(t, err, "Error result should be nil.")
}
