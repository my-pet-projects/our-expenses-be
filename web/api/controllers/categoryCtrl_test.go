package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"our-expenses-server/db/repositories"
	"our-expenses-server/logger"
	"our-expenses-server/models"
	"our-expenses-server/testing/mocks"
	"our-expenses-server/validators"
	"our-expenses-server/web/requests"
	"our-expenses-server/web/responses"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestProvideCategoryController_ReturnsController(t *testing.T) {
	logger := new(logger.AppLogger)
	repo := new(repositories.CategoryRepository)
	validator := new(validators.Validator)

	results := ProvideCategoryController(repo, logger, validator)

	assert.NotNil(t, results, "Controller should not be nil.")
}

func TestGetAllCategories_ReturnsCategoriesFromDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	category1ID := primitive.NewObjectID()
	category2ID := primitive.NewObjectID()
	categories := []models.Category{
		{
			ID:   &category1ID,
			Name: "category 1",
			Path: "/path/to/category/1",
		},
		{
			ID:   &category2ID,
			Name: "category 2",
			Path: "/path/to/category/2",
		},
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter models.CategoryFilter) bool {
		return filter.ParentID == ""
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetAll", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(categories, nil)

	// Act
	ctrl.GetAllCategories(response, request)

	var responseCategories []models.Category
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseCategories)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.Equal(t, categories, responseCategories, "Categories in the response should be the same")
}

func TestGetAllCategories_Throws500ErrorOnDatabaseError(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter models.CategoryFilter) bool {
		return filter.ParentID == ""
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetAll", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("error"))

	// Act
	ctrl.GetAllCategories(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
}

func TestCreateCategory_SavesCategoryInDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	categoryRequest := &requests.CreateCategoryRequest{
		Name:     "category",
		Path:     "/path/to/category",
		ParentID: &parentID,
	}
	category := &models.Category{
		ID:       &categoryID,
		Name:     categoryRequest.Name,
		Path:     categoryRequest.Path,
		ParentID: categoryRequest.ParentID,
	}
	jsonCategory, _ := json.Marshal(categoryRequest)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchCategoryFn := func(cat *models.Category) bool {
		return cat.Name == category.Name &&
			cat.ParentID.Hex() == category.ParentID.Hex() && cat.Path == category.Path
	}

	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Insert", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchCategoryFn)).Return(category, nil)

	// Act
	ctrl.CreateCategory(response, request)

	var responseCategory models.Category
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseCategory)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusCreated, response.Code, "HTTP status should be 201.")
	assert.Equal(t, responseCategory.ID, category.ID, "Category ID should be set after database call")
	assert.Equal(t, responseCategory.Name, category.Name, "Category Name should be the same")
	assert.Equal(t, responseCategory.Path, category.Path, "Category Path should be the same")
	assert.Equal(t, responseCategory.ParentID, category.ParentID, "Category ParentID should be the same")
}

func TestCreateCategory_Returns400_WhenInvalidJson(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	jsonCategory, _ := json.Marshal(`{"invalid json":`)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	ctrl.CreateCategory(response, request)

	var responseError map[string]interface{}
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseError)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, responseError["message"], "Response message property should not be empty.")
}

func TestCreateCategory_Returns400_WhenInvalidPayload(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	validationErrors := []validators.ValidationError{
		{Details: "required", Field: "field1"},
	}

	category := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(category)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	validator.On("ValidateStruct", mock.Anything).Once().Return(validationErrors)

	// Act
	ctrl.CreateCategory(response, request)

	var errorResponse responses.ErrorResponse
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &errorResponse)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}
	errorResponseMessageBytes, _ := json.Marshal(errorResponse.Message)
	var errorResponseValidationMessages []validators.ValidationError
	json.Unmarshal(errorResponseMessageBytes, &errorResponseValidationMessages)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.Equal(t, validationErrors, errorResponseValidationMessages, "Should return validation error messages.")
}

func TestCreateCategory_Returns500_WhenDatabaseFailed(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryRequest := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(categoryRequest)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Insert", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx == request.Context()
	}), mock.Anything).Return(nil, errors.New("database error"))

	// Act
	ctrl.CreateCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestGetCategory_ReturnsCategoryFromDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	category := &models.Category{
		ID:       &categoryID,
		Name:     "category",
		ParentID: &parentID,
		Path:     "/path/to/category",
	}

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(category, nil)

	// Act
	ctrl.GetCategory(response, request)

	var responseCategory models.Category
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseCategory)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.Equal(t, category, &responseCategory, "Categories in the response should be the same")
}

