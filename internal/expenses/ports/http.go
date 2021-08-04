package ports

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/server/httperr"
)

var tracer trace.Tracer

// HTTPServer represents HTTP server with application dependency.
type HTTPServer struct {
	app app.Application
}

// NewHTTPServer instantiates http server with application.
func NewHTTPServer(app *app.Application) HTTPServer {
	tracer = otel.Tracer("ports.http")
	return HTTPServer{
		app: *app,
	}
}

// AddExpense adds a new expense.
func (h HTTPServer) AddExpense(echoCtx echo.Context) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get expense http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling add expense HTTP request")

	var newExpense NewExpense
	bindErr := echoCtx.Bind(&newExpense)
	if bindErr != nil {
		expenseErr := Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid format for a new expense",
		}
		h.app.Logger.Error(ctx, "Invalid expense format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest, expenseErr)
	}

	cmdArgs := command.AddExpenseCommand{
		CategoryID: newExpense.CategoryId,
		Price:      fmt.Sprint(newExpense.Price),
		Currency:   newExpense.Currency,
		Quantity:   fmt.Sprint(newExpense.Quantity),
		Comment:    newExpense.Comment,
		Date:       newExpense.Date,
	}
	expenseID, expenseCrtErr := h.app.Commands.AddExpense.Handle(ctx, cmdArgs)
	if expenseCrtErr != nil {
		h.app.Logger.Error(ctx, "Failed to create expense", expenseCrtErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(expenseCrtErr))
	}

	response := NewExpenseResponse{
		Id: *expenseID,
	}

	return echoCtx.JSON(http.StatusCreated, response)
}
