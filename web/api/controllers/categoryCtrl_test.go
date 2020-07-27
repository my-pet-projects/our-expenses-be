package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProvideCategoryController_ReturnsController(t *testing.T) {
	logger := new(logger.AppLogger)
	repo := new(repositories.CategoryRepository)
	validator := new(validators.Validator)

	results := ProvideCategoryController(repo, logger, validator)

	assert.NotNil(t, results, "Controller should not be nil.")
}

func TestGetAllCategories_ReturnsCategoriesFromDatabase(t *testing.T) {
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

	categories := []models.Category{
		{
			ID:   "1",
			Name: "category 1",
			Path: "/path/to/category/1",
		},
		{
			ID:   "2",
			Name: "category 2",
			Path: "/path/to/category/2",
		},
	}

	repo.On("GetAll", mock.MatchedBy(func(ctx context.Context) bool { return ctx == request.Context() })).Return(categories, nil)

	ctrl.GetAllCategories(response, request)

	var responseCategories []models.Category
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseCategories)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.Equal(t, categories, responseCategories, "Categories in the response should be the same")
}

func TestGetAllCategories_Throws500ErrorOnDatabaseError(t *testing.T) {
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

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	repo.On("GetAll", mock.MatchedBy(func(ctx context.Context) bool { return ctx == request.Context() })).Return(nil, errors.New("error"))

	ctrl.GetAllCategories(response, request)

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
}

func TestCreateCategory_SavesCategoryInDatabase(t *testing.T) {
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	newID := "id"
	category := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(category)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Save", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx == request.Context()
	}), mock.Anything).Return(newID, nil)

	ctrl.CreateCategory(response, request)

	var responseCategory models.Category
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseCategory)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusCreated, response.Code, "HTTP status should be 201.")
	assert.Equal(t, responseCategory.ID, newID, "Category ID should be set after database call")
	assert.Equal(t, responseCategory.Name, category.Name, "Category Name should be the same")
	assert.Equal(t, responseCategory.Path, category.Path, "Category Path should be the same")
}

func TestCreateCategory_Returns400_WhenInvalidJson(t *testing.T) {
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

	ctrl.CreateCategory(response, request)

	var responseError map[string]interface{}
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &responseError)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, responseError["message"], "Response message property should not be empty.")
}

func TestCreateCategory_Returns400_WhenInvalidPayload(t *testing.T) {
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

	ctrl.CreateCategory(response, request)

	var errorResponse responses.ErrorResponse
	jsonErr := json.Unmarshal([]byte(response.Body.String()), &errorResponse)
	if jsonErr != nil {
		t.Errorf("Cannot convert to json: %v", jsonErr)
	}
	errorResponseMessageBytes, _ := json.Marshal(errorResponse.Message)
	var errorResponseValidationMessages []validators.ValidationError
	json.Unmarshal(errorResponseMessageBytes, &errorResponseValidationMessages)

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.Equal(t, validationErrors, errorResponseValidationMessages, "Should return validation error messages.")
}

func TestCreateCategory_Returns500_WhenDatabaseFailed(t *testing.T) {
	logger := new(mocks.AppLoggerInterface)
	repo := new(mocks.CategoryRepoInterface)
	validator := new(mocks.ValidatorInterface)
	ctrl := &CategoryController{
		repo:      repo,
		logger:    logger,
		validator: validator,
	}

	category := &requests.CreateCategoryRequest{
		Name: "category",
		Path: "/path/to/category",
	}
	jsonCategory, _ := json.Marshal(category)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonCategory))

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	validator.On("ValidateStruct", mock.Anything).Return(nil)
	repo.On("Save", mock.MatchedBy(func(ctx context.Context) bool {
		return ctx == request.Context()
	}), mock.Anything).Return("", errors.New("database error"))

	ctrl.CreateCategory(response, request)

	repo.AssertExpectations(t)
	validator.AssertExpectations(t)
	logger.AssertExpectations(t)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}
