package query

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var findCategoriesTracer trace.Tracer

// FindCategoriesQuery defines a category query.
type FindCategoriesQuery struct {
	ParentID        *string
	FindAllChildren bool
	FindAll         bool
}

// FindCategoriesHandler defines a handler to fetch categories.
type FindCategoriesHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// FindCategoriesHandlerInterface defines a contract to handle query.
type FindCategoriesHandlerInterface interface {
	Handle(ctx context.Context, query FindCategoriesQuery) ([]domain.Category, error)
}

// NewFindCategoriesHandler returns a query handler.
func NewFindCategoriesHandler(
	repo adapters.CategoryRepoInterface,
	logger logger.LogInterface,
) FindCategoriesHandler {
	findCategoriesTracer = otel.Tracer("app.query.find_categories")
	return FindCategoriesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles query to find categories.
func (h FindCategoriesHandler) Handle(
	ctx context.Context,
	query FindCategoriesQuery,
) ([]domain.Category, error) {
	ctx, span := findCategoriesTracer.Start(ctx, "execute find categories query")
	defer span.End()

	filter := domain.CategoryFilter{
		ParentID:     query.ParentID,
		FindChildren: query.FindAllChildren,
		FindAll:      query.FindAll,
	}
	return h.repo.GetAll(ctx, filter)
}
