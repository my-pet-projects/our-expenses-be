package command

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// MoveCategoryCommand defines a category move command.
type MoveCategoryCommand struct {
	CategoryID    string
	DestinationID string
}

// MoveCategoryHandler defines a handler to move category.
type MoveCategoryHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// MoveCategoryHandlerInterface defines a contract to handle command.
type MoveCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd MoveCategoryCommand) (*domain.UpdateResult, error)
}

// NewMoveCategoryHandler returns command handler.
func NewMoveCategoryHandler(
	repo adapters.CategoryRepoInterface,
	logger logger.LogInterface,
) MoveCategoryHandler {
	return MoveCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles move category command.
func (h MoveCategoryHandler) Handle(ctx context.Context, cmd MoveCategoryCommand) (*domain.UpdateResult, error) {
	ctx, span := tracer.NewSpan(ctx, "execute move category command")
	defer span.End()

	targetCat, targetCatErr := h.repo.GetOne(ctx, cmd.CategoryID)
	if targetCatErr != nil {
		return nil, errors.Wrap(targetCatErr, "get target category to move")
	}

	if targetCat == nil {
		return nil, nil
	}

	pathFilter := domain.CategoryFilter{
		CategoryID:   cmd.CategoryID,
		FindChildren: true,
	}

	categoryUsages, categoryUsagesErr := h.repo.GetAll(ctx, pathFilter)
	if categoryUsagesErr != nil {
		return nil, errors.Wrap(categoryUsagesErr, "get category children")
	}

	var destinationCat *domain.Category
	if cmd.DestinationID != "root" {
		var destinationCatErr error
		destinationCat, destinationCatErr = h.repo.GetOne(ctx, cmd.DestinationID)
		if destinationCatErr != nil {
			return nil, errors.Wrap(destinationCatErr, "get destination category to move")
		}

		if destinationCat == nil {
			return nil, nil
		}
	}

	oldPath := targetCat.Path()
	targetCat.AssignNewParent(destinationCat)
	newPath := targetCat.Path()

	categoriesToUpdate := []domain.Category{*targetCat}

	for _, category := range categoryUsages {
		category.ReplaceAncestor(oldPath, newPath)
		categoriesToUpdate = append(categoriesToUpdate, category)
	}

	moveResult := &domain.UpdateResult{}
	for _, category := range categoriesToUpdate {
		updateResult, updateErr := h.repo.Update(ctx, category)
		if updateErr != nil {
			return nil, errors.Wrap(updateErr, "target category")
		}
		moveResult.UpdateCount += updateResult.UpdateCount
	}

	return moveResult, nil
}
