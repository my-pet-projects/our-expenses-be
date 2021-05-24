package command

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var deleteCategoryTracer trace.Tracer

// DeleteCategoryCommand defines a category delete command.
type DeleteCategoryCommand struct {
	CategoryID string
}

// DeleteCategoryHandler defines a handler to delete category.
type DeleteCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// DeleteCategoryHandlerInterface defines a contract to handle command.
type DeleteCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd DeleteCategoryCommand) (*domain.DeleteResult, error)
}

// NewDeleteCategoryHandler returns command handler.
func NewDeleteCategoryHandler(
	repo repository.CategoryRepoInterface,
	logger logger.LogInterface,
) DeleteCategoryHandler {
	deleteCategoryTracer = otel.Tracer("app.command.delete_category")
	return DeleteCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles delete category command.
func (h DeleteCategoryHandler) Handle(ctx context.Context, cmd DeleteCategoryCommand) (*domain.DeleteResult, error) {
	ctx, span := deleteCategoryTracer.Start(ctx, "execute delete category command")
	defer span.End()

	category, categoryErr := h.repo.GetOne(ctx, cmd.CategoryID)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "get category for deletion")
	}

	if category == nil {
		return nil, nil
	}

	deleteFilter := domain.CategoryFilter{
		Path:         category.Path(),
		FindChildren: true,
	}

	deleteCmdResult, deleteCmdErr := h.repo.DeleteAll(ctx, deleteFilter)
	if deleteCmdErr != nil {
		return nil, errors.Wrap(deleteCmdErr, "delete category command")
	}

	return deleteCmdResult, nil
}
