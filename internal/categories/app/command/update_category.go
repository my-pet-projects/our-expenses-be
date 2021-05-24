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

// UpdateCategoryCommandArgs defines category command arguments.
type UpdateCategoryCommandArgs struct {
	ID       string
	ParentID *string
	Name     string
	Path     string
	Level    int
}

// UpdateCategoryHandler defines handler to add category.
type UpdateCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// UpdateCategoryHandlerInterface defines a contract to handle command.
type UpdateCategoryHandlerInterface interface {
	Handle(ctx context.Context, category UpdateCategoryCommandArgs) error
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
func (h UpdateCategoryHandler) Handle(ctx context.Context, args UpdateCategoryCommandArgs) error {
	ctx, span := updateCategoryTracer.Start(ctx, "execute update category command")
	defer span.End()

	category, categoryErr := domain.NewCategory(args.ID, args.Name, args.ParentID, args.Path, args.Level, time.Now(), nil)
	if categoryErr != nil {
		return errors.Wrap(categoryErr, "prepare category failed")
	}

	_, updateCmdErr := h.repo.Update(ctx, category)
	if updateCmdErr != nil {
		return errors.Wrap(updateCmdErr, "update category command failed")
	}

	return nil
}
