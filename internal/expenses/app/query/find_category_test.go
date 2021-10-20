package query_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewFindCategoryHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseCategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := query.NewFindCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryHandle_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseCategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryId,
	}

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryHandle_RepoSuccess_CategoryHasNoPath_ReturnsCategory(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseCategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	parentId1 := "parentId1"
	path := ""
	icon := "icon"
	category, _ := domain.NewCategory(categoryId, &parentId1, "name", &icon, 1, path)
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryId,
	}

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, category, query, "Should return category.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestFindCategoryHandle_RepoSuccess_ReturnsCategory(t *testing.T) {
	// Arrange
	repo := new(mocks.ExpenseCategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	parentId1 := "parentId1"
	parentId2 := "parentId2"
	icon := "icon"
	path := fmt.Sprintf("|%s|%s", parentId1, parentId2)
	category, _ := domain.NewCategory(categoryId, &parentId1, "name", &icon, 1, path)
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryId,
	}

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, category, query, "Should return category.")
	assert.Nil(t, err, "Error result should be nil.")
}
