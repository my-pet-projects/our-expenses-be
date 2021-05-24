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

	cmd := command.UpdateCategoryCommand{}

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestUpdateCategoryHandler_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID := "parentId"
	cmd := command.UpdateCategoryCommand{
		ID:       categoryID,
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.ID() == cmd.ID && cat.Name() == cmd.Name && cat.Path() == cmd.Path &&
			cat.Level() == cmd.Level && cat.ParentID() == cmd.ParentID
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestUpdateCategoryHandler_RepoSuccess_ReturnsResult(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID := "parentId"
	cmd := command.UpdateCategoryCommand{
		ID:       categoryID,
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}
	updateResult := &domain.UpdateResult{UpdateCount: 10}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.ID() == cmd.ID && cat.Name() == cmd.Name && cat.Path() == cmd.Path &&
			cat.Level() == cmd.Level && cat.ParentID() == cmd.ParentID
	}
	repo.On("Update", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return(updateResult, nil)

	// SUT
	sut := command.NewUpdateCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, result, "Result should be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}
