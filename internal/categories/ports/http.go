package ports

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
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
	h.app.Logger.Info(ctx, "Handling get categories HTTP request")

	query := query.FindCategoriesQuery{
		ParentID:        params.ParentId,
		FindAllChildren: false,
	}
	if params.AllChildren != nil {
		query.FindAllChildren = *params.AllChildren
	}
	if params.All != nil {
		query.FindAll = *params.All
	}
	queryRes, queryErr := h.app.Queries.FindCategories.Handle(ctx, query)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to fetch categories", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	response := categoriesToResponse(queryRes)
	return echoCtx.JSON(http.StatusOK, response)
}

// FindCategoryByID returns categories.
func (h HTTPServer) FindCategoryByID(echoCtx echo.Context, id string) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get category http request")
	span.SetAttributes(attribute.Any("id", id))
	defer span.End()
	h.app.Logger.Infof(ctx, "Handling get %s category HTTP request", id)

	query := query.FindCategoryQuery{
		CategoryID: id,
	}
	queryRes, queryErr := h.app.Queries.FindCategory.Handle(ctx, query)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to find category", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	if queryRes == nil {
		catErr := Error{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Could not find category with ID %s", id),
		}
		return echoCtx.JSON(http.StatusNotFound, catErr)
	}

	categoryRes := categoryToResponse(queryRes)
	return echoCtx.JSON(http.StatusOK, categoryRes)
}

// AddCategory adds a new category.
func (h HTTPServer) AddCategory(echoCtx echo.Context) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get category http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling add categories HTTP request")

	var newCategory Category
	bindErr := echoCtx.Bind(&newCategory)
	if bindErr != nil {
		catErr := Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid format for a new category",
		}
		h.app.Logger.Error(ctx, "Invalid category format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest, catErr)
	}

	cmdArgs := command.AddCategoryCommand{
		ParentID: newCategory.ParentId,
		Name:     newCategory.Name,
		Path:     newCategory.Path,
		Level:    newCategory.Level,
	}
	categoryID, categoryCrtErr := h.app.Commands.AddCategory.Handle(ctx, cmdArgs)
	if categoryCrtErr != nil {
		h.app.Logger.Error(ctx, "Failed to create category", categoryCrtErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(categoryCrtErr))
	}

	return echoCtx.JSON(http.StatusCreated, categoryID)
}

// UpdateCategory updates a category.
func (h HTTPServer) UpdateCategory(echoCtx echo.Context, id string) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle update category http request")
	defer span.End()
	h.app.Logger.Info(ctx, "Handling update category HTTP request")

	var category Category
	bindErr := echoCtx.Bind(&category)
	if bindErr != nil {
		catErr := Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid format for a category to update",
		}
		h.app.Logger.Error(ctx, "Invalid category format", bindErr)
		return echoCtx.JSON(http.StatusBadRequest, catErr)
	}

	cmd := command.UpdateCategoryCommand{
		ID:       category.Id,
		ParentID: category.ParentId,
		Name:     category.Name,
		Path:     category.Path,
		Level:    category.Level,
	}
	_, categoryUpdErr := h.app.Commands.UpdateCategory.Handle(ctx, cmd)
	if categoryUpdErr != nil {
		h.app.Logger.Error(ctx, "Failed to update category", categoryUpdErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(categoryUpdErr))
	}

	return echoCtx.NoContent(http.StatusOK)
}

// DeleteCategory deletes a category.
func (h HTTPServer) DeleteCategory(echoCtx echo.Context, id string) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle delete category http request")
	span.SetAttributes(attribute.Any("id", id))
	defer span.End()
	h.app.Logger.Info(ctx, "Handling delete categories HTTP request")

	cmd := command.DeleteCategoryCommand{
		CategoryID: id,
	}
	cmdRes, cmdErr := h.app.Commands.DeleteCategory.Handle(ctx, cmd)
	if cmdErr != nil {
		h.app.Logger.Error(ctx, "Failed to delete category", cmdErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(cmdErr))
	}

	if cmdRes == nil {
		catErr := Error{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Could not find category with ID %s", id),
		}
		return echoCtx.JSON(http.StatusNotFound, catErr)
	}

	return echoCtx.NoContent(http.StatusNoContent)
}

// FindCategoryUsages returns categories.
func (h HTTPServer) FindCategoryUsages(echoCtx echo.Context, id string) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get category usages http request")
	span.SetAttributes(attribute.Any("id", id))
	defer span.End()
	h.app.Logger.Info(ctx, "Handling get category usages HTTP request")

	query := query.FindCategoryUsagesQuery{
		CategoryID: id,
	}
	queryRes, queryErr := h.app.Queries.FindCategoryUsages.Handle(ctx, query)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to find category usages", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	response := categoriesToResponse(queryRes)
	return echoCtx.JSON(http.StatusOK, response)
}

// MoveCategory moves category.
func (h HTTPServer) MoveCategory(echoCtx echo.Context, id string, params MoveCategoryParams) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle move category http request")
	span.SetAttributes(attribute.Any("id", id))
	span.SetAttributes(attribute.Any("destinationId", params.DestinationId))
	defer span.End()
	h.app.Logger.Info(ctx,
		fmt.Sprintf("Handling HTTP request to move category %s to %s", id, params.DestinationId))

	cmd := command.MoveCategoryCommand{
		CategoryID:    id,
		DestinationID: params.DestinationId,
	}
	cmdResult, cmdErr := h.app.Commands.MoveCategory.Handle(ctx, cmd)
	if cmdErr != nil {
		h.app.Logger.Error(ctx, "Failed to move category", cmdErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(cmdErr))
	}

	return echoCtx.JSON(http.StatusOK, cmdResult)
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
		Id: domainCategory.ID(),
		NewCategory: NewCategory{
			Name:     domainCategory.Name(),
			ParentId: domainCategory.ParentID(),
			Path:     domainCategory.Path(),
			Level:    domainCategory.Level(),
			Parents:  parents,
		},
	}
	return category
}
