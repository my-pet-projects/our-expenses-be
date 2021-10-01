package adapters

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

var categoriesRepoTracer trace.Tracer

const collectionName string = "expenses"

type categoryDbModel struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	ParentID  *primitive.ObjectID `bson:"parentId,omitempty"`
	Parents   []categoryDbModel   `bson:"parents,omitempty"`
	Path      string              `bson:"path"`
	Icon      *string             `bson:"icon,omitempty"`
	Level     int                 `bson:"level"`
	CreatedAt time.Time           `bson:"createdAt"`
	UpdatedAt *time.Time          `bson:"updatedAt,omitempty"`
}

type expenseDbModel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID primitive.ObjectID `bson:"categoryId"`
	Category   *categoryDbModel   `bson:"category,omitempty"`
	Price      float64            `bson:"price"`
	Currency   string             `bson:"currency"`
	Quantity   float64            `bson:"quantity"`
	Date       time.Time          `bson:"date"`
	Comment    *string            `bson:"comment,omitempty"`
	Trip       *string            `bson:"trip,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty"`
	CreatedBy  string             `bson:"createdBy,omitempty"`
	UpdatedAt  *time.Time         `bson:"updatedAt,omitempty"`
	UpdatedBy  *string            `bson:"updatedBy,omitempty"`
}

// ExpenseRepository represents a struct to access expenses MongoDB collection.
type ExpenseRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// ExpenseRepoInterface defines a contract to persist expenses in the database.
type ExpenseRepoInterface interface {
	Insert(ctx context.Context, expense domain.Expense) (*string, error)
	DeleteAll(ctx context.Context) (*domain.DeleteResult, error)
}

// NewExpenseRepo returns a Expenseadapters.
func NewExpenseRepo(client *database.MongoClient, logger logger.LogInterface) *ExpenseRepository {
	categoriesRepoTracer = otel.Tracer("app.adapters.expenses")
	return &ExpenseRepository{
		logger: logger,
		client: client,
	}
}

// collection returns collection handle.
func (r *ExpenseRepository) collection() *mongo.Collection {
	return r.client.Collection(collectionName)
}

// Insert insert a new record into database.
func (r *ExpenseRepository) Insert(ctx context.Context, category domain.Expense) (*string, error) {
	ctx, span := categoriesRepoTracer.Start(ctx, "add expense to the database")
	defer span.End()

	dbModel := r.marshalExpense(category)
	tempUser := "kot"
	dbModel.CreatedBy = tempUser
	dbModel.CreatedAt = time.Now()

	insRes, insErr := r.collection().InsertOne(ctx, dbModel)
	if insErr != nil {
		return nil, errors.Wrap(insErr, "mongodb insert expense")
	}

	objID, _ := insRes.InsertedID.(primitive.ObjectID)
	objIDString := objID.Hex()

	return &objIDString, nil
}

// DeleteAll deletes all expenses in the database.
func (r *ExpenseRepository) DeleteAll(ctx context.Context) (*domain.DeleteResult, error) {
	query := bson.M{}
	mongoDelResult, mongoDelErr := r.collection().DeleteMany(ctx, query)
	if mongoDelErr != nil {
		return nil, errors.Wrap(mongoDelErr, "mongo delete expenses")
	}

	result := &domain.DeleteResult{
		DeleteCount: int(mongoDelResult.DeletedCount),
	}

	return result, nil
}

// marshalExpense marshalls expense domain object into MongoDB model.
func (r ExpenseRepository) marshalExpense(expense domain.Expense) expenseDbModel {
	id, _ := primitive.ObjectIDFromHex(expense.ID())
	categoryID, _ := primitive.ObjectIDFromHex(expense.CategoryID())

	return expenseDbModel{
		ID:         id,
		CategoryID: categoryID,
		Price:      expense.Price(),
		Currency:   expense.Currency(),
		Quantity:   expense.Quantity(),
		Comment:    expense.Comment(),
		Trip:       expense.Trip(),
		Date:       expense.Date(),
	}
}

func (r ExpenseRepository) unmarshalExpense(expenseModel expenseDbModel) (*domain.Expense, error) {
	exp, expErr := domain.NewExpense(expenseModel.ID.Hex(), expenseModel.CategoryID.Hex(),
		expenseModel.Price, expenseModel.Currency, expenseModel.Quantity,
		expenseModel.Comment, expenseModel.Trip, expenseModel.Date)
	if expErr != nil {
		return nil, errors.Wrap(expErr, "unmarshal expense")
	}
	exp.SetMetadata(expenseModel.CreatedBy, expenseModel.CreatedAt, expenseModel.UpdatedBy, expenseModel.UpdatedAt)
	return exp, nil
}
