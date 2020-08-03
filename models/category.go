package models

// Category struct represents a category.
type Category struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name"`
	Parent string `json:"parent"`
	Path   string `json:"path"`
}
