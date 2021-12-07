package domain

import (
	"fmt"
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
	icon      *string
	level     int
	parents   []Category
	createdBy string
	createdAt time.Time
	updatedBy *string
	updatedAt *time.Time
}

// NewCategory creates a new category domain object.
func NewCategory(id string, name string, parentID *string, path string, icon *string, level int) (*Category, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	if icon != nil {
		trimmedIcon := strings.TrimSpace(*icon)
		icon = &trimmedIcon
	}

	return &Category{
		id:       id,
		name:     strings.TrimSpace(name),
		parentID: parentID,
		path:     strings.TrimSpace(path),
		icon:     icon,
		level:    level,
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

// Icon returns category icon.
func (c Category) Icon() *string {
	return c.icon
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

	parentIDs := make([]string, 0)
	for _, parentID := range strings.Split(c.path, "|") {
		if parentID == c.id || len(parentID) == 0 {
			continue
		}
		parentIDs = append(parentIDs, parentID)
	}
	return parentIDs
}

// CreatedBy returns category creator.
func (c Category) CreatedBy() string {
	return c.createdBy
}

// UpdatedBy returns category updater.
func (c Category) UpdatedBy() *string {
	return c.updatedBy
}

// CreatedAt returns category creation date.
func (c Category) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns category update date.
func (c Category) UpdatedAt() *time.Time {
	return c.updatedAt
}

// SetMetadata sets category metadata.
func (c *Category) SetMetadata(createdBy string, createdAt time.Time, updatedBy *string, updatedAt *time.Time) {
	c.createdBy = createdBy
	c.createdAt = createdAt
	c.updatedBy = updatedBy
	c.updatedAt = updatedAt
}

// SetParents sets category parents.
func (c *Category) SetParents(parents []Category) {
	c.parents = parents
}

// AssignNewParent assigns a new parent.
func (c *Category) AssignNewParent(parent *Category) {
	var newParentID *string
	var newLevel int
	var newPath string

	if parent == nil {
		newParentID = nil
		newLevel = 1
		newPath = fmt.Sprintf("|%s", c.id)
	} else {
		newParentID = &parent.id
		newLevel = parent.level + 1
		newPath = fmt.Sprintf("%s|%s", parent.path, c.id)
	}

	c.parentID = newParentID
	c.path = newPath
	c.level = newLevel
}

// ReplaceAncestor replaces all ancestors.
func (c *Category) ReplaceAncestor(oldPath string, newPath string) {
	path := strings.ReplaceAll(c.path, oldPath, newPath)
	level := strings.Count(path, "|")
	c.path = path
	c.level = level
}
