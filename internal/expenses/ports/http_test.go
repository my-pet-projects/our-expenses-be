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

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/ports"
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

func TestAddExpense_SuccessfulCommand_Returns201(t *testing.T) {
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
	expenseJSON := `{"categoryId":"123"}`
	expenseId := "expenseId"

	matchCatFn := func(command command.AddExpenseCommand) bool {
		return command.CategoryID == "123"
	}
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.MatchedBy(matchCatFn)).Return(&expenseId, nil)

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
	assert.Equal(t, http.StatusCreated, response.Code, "HTTP status should be 201.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_FailedCommand_Returns500(t *testing.T) {
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
	expenseJSON := `{"name":"expense"}`

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	handler.On("Handle", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
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
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestAddExpense_InvalidPayload_Returns400(t *testing.T) {
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
