package domain

import (
	"fmt"
	"strings"
	"time"
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
	updatedAt time.Time
}

// NewCategory creates a new category domain object.
func NewCategory(id string, name string, parentID *string, path string, level int) (*Category, error) {

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
	fmt.Printf("\n\n %+v \n\n", parentIDs)
	fmt.Printf("\n\n %+v \n\n", strings.Split(c.path, "|"))
	return parentIDs
}

// SetParents sets category parents.
func (c *Category) SetParents(parents []Category) {
	c.parents = parents
}
