package models

// Category struct represents a category.
type Category struct {
	ID   string `json:"id,omitempty" bson:"_id"`
	Name string `json:"name"`
	Path string `json:"path"`
}
