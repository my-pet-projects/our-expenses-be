package repository

import (
	"context"
	"fmt"
	"our-expenses-server/entity"
	"our-expenses-server/logger"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName string = "categories"
const loggerCategory string = "repository/category"

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	logger logger.AppLoggerInterface
	db     *mongo.Database
}

// CategoryRepoInterface defines a contract to persist categories in the database.
type CategoryRepoInterface interface {
	GetAll(ctx context.Context, filter entity.CategoryFilter) ([]entity.Category, error)
	GetOne(ctx context.Context, id string) (*entity.Category, error)
	Insert(ctx context.Context, category *entity.Category) (*entity.ID, error)
	Update(ctx context.Context, category *entity.Category) (string, error)
	DeleteAll(ctx context.Context, filter entity.CategoryFilter) (int, error)
	DeleteOne(ctx context.Context, id string) (int, error)
}

// ProvideCategoryRepository returns a CategoryRepository.
func ProvideCategoryRepo(logger *logger.AppLogger, db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		logger: logger,
		db:     db,
	}
}

// collection returns collection handle.
func (repo *CategoryRepository) collection() *mongo.Collection {
	return repo.db.Collection(collectionName)
}

// GetAll returns all categories from the database that matches the filter.
func (repo *CategoryRepository) GetAll(ctx context.Context, filter entity.CategoryFilter) ([]entity.Category, error) {
	start := time.Now()
	categories := []entity.Category{}

	query := bson.M{}
	if filter.ParentID == "" {
		query["parentId"] = bson.M{"$exists": false}
	} else {
		parentID, _ := primitive.ObjectIDFromHex(filter.ParentID)
		query["parentId"] = parentID
	}

	parantCategoriesIDs := make([]primitive.ObjectID, len(filter.CategoryIDs))
	for index := range filter.CategoryIDs {
		parantCategoriesID, err := primitive.ObjectIDFromHex(filter.CategoryIDs[index])
		if err == nil {
			parantCategoriesIDs[index] = parantCategoriesID
		}
	}

	if len(parantCategoriesIDs) != 0 {
		query = bson.M{}
		query["_id"] = bson.M{"$in": parantCategoriesIDs}
	}

	if filter.FindChildren {
		query = bson.M{}
		query["path"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: fmt.Sprintf(".*\\|%s\\|.*", filter.CategoryID),
				Options: "i",
			},
		}
	}

	if filter.FindAll {
		query = bson.M{}
	}

	repo.logger.InfoWithFields(ctx, "Fetching categories from database ...", logger.FieldsSet{
		"component": "db/start",
		"payload":   fmt.Sprintf("%+v", query),
	})

	cursor, findError := repo.collection().Find(ctx, query)
	if findError != nil {
		return nil, errors.Wrap(findError, "find command")
	}

	allError := cursor.All(ctx, &categories)
	if allError != nil {
		fmt.Printf("%+v", allError)
		return nil, errors.Wrap(allError, "cursor iteration")
	}

	repo.logger.InfoWithFields(ctx, fmt.Sprintf("Found %d items", len(categories)), logger.FieldsSet{
		"component": "db/end",
		"duration":  time.Since(start),
	})

	return categories, nil
}

// GetOne returns a single category from the database.
func (repo *CategoryRepository) GetOne(ctx context.Context, id string) (*entity.Category, error) {
	category := entity.Category{}

	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	findError := repo.collection().FindOne(ctx, filter).Decode(&category)
	if findError != nil {
		if findError == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, findError
	}

	return &category, nil
}

// Insert a category into the database.
func (repo *CategoryRepository) Insert(ctx context.Context, category *entity.Category) (*entity.ID, error) {
	start := time.Now()

	// Datadase trace object?
	repo.logger.InfoWithFields(ctx, "Inserting category into database ...", logger.FieldsSet{
		"component": "db/start",
		"payload":   fmt.Sprintf("%+v", category),
	})

	insRes, insErr := repo.collection().InsertOne(ctx, category)
	if insErr != nil {
		return nil, errors.Wrap(insErr, "failed to insert category into db")
	}
	objID, _ := insRes.InsertedID.(primitive.ObjectID)

	repo.logger.InfoWithFields(ctx, fmt.Sprintf("Inserted category with ID %s", objID.Hex()), logger.FieldsSet{
		"component": "db/end",
		"duration":  time.Since(start),
	})

	return &objID, nil
}

// Update updates a category in the database.
func (repo *CategoryRepository) Update(ctx context.Context, category *entity.Category) (string, error) {
	filter := bson.M{"_id": category.ID}

	updater := bson.M{"$set": category}

	if &category.ParentID == nil {
		updater["$unset"] = bson.M{
			"parentId": "",
		}
	}

	fmt.Print(updater)

	updateResult, updateError := repo.collection().UpdateOne(ctx, filter, updater)
	if updateError != nil {
		return "", updateError
	}

	id, _ := updateResult.UpsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

// DeleteOne deletes a category in the database.
func (repo *CategoryRepository) DeleteOne(ctx context.Context, id string) (int, error) {
	categoryID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": categoryID}

	deleteResult, deleteError := repo.collection().DeleteOne(ctx, filter)
	if deleteError != nil {
		return 0, deleteError
	}

	deletedDocuments := deleteResult.DeletedCount

	return int(deletedDocuments), nil
}

// DeleteAll deletes all categories in the database.
func (repo *CategoryRepository) DeleteAll(ctx context.Context, filter entity.CategoryFilter) (int, error) {
	query := bson.M{}

	if len(filter.Path) != 0 && filter.FindChildren {
		query["path"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: fmt.Sprintf("^%s.*", strings.ReplaceAll(filter.Path, "|", "\\|")),
				Options: "i",
			},
		}
	}

	deleteResult, deleteError := repo.collection().DeleteMany(ctx, query)
	if deleteError != nil {
		return 0, deleteError
	}

	deletedDocuments := deleteResult.DeletedCount

	return int(deletedDocuments), nil
}
