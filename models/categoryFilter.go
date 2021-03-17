package models

// CategoryFilter struct represents a category filter.
type CategoryFilter struct {
	CategoryID   string
	ParentID     string
	CategoryIDs  []string
	FindChildren bool
}
