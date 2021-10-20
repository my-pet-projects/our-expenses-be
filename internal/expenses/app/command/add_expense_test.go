package command_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestAddExpenseHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewAddExpenseHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddExpenseHandler_ExpenseError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.AddExpenseCommand{}

	// SUT
	sut := command.NewAddExpenseHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddExpenseHandler_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	comment := "comment"
	parentID := "parentID"
	category, _ := domain.NewCategory("categoryID", &parentID, "category", nil, 1, "path")
	cmd := command.AddExpenseCommand{
		Category: *category,
		Price:    12.55,
		Currency: "EUR",
		Quantity: 2,
		Comment:  &comment,
		Date:     time.Now(),
	}

	matchExpenseFn := func(cat domain.Expense) bool {
		return cat.ID() == "" && reflect.DeepEqual(cat.Category(), cmd.Category) &&
			cat.Price() == cmd.Price && cat.Currency() == cmd.Currency && cat.Quantity() == cmd.Quantity &&
			cat.Comment() == cmd.Comment && cat.Date() == cmd.Date
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchExpenseFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewAddExpenseHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddExpenseHandler_RepoSuccess_ReturnsNewId(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	expenseID := "expenseId"
	comment := "comment"
	parentID := "parentID"
	category, _ := domain.NewCategory("categoryID", &parentID, "category", nil, 1, "path")
	cmd := command.AddExpenseCommand{
		Category: *category,
		Price:    12.55,
		Currency: "EUR",
		Quantity: 2,
		Comment:  &comment,
		Date:     time.Now(),
	}

	matchExpenseFn := func(cat domain.Expense) bool {
		return cat.ID() == "" && reflect.DeepEqual(cat.Category(), cmd.Category) &&
			cat.Price() == cmd.Price && cat.Currency() == cmd.Currency && cat.Quantity() == cmd.Quantity &&
			cat.Comment() == cmd.Comment && cat.Date() == cmd.Date
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchExpenseFn)).Return(&expenseID, nil)

	// SUT
	sut := command.NewAddExpenseHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, &expenseID, query, "Should return expense id.")
	assert.Nil(t, err, "Error result should be nil.")
}
