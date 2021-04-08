package entity

// CategoryFilter struct represents a category filter.
type CategoryFilter struct {
	CategoryID   string
	ParentID     string
	CategoryIDs  []string
	Path         string
	FindChildren bool
	FindAll      bool
}
