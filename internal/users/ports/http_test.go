package ports_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/ports"
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

// nolint:dupl
func TestSignUp_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	signUpHandler := new(mocks.SignUpHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			SignUp: signUpHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	expenseJSON := fmt.Sprintf(`invalid json`)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Signup(ctx)

	// Assert
	logger.AssertExpectations(t)
	signUpHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestSignUp_CommandFails_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	signUpHandler := new(mocks.SignUpHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			SignUp: signUpHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	username := "user1"
	password := "pass1"
	credentialsJSON := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	matchCredFn := func(cmd command.SignUpCommand) bool {
		return cmd.Username == username && cmd.Password == password
	}
	signUpHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCredFn)).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(credentialsJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Signup(ctx)

	// Assert
	logger.AssertExpectations(t)
	signUpHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

// nolint:dupl
func TestSignUp_HappyPath_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	signUpHandler := new(mocks.SignUpHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			SignUp: signUpHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	username := "user2"
	password := "pass2"
	credentialsJSON := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	user, _ := domain.NewUser("id", username, "hashedpass", "token", "refreshtoken")

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	matchCredFn := func(cmd command.SignUpCommand) bool {
		return cmd.Username == username && cmd.Password == password
	}
	signUpHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCredFn)).Return(user, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(credentialsJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Signup(ctx)

	// Assert
	logger.AssertExpectations(t)
	signUpHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
	assert.NotContains(t, response.Body.String(), password)
	assert.Contains(t, response.Body.String(), user.ID())
	assert.Contains(t, response.Body.String(), user.Token())
	assert.Contains(t, response.Body.String(), user.RefreshToken())
}

// nolint:dupl
func TestLogin_InvalidPayload_Returns400(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	loginHandler := new(mocks.LoginHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			Login: loginHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	expenseJSON := fmt.Sprintf(`invalid json`)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(expenseJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Login(ctx)

	// Assert
	logger.AssertExpectations(t)
	loginHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusBadRequest, response.Code, "HTTP status should be 400.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestLogin_CommandFails_Returns500(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	loginHandler := new(mocks.LoginHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			Login: loginHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	username := "user3"
	password := "pass3"
	credentialsJSON := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	matchCredFn := func(cmd command.LoginCommand) bool {
		return cmd.Username == username && cmd.Password == password
	}
	loginHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCredFn)).Return(nil, errors.New("error"))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(credentialsJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Login(ctx)

	// Assert
	logger.AssertExpectations(t)
	loginHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusInternalServerError, response.Code, "HTTP status should be 500.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

func TestLogin_CommandFailsWithWrongPassword_Returns401(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	loginHandler := new(mocks.LoginHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			Login: loginHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	username := "user3"
	password := "pass3"
	credentialsJSON := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()

	matchCredFn := func(cmd command.LoginCommand) bool {
		return cmd.Username == username && cmd.Password == password
	}
	loginHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCredFn)).Return(nil, domain.ErrWrongPassword)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(credentialsJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Login(ctx)

	// Assert
	logger.AssertExpectations(t)
	loginHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusUnauthorized, response.Code, "HTTP status should be 401.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
}

// nolint:dupl
func TestLogin_HappyPath_Returns200(t *testing.T) {
	t.Parallel()
	// Arrange
	e := echo.New()
	logger := new(mocks.LogInterface)
	loginHandler := new(mocks.LoginHandlerInterface)
	app := &app.Application{
		Commands: app.Commands{
			Login: loginHandler,
		},
		Queries: app.Queries{},
		Logger:  logger,
	}
	username := "user4"
	password := "pass4"
	credentialsJSON := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	user, _ := domain.NewUser("id", username, "hashedpass", "token", "refreshtoken")

	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Return()

	matchCredFn := func(cmd command.LoginCommand) bool {
		return cmd.Username == username && cmd.Password == password
	}
	loginHandler.On("Handle", mock.Anything, mock.MatchedBy(matchCredFn)).Return(user, nil)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/users", strings.NewReader(credentialsJSON))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx := e.NewContext(request, response)

	// SUT
	server := ports.NewHTTPServer(app)

	// Act
	server.Login(ctx)

	// Assert
	logger.AssertExpectations(t)
	loginHandler.AssertExpectations(t)
	assert.Equal(t, http.StatusOK, response.Code, "HTTP status should be 200.")
	assert.NotEmpty(t, response.Body.String(), "Should not return empty body.")
	assert.NotContains(t, response.Body.String(), password)
	assert.Contains(t, response.Body.String(), user.ID())
	assert.Contains(t, response.Body.String(), user.Token())
	assert.Contains(t, response.Body.String(), user.RefreshToken())
}
