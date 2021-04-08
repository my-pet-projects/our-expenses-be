package entity

import (
	"time"

	"github.com/pkg/errors"
)

// Category data
type Category struct {
	ID               ID         `bson:"_id,omitempty"`
	Name             string     `bson:"name"`
	ParentID         ID         `bson:"parentId,omitempty"`
	Path             string     `bson:"path"`
	Level            int        `bson:"level"`
	ParentCategories []Category `bson:"-"`
	CreatedAt        time.Time  `bson:"createdAt"`
	UpdatedAt        time.Time  `bson:"updatedAt,omitempty"`
}

// ID               *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// Name             string              `json:"name"`
// ParentID         *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`
// Path             string              `json:"path"`
// Level            int32               `json:"level"`
// ParentCategories []Category          `json:"parents,omitempty" bson:"-"`

// NewCategory creates a new category
func NewCategory(name string, parentID string, path string, level int) (*Category, error) {
	c := &Category{
		ID:        NewID(),
		Name:      name,
		ParentID:  IDFromString(parentID),
		Path:      path,
		Level:     level,
		CreatedAt: time.Now(),
	}
	err := c.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "wrap 2")
	}
	return c, nil
}

//Validate validate book
func (b *Category) Validate() error {
	// if b.Title == "" || b.Author == "" || b.Pages <= 0 || b.Quantity <= 0 {
	// 	return ErrInvalidEntity
	// }
	err := b.Validate1()
	if err != nil {
		return errors.Wrap(err, "wrap 1")
	}

	return nil
}

//Validate validate book
func (b *Category) Validate1() error {
	// if b.Title == "" || b.Author == "" || b.Pages <= 0 || b.Quantity <= 0 {
	// 	return ErrInvalidEntity
	// }
	return errors.New("not enough arguments, expected at least 3, got %d")
}
