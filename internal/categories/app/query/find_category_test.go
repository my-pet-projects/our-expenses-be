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
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := query.NewFindCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestFindCategoryHandle_RepoError_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return id == categoryID
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(nil, errors.New("error"))

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
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID1 := "parentId1"
	path := ""
	icon := "icon"
	category, _ := domain.NewCategory(categoryID, "name", &parentID1, path, &icon, 1)
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return id == categoryID
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(category, nil)

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

func TestFindCategoryHandle_RepoSuccess_AndParentsCategories_RepoSuccess_ReturnsCategoryWithParents(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID1 := "parentId1"
	parentID2 := "parentId2"
	icon := "icon"
	path := fmt.Sprintf("|%s|%s", parentID1, parentID2)
	category, _ := domain.NewCategory(categoryID, "name", &parentID1, path, &icon, 1)
	parentCategory1, _ := domain.NewCategory(parentID1, "name1", nil, path, &icon, 1)
	parentCategory2, _ := domain.NewCategory(parentID2, "name1", nil, path, &icon, 1)
	parentCategories := []domain.Category{*parentCategory1, *parentCategory2}
	parentFilter := domain.CategoryFilter{CategoryIDs: []string{parentID1, parentID2}}
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return id == categoryID
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, parentFilter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(parentCategories, nil)

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.NotNil(t, query, "Result should not be nil.")
	assert.Equal(t, category, query, "Should return category.")
	assert.Equal(t, parentCategories, query.Parents(), "Parents should match.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestFindCategoryHandle_RepoSuccess_AndParentsCategories_RepoError_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	parentID1 := "parentId1"
	parentID2 := "parentId2"
	icon := "icon"
	path := fmt.Sprintf("|%s|%s", parentID1, parentID2)
	category, _ := domain.NewCategory(categoryID, "name", &parentID1, path, &icon, 1)
	parentFilter := domain.CategoryFilter{CategoryIDs: []string{parentID1, parentID2}}
	findQuery := query.FindCategoryQuery{
		CategoryID: categoryID,
	}

	matchIDFn := func(id string) bool {
		return id == categoryID
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIDFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, parentFilter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := query.NewFindCategoryHandler(repo, log)

	// Act
	query, err := sut.Handle(ctx, findQuery)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, query, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}
