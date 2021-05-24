package command_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestUpdateCategoryHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewUpdateCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestUpdateCategoryHandler_CategoryError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	args := command.UpdateCategoryCommandArgs{}

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	err := sut.Handle(ctx, args)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestUpdateCategoryHandler_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID := "parentId"
	args := command.UpdateCategoryCommandArgs{
		ID:       categoryID,
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.ID() == args.ID && cat.Name() == args.Name && cat.Path() == args.Path &&
			cat.Level() == args.Level && cat.ParentID() == args.ParentID
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return("", errors.New("error"))

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	err := sut.Handle(ctx, args)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestUpdateCategoryHandler_RepoSuccess_ReturnsResult(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID := "parentId"
	args := command.UpdateCategoryCommandArgs{
		ID:       categoryID,
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.ID() == args.ID && cat.Name() == args.Name && cat.Path() == args.Path &&
			cat.Level() == args.Level && cat.ParentID() == args.ParentID
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return(categoryID, nil)

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	err := sut.Handle(ctx, args)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, err, "Error result should be nil.")
}
