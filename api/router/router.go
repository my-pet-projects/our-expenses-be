package router

import (
	"net/http"
	"our-expenses-server/api/handler"

	"github.com/gorilla/mux"
)

// Router defines available application controllers.
type Router struct {
	category handler.CategoryControllerInterface
}

// InitializeRoutes returns HTTP handler with defined application routes.
func (router *Router) InitializeRoutes() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/categories", router.category.Create).Methods(http.MethodPost)
	apiRouter.HandleFunc("/categories", router.category.GetAll).Methods(http.MethodGet)
	apiRouter.HandleFunc("/categories/{id}", router.category.GetOne).Methods(http.MethodGet)
	apiRouter.HandleFunc("/categories/{id}", router.category.Update).Methods(http.MethodPut)
	apiRouter.HandleFunc("/categories/{id}", router.category.Delete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/categories/{id}/usages", router.category.GetUsages).Methods(http.MethodGet)
	apiRouter.HandleFunc("/categories/{id}/move", router.category.Move).Methods(http.MethodPut)

	return apiRouter
}

// ProvideRouter returns Router.
func ProvideRouter(category *handler.CategoryController) *Router {
	return &Router{
		category: category,
	}
}
