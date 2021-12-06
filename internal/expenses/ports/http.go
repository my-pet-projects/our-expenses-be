package ports

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
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

// AddExpense adds a new expense.
func (h HTTPServer) AddExpense(echoCtx echo.Context) error {
	ctx, span := tracer.NewSpan(echoCtx.Request().Context(), "handle add expense http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling add expense HTTP request")

	var newExpense NewExpense
	bindErr := echoCtx.Bind(&newExpense)
	if bindErr != nil {
		tracer.AddSpanError(span, bindErr)
		h.app.Logger.Error(ctx, "Invalid expense format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest,
			httperr.BadRequest("Invalid expense format"))
	}

	catQuery := query.FindCategoryQuery{
		CategoryID: newExpense.CategoryId,
	}
	category, categoryErr := h.app.Queries.FindCategory.Handle(ctx, catQuery)
	if categoryErr != nil {
		tracer.AddSpanError(span, categoryErr)
		h.app.Logger.Error(ctx, "Failed to get category", categoryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(categoryErr))
	}

	if category == nil {
		return echoCtx.JSON(http.StatusBadRequest,
			httperr.BadRequest(fmt.Sprintf("Invalid provided category with ID %s", newExpense.CategoryId)))
	}

	cmdArgs := command.AddExpenseCommand{
		Category: *category,
		Price:    newExpense.Price,
		Currency: newExpense.Currency,
		Quantity: newExpense.Quantity,
		Comment:  newExpense.Comment,
		Trip:     newExpense.Trip,
		Date:     newExpense.Date,
	}
	expenseID, expenseCrtErr := h.app.Commands.AddExpense.Handle(ctx, cmdArgs)
	if expenseCrtErr != nil {
		tracer.AddSpanError(span, expenseCrtErr)
		h.app.Logger.Error(ctx, "Failed to create expense", expenseCrtErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(expenseCrtErr))
	}

	response := NewExpenseResponse{
		Id: *expenseID,
	}

	return echoCtx.JSON(http.StatusCreated, response)
}

// GenerateReport generates a new expense report.
func (h HTTPServer) GenerateReport(echoCtx echo.Context, params GenerateReportParams) error {
	ctx, span := tracer.NewSpan(echoCtx.Request().Context(), "handle generate report http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling generate report HTTP request")

	dateRange, dateRangeErr := domain.NewDateRange(params.From, params.To)
	if dateRangeErr != nil {
		tracer.AddSpanError(span, dateRangeErr)
		return echoCtx.JSON(http.StatusBadRequest,
			httperr.BadRequest("Date range has invalid format"))
	}

	fetchCmdArgs := command.FetchExchangeRatesCommand{
		DateRange: *dateRange,
	}
	rates, ratesErr := h.app.Commands.FetchExchangeRates.Handle(ctx, fetchCmdArgs)
	if ratesErr != nil {
		tracer.AddSpanError(span, ratesErr)
		h.app.Logger.Error(ctx, "Failed to fetch exchange rates", ratesErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(ratesErr))
	}

	queryArgs := query.FindExpensesQuery{
		DateRange:     *dateRange,
		Interval:      string(params.Interval),
		ExchangeRates: rates,
	}

	expenseRpt, expenseRptErr := h.app.Queries.FindExpenses.Handle(ctx, queryArgs)
	if expenseRptErr != nil {
		tracer.AddSpanError(span, expenseRptErr)
		h.app.Logger.Error(ctx, "Failed to create expense report", expenseRptErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(expenseRptErr))
	}

	response := reportToResponse(*expenseRpt)
	return echoCtx.JSON(http.StatusOK, response)
}
