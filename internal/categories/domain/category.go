package domain

import (
	"time"
)

// Category represents a domain object.
type Category struct {
	id        string
	name      string
	parentID  string
	path      string
	level     int
	parents   []Category
	createdAt time.Time
	updatedAt time.Time
}

// NewCategory creates a new category domain object.
func NewCategory(id string, name string, parentID string, path string, level int) (*Category, error) {

	return &Category{
		id:        id,
		name:      name,
		parentID:  parentID,
		path:      path,
		level:     level,
		createdAt: time.Now(),
	}, nil
}

// ID returns category id.
func (c Category) ID() string {
	return c.id
}

// Name returns category name.
func (c Category) Name() string {
	return c.name
}

// ParentID returns category parent id.
func (c Category) ParentID() string {
	return c.parentID
}

// Path returns category path.
func (c Category) Path() string {
	return c.path
}

// Level returns category level.
func (c Category) Level() int {
	return c.level
}

// Parents returns category parents.
func (c Category) Parents() []Category {
	return c.parents
}
