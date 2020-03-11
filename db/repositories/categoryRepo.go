package repositories

import (
	"context"
	"our-expenses-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Category represents a category database entity.
type Category struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	db *mongo.Database
}

// ProvideCategoryRepository returns a CategoryRepository.
func ProvideCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// collectionName returns category collection name.
func (repo *CategoryRepository) collectionName() string {
	return "categories"
}

// collection returns category collection handle.
func (repo *CategoryRepository) collection() *mongo.Collection {
	return repo.db.Collection(repo.collectionName())
}

// GetAll returns all categories from the database.
func (repo *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	categories := []models.Category{}

	findOpts := options.Find()
	cursor, findError := repo.collection().Find(ctx, findOpts)
	if findError != nil {
		return nil, findError
	}

	allError := cursor.All(ctx, &categories)
	if allError != nil {
		return nil, allError
	}

	return categories, nil
}

// GetOne returns a single category from the database.
func (repo *CategoryRepository) GetOne(ctx context.Context, id string) (*Category, error) {
	category := Category{}

	objID, objError := primitive.ObjectIDFromHex(id)
	if objError != nil {
		return nil, objError
	}

	filter := bson.M{"_id": objID}
	findError := repo.collection().FindOne(ctx, filter).Decode(&category)
	if findError != nil {
		return nil, findError
	}

	return &category, nil
}

// Save inserts a category to the database.
func (repo *CategoryRepository) Save(ctx context.Context, category *models.Category) (string, error) {
	insertResult, insertError := repo.collection().InsertOne(ctx, category)
	if insertError != nil {
		return "", insertError
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}
