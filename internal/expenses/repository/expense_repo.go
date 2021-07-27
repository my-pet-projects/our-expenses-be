package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
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

type expenseDbModel struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	CategoryID string               `bson:"categoryId"`
	Price      primitive.Decimal128 `bson:"price"`
	Currency   string               `bson:"currency"`
	Quantity   int                  `bson:"quantity"`
	Comment    *string              `bson:"comment,omitempty"`
	CreatedAt  time.Time            `bson:"createdAt"`
	UpdatedAt  *time.Time           `bson:"updatedAt,omitempty"`
}

// ExpenseRepository represents a struct to access expenses MongoDB collection.
type ExpenseRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// ExpenseRepoInterface defines a contract to persist expenses in the database.
type ExpenseRepoInterface interface {
	Insert(ctx context.Context, expense domain.Expense) (*string, error)
}

// NewExpenseRepo returns a ExpenseRepository.
func NewExpenseRepo(client *database.MongoClient, logger logger.LogInterface) *ExpenseRepository {
	categoriesRepoTracer = otel.Tracer("app.repository.expenses")
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

	insRes, insErr := r.collection().InsertOne(ctx, dbModel)
	if insErr != nil {
		return nil, errors.Wrap(insErr, "mongodb insert expense")
	}

	objID, _ := insRes.InsertedID.(primitive.ObjectID)
	objIDString := objID.Hex()

	return &objIDString, nil
}

// marshalExpense marshalls expense domain object into MongoDB model.
func (r ExpenseRepository) marshalExpense(expense domain.Expense) expenseDbModel {
	id, _ := primitive.ObjectIDFromHex(expense.ID())
	price, _ := primitive.ParseDecimal128(expense.Price())

	return expenseDbModel{
		ID:         id,
		CategoryID: expense.CategoryID(),
		Price:      price,
		Currency:   expense.Currency(),
		Quantity:   expense.Quantity(),
		Comment:    expense.Comment(),
		CreatedAt:  expense.CreatedAt(),
		UpdatedAt:  expense.UpdatedAt(),
	}
}
