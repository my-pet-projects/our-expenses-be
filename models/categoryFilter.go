package models

// CategoryFilter struct represents a category filter.
type CategoryFilter struct {
	Path   string `bson:"path,omitempty"`
	Parent string `bson:"parent,omitempty"`
}
