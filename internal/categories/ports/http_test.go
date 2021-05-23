package ports_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/ports"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewHTTPServer_ReturnsServer(t *testing.T) {
	// Arrange
	app := &app.Application{}

	// Act
	result := ports.NewHTTPServer(app)

	// Assert
	assert.NotNil(t, result, "Result result should not be nil.")
}

func TestFindCategories_SuccessfulQuery_Returns200(t *testing.T) {
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	query := new(mocks.FindCategoriesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategories: query,
		},
		Logger: logger,
	}
	categories := []domain.Category{{}}

	query.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(categories, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)
	parentID := "parentID"
	allChildren := true
	params := ports.FindCategoriesParams{
		ParentId:    &parentID,
		AllChildren: &allChildren,
	}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategories(ctx, params)

	// Assert
	logger.AssertExpectations(t)
	query.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategories_FailedQuery_Returns500(t *testing.T) {
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	query := new(mocks.FindCategoriesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategories: query,
		},
		Logger: logger,
	}

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	query.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)
	params := ports.FindCategoriesParams{}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategories(ctx, params)

	// Assert
	logger.AssertExpectations(t)
	query.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_SuccessfulQuery_Returns200(t *testing.T) {
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	query := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: query,
		},
		Logger: logger,
	}
	categoryId := "categoryId"
	category := domain.Category{}
	category.SetParents([]domain.Category{{}})

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	query.On("Handle", mock.Anything, mock.MatchedBy(matchIdFn)).Return(&category, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryId)

	// Assert
	logger.AssertExpectations(t)
	query.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_FailedQuery_Returns500(t *testing.T) {
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	query := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: query,
		},
		Logger: logger,
	}
	categoryID := "categoryId"

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	query.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	query.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_NilQueryResult_Returns404(t *testing.T) {
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	query := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: query,
		},
		Logger: logger,
	}
	categoryId := "categoryId"

	matchIdFn := func(id string) bool {
		return id == categoryId
	}
	query.On("Handle", mock.Anything, mock.MatchedBy(matchIdFn)).Return(nil, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryId)

	// Assert
	logger.AssertExpectations(t)
	query.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, response.Code, "HTTP status should be 404.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}
