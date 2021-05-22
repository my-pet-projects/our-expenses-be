package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/entity"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// categoryModel defines category structure in MongoDB.
type categoryModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ParentID  primitive.ObjectID `bson:"parentId,omitempty"`
	Path      string             `bson:"path"`
	Level     int                `bson:"level"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty"`
}

const collectionName string = "categories"
const loggerCategory string = "repository/category"

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

var tracer trace.Tracer

// CategoryRepoInterface defines a contract to persist categories in the database.
type CategoryRepoInterface interface {
	GetAll(ctx context.Context, filter entity.CategoryFilter) ([]domain.Category, error)
	GetOne(ctx context.Context, id string) (*entity.Category, error)
	Insert(ctx context.Context, category *entity.Category) (*entity.ID, error)
	Update(ctx context.Context, category *entity.Category) (string, error)
	DeleteAll(ctx context.Context, filter entity.CategoryFilter) (int, error)
	DeleteOne(ctx context.Context, id string) (int, error)
}

// NewCategoryRepo returns a CategoryRepository.
func NewCategoryRepo(client *database.MongoClient, logger logger.LogInterface) *CategoryRepository {
	tracer = otel.Tracer("categories repository")
	return &CategoryRepository{
		logger: logger,
		client: client,
	}
}

// collection returns collection handle.
func (r *CategoryRepository) collection() *mongo.Collection {
	return r.client.Collection(collectionName)
}

// GetAll returns all categories from the database that matches the filter.
func (r *CategoryRepository) GetAll(ctx context.Context, filter entity.CategoryFilter) ([]domain.Category, error) {
	ctx, span := tracer.Start(ctx, "find categories in the database")
	// span.SetAttributes(attribute.Any("from", from), attribute.Any("to", to))
	defer span.End()

	start := time.Now()

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

	r.logger.InfoWithFields(ctx, "Fetching categories from database ...", logger.FieldsSet{
		"component": "db/start",
		"payload":   fmt.Sprintf("%+v", query),
	})

	span.AddEvent("start query", trace.WithAttributes(attribute.Any("filter", query)))

	cursor, findError := r.collection().Find(ctx, query)
	if findError != nil {
		return nil, errors.Wrap(findError, "find command")
	}

	span.AddEvent("cursor iteration")

	categoryModels := []categoryModel{}
	allError := cursor.All(ctx, &categoryModels)
	if allError != nil {
		return nil, errors.Wrap(allError, "cursor iteration")
	}

	span.AddEvent("fetched finished", trace.WithAttributes(attribute.Any("items", len(categoryModels))))

	r.logger.InfoWithFields(ctx, fmt.Sprintf("Found %d items", len(categoryModels)), logger.FieldsSet{
		"component": "db/end",
		"duration":  time.Since(start),
	})

	categories := []domain.Category{}
	for _, categoryModel := range categoryModels {
		cat, catErr := r.unmarshalCategory(categoryModel)
		if catErr != nil {
			return nil, catErr
		}
		categories = append(categories, *cat)
	}

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

func (r CategoryRepository) unmarshalCategory(categoryModel categoryModel) (*domain.Category, error) {
	cat, catErr := domain.NewCategory(categoryModel.ID.Hex(), categoryModel.Name,
		categoryModel.ParentID.Hex(), categoryModel.Path, categoryModel.Level)
	if catErr != nil {
		return nil, errors.Wrap(catErr, "unmarshal category")
	}
	return cat, nil
}
