package ports

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
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
		Price:      newExpense.Price,
		Currency:   newExpense.Currency,
		Quantity:   newExpense.Quantity,
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

// GenerateReport generates a new expense report.
func (h HTTPServer) GenerateReport(echoCtx echo.Context) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get report http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling get report HTTP request")

	from := time.Date(2021, time.July, 3, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, time.August, 3, 0, 0, 0, 0, time.UTC)
	queryArgs := query.FindExpensesQuery{
		From: from,
		To:   to,
	}

	expenseRpt, expenseRptErr := h.app.Queries.FindExpenses.Handle(ctx, queryArgs)
	if expenseRptErr != nil {
		h.app.Logger.Error(ctx, "Failed to create expense report", expenseRptErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(expenseRptErr))
	}

	response := ExpenseReport{
		Expenses: expensesToResponse(expenseRpt),
	}
	return echoCtx.JSON(http.StatusOK, response)
}

func expensesToResponse(domainExpenses []domain.Expense) []Expense {
	expenses := []Expense{}
	for _, exp := range domainExpenses {
		e := expenseToResponse(exp)
		expenses = append(expenses, e)
	}
	return expenses
}

func expenseToResponse(domainExpense domain.Expense) Expense {
	expense := Expense{
		Id: domainExpense.ID(),
		NewExpense: NewExpense{
			CategoryId: domainExpense.CategoryID(),
			Comment:    domainExpense.Comment(),
			Currency:   domainExpense.Currency(),
			Date:       domainExpense.Date(),
			Price:      domainExpense.Price(),
			Quantity:   domainExpense.Quantity(),
		},
	}
	return expense
}
