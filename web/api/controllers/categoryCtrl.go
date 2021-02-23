package controllers

import (
	"encoding/json"
	"net/http"
	"our-expenses-server/db/repositories"
	"our-expenses-server/logger"
	"our-expenses-server/models"
	"our-expenses-server/validators"
	"our-expenses-server/web/requests"
	"our-expenses-server/web/responses"
)

// CategoryController defines a category API endpoint.
type CategoryController struct {
	repo      repositories.CategoryRepoInterface
	logger    logger.AppLoggerInterface
	validator validators.ValidatorInterface
}

// CategoryControllerInterface defines a contract to expose category API.
type CategoryControllerInterface interface {
	CreateCategory(w http.ResponseWriter, req *http.Request)
	GetAllCategories(w http.ResponseWriter, req *http.Request)
}

// ProvideCategoryController returns a CategoryController.
func ProvideCategoryController(repo *repositories.CategoryRepository, logger *logger.AppLogger, validator *validators.Validator) *CategoryController {
	return &CategoryController{repo: repo, logger: logger, validator: validator}
}

// GetAllCategories returns a list of all categories to the response.
func (ctrl *CategoryController) GetAllCategories(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	parentIDParam := req.URL.Query().Get("parentId")
	loggerTags := logger.Fields{"api/categories": "getAll", "query": req.URL.Query()}
	ctrl.logger.Info("Http request", loggerTags)

	filter := models.CategoryFilter{
		ParentID: parentIDParam,
	}

	categories, categoriesError := ctrl.repo.GetAll(ctx, filter)
	if categoriesError != nil {
		ctrl.logger.Error("Failed to load categories from the database", categoriesError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

// CreateCategory inserts a category into the database.
func (ctrl *CategoryController) CreateCategory(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := req.Context()
	request := &requests.CreateCategoryRequest{}
	loggerTags := logger.Fields{"api/categories": "create"}

	decodeError := json.NewDecoder(req.Body).Decode(request)
	if decodeError != nil {
		ctrl.logger.Error("Error while decoding request body", decodeError, loggerTags)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.ErrorResponse{Message: "Invalid resquest payload"})
		return
	}
	defer req.Body.Close()

	validationError := ctrl.validator.ValidateStruct(request)
	if validationError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responses.ErrorResponse{Message: validationError})
		return
	}

	category := &models.Category{
		Name: request.Name,
		Path: request.Path,
	}

	savedCategory, saveError := ctrl.repo.Insert(ctx, category)
	if saveError != nil {
		ctrl.logger.Error("Failed to insert a category", saveError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedCategory)
}
