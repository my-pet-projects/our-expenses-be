package command

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// AddExpenseCommand defines an expense command.
type AddExpenseCommand struct {
	Category domain.Category
	Price    float64
	Quantity float64
	Currency string
	Date     time.Time
	Comment  *string
	Trip     *string
}

// AddExpenseHandler defines a handler to add expense.
type AddExpenseHandler struct {
	repo   adapters.ExpenseRepoInterface
	logger logger.LogInterface
}

// AddExpenseHandlerInterface defines a contract to handle command.
type AddExpenseHandlerInterface interface {
	Handle(ctx context.Context, cmd AddExpenseCommand) (*string, error)
}

// NewAddExpenseHandler returns command handler.
func NewAddExpenseHandler(
	repo adapters.ExpenseRepoInterface,
	logger logger.LogInterface,
) AddExpenseHandler {
	return AddExpenseHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles add expense command.
func (h AddExpenseHandler) Handle(ctx context.Context, cmd AddExpenseCommand) (*string, error) {
	ctx, span := tracer.NewSpan(ctx, "execute add expense command")
	defer span.End()

	expense, expenseErr := domain.NewExpense("", cmd.Category, cmd.Price, cmd.Currency, cmd.Quantity,
		cmd.Comment, cmd.Trip, cmd.Date)
	if expenseErr != nil {
		return nil, errors.Wrap(expenseErr, "prepare expense failed")
	}

	return h.repo.Insert(ctx, *expense)
}
