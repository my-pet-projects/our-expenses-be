package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"our-expenses-server/api/httperr"
	"our-expenses-server/api/presenter"
	"our-expenses-server/entity"
	"our-expenses-server/logger"
	"our-expenses-server/service/category"
	"our-expenses-server/validator"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const loggerCategory = "api/category"

// CategoryController defines a category API endpoint.
type CategoryController struct {
	service   category.CategoryServiceInterface
	logger    logger.AppLoggerInterface
	validator validator.ValidatorInterface
}

// CategoryControllerInterface defines a contract to expose category API.
type CategoryControllerInterface interface {
	GetAll(w http.ResponseWriter, req *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
	GetOne(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
	GetUsages(w http.ResponseWriter, req *http.Request)
	Move(w http.ResponseWriter, req *http.Request)
}

// ProvideCategoryController returns a CategoryController.
func ProvideCategoryController(service *category.CategoryService, logger *logger.AppLogger, validator *validator.Validator) *CategoryController {
	return &CategoryController{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}

// GetAllCategories returns a list of all categories.
func (ctrl *CategoryController) GetAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	parentIDParam := req.URL.Query().Get("parentId")
	allParam := req.URL.Query().Get("all")
	ctrl.logger.Info(ctx, "Handling HTTP request to get all categories")

	isAll, isAllError := strconv.ParseBool(allParam)
	if isAllError != nil {
		isAll = false
	}

	filter := entity.CategoryFilter{
		ParentID: parentIDParam,
		FindAll:  isAll,
	}

	categories, categoriesErr := ctrl.service.GetAll(ctx, filter)
	if categoriesErr != nil {
		ctrl.logger.Error(ctx, "Failed to fetch categories from database", categoriesErr)
		httperr.RespondWithError(categoriesErr, w, req)
		return
	}

	var toJson []*presenter.Category
	for _, c := range categories {
		toJson = append(toJson, &presenter.Category{
			ID:        c.ID.Hex(),
			Name:      c.Name,
			ParentID:  c.ParentID.Hex(),
			Path:      c.Path,
			Level:     c.Level,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toJson)
	return
}

// GetCategory returns a single category found by id.
func (ctrl *CategoryController) GetOne(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	// loggerTags := logger.FieldsSet{loggerCategory: "get", "query": req.URL.Query(), "routeVars": vars}
	ctrl.logger.Info(ctx, "Http request")

	category, categoryError := ctrl.service.GetOne(ctx, categoryID)
	if categoryError != nil {
		ctrl.logger.Error(ctx, "Failed to get a category from the database", categoryError)
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

	var toJson = presenter.Category{
		ID:        category.ID.Hex(),
		Name:      category.Name,
		ParentID:  category.ParentID.Hex(),
		Path:      category.Path,
		Level:     category.Level,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	if len(parentCategoryIDs) != 0 {
		parentCategoriesFilter := entity.CategoryFilter{
			CategoryIDs: parentCategoryIDs,
		}

		parentCategories, parentCategoriesError := ctrl.service.GetAll(ctx, parentCategoriesFilter)
		if parentCategoriesError != nil {
			ctrl.logger.Error(ctx, "Failed to get parent categories from the database", parentCategoriesError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, c := range parentCategories {
			toJson.ParentCategories = append(toJson.ParentCategories, presenter.Category{
				ID:        c.ID.Hex(),
				Name:      c.Name,
				ParentID:  c.ParentID.Hex(),
				Path:      c.Path,
				Level:     c.Level,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(toJson)
}

// CreateCategory inserts a category into the database.
func (ctrl *CategoryController) Create(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctrl.logger.Info(ctx, "Handling create category request")

	var payload struct {
		Name     string `json:"name" validate:"required"`
		ParentID string `json:"parentId"`
		Path     string `json:"path"` // TODO: maybe require path, but not for the 1st level categories
		Level    int32  `json:"level" validate:"required,gt=0"`
	}

	decodeErr := json.NewDecoder(req.Body).Decode(&payload)
	if decodeErr != nil {
		// ctrl.logger.Error("Error while decoding request body", decodeErr, loggerTags)
		httperr.RespondWithBadRequest("Invalid resquest payload", w, req)
		return
	}
	defer req.Body.Close()

	validatorErr := ctrl.validator.ValidateStruct(payload)
	if validatorErr != nil {
		// httperr.RespondWithBadRequest(validatorErr, w, req)
		return
	}

	id, saveErr := ctrl.service.Create(ctx, payload.Name, payload.ParentID, payload.Path, int(payload.Level))
	if saveErr != nil {
		fmt.Printf("\n\nFATAL: %+v \n\n", saveErr)

		ctrl.logger.Error(ctx, "-------> Failed to save category", saveErr)

		var perr *entity.AppError
		if errors.As(saveErr, &perr) {

			// ctrl.logger.Error(ctx, "-------> errors.As", perr)

			// w.WriteHeader(http.StatusBadRequest)

			// json.NewEncoder(w).Encode(responses.ErrorResponse{Message: perr.Msg})

		}

		// // switch err := errors.Cause(err).(type) {
		// if originalErr, ok := errors.Cause(twiceWrappedError).(*CustomError); ok {
		// switch err := errors.Cause(saveErr).(type) {
		//

		switch err := errors.Cause(saveErr).(type) {
		case *entity.AppError:
			// ctrl.logger.Error(ctx, "-------> *entity.ErrZeroDivision", err)

			httperr.RespondWithBadRequest(err.Msg, w, req)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// if errVal, ok := saveErr.(*entity.ErrInvalidEntity); ok {
		// 	fmt.Println(errVal)
		// 	fmt.Println(errVal)
		// 	fmt.Println(errVal)
		// }

		// // ctrl.logger.Error("Failed to insert a category", saveErr, loggerTags)
		// w.WriteHeader(http.StatusInternalServerError)
		// return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

// UpdateCategory updates a category in the database.
func (ctrl *CategoryController) Update(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// vars := mux.Vars(req)
	// categoryID := vars["id"]

	var payload struct {
		Name     string `json:"name" validate:"required"`
		ParentID string `json:"parentId"`
		Path     string `json:"path"` // TODO: maybe require path, but not for the 1st level categories
		Level    int32  `json:"level" validate:"required,gt=0"`
	}

	decodeError := json.NewDecoder(req.Body).Decode(payload)
	if decodeError != nil {
		ctrl.logger.Error(ctx, "Error while decoding request body", decodeError)
		httperr.RespondWithBadRequest("Invalid resquest payload", w, req)
		return
	}
	defer req.Body.Close()

	validationError := ctrl.validator.ValidateStruct(payload)
	if validationError != nil {
		// httperr.RespondWithBadRequest(validationError, w, req)
		return
	}

	// objID, _ := primitive.ObjectIDFromHex(categoryID)

	category := &entity.Category{
		// ID:       &objID,
		// Name:     request.Name,
		// Path:     request.Path,
		// ParentID: request.ParentID,
		// Level:    request.Level,
	}

	_, updateError := ctrl.service.Update(ctx, category)
	if updateError != nil {
		ctrl.logger.Error(ctx, "Failed to update a category", updateError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCategory deletes a category found by id.
func (ctrl *CategoryController) Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	ctrl.logger.Info(ctx, "Http request")

	category, categoryError := ctrl.service.GetOne(ctx, categoryID)
	if categoryError != nil {
		ctrl.logger.Error(ctx, "Failed to get a category from the database for deletion", categoryError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if category == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	categoryFilter := entity.CategoryFilter{
		Path:         category.Path,
		FindChildren: true,
	}

	deleteResult, deleteError := ctrl.service.DeleteAll(ctx, categoryFilter)
	if deleteError != nil {
		ctrl.logger.Error(ctx, "Failed to delete a category from the database", deleteError)
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

// GetCategoryUsages returns category usages.
func (ctrl *CategoryController) GetUsages(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	ctrl.logger.Info(ctx, "Http request")

	pathFilter := entity.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}

	categoryUsages, categoryUsagesError := ctrl.service.GetAll(ctx, pathFilter)
	if categoryUsagesError != nil {
		ctrl.logger.Error(ctx, "Failed to get category usages from the database", categoryUsagesError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categoryUsages)
}

// MoveCategory moves category.
func (ctrl *CategoryController) Move(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	categoryID := vars["id"]
	destinationID := req.URL.Query().Get("destinationId")
	ctrl.logger.Info(ctx, "Http request")

	targetCategory, targetCategoryError := ctrl.service.GetOne(ctx, categoryID)
	if targetCategoryError != nil {
		ctrl.logger.Error(ctx, "Failed to get destination category from the database for move", targetCategoryError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if targetCategory == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println(destinationID != "root")

	var destinationCategory *entity.Category
	if destinationID != "root" {
		var destinationCategoryError error
		destinationCategory, destinationCategoryError = ctrl.service.GetOne(ctx, destinationID)
		if destinationCategoryError != nil {
			ctrl.logger.Error(ctx, "Failed to get destination category from the database for move", destinationCategoryError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if destinationCategory == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	pathFilter := entity.CategoryFilter{
		CategoryID:   categoryID,
		FindChildren: true,
	}

	categoryUsages, categoryUsagesError := ctrl.service.GetAll(ctx, pathFilter)
	if categoryUsagesError != nil {
		ctrl.logger.Error(ctx, "Failed to get category usages from the database", categoryUsagesError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var targetPath string
	var destinationPath string
	var newPath string
	var levelDiff int
	var parentID *primitive.ObjectID

	parentID = nil
	destinationPath = ""
	levelDiff = -targetCategory.Level

	if destinationCategory != nil {
		fmt.Println("destinationCategory not nil")
		parentID = &destinationCategory.ID
		destinationPath = destinationCategory.Path
		levelDiff = destinationCategory.Level - targetCategory.Level
	}

	targetPath = targetCategory.Path
	newPath = strings.Replace(targetCategory.Path, targetPath, destinationPath, -1) + "|" + targetCategory.ID.Hex()

	targetCategory.Path = newPath
	targetCategory.ParentID = *parentID
	targetCategory.Level = targetCategory.Level + levelDiff + 1
	_, updateError := ctrl.service.Update(ctx, targetCategory)
	if updateError != nil {
		ctrl.logger.Error(ctx, "Failed to update a category from the database for move", updateError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, category := range categoryUsages {
		newChildPath := strings.Replace(category.Path, targetPath, newPath, -1)
		newChildLevel := category.Level + levelDiff + 1

		category.Path = newChildPath
		category.Level = newChildLevel
		_, updateError := ctrl.service.Update(ctx, &category)
		if updateError != nil {
			ctrl.logger.Error(ctx, "Failed to update a category from the database for move", updateError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(0)
}
