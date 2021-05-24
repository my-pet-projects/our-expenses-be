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

func TestAddCategoryHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewAddCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddCategoryHandler_CategoryError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()

	cmd := command.NewCategoryCommandArgs{}

	// SUT
	sut := command.NewAddCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddCategoryHandler_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	parentID := "parentId"
	cmd := command.NewCategoryCommandArgs{
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.Name() == cmd.Name && cat.Path() == cmd.Path &&
			cat.Level() == cmd.Level && cat.ParentID() == cmd.ParentID
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewAddCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestAddCategoryHandler_RepoSuccess_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID := "parentId"
	cmd := command.NewCategoryCommandArgs{
		Name:     "name",
		ParentID: &parentID,
		Path:     "path",
		Level:    1,
	}

	matchCategoryFn := func(cat *domain.Category) bool {
		return cat.Name() == cmd.Name && cat.Path() == cmd.Path &&
			cat.Level() == cmd.Level && cat.ParentID() == cmd.ParentID
	}
	repo.On("Insert", mock.Anything,
		mock.MatchedBy(matchCategoryFn)).Return(&categoryID, nil)

	// SUT
	sut := command.NewAddCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, &categoryID, query, "Should return category id.")
	assert.Nil(t, err, "Error result should be nil.")
}
