package ports

import (
	"net/http"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/server/httperr"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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

// GetDeviceLocations returns device locations.
func (h HTTPServer) GetCategories(echoCtx echo.Context) error {
	ctx, span := tracer.Start(echoCtx.Request().Context(), "handle get device locations http request")
	defer span.End()

	query, queryErr := h.app.Queries.FindCategories.Handle(ctx)
	if queryErr != nil {
		h.app.Logger.Error(ctx, "Failed to fetch categories", queryErr)
		return echoCtx.JSON(http.StatusInternalServerError, httperr.InternalError(queryErr))
	}

	categoriesRes := categoriesToResponse(query)
	return echoCtx.JSON(http.StatusOK, categoriesRes)
}

func categoriesToResponse(domainCategories []domain.Category) []Category {
	categories := []Category{}
	for _, cat := range domainCategories {
		c := Category{
			Id:       cat.ID(),
			Name:     cat.Name(),
			ParentId: cat.ParentID(),
			Path:     cat.Path(),
			Level:    cat.Level(),
			// Parents:  cat.Parents(),
		}
		categories = append(categories, c)
	}
	return categories
}
