package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category struct represents a category.
type Category struct {
	ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string              `json:"name"`
	ParentID *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Path     string              `json:"path"`
	Level    int32               `json:"level"`
}
