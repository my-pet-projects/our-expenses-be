package query

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var findCategoriesTracer trace.Tracer

// FindCategoriesHandler defines handler to fetch categories.
type FindCategoriesHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// FindCategoriesHandlerInterface defines a contract to handle query.
type FindCategoriesHandlerInterface interface {
	Handle(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, error)
}

// NewFindCategoriesHandler returns query handler.
func NewFindCategoriesHandler(repo repository.CategoryRepoInterface, logger logger.LogInterface) FindCategoriesHandler {
	findCategoriesTracer = otel.Tracer("app.query.find_categories")
	return FindCategoriesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find categories query.
func (h FindCategoriesHandler) Handle(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, error) {
	ctx, span := findCategoriesTracer.Start(ctx, "execute find categories query")
	defer span.End()

	return h.repo.GetAll(ctx, filter)
}