func TestGetCategory_Returns404_WhenCategoryNotFoundInDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(nil, nil)

	// Act
	ctrl.GetCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusNotFound, response.Code, "HTTP status should be 404.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestGetCategory_Returns500_WhenDatabaseThrowsError(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(nil, errors.New("database error"))

	// Act
	ctrl.GetCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestUpdateCategory_UpdatesCategoryInDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	parentID := primitive.NewObjectID()
	categoryRequest := &requests.CreateCategoryRequest{
		Name:     "category",
		Path:     "/path/to/category",
		ParentID: &parentID,
	}
	category := &models.Category{
		ID:       &categoryID,
		Name:     categoryRequest.Name,
		Path:     categoryRequest.Path,
		ParentID: categoryRequest.ParentID,
	}
	jsonCategory, _ := json.Marshal(categoryRequest)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%s", categoryID.Hex()), bytes.NewBuffer(jsonCategory))
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchCategoryFn := func(cat *models.Category) bool {
		return cat.ID.Hex() == category.ID.Hex() && cat.Name == category.Name &&
			cat.ParentID.Hex() == category.ParentID.Hex() && cat.Path == category.Path
	}

	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Update", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchCategoryFn)).Return(categoryID.Hex(), nil)

	// Act
	ctrl.UpdateCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusNoContent, response.Code, "HTTP status should be 204.")
	assert.Equal(t, response.Body.String(), "", "Response should be empty")
}

func TestUpdateCategory_Returns400_WhenInvalidJson(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	jsonCategory, _ := json.Marshal(`{"invalid json":`)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%s", categoryID.Hex()), bytes.NewBuffer(jsonCategory))
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	// Act
	ctrl.UpdateCategory(response, request)

	var responseError map[string]interface{}
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseError)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, responseError["message"], "Response message property should not be empty.")
}

func TestUpdateCategory_Returns400_WhenInvalidPayload(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	validationErrors := []validators.ValidationError{
		{Details: "required", Field: "field1"},
	}

	categoryID := primitive.NewObjectID()
	category := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(category)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%s", categoryID.Hex()), bytes.NewBuffer(jsonCategory))
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	validator.On("ValidateStruct", mock.Anything).Once().Return(validationErrors)

	// Act
	ctrl.UpdateCategory(response, request)

	var errorResponse responses.ErrorResponse
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &errorResponse)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}
	errorResponseMessageBytes, _ := json.Marshal(errorResponse.Message)
	var errorResponseValidationMessages []validators.ValidationError
	json.Unmarshal(errorResponseMessageBytes, &errorResponseValidationMessages)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.Equal(t, validationErrors, errorResponseValidationMessages, "Should return validation error messages.")
}

func TestUpdateCategory_Returns500_WhenDatabaseFailed(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	categoryRequest := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(categoryRequest)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("PUT", fmt.Sprintf("/categories/%s", categoryID.Hex()), bytes.NewBuffer(jsonCategory))
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx == request.Context()
	}), mock.Anything).Return("", errors.New("database error"))

	// Act
	ctrl.UpdateCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestDeleteCategory_ReturnsCategoryFromDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	deletionsNumber := 15

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("DeleteOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(int64(deletionsNumber), nil)

	// Act
	ctrl.DeleteCategory(response, request)

	var responseDeleted int
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseDeleted)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.Equal(t, deletionsNumber, responseDeleted, "Number of deletions should be the same")
}

func TestDeleteCategory_Returns404_WhenCategoryNotFoundInDatabase(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	deletionsNumber := 0

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("DeleteOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(int64(deletionsNumber), nil)

	// Act
	ctrl.DeleteCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusNotFound, response.Code, "HTTP status should be 404.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestDeleteCategory_Returns500_WhenDatabaseThrowsError(t *testing.T) {
	// Arrange
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	categoryID := primitive.NewObjectID()
	deletionsNumber := 0

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", fmt.Sprintf("/categories/%s", categoryID.Hex()), nil)
	vars := map[string]string{
		"id": categoryID.Hex(),
	}
	request = mux.SetURLVars(request, vars)

	matchCtxFn := func(ctx context.Context) bool {
		return ctx == request.Context()
	}
	matchFilterFn := func(filter string) bool {
		return filter == categoryID.Hex()
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("DeleteOne", mock.MatchedBy(matchCtxFn), mock.MatchedBy(matchFilterFn)).Return(int64(deletionsNumber), errors.New("database error"))

	// Act
	ctrl.DeleteCategory(response, request)

	// Assert
	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}
