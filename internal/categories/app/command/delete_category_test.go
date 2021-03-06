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

func TestDeleteCategoryHandler_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewDeleteCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestDeleteCategoryHandler_FailedToGetCategory_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	cmd := command.DeleteCategoryCommand{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewDeleteCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestDeleteCategoryHandler_NoCategoryFound_ReturnsEmptyResult(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	cmd := command.DeleteCategoryCommand{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(nil, nil)

	// SUT
	sut := command.NewDeleteCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestDeleteCategoryHandler_FailedDeleteCategory_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	icon := "icon"
	cmd := command.DeleteCategoryCommand{
		CategoryID: categoryID,
	}
	category, _ := domain.NewCategory(categoryID, "name", nil, "path", &icon, 1)

	matchIDFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return filter.Path == category.Path() && filter.FindChildren == true
	}
	repo.On("DeleteAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewDeleteCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestDeleteCategoryHandler_DeletesCategory_ReturnsResult(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	icon := "icon"
	cmd := command.DeleteCategoryCommand{
		CategoryID: categoryID,
	}
	category, _ := domain.NewCategory(categoryID, "name", nil, "path", &icon, 1)
	deleteResult := &domain.DeleteResult{DeleteCount: 10}

	matchIDFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return filter.Path == category.Path() && filter.FindChildren == true
	}
	repo.On("DeleteAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(deleteResult, nil)

	// SUT
	sut := command.NewDeleteCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, result, "Result should not be nil.")
	assert.Equal(t, deleteResult, result, "Result should match.")
	assert.Nil(t, err, "Error result should be nil.")
}
