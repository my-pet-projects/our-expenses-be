package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

// CreateCategoryRequest provides the schema definition for create category API request body.
type CreateCategoryRequest struct {
	Name     string              `json:"name" validate:"required"`
	ParentID *primitive.ObjectID `json:"parentId"`
	Path     string              `json:"path" validate:"required"`
	Level    int32               `json:"level" validate:"required,gt=0"`
}

// UpdateCategoryRequest provides the schema definition for update category API request body.
type UpdateCategoryRequest struct {
	CreateCategoryRequest
}
