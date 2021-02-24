package api

import (
	"net/http"
	"our-expenses-server/web/api/controllers"

	"github.com/gorilla/mux"
)

// Router defines available application controllers.
type Router struct {
	categoryCtrl controllers.CategoryControllerInterface
}

// InitializeRoutes returns HTTP handler with defined application routes.
func (router *Router) InitializeRoutes() *mux.Router {
	apiRouter := mux.NewRouter().StrictSlash(true).PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/categories", router.categoryCtrl.CreateCategory).Methods(http.MethodPost)
	apiRouter.HandleFunc("/categories", router.categoryCtrl.GetAllCategories).Methods(http.MethodGet)
	apiRouter.HandleFunc("/categories/{id}", router.categoryCtrl.GetCategory).Methods(http.MethodGet)
	apiRouter.HandleFunc("/categories/{id}", router.categoryCtrl.UpdateCategory).Methods(http.MethodPut)
	apiRouter.HandleFunc("/categories/{id}", router.categoryCtrl.DeleteCategory).Methods(http.MethodDelete)

	return apiRouter
}

// ProvideRouter returns Router.
func ProvideRouter(categoryCtrl *controllers.CategoryController) *Router {
	return &Router{categoryCtrl: categoryCtrl}
}
