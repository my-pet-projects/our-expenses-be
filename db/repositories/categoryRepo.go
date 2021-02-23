package repositories

import (
	"context"
	"our-expenses-server/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	db *mongo.Database
}

// CategoryRepoInterface defines a contract to persist categories in the database.
type CategoryRepoInterface interface {
	GetAll(ctx context.Context, filter models.CategoryFilter) ([]models.Category, error)
	GetOne(ctx context.Context, id string) (*models.Category, error)
	Insert(ctx context.Context, category *models.Category) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) (string, error)
	DeleteAll(ctx context.Context) (int64, error)
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

// GetAll returns all categories from the database that matches the filter.
func (repo *CategoryRepository) GetAll(ctx context.Context, filter models.CategoryFilter) ([]models.Category, error) {
	categories := []models.Category{}

	query := bson.M{}
	if filter.ParentID == "" {
		query["parentId"] = bson.M{"$exists": false}
	} else {
		parentID, _ := primitive.ObjectIDFromHex(filter.ParentID)
		query["parentId"] = parentID
	}

	cursor, findError := repo.collection().Find(ctx, query)
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
func (repo *CategoryRepository) GetOne(ctx context.Context, id string) (*models.Category, error) {
	category := models.Category{}

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

// Insert inserts a category to the database.
func (repo *CategoryRepository) Insert(ctx context.Context, category *models.Category) (*models.Category, error) {
	insertResult, insertError := repo.collection().InsertOne(ctx, category)
	if insertError != nil {
		return nil, insertError
	}

	id, _ := insertResult.InsertedID.(primitive.ObjectID)

	category.ID = &id

	return category, nil
}

// Update updates a category in the database.
func (repo *CategoryRepository) Update(ctx context.Context, category *models.Category) (string, error) {
	updateResult, updateError := repo.collection().UpdateOne(ctx, bson.M{"_id": category.ID}, category)
	if updateError != nil {
		return "", updateError
	}

	id, _ := updateResult.UpsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

// DeleteAll deletes all categories in the database.
func (repo *CategoryRepository) DeleteAll(ctx context.Context) (int64, error) {
	deleteResult, deleteError := repo.collection().DeleteMany(ctx, bson.M{})
	if deleteError != nil {
		return 0, deleteError
	}

	deletedDocuments := deleteResult.DeletedCount

	return deletedDocuments, nil
}
