package command

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var addExpenseTracer trace.Tracer

// AddExpenseCommand defines an expense command.
type AddExpenseCommand struct {
	CategoryID string
	Price      float64
	Quantity   float64
	Currency   string
	Date       time.Time
	Comment    *string
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
	addExpenseTracer = otel.Tracer("app.command.add_expense")
	return AddExpenseHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles add expense command.
func (h AddExpenseHandler) Handle(ctx context.Context, cmd AddExpenseCommand) (*string, error) {
	ctx, span := addExpenseTracer.Start(ctx, "execute add expense command")
	defer span.End()

	expense, expenseErr := domain.NewExpense("", cmd.CategoryID, cmd.Price, cmd.Currency, cmd.Quantity,
		cmd.Comment, cmd.Date, time.Now(), nil)
	if expenseErr != nil {
		return nil, errors.Wrap(expenseErr, "prepare expense failed")
	}

	return h.repo.Insert(ctx, *expense)
}
