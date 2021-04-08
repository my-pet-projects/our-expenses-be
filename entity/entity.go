package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ID entity ID.
type ID = primitive.ObjectID

// NewID creates a new entity ID.
func NewID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// IDFromString creates entity ID from string value.
func IDFromString(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}
