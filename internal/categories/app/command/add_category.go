package command

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var addCategoryTracer trace.Tracer

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
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// AddCategoryHandlerInterface defines a contract to handle command.
type AddCategoryHandlerInterface interface {
	Handle(ctx context.Context, cmd AddCategoryCommand) (*string, error)
}

// NewAddCategoryHandler returns command handler.
func NewAddCategoryHandler(
	repo repository.CategoryRepoInterface,
	logger logger.LogInterface,
) AddCategoryHandler {
	addCategoryTracer = otel.Tracer("app.command.add_category")
	return AddCategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// Handle handles add category command.
func (h AddCategoryHandler) Handle(ctx context.Context, cmd AddCategoryCommand) (*string, error) {
	ctx, span := addCategoryTracer.Start(ctx, "execute add category command")
	defer span.End()

	if len(cmd.ID) == 0 {
		cmd.ID = primitive.NewObjectID().Hex()
	}
	path := fmt.Sprintf("%s|%s", cmd.Path, cmd.ID)
	category, categoryErr := domain.NewCategory(cmd.ID, cmd.Name, cmd.ParentID, path, cmd.Icon, cmd.Level, time.Now(), nil)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "prepare category failed")
	}

	return h.repo.Insert(ctx, *category)
}
