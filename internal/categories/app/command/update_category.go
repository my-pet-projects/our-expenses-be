package command

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// UpdateCategoryCommand defines a category update command.
type UpdateCategoryCommand struct {
	ID       string
	ParentID *string
	Name     string
	Path     string
	Icon     *string
	Level    int
}

// UpdateCategoryHandler defines a handler to update category.
type UpdateCategoryHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// UpdateCategoryHandlerInterface defines a contract to handle command.
type UpdateCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd UpdateCategoryCommand) (*domain.UpdateResult, error)
}

// NewUpdateCategoryHandler returns command handler.
func NewUpdateCategoryHandler(
	repo adapters.CategoryRepoInterface,
	logger logger.LogInterface,
) UpdateCategoryHandler {
	return UpdateCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles update category command.
func (h UpdateCategoryHandler) Handle(
	ctx context.Context,
	cmd UpdateCategoryCommand,
) (*domain.UpdateResult, error) {
	ctx, span := tracer.NewSpan(ctx, "execute update category command")
	defer span.End()

	// TODO: get a category from the database and create a new category based on that object.

	category, categoryErr := domain.NewCategory(cmd.ID, cmd.Name, cmd.ParentID, cmd.Path,
		cmd.Icon, cmd.Level)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "prepare category failed")
	}

	updateCmdResult, updateCmdErr := h.repo.Update(ctx, *category)
	if updateCmdErr != nil {
		return nil, errors.Wrap(updateCmdErr, "update category command failed")
	}

	return updateCmdResult, nil
}
