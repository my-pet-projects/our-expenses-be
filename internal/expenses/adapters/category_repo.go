package adapters

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var categoriesRepoTracer trace.Tracer

const categoriesCollectionName string = "categories"

// categoryDbModel defines category structure in MongoDB.
type categoryDbModel struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	Name     string              `bson:"name"`
	ParentID *primitive.ObjectID `bson:"parentId,omitempty"`
	Parents  []categoryDbModel   `bson:"parents,omitempty"`
	Path     string              `bson:"path"`
	Icon     *string             `bson:"icon,omitempty"`
	Level    int                 `bson:"level"`
}

// CategoryRepository represents a struct to access categories MongoDB collection.
type CategoryRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// ExpenseCategoryRepoInterface defines a contract to persist categories in the database.
type ExpenseCategoryRepoInterface interface {
	GetOne(ctx context.Context, id string) (*domain.Category, error)
}

// NewCategoryRepo returns a CategoryRepository.
func NewCategoryRepo(client *database.MongoClient, logger logger.LogInterface) *CategoryRepository {
	categoriesRepoTracer = otel.Tracer("app.repository.categories")
	return &CategoryRepository{
		logger: logger,
		client: client,
	}
}

// collection returns collection handle.
func (r *CategoryRepository) collection() *mongo.Collection {
	return r.client.Collection(categoriesCollectionName)
}

// GetOne returns a single category from the database.
func (r *CategoryRepository) GetOne(ctx context.Context, id string) (*domain.Category, error) {
	ctx, span := categoriesRepoTracer.Start(ctx, "find categories in the database")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": objID}
	categoryDbModel := categoryDbModel{}
	findError := r.collection().FindOne(ctx, filter).Decode(&categoryDbModel)
	if findError != nil {
		if findError == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, findError
	}

	category, categoryErr := r.unmarshalCategory(categoryDbModel)
	if categoryErr != nil {
		return nil, categoryErr
	}

	return category, nil
}

func (r CategoryRepository) unmarshalCategory(categoryModel categoryDbModel) (*domain.Category, error) {
	var parentID *string
	if categoryModel.ParentID != nil {
		parentIDHex := categoryModel.ParentID.Hex()
		parentID = &parentIDHex
	}
	cat, catErr := domain.NewCategory(categoryModel.ID.Hex(), parentID, categoryModel.Name,
		categoryModel.Icon, categoryModel.Level, categoryModel.Path)
	if catErr != nil {
		return nil, errors.Wrap(catErr, "unmarshal category")
	}
	return cat, nil
}
