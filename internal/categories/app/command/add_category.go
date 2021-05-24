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

var addCategoryTracer trace.Tracer

// NewCategoryCommandArgs defines category command object.
type NewCategoryCommandArgs struct {
	ParentID *string
	Name     string
	Path     string
	Level    int
}

// AddCategoryHandler defines handler to add category.
type AddCategoryHandler struct {
	repo   repository.CategoryRepoInterface
	logger logger.LogInterface
}

// AddCategoryHandlerInterface defines a contract to handle command.
type AddCategoryHandlerInterface interface {
	Handle(ctx context.Context, args NewCategoryCommandArgs) (*string, error)
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
func (h AddCategoryHandler) Handle(ctx context.Context, args NewCategoryCommandArgs) (*string, error) {
	ctx, span := addCategoryTracer.Start(ctx, "execute add category command")
	defer span.End()

	category, categoryErr := domain.NewCategory("", args.Name, args.ParentID, args.Path, args.Level, time.Now(), nil)
	if categoryErr != nil {
		return nil, errors.Wrap(categoryErr, "prepare category failed")
	}

	return h.repo.Insert(ctx, category)
}
