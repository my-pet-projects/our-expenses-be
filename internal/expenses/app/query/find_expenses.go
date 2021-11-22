package query

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// FindExpensesQuery defines an expense query.
type FindExpensesQuery struct {
	DateRange domain.DateRange
	Interval  string
}

// FindExpensesHandler defines a handler to fetch expenses.
type FindExpensesHandler struct {
	repo   adapters.ReportRepoInterface
	logger logger.LogInterface
}

// FindExpensesHandlerInterface defines a contract to handle query.
type FindExpensesHandlerInterface interface {
	Handle(ctx context.Context, query FindExpensesQuery) (*domain.ReportByDate, error)
}

// NewFindExpensesHandler returns a query handler.
func NewFindExpensesHandler(
	repo adapters.ReportRepoInterface,
	logger logger.LogInterface,
) FindExpensesHandler {
	return FindExpensesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles query to find expenses.
func (h FindExpensesHandler) Handle(
	ctx context.Context,
	query FindExpensesQuery,
) (*domain.ReportByDate, error) {
	ctx, span := tracer.NewSpan(ctx, "execute find expenses query")
	defer span.End()

	filter, filterErr := domain.NewExpenseFilter(query.DateRange.From(), query.DateRange.To(), query.Interval)
	if filterErr != nil {
		tracer.AddSpanError(span, filterErr)
		return nil, errors.Wrap(filterErr, "prepare filter")
	}

	expenses, expensesErr := h.repo.GetAll(ctx, *filter)
	if expensesErr != nil {
		tracer.AddSpanError(span, expensesErr)
		return nil, errors.Wrap(expensesErr, "fetch expenses")
	}

	reportGenerator := domain.NewReportGenerator(expenses, *filter)
	report := reportGenerator.GenerateByDateReport()

	return &report, nil
}
