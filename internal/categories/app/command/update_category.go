package command

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var updateCategoryTracer trace.Tracer

// UpdateCategoryCommand defines a category update command.
type UpdateCategoryCommand struct {
	ID       string
	ParentID *string
	Name     string
	Path     string
	Level    int
}

// UpdateCategoryHandler defines a handler to update category.
type UpdateCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// UpdateCategoryHandlerInterface defines a contract to handle command.
type UpdateCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd UpdateCategoryCommand) (*domain.UpdateResult, error)
}

// NewUpdateCategoryHandler returns command handler.
func NewUpdateCategoryHandler(
	repo repository.CategoryRepoInterface,
	logger logger.LogInterface,
) UpdateCategoryHandler {
	updateCategoryTracer = otel.Tracer("app.command.update_category")
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
	ctx, span := updateCategoryTracer.Start(ctx, "execute update category command")
	defer span.End()

	// TODO: get a category from the database and create a new category based on that object.

	category, categoryErr := domain.NewCategory(cmd.ID, cmd.Name, cmd.ParentID, cmd.Path, cmd.Level, time.Now(), nil)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "prepare category failed")
	}

	updateCmdResult, updateCmdErr := h.repo.Update(ctx, *category)
	if updateCmdErr != nil {
		return nil, errors.Wrap(updateCmdErr, "update category command failed")
	}

	return updateCmdResult, nil
}
