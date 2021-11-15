package ports

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/server/httperr"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// HTTPServer represents HTTP server with application dependency.
type HTTPServer struct {
	app app.Application
}

// NewHTTPServer instantiates http server with application.
func NewHTTPServer(app *app.Application) HTTPServer {
	return HTTPServer{
		app: *app,
	}
}

// Signup creates a user.
func (h HTTPServer) Signup(echoCtx echo.Context) error {
	ctx, span := tracer.NewSpan(echoCtx.Request().Context(), "handle sign-up http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling sign-up HTTP request")

	var userCredentials UserCredentials
	bindErr := echoCtx.Bind(&userCredentials)
	if bindErr != nil {
		tracer.AddSpanError(span, bindErr)
		h.app.Logger.Error(ctx, "Invalid user credentials format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest,
			httperr.BadRequest("Credentials has invalid format"))
	}

	cmdArgs := command.SignUpCommand{
		Username: userCredentials.Username,
		Password: userCredentials.Password,
	}

	user, userErr := h.app.Commands.SignUp.Handle(ctx, cmdArgs)
	if userErr != nil {
		tracer.AddSpanError(span, userErr)
		h.app.Logger.Error(ctx, "Failed to sign-up a user", userErr)
		return echoCtx.JSON(http.StatusInternalServerError,
			httperr.InternalError(userErr))
	}

	return echoCtx.JSON(http.StatusOK, userToResponse(user))
}

// Login authenticates a user.
func (h HTTPServer) Login(echoCtx echo.Context) error {
	ctx, span := tracer.NewSpan(echoCtx.Request().Context(), "handle login http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling login HTTP request")

	var userCredentials UserCredentials
	bindErr := echoCtx.Bind(&userCredentials)
	if bindErr != nil {
		tracer.AddSpanError(span, bindErr)
		h.app.Logger.Error(ctx, "Invalid user credentials format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest,
			httperr.BadRequest("Credentials has invalid format"))
	}

	cmdArgs := command.LoginCommand{
		Username: userCredentials.Username,
		Password: userCredentials.Password,
	}

	user, userErr := h.app.Commands.Login.Handle(ctx, cmdArgs)
	if userErr != nil {
		if errors.Is(userErr, domain.ErrWrongPassword) {
			return echoCtx.JSON(http.StatusUnauthorized,
				httperr.Unauthorized("Wrong user or password"))
		}
		h.app.Logger.Error(ctx, "Failed to login", userErr)
		return echoCtx.JSON(http.StatusInternalServerError,
			httperr.InternalError(userErr))
	}

	return echoCtx.JSON(http.StatusOK, userToResponse(user))
}
