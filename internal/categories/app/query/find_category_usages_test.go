package query_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewFindCategoryUsagesHandler_ReturnsHandler(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := query.NewFindCategoryUsagesHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryUsagesHandle_RepoError_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	filter := domain.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}
	findQuery := query.FindCategoryUsagesQuery{
		CategoryID: categoryID,
	}

	matchFilterFn := func(f domain.CategoryFilter) bool {
		return reflect.DeepEqual(f, filter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindCategoryUsagesHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryUsagesHandle_RepoSuccess_ReturnsCategories(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	filter := domain.CategoryFilter{
		CategoryID:   "categoryId",
		FindChildren: true,
	}
	findQuery := query.FindCategoryUsagesQuery{
		CategoryID: categoryID,
	}
	categories := []domain.Category{{}}

	matchFilterFn := func(f domain.CategoryFilter) bool {
		return reflect.DeepEqual(f, filter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(categories, nil)

	// SUT
	sut := query.NewFindCategoryUsagesHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, categories, query, "Should return categories.")
	assert.Nil(t, err, "Error result should be nil.")
}
