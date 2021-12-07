package command

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// AddCategoryCommand defines a category command.
type AddCategoryCommand struct {
	ID       string
	ParentID *string
	Name     string
	Path     string
	Icon     *string
	Level    int
}

// AddCategoryHandler defines a handler to add category.
type AddCategoryHandler struct {
	repo   adapters.CategoryRepoInterface
	logger logger.LogInterface
}

// AddCategoryHandlerInterface defines a contract to handle command.
type AddCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd AddCategoryCommand) (*string, error)
}

// NewAddCategoryHandler returns command handler.
func NewAddCategoryHandler(
	repo adapters.CategoryRepoInterface,
	logger logger.LogInterface,
) AddCategoryHandler {
	return AddCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles add category command.
func (h AddCategoryHandler) Handle(ctx context.Context, cmd AddCategoryCommand) (*string, error) {
	ctx, span := tracer.NewSpan(ctx, "execute add category command")
	defer span.End()

	if len(cmd.ID) == 0 {
		cmd.ID = primitive.NewObjectID().Hex()
	}
	path := fmt.Sprintf("%s|%s", cmd.Path, cmd.ID)
	category, categoryErr := domain.NewCategory(cmd.ID, cmd.Name, cmd.ParentID, path, cmd.Icon, cmd.Level)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "prepare category failed")
	}

	insRes, insResErr := h.repo.Insert(ctx, *category)
	if insResErr != nil {
		return nil, errors.Wrap(insResErr, "insert category")
	}
	return insRes, nil
}
