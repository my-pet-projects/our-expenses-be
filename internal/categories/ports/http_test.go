package ports_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/ports"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewHTTPServer_ReturnsServer(t *testing.T) {
	t.Parallel()
	// Arrange
	app := &app.Application{}

	// Act
	result := ports.NewHTTPServer(app)

	// Assert
	assert.NotNil(t, result, "Result result should not be nil.")
}

func TestFindCategories_SuccessfulQuery_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoriesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategories: handler,
		},
		Logger: logger,
	}
	categories := []domain.Category{{}}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(categories, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)
	parentID := "parentID"
	allChildren := true
	all := true
	params := ports.FindCategoriesParams{
		ParentId:    &parentID,
		AllChildren: &allChildren,
		All:         &all,
	}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategories(ctx, params)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategories_FailedQuery_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoriesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategories: handler,
		},
		Logger: logger,
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

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
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_SuccessfulQuery_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: handler,
		},
		Logger: logger,
	}
	categoryID := "categoryID"
	category := domain.Category{}
	category.SetParents([]domain.Category{{}})

	matchIDFn := func(q query.FindCategoryQuery) bool {
		return q.CategoryID == categoryID
	}
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchIDFn)).Return(&category, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_FailedQuery_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: handler,
		},
		Logger: logger,
	}
	categoryID := "categoryID"

	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategory_NilQueryResult_Returns404(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategory: handler,
		},
		Logger: logger,
	}
	categoryID := "categoryID"

	matchIDFn := func(q query.FindCategoryQuery) bool {
		return q.CategoryID == categoryID
	}
	logger.On("Infof", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchIDFn)).Return(nil, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryByID(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, response.Code, "HTTP status should be 404.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddCategory_SuccessfulCommand_Returns201(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.AddCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := `{"name":"category"}`
	categoryID := "categoryID"

	matchCatFn := func(command command.AddCategoryCommand) bool {
		return command.Name == "category"
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(&categoryID, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddCategory(ctx)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusCreated, response.Code, "HTTP status should be 201.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddCategory_FailedCommand_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.AddCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := `{"name":"category"}`

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddCategory(ctx)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddCategory_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.AddCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := "invalid"

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddCategory(ctx)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestUpdateCategory_SuccessfulCommand_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.UpdateCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			UpdateCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := `{"name":"category"}`
	categoryID := "categoryID"
	updateResult := &domain.UpdateResult{UpdateCount: 10}

	matchCatFn := func(args command.UpdateCategoryCommand) bool {
		return args.Name == "category"
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(updateResult, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.UpdateCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestUpdateCategory_FailedCommand_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.UpdateCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			UpdateCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := `{"name":"category"}`
	categoryID := "categoryID"

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.UpdateCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestUpdateCategory_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.UpdateCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			UpdateCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryJSON := "invalid"
	categoryID := "categoryID"

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", strings.NewReader(categoryJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.UpdateCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestDeleteCategory_SuccessfulCommand_Returns204(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.DeleteCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			DeleteCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryID := "categoryID"
	deleteResult := &domain.DeleteResult{DeleteCount: 10}

	matchCatFn := func(cmd command.DeleteCategoryCommand) bool {
		return cmd.CategoryID == categoryID
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(deleteResult, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.DeleteCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusNoContent, response.Code, "HTTP status should be 204.")
	assert.Empty(t, response.Body.String(), "Should return empty body.")
}

func TestDeleteCategory_FailedCommand_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.DeleteCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			DeleteCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryID := "categoryID"

	matchCatFn := func(cmd command.DeleteCategoryCommand) bool {
		return cmd.CategoryID == categoryID
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(nil, errors.New("error"))
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.DeleteCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should return not empty body.")
}

func TestDeleteCategory_EmptyDeleteResult_Returns404(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.DeleteCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			DeleteCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryID := "categoryID"

	matchCatFn := func(cmd command.DeleteCategoryCommand) bool {
		return cmd.CategoryID == categoryID
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(nil, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.DeleteCategory(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusNotFound, response.Code, "HTTP status should be 404.")
	assert.NotEmpty(t, response.Body.String(), "Should return not empty body.")
}

func TestFindCategoryUsages_SuccessfulQuery_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoryUsagesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategoryUsages: handler,
		},
		Logger: logger,
	}
	categories := []domain.Category{{}}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(categories, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)
	categoryID := "categoryID"

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryUsages(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestFindCategoryUsages_FailedQuery_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindCategoryUsagesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindCategoryUsages: handler,
		},
		Logger: logger,
	}
	categoryID := "categoryID"

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.FindCategoryUsages(ctx, categoryID)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestMoveCategory_SuccessfulCommand_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.MoveCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			MoveCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	updateResult := &domain.UpdateResult{
		UpdateCount: 10,
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(updateResult, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)
	categoryID := "categoryID"
	destinationID := "destinationId"
	params := ports.MoveCategoryParams{
		DestinationId: destinationID,
	}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.MoveCategory(ctx, categoryID, params)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestMoveCategory_FailedCommand_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.MoveCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			MoveCategory: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	categoryID := "categoryID"
	destinationID := "destinationId"
	params := ports.MoveCategoryParams{
		DestinationId: destinationID,
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/categories", nil)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.MoveCategory(ctx, categoryID, params)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}
