package command_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestMoveCategoryHandler_ReturnsHandler(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)

	// Act
	err := command.NewMoveCategoryHandler(repo, log)

	// Assert
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestMoveCategoryHandler_FailedToGetCategory_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestMoveCategoryHandler_NoCategoryFound_ReturnsEmptyResult(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(nil, nil)

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestMoveCategoryHandler_FailedToGetCategoryUsages_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}
	pathFilter := domain.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}
	parentID := "parentId"
	path := fmt.Sprintf("|%s", parentID)
	category, _ := domain.NewCategory(categoryID, "name", &parentID, path, nil, 1)

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything,
		mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, pathFilter)
	}
	repo.On("GetAll", mock.Anything,
		mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestMoveCategoryHandler_FailedToGetDestinationCategory_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}
	pathFilter := domain.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}
	parentID := "parentId"
	path := fmt.Sprintf("|%s", parentID)
	category, _ := domain.NewCategory(categoryID, "name", &parentID, path, nil, 1)
	categories := []domain.Category{{}}

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, pathFilter)
	}
	repo.On("GetAll", mock.Anything, mock.MatchedBy(matchFilterFn)).Return(categories, nil)
	matchDestIdFn := func(id string) bool {
		return id == destinationID
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchDestIdFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestMoveCategoryHandler_NoDestinationCategory_ReturnsEmptyResult(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}
	pathFilter := domain.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}
	parentID := "parentId"
	path := fmt.Sprintf("|%s", parentID)
	category, _ := domain.NewCategory(categoryID, "name", &parentID, path, nil, 1)
	categories := []domain.Category{{}}

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, pathFilter)
	}
	repo.On("GetAll", mock.Anything, mock.MatchedBy(matchFilterFn)).Return(categories, nil)
	matchDestIdFn := func(id string) bool {
		return id == destinationID
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchDestIdFn)).Return(nil, nil)

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}

func TestMoveCategoryHandler_FailedToUpdateCategory_ThrowsError(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	categoryID := "categoryId"
	destinationID := "destinationID"
	cmd := command.MoveCategoryCommand{
		CategoryID:    categoryID,
		DestinationID: destinationID,
	}
	pathFilter := domain.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}
	parentID := "parentId"
	path := fmt.Sprintf("|%s", parentID)
	category, _ := domain.NewCategory(categoryID, "name", &parentID, path, nil, 1)
	destCategory, _ := domain.NewCategory(destinationID, "dest name", &parentID, path, nil, 1)
	categories := []domain.Category{{}}

	matchIdFn := func(id string) bool {
		return categoryID == id
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchIdFn)).Return(category, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, pathFilter)
	}
	repo.On("GetAll", mock.Anything, mock.MatchedBy(matchFilterFn)).Return(categories, nil)
	matchDestIdFn := func(id string) bool {
		return id == destinationID
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchDestIdFn)).Return(destCategory, nil)
	matchUpdFn := func(category domain.Category) bool {
		return category.ID() == categoryID
	}
	repo.On("Update", mock.Anything, mock.MatchedBy(matchUpdFn)).Return(nil, errors.New("error"))

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error result should not be nil.")
}

func TestMoveCategoryHandler_UpdateCategory_ReturnsResult(t *testing.T) {
	// Arrange
	repo := new(mocks.CategoryRepoInterface)
	log := new(mocks.LogInterface)
	ctx := context.Background()
	targetID := "targetID"
	targetParentID := "targetParentID"
	targetPath := fmt.Sprintf("|%s|%s", targetParentID, targetID)
	targetLevel := 2
	destinationID := "destinationID"
	destinationParentID := "destinationParentID"
	destinationPath := fmt.Sprintf("|level1|%s|%s", destinationParentID, destinationID)
	destinationLevel := 3
	childID := "childID"
	childParentID := "childParentID"
	childPath := fmt.Sprintf("%s|%s|%s", targetPath, childParentID, childID)
	childLevel := 4
	icon := "icon"
	cmd := command.MoveCategoryCommand{
		CategoryID:    targetID,
		DestinationID: destinationID,
	}
	pathFilter := domain.CategoryFilter{
		CategoryID:   targetID,
		FindChildren: true,
	}
	targetCat, _ := domain.NewCategory(targetID, "target name", &targetParentID, targetPath, &icon, targetLevel)
	destCategory, _ := domain.NewCategory(destinationID, "dest name", &destinationParentID, destinationPath, &icon, destinationLevel)
	childCategory, _ := domain.NewCategory(childID, "child name", &childParentID, childPath, &icon, childLevel)
	categories := []domain.Category{*childCategory}
	updateResult := &domain.UpdateResult{
		UpdateCount: 5,
	}

	matchIdFn := func(id string) bool {
		return targetID == id
	}
	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchIdFn)).Return(targetCat, nil)
	matchFilterFn := func(filter domain.CategoryFilter) bool {
		return reflect.DeepEqual(filter, pathFilter)
	}
	repo.On("GetAll", mock.Anything, mock.MatchedBy(matchFilterFn)).Return(categories, nil)
	matchDestIdFn := func(id string) bool {
		return id == destinationID
	}

	repo.On("GetOne", mock.Anything, mock.MatchedBy(matchDestIdFn)).Return(destCategory, nil)
	matchUpdTargetFn := func(category domain.Category) bool {
		return category.ID() == targetID && category.Path() == fmt.Sprintf("%s|%s", destinationPath, targetID) &&
			category.Level() == destinationLevel+1
	}
	repo.On("Update", mock.Anything, mock.MatchedBy(matchUpdTargetFn)).Return(updateResult, nil)
	matchUpdChildFn := func(category domain.Category) bool {
		return category.ID() == childID && category.Path() == fmt.Sprintf("%s|%s|%s|%s", destinationPath, targetID, childParentID, childID) &&
			category.Level() == 6
	}
	repo.On("Update", mock.Anything, mock.MatchedBy(matchUpdChildFn)).Return(updateResult, nil)

	// SUT
	sut := command.NewMoveCategoryHandler(repo, log)

	// Act
	result, err := sut.Handle(ctx, cmd)

	// Assert
	repo.AssertExpectations(t)
	assert.Nil(t, result, "Result should be nil.")
	assert.Nil(t, err, "Error result should be nil.")
}
