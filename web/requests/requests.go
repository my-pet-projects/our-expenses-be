package requests

// CreateCategoryRequest provides the schema definition for create category API request body.
type CreateCategoryRequest struct {
	Name   string `json:"name" validate:"required"`
	Parent string `json:"parent"`
	Path   string `json:"path" validate:"required"`
}

// UpdateCategoryRequest provides the schema definition for update category API request body.
type UpdateCategoryRequest struct {
	ID string `json:"id" validate:"required,len=24"`
	CreateCategoryRequest
}
