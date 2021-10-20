package ports

import (
	"fmt"
	"net/http"

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

	catQuery := query.FindCategoryQuery{
		CategoryID: newExpense.CategoryId,
	}

	category, categoryErr := h.app.Queries.FindCategory.Handle(ctx, catQuery)
	if categoryErr != nil {
		h.app.Logger.Error(ctx, "Failed to get category", categoryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(categoryErr))
	}

	if category == nil {
		catErr := Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid provided category with ID %s", newExpense.CategoryId),
		}
		return echoCtx.JSON(http.StatusBadRequest, catErr)
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
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get report http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling get report HTTP request")

	queryArgs := query.FindExpensesQuery{
		From: params.From,
		To:   params.To,
	}

	expenseRpt, expenseRptErr := h.app.Queries.FindExpenses.Handle(ctx, queryArgs)
	if expenseRptErr != nil {
		h.app.Logger.Error(ctx, "Failed to create expense report", expenseRptErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(expenseRptErr))
	}

	response := reportToResponse(*expenseRpt)
	return echoCtx.JSON(http.StatusOK, response)
}

func reportToResponse(domainReport domain.ReportByDate) ExpenseReport {
	dateCategoryReport := []DateCategoryReport{}
	for _, categoryByDate := range domainReport.CategoryByDate {
		someResponses := make([]CategoryReport, 0)
		for _, category := range categoryByDate.SubCategories {
			someResponses = append(someResponses, toResponse(*category))
		}

		dateCategoryReport = append(dateCategoryReport, DateCategoryReport{
			Date:            categoryByDate.Date,
			CategoryReports: someResponses,
			Total:           totalToTotalResponse(categoryByDate.Total),
		})
	}

	report := ExpenseReport{
		DateReports: dateCategoryReport,
		Total:       totalToTotalResponse(domainReport.Total),
	}
	return report
}

func totalToTotalResponse(d domain.Total) Total {
	return Total{
		Debug: d.SumDebug,
		Sum:   d.Sum.String(),
	}
}

func toResponse(d domain.CategoryExpenses) CategoryReport {
	response := CategoryReport{
		Category: categoryToResponse(d.Category),
		Total:    totalToTotalResponse(d.Total),
	}

	es := []Expense{}

	if d.Expenses != nil {
		for _, e := range *d.Expenses {
			es = append(es, expenseToResponse(e))
		}
		response.Expenses = &es
	}

	c := []CategoryReport{}
	for _, s := range d.SubCategories {
		c = append(c, toResponse(*s))
	}

	response.Children = &c

	return response
}

func categoryToResponse(domainCategory domain.Category) Category {
	// parents := []Category{}
	// p := domainCategory.Parents()
	// if p != nil {
	// 	for _, parentCategory := range *p {
	// 		parents = append(parents, categoryToResponse(parentCategory))
	// 	}
	// }

	category := Category{
		Id:    domainCategory.ID(),
		Name:  domainCategory.Name(),
		Icon:  domainCategory.Icon(),
		Level: domainCategory.Level(),
		// Parents: parents,
	}
	return category
}

func expenseToResponse(domainExpense domain.Expense) Expense {
	category := categoryToResponse(domainExpense.Category())
	expense := Expense{
		Id: domainExpense.ID(),
		NewExpense: NewExpense{
			// CategoryId: domainExpense.CategoryID(),
			Comment:  domainExpense.Comment(),
			Currency: domainExpense.Currency(),
			Date:     domainExpense.Date(),
			Price:    domainExpense.Price(),
			Quantity: domainExpense.Quantity(),
		},
		Category: category,
	}
	return expense
}
