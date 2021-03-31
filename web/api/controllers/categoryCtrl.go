package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"our-expenses-server/db/repositories"
	"our-expenses-server/logger"
	"our-expenses-server/models"
	"our-expenses-server/validators"
	"our-expenses-server/web/requests"
	"our-expenses-server/web/responses"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	GetCategory(w http.ResponseWriter, req *http.Request)
	UpdateCategory(w http.ResponseWriter, req *http.Request)
	DeleteCategory(w http.ResponseWriter, req *http.Request)
	GetCategoryUsages(w http.ResponseWriter, req *http.Request)
	MoveCategory(w http.ResponseWriter, req *http.Request)
}

// ProvideCategoryController returns a CategoryController.
func ProvideCategoryController(repo *repositories.CategoryRepository, logger *logger.AppLogger, validator *validators.Validator) *CategoryController {
	return &CategoryController{repo: repo, logger: logger, validator: validator}
}

const loggerCategory = "api/categories"

// GetAllCategories returns a list of all categories.
func (ctrl *CategoryController) GetAllCategories(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	parentIDParam := req.URL.Query().Get("parentId")
	allParam := req.URL.Query().Get("all")
	loggerTags := logger.Fields{loggerCategory: "getAll", "query": req.URL.Query()}
	ctrl.logger.Info("Http request", loggerTags)

	isAll, isAllError := strconv.ParseBool(allParam)
	if isAllError != nil {
		isAll = false
	}

	filter := models.CategoryFilter{
		ParentID: parentIDParam,
		FindAll:  isAll,
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

// GetCategory returns a single category found by id.
func (ctrl *CategoryController) GetCategory(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	loggerTags := logger.Fields{loggerCategory: "get", "query": req.URL.Query(), "routeVars": vars}
	ctrl.logger.Info("Http request", loggerTags)

	category, categoryError := ctrl.repo.GetOne(ctx, categoryID)
	if categoryError != nil {
		ctrl.logger.Error("Failed to get a category from the database", categoryError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rawParentCategories := strings.Split(category.Path, "|")
	var parentCategoryIDs []string
	for _, str := range rawParentCategories {
		if str != "" && str != category.ID.Hex() {
			parentCategoryIDs = append(parentCategoryIDs, str)
		}
	}

	if len(parentCategoryIDs) != 0 {
		parentCategoriesFilter := models.CategoryFilter{
			CategoryIDs: parentCategoryIDs,
		}

		parentCategories, parentCategoriesError := ctrl.repo.GetAll(ctx, parentCategoriesFilter)
		if parentCategoriesError != nil {
			ctrl.logger.Error("Failed to get parent categories from the database", parentCategoriesError, loggerTags)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		category.ParentCategories = parentCategories
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category found by id.
func (ctrl *CategoryController) DeleteCategory(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	loggerTags := logger.Fields{loggerCategory: "delete", "query": req.URL.Query(), "routeVars": vars}
	ctrl.logger.Info("Http request", loggerTags)

	category, categoryError := ctrl.repo.GetOne(ctx, categoryID)
	if categoryError != nil {
		ctrl.logger.Error("Failed to get a category from the database for deletion", categoryError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	categoryFilter := models.CategoryFilter{
		Path:         category.Path,
		FindChildren: true,
	}

	deleteResult, deleteError := ctrl.repo.DeleteAll(ctx, categoryFilter)
	if deleteError != nil {
		ctrl.logger.Error("Failed to delete a category from the database", deleteError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deleteResult == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deleteResult)
}

// CreateCategory inserts a category into the database.
func (ctrl *CategoryController) CreateCategory(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	request := &requests.CreateCategoryRequest{}
	loggerTags := logger.Fields{loggerCategory: "create"}

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

	categoryID := primitive.NewObjectID()
	category := &models.Category{
		ID:       &categoryID,
		Name:     request.Name,
		Path:     fmt.Sprintf("%s|%s", request.Path, categoryID.Hex()),
		ParentID: request.ParentID,
		Level:    request.Level,
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

// UpdateCategory updates a category in the database.
func (ctrl *CategoryController) UpdateCategory(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	request := &requests.UpdateCategoryRequest{}
	loggerTags := logger.Fields{loggerCategory: "update"}

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

	objID, _ := primitive.ObjectIDFromHex(categoryID)

	category := &models.Category{
		ID:       &objID,
		Name:     request.Name,
		Path:     request.Path,
		ParentID: request.ParentID,
		Level:    request.Level,
	}

	_, updateError := ctrl.repo.Update(ctx, category)
	if updateError != nil {
		ctrl.logger.Error("Failed to update a category", updateError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCategoryUsages returns category usages.
func (ctrl *CategoryController) GetCategoryUsages(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	loggerTags := logger.Fields{loggerCategory: "usages", "query": req.URL.Query()}
	ctrl.logger.Info("Http request", loggerTags)

	pathFilter := models.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}

	categoryUsages, categoryUsagesError := ctrl.repo.GetAll(ctx, pathFilter)
	if categoryUsagesError != nil {
		ctrl.logger.Error("Failed to get category usages from the database", categoryUsagesError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categoryUsages)
}

// MoveCategory moves category.
func (ctrl *CategoryController) MoveCategory(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	destinationID := req.URL.Query().Get("destinationId")
	loggerTags := logger.Fields{loggerCategory: "move", "query": req.URL.Query()}
	ctrl.logger.Info("Http request", loggerTags)

	targetCategory, targetCategoryError := ctrl.repo.GetOne(ctx, categoryID)
	if targetCategoryError != nil {
		ctrl.logger.Error("Failed to get destination category from the database for move", targetCategoryError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if targetCategory == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(destinationID != "root")

	var destinationCategory *models.Category
	if destinationID != "root" {
		var destinationCategoryError error
		destinationCategory, destinationCategoryError = ctrl.repo.GetOne(ctx, destinationID)
		if destinationCategoryError != nil {
			ctrl.logger.Error("Failed to get destination category from the database for move", destinationCategoryError, loggerTags)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if destinationCategory == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	pathFilter := models.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}

	categoryUsages, categoryUsagesError := ctrl.repo.GetAll(ctx, pathFilter)
	if categoryUsagesError != nil {
		ctrl.logger.Error("Failed to get category usages from the database", categoryUsagesError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var targetPath string
	var destinationPath string
	var newPath string
	var levelDiff int32
	var parentID *primitive.ObjectID

	parentID = nil
	destinationPath = ""
	levelDiff = -targetCategory.Level

	if destinationCategory != nil {
		fmt.Println("destinationCategory not nil")
		parentID = destinationCategory.ID
		destinationPath = destinationCategory.Path
		levelDiff = destinationCategory.Level - targetCategory.Level
	}

	targetPath = targetCategory.Path
	newPath = strings.Replace(targetCategory.Path, targetPath, destinationPath, -1) + "|" + targetCategory.ID.Hex()

	targetCategory.Path = newPath
	targetCategory.ParentID = parentID
	targetCategory.Level = targetCategory.Level + levelDiff + 1
	_, updateError := ctrl.repo.Update(ctx, targetCategory)
	if updateError != nil {
		ctrl.logger.Error("Failed to update a category from the database for move", updateError, loggerTags)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, category := range categoryUsages {
		newChildPath := strings.Replace(category.Path, targetPath, newPath, -1)
		newChildLevel := category.Level + levelDiff + 1

		category.Path = newChildPath
		category.Level = newChildLevel
		_, updateError := ctrl.repo.Update(ctx, &category)
		if updateError != nil {
			ctrl.logger.Error("Failed to update a category from the database for move", updateError, loggerTags)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(0)
}
