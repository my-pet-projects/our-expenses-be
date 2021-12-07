package query

import (
	"context"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
	"github.com/pkg/errors"
)

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
	ctx, span := tracer.NewSpan(ctx, "execute find categories query")
	defer span.End()

	filter := domain.CategoryFilter{
		ParentID:     query.ParentID,
		FindChildren: query.FindAllChildren,
		FindAll:      query.FindAll,
	}
	res, resErr := h.repo.GetAll(ctx, filter)
	if resErr != nil {
		return nil, errors.Wrap(resErr, "get categories")
	}
	return res, nil
}
