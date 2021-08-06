package query

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var findExpensesTracer trace.Tracer

// FindExpensesQuery defines a expense query.
type FindExpensesQuery struct {
	From time.Time
	To   time.Time
}

// FindExpensesHandler defines a handler to fetch expenses.
type FindExpensesHandler struct {
	repo   adapters.ReportRepoInterface
	logger logger.LogInterface
}

// FindExpensesHandlerInterface defines a contract to handle query.
type FindExpensesHandlerInterface interface {
	Handle(ctx context.Context, query FindExpensesQuery) ([]domain.Expense, error)
}

// NewFindExpensesHandler returns a query handler.
func NewFindExpensesHandler(
	repo adapters.ReportRepoInterface,
	logger logger.LogInterface,
) FindExpensesHandler {
	findExpensesTracer = otel.Tracer("app.query.find_expenses")
	return FindExpensesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles query to find expenses.
func (h FindExpensesHandler) Handle(
	ctx context.Context,
	query FindExpensesQuery,
) ([]domain.Expense, error) {
	ctx, span := findExpensesTracer.Start(ctx, "execute find expenses query")
	defer span.End()

	filter := domain.ExpenseFilter{
		From: query.From,
		To:   query.To,
	}
	return h.repo.GetAll(ctx, filter)
}
