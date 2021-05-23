package ports

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
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

// FindCategories returns categories.
func (h HTTPServer) FindCategories(echoCtx echo.Context, params FindCategoriesParams) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get categories http request")
	defer span.End()

	filter := domain.CategoryFilter{}

	query, queryErr := h.app.Queries.FindCategories.Handle(ctx, filter)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to fetch categories", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	categoriesRes := categoriesToResponse(query)
	return echoCtx.JSON(http.StatusOK, categoriesRes)
}

// FindCategoryByID returns categories.
func (h HTTPServer) FindCategoryByID(echoCtx echo.Context, id string) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get category http request")
	span.SetAttributes(attribute.Any("id", id))
	defer span.End()

	query, queryErr := h.app.Queries.FindCategory.Handle(ctx, id)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to find category", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	if query == nil {
		catErr := Error{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Could not find category with ID %s", id),
		}
		return echoCtx.JSON(http.StatusNotFound, catErr)
	}

	categoryRes := categoryToResponse(query)
	return echoCtx.JSON(http.StatusOK, categoryRes)
}

func categoriesToResponse(domainCategories []domain.Category) []Category {
	categories := []Category{}
	for _, cat := range domainCategories {
		c := categoryToResponse(&cat)
		categories = append(categories, c)
	}
	return categories
}

func categoryToResponse(domainCategory *domain.Category) Category {
	var parents *[]Category
	if len(domainCategory.Parents()) != 0 {
		categoryParents := categoriesToResponse(domainCategory.Parents())
		parents = &categoryParents
	}
	category := Category{
		Id:       domainCategory.ID(),
		Name:     domainCategory.Name(),
		ParentId: domainCategory.ParentID(),
		Path:     domainCategory.Path(),
		Level:    domainCategory.Level(),
		Parents:  parents,
	}
	return category
}
