package query

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var findCategoryTracer trace.Tracer

// FindCategoryQuery defines a category query.
type FindCategoryQuery struct {
	CategoryID string
}

// FindCategoryHandler defines handler to fetch category.
type FindCategoryHandler struct {
	repo   adapters.ExpenseCategoryRepoInterface
	logger logger.LogInterface
}

// FindExpenseCategoryHandlerInterface defines a contract to handle query.
type FindExpenseCategoryHandlerInterface interface {
	Handle(ctx context.Context, query FindCategoryQuery) (*domain.Category, error)
}

// NewFindCategoryHandler returns query handler.
func NewFindCategoryHandler(repo adapters.ExpenseCategoryRepoInterface, logger logger.LogInterface) FindCategoryHandler {
	findCategoryTracer = otel.Tracer("app.query.find_category")
	return FindCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find category query.
func (h FindCategoryHandler) Handle(ctx context.Context, query FindCategoryQuery) (*domain.Category, error) {
	ctx, span := findCategoryTracer.Start(ctx, "execute find category query")
	defer span.End()

	category, categoryErr := h.repo.GetOne(ctx, query.CategoryID)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "get category")
	}

	return category, nil
}
