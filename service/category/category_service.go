package category

import (
	"context"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/entity"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/infrastructure/db/repository"
	"go.opentelemetry.io/otel"

	"github.com/pkg/errors"
)

// CategoryService represents a struct to
type CategoryService struct {
	repo repository.CategoryRepoInterface
}

// CategoryServiceInterface defines a contract to persist categories in the database.
type CategoryServiceInterface interface {
	GetAll(ctx context.Context, filter entity.CategoryFilter) ([]entity.Category, error)
	GetOne(ctx context.Context, id string) (*entity.Category, error)
	Create(ctx context.Context, name string, parentID string, path string, level int) (*entity.ID, error)
	Update(ctx context.Context, category *entity.Category) (string, error)
	DeleteAll(ctx context.Context, filter entity.CategoryFilter) (int, error)
	DeleteOne(ctx context.Context, id string) (int, error)
}

// ProvideCategoryService create new service.
func ProvideCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetAll fetches categories.
func (s *CategoryService) GetAll(ctx context.Context, filter entity.CategoryFilter) ([]entity.Category, error) {
	tracer := otel.Tracer("category.service.GetAll")
	ctx, span := tracer.Start(ctx, "service get all")
	defer span.End()

	categories, categoriesErr := s.repo.GetAll(ctx, filter)
	if categoriesErr != nil {
		return nil, entity.NewAppDbError(errors.Wrap(categoriesErr, "fetch categories"))
	}
	return categories, nil
}

// GetOne fetches categories.
func (s *CategoryService) GetOne(ctx context.Context, id string) (*entity.Category, error) {
	return s.repo.GetOne(ctx, id)
}

// Create creates a category.
func (s *CategoryService) Create(ctx context.Context, name string, parentID string, path string, level int) (*entity.ID, error) {
	c, err := entity.NewCategory(name, parentID, path, level)
	if err != nil {
		return nil, entity.NewAppError(errors.Wrap(err, "create a category"))

	}
	return s.repo.Insert(ctx, c)
}

// GetOne fetches categories.
func (s *CategoryService) Update(ctx context.Context, category *entity.Category) (string, error) {
	return s.repo.Update(ctx, category)
}

// GetOne fetches categories.
func (s *CategoryService) DeleteAll(ctx context.Context, filter entity.CategoryFilter) (int, error) {
	return s.repo.DeleteAll(ctx, filter)
}

// GetOne fetches categories.
func (s *CategoryService) DeleteOne(ctx context.Context, id string) (int, error) {
	return s.repo.DeleteOne(ctx, id)
}
