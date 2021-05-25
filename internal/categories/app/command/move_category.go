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

var moveCategoryTracer trace.Tracer

// MoveCategoryCommand defines a category move command.
type MoveCategoryCommand struct {
	CategoryID    string
	DestinationID string
}

// MoveCategoryHandler defines a handler to move category.
type MoveCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// MoveCategoryHandlerInterface defines a contract to handle command.
type MoveCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd MoveCategoryCommand) (*domain.UpdateResult, error)
}

// NewMoveCategoryHandler returns command handler.
func NewMoveCategoryHandler(
	repo repository.CategoryRepoInterface,
	logger logger.LogInterface,
) MoveCategoryHandler {
	moveCategoryTracer = otel.Tracer("app.command.move_category")
	return MoveCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles move category command.
func (h MoveCategoryHandler) Handle(ctx context.Context, cmd MoveCategoryCommand) (*domain.UpdateResult, error) {
	ctx, span := moveCategoryTracer.Start(ctx, "execute move category command")
	defer span.End()

	targetCat, targetCatErr := h.repo.GetOne(ctx, cmd.CategoryID)
	if targetCatErr != nil {
		return nil, errors.Wrap(targetCatErr, "get target category to move")
	}

	if targetCat == nil {
		return nil, nil
	}

	// moveFilter := domain.CategoryFilter{
	// 	Path:         targetCat.Path(),
	// 	FindChildren: true,
	// }

	// moveCmdResult, moveCmdErr := h.repo.MoveAll(ctx, moveFilter)
	// if moveCmdErr != nil {
	// 	return nil, errors.Wrap(moveCmdErr, "move category command")
	// }

	return nil, nil
}
