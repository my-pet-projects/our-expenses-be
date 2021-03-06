package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

const collectionName string = "categories"

// categoryModel defines category structure in MongoDB.
type categoryModel struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	ParentID  *primitive.ObjectID `bson:"parentId,omitempty"`
	Path      string              `bson:"path"`
	Icon      *string             `bson:"icon,omitempty"`
	Level     int                 `bson:"level"`
	CreatedBy string              `bson:"createdBy,omitempty"`
	UpdatedBy *string             `bson:"updatedBy,omitempty"`
	CreatedAt time.Time           `bson:"createdAt,omitempty"`
	UpdatedAt *time.Time          `bson:"updatedAt,omitempty"`
}

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// CategoryRepoInterface defines a contract to persist categories in the database.
type CategoryRepoInterface interface {
	GetAll(ctx context.Context, filter domain.CategoryFilter) ([]domain.Category, error)
	GetOne(ctx context.Context, id string) (*domain.Category, error)
	Insert(ctx context.Context, category domain.Category) (*string, error)
	Update(ctx context.Context, category domain.Category) (*domain.UpdateResult, error)
	DeleteAll(ctx context.Context, filter domain.CategoryFilter) (*domain.DeleteResult, error)
	DeleteOne(ctx context.Context, id string) (*domain.DeleteResult, error)
}

// NewCategoryRepo returns a CategoryRepository.
func NewCategoryRepo(client *database.MongoClient, logger logger.LogInterface) *CategoryRepository {
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
func (r *CategoryRepository) GetAll(
	ctx context.Context,
	filter domain.CategoryFilter,
) ([]domain.Category, error) {
	ctx, span := tracer.NewSpan(ctx, "find categories in the database")
	// span.SetAttributes(attribute.Any("filter", filter)) TODO: trace attributes?
	defer span.End()

	query := bson.M{}
	if filter.ParentID == nil {
		query["parentId"] = bson.M{"$exists": false}
	} else {
		parentID, _ := primitive.ObjectIDFromHex(*filter.ParentID)
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

	// span.AddEvent("start query", trace.WithAttributes(attribute.Any("filter", query)))

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

	span.AddEvent("fetched finished", trace.WithAttributes(attribute.Int("items", len(categoryModels))))

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
func (r *CategoryRepository) GetOne(ctx context.Context, id string) (*domain.Category, error) {
	ctx, span := tracer.NewSpan(ctx, "find categories in the database")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	categoryDbModel := categoryModel{}
	findError := r.collection().FindOne(ctx, filter).Decode(&categoryDbModel)
	if findError != nil {
		if errors.Is(findError, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.Wrap(findError, "find category")
	}

	category, categoryErr := r.unmarshalCategory(categoryDbModel)
	if categoryErr != nil {
		return nil, categoryErr
	}

	return category, nil
}

// Insert a category into the database.
func (r *CategoryRepository) Insert(ctx context.Context, category domain.Category) (*string, error) {
	ctx, span := tracer.NewSpan(ctx, "add category to the database")
	defer span.End()

	categoryDbModel := r.marshalCategory(category)
	tempUser := "kot"
	categoryDbModel.CreatedBy = tempUser
	categoryDbModel.CreatedAt = time.Now()

	insRes, insErr := r.collection().InsertOne(ctx, categoryDbModel)
	if insErr != nil {
		return nil, errors.Wrap(insErr, "mongodb insert category")
	}

	objID, _ := insRes.InsertedID.(primitive.ObjectID)
	objIDString := objID.Hex()

	return &objIDString, nil
}

// Update updates a category in the database.
func (r *CategoryRepository) Update(ctx context.Context, category domain.Category) (*domain.UpdateResult, error) {
	ctx, span := tracer.NewSpan(ctx, "update category in the database")
	defer span.End()

	categoryDbModel := r.marshalCategory(category)
	tempUser := "kot"
	now := time.Now()
	categoryDbModel.UpdatedBy = &tempUser
	categoryDbModel.UpdatedAt = &now

	filter := bson.M{"_id": categoryDbModel.ID}
	updater := bson.M{"$set": categoryDbModel}

	if categoryDbModel.ParentID == nil {
		updater["$unset"] = bson.M{
			"parentId": "",
		}
	}

	opts := options.Update().SetUpsert(false)

	// if &category.ParentID == nil {
	// 	updater["$unset"] = bson.M{
	// 		"parentId": "",
	// 	}
	// }

	mongoUpdResult, mongoUpdErr := r.collection().UpdateOne(ctx, filter, updater, opts)
	if mongoUpdErr != nil {
		return nil, errors.Wrap(mongoUpdErr, "mongo update category")
	}

	result := &domain.UpdateResult{
		UpdateCount: int(mongoUpdResult.ModifiedCount),
	}

	return result, nil
}

// DeleteOne deletes a category in the database.
func (r *CategoryRepository) DeleteOne(ctx context.Context, id string) (*domain.DeleteResult, error) {
	categoryID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": categoryID}

	mongoDelResult, mongoDelErr := r.collection().DeleteOne(ctx, filter)
	if mongoDelErr != nil {
		return nil, errors.Wrap(mongoDelErr, "mongo delete categories")
	}

	result := &domain.DeleteResult{
		DeleteCount: int(mongoDelResult.DeletedCount),
	}

	return result, nil
}

// DeleteAll deletes all categories in the database.
func (r *CategoryRepository) DeleteAll(
	ctx context.Context,
	filter domain.CategoryFilter,
) (*domain.DeleteResult, error) {
	query := bson.M{}

	// TODO: little safety net to prevent deleting all items

	if len(filter.Path) != 0 && filter.FindChildren {
		query["path"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: fmt.Sprintf("^%s.*", strings.ReplaceAll(filter.Path, "|", "\\|")),
				Options: "i",
			},
		}
	}

	mongoDelResult, mongoDelErr := r.collection().DeleteMany(ctx, query)
	if mongoDelErr != nil {
		return nil, errors.Wrap(mongoDelErr, "mongo delete categories")
	}

	result := &domain.DeleteResult{
		DeleteCount: int(mongoDelResult.DeletedCount),
	}

	return result, nil
}

func (r CategoryRepository) marshalCategory(category domain.Category) categoryModel {
	id, _ := primitive.ObjectIDFromHex(category.ID())
	var parentID *primitive.ObjectID
	if category.ParentID() != nil {
		parentIDObj, _ := primitive.ObjectIDFromHex(*category.ParentID())
		parentID = &parentIDObj
	}

	return categoryModel{
		ID:       id,
		Name:     category.Name(),
		Icon:     category.Icon(),
		ParentID: parentID,
		Path:     category.Path(),
		Level:    category.Level(),
	}
}

func (r CategoryRepository) unmarshalCategory(categoryModel categoryModel) (*domain.Category, error) {
	var parentID *string
	if categoryModel.ParentID != nil {
		parentIDHex := categoryModel.ParentID.Hex()
		parentID = &parentIDHex
	}
	cat, catErr := domain.NewCategory(categoryModel.ID.Hex(), categoryModel.Name,
		parentID, categoryModel.Path, categoryModel.Icon, categoryModel.Level)
	cat.SetMetadata(categoryModel.CreatedBy, categoryModel.CreatedAt, categoryModel.UpdatedBy, categoryModel.UpdatedAt)
	if catErr != nil {
		return nil, errors.Wrap(catErr, "unmarshal category")
	}
	return cat, nil
}
