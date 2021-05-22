package query

import (
	"context"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/entity"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// FindCategoriesHandler defines handler to fetch categories.
type FindCategoriesHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

var locationsTracer trace.Tracer

// FindCategoriesHandlerInterface defines a contract to handle query.
type FindCategoriesHandlerInterface interface {
	Handle(ctx context.Context) ([]domain.Category, error)
}

// NewFindCategoriesHandler returns query handler.
func NewFindCategoriesHandler(
	repo repository.CategoryRepoInterface,
	logger logger.LogInterface,
) FindCategoriesHandler {
	locationsTracer = otel.Tracer("app.query.find_categories")
	return FindCategoriesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find locations query.
func (h FindCategoriesHandler) Handle(ctx context.Context) ([]domain.Category, error) {
	ctx, span := locationsTracer.Start(ctx, "execute find categories query")
	defer span.End()

	filter := entity.CategoryFilter{}

	return h.repo.GetAll(ctx, filter)
}
