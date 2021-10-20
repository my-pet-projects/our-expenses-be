package ports_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/ports"
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

func TestAddExpense_SuccessfulCommand_Returns201(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	addExpenseHandler := new(mocks.AddExpenseHandlerInterface)
	findCategoryHandler := new(mocks.FindExpenseCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddExpense: addExpenseHandler,
		},
		Queries: app.Queries{
			FindCategory: findCategoryHandler,
		},
		Logger: logger,
	}
	categoryID := "123"
	expenseJSON := fmt.Sprintf(`{"categoryId":"%s"}`, categoryID)
	expenseId := "expenseId"
	category, _ := domain.NewCategory(categoryID, nil, "category", nil, 1, "path")

	matchCatFn := func(query query.FindCategoryQuery) bool {
		return query.CategoryID == categoryID
	}
	findCategoryHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(category, nil)

	matchExpFn := func(command command.AddExpenseCommand) bool {
		return reflect.DeepEqual(command.Category, *category)
	}
	addExpenseHandler.On("Handle", mock.Anything, mock.MatchedBy(matchExpFn)).Return(&expenseId, nil)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/expenses", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddExpense(ctx)

	// Assert
	logger.AssertExpectations(t)
	addExpenseHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusCreated, response.Code, "HTTP status should be 201.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_FailedCategoryQuery_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	expenseHandler := new(mocks.AddExpenseHandlerInterface)
	findCategoryHandler := new(mocks.FindExpenseCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddExpense: expenseHandler,
		},
		Queries: app.Queries{
			FindCategory: findCategoryHandler,
		},
		Logger: logger,
	}
	categoryID := "123"
	expenseJSON := fmt.Sprintf(`{"categoryId":"%s"}`, categoryID)

	findCategoryHandler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/expenses", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddExpense(ctx)

	// Assert
	logger.AssertExpectations(t)
	expenseHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_CategoryQueryNotFound_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	expenseHandler := new(mocks.AddExpenseHandlerInterface)
	findCategoryHandler := new(mocks.FindExpenseCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddExpense: expenseHandler,
		},
		Queries: app.Queries{
			FindCategory: findCategoryHandler,
		},
		Logger: logger,
	}
	categoryID := "123"
	expenseJSON := fmt.Sprintf(`{"categoryId":"%s"}`, categoryID)

	findCategoryHandler.On("Handle", mock.Anything, mock.Anything).Return(nil, nil)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/expenses", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddExpense(ctx)

	// Assert
	logger.AssertExpectations(t)
	expenseHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_FailedExpenseCommand_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	expenseHandler := new(mocks.AddExpenseHandlerInterface)
	findCategoryHandler := new(mocks.FindExpenseCategoryHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddExpense: expenseHandler,
		},
		Queries: app.Queries{
			FindCategory: findCategoryHandler,
		},
		Logger: logger,
	}
	categoryID := "123"
	expenseJSON := fmt.Sprintf(`{"categoryId":"%s"}`, categoryID)
	category, _ := domain.NewCategory(categoryID, nil, "category", nil, 1, "path")

	matchCatFn := func(query query.FindCategoryQuery) bool {
		return query.CategoryID == categoryID
	}
	findCategoryHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(category, nil)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	expenseHandler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/expenses", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddExpense(ctx)

	// Assert
	logger.AssertExpectations(t)
	expenseHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.AddExpenseHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			AddExpense: handler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	expenseJSON := "invalid"

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/expenses", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.AddExpense(ctx)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestGenerateReport_SuccessfulQuery_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindExpensesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindExpenses: handler,
		},
		Logger: logger,
	}
	from := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, time.August, 3, 0, 0, 0, 0, time.UTC)
	report := &domain.ReportByDate{}

	matchFn := func(query query.FindExpensesQuery) bool {
		return query.From == from && query.To == to
	}
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchFn)).Return(report, nil)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/reports", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)
	params := ports.GenerateReportParams{
		To:   to,
		From: from,
	}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.GenerateReport(ctx, params)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestGenerateReport_FailedQuery_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	handler := new(mocks.FindExpensesHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			FindExpenses: handler,
		},
		Logger: logger,
	}

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/reports", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)
	params := ports.GenerateReportParams{
		To:   time.Now(),
		From: time.Now(),
	}

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.GenerateReport(ctx, params)

	// Assert
	logger.AssertExpectations(t)
	handler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}
