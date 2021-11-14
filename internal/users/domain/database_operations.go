package domain

const SystemUser string = "system"

// InsertResult represents a struct with new entity ID.
type InsertResult struct {
	ID string
}

// UpdateResult represents a struct with update operation result details.
type UpdateResult struct {
	UpdateCount int
}
