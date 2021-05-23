package query_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewFindCategoryHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := query.NewFindCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryHandle_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, categoryId)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryHandle_RepoSuccess_CategoryHasNoPath_ReturnsCategory(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	parentId1 := "parentId1"
	path := ""
	category, _ := domain.NewCategory(categoryId, "name", &parentId1, path, 1)

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, categoryId)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestFindCategoryHandle_RepoSuccess_AndParentsCategories_RepoSuccess_ReturnsCategoryWithParents(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	parentId1 := "parentId1"
	parentId2 := "parentId2"
	path := fmt.Sprintf("|%s|%s", parentId1, parentId2)
	category, _ := domain.NewCategory(categoryId, "name", &parentId1, path, 1)
	parentCategory1, _ := domain.NewCategory(parentId1, "name1", nil, path, 1)
	parentCategory2, _ := domain.NewCategory(parentId2, "name1", nil, path, 1)
	parentCategories := []domain.Category{*parentCategory1, *parentCategory2}
	parentFilter := domain.CategoryFilter{CategoryIDs: []string{parentId1, parentId2}}

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, parentFilter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(parentCategories, nil)

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, categoryId)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, parentCategories, query.Parents(), "Parents should match.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestFindCategoryHandle_RepoSuccess_AndParentsCategories_RepoError_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryId := "categoryId"
	parentId1 := "parentId1"
	parentId2 := "parentId2"
	path := fmt.Sprintf("|%s|%s", parentId1, parentId2)
	category, _ := domain.NewCategory(categoryId, "name", &parentId1, path, 1)
	parentFilter := domain.CategoryFilter{CategoryIDs: []string{parentId1, parentId2}}

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, parentFilter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, categoryId)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}
