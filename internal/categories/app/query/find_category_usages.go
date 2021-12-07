package query

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// FindCategoryUsagesQuery defines a category usages query.
type FindCategoryUsagesQuery struct {
	CategoryID string
}

// FindCategoryUsagesHandler defines handler to fetch category usages.
type FindCategoryUsagesHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// FindCategoryUsagesHandlerInterface defines a contract to handle query.
type FindCategoryUsagesHandlerInterface interface {
	Handle(ctx context.Context, query FindCategoryUsagesQuery) ([]domain.Category, error)
}

// NewFindCategoryUsagesHandler returns query handler.
func NewFindCategoryUsagesHandler(
	repo adapters.CategoryRepoInterface,
	logger logger.LogInterface,
) FindCategoryUsagesHandler {
	return FindCategoryUsagesHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles find category usages query.
func (h FindCategoryUsagesHandler) Handle(
	ctx context.Context,
	query FindCategoryUsagesQuery,
) ([]domain.Category, error) {
	ctx, span := tracer.NewSpan(ctx, "execute find category usages")
	span.SetAttributes(attribute.String("id", query.CategoryID))
	defer span.End()

	filter := domain.CategoryFilter{
		CategoryID:   query.CategoryID,
		FindChildren: true,
	}

	queryResult, queryErr := h.repo.GetAll(ctx, filter)
	if queryErr != nil {
		return nil, errors.Wrap(queryErr, "fetch category usages")
	}

	return queryResult, nil
}
