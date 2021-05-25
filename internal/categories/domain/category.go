package domain

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Category represents a domain object.
type Category struct {
	id        string
	name      string
	parentID  *string
	path      string
	level     int
	parents   []Category
	createdAt time.Time
	updatedAt *time.Time
}

// NewCategory creates a new category domain object.
func NewCategory(id string, name string, parentID *string, path string, level int,
	createdAt time.Time, updatedAt *time.Time) (*Category, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	return &Category{
		id:        id,
		name:      name,
		parentID:  parentID,
		path:      path,
		level:     level,
		createdAt: createdAt,
		updatedAt: updatedAt,
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
func (c Category) ParentID() *string {
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

// ParentIDs returns parent IDs.
func (c Category) ParentIDs() []string {
	if len(c.path) == 0 {
		return []string{}
	}

	var parentIDs []string
	for _, parentID := range strings.Split(c.path, "|") {
		if parentID == c.id || len(parentID) == 0 {
			continue
		}
		parentIDs = append(parentIDs, parentID)
	}
	return parentIDs
}

// CreatedAt returns category creation date.
func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns category update date.
func (c Category) UpdatedAt() *time.Time {
	return c.updatedAt
}

// SetParents sets category parents.
func (c *Category) SetParents(parents []Category) {
	c.parents = parents
}
