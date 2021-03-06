package query

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// FindCategoryQuery defines a category query.
type FindCategoryQuery struct {
	CategoryID string
}

// FindCategoryHandler defines handler to fetch category.
type FindCategoryHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// FindCategoryHandlerInterface defines a contract to handle query.
type FindCategoryHandlerInterface interface {
	Handle(ctx context.Context, query FindCategoryQuery) (*domain.Category, error)
}

// NewFindCategoryHandler returns query handler.
func NewFindCategoryHandler(repo adapters.CategoryRepoInterface, logger logger.LogInterface) FindCategoryHandler {
	return FindCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find category query.
func (h FindCategoryHandler) Handle(ctx context.Context, query FindCategoryQuery) (*domain.Category, error) {
	ctx, span := tracer.NewSpan(ctx, "execute find category query")
	defer span.End()

	category, categoryErr := h.repo.GetOne(ctx, query.CategoryID)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "get category")
	}

	if category != nil && len(category.ParentIDs()) != 0 {
		parentCategoriesFilter := domain.CategoryFilter{
			CategoryIDs: category.ParentIDs(),
		}

		parentCategories, parentCategoriesError := h.repo.GetAll(ctx, parentCategoriesFilter)
		if parentCategoriesError != nil {
			return nil, errors.Wrapf(parentCategoriesError, "fetch parents for %s category", query.CategoryID)
		}

		category.SetParents(parentCategories)
	}

	return category, nil
}
