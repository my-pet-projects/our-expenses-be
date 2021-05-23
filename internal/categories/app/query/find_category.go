package query

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var findCategoryTracer trace.Tracer

// FindCategoryHandler defines handler to fetch category.
type FindCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// FindCategoryHandlerInterface defines a contract to handle query.
type FindCategoryHandlerInterface interface {
	Handle(ctx context.Context, id string) (*domain.Category, error)
}

// NewFindCategoryHandler returns query handler.
func NewFindCategoryHandler(repo repository.CategoryRepoInterface, logger logger.LogInterface) FindCategoryHandler {
	findCategoryTracer = otel.Tracer("app.query.find_category")
	return FindCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find category query.
func (h FindCategoryHandler) Handle(ctx context.Context, id string) (*domain.Category, error) {
	ctx, span := findCategoryTracer.Start(ctx, "execute find category query")
	defer span.End()

	category, categoryErr := h.repo.GetOne(ctx, id)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "get category")
	}

	if category != nil && len(category.ParentIDs()) != 0 {
		parentCategoriesFilter := domain.CategoryFilter{
			CategoryIDs: category.ParentIDs(),
		}

		parentCategories, parentCategoriesError := h.repo.GetAll(ctx, parentCategoriesFilter)
		if parentCategoriesError != nil {
			return nil, errors.Wrapf(parentCategoriesError, "fetch parents for %s category", id)
		}

		category.SetParents(parentCategories)
	}

	return category, nil
}
