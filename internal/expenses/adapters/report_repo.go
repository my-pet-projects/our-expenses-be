package adapters

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

// ReportRepository represents a struct to access expenses MongoDB collection.
type ReportRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
	tracer trace.Tracer
}

// ReportRepoInterface defines a contract to persist expenses in the database.
type ReportRepoInterface interface {
	GetAll(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, error)
}

// NewReportRepo returns a report repository.
func NewReportRepo(client *database.MongoClient, logger logger.LogInterface) *ReportRepository {
	return &ReportRepository{
		logger: logger,
		client: client,
		tracer: otel.Tracer("app.repository.report"),
	}
}

// collection returns collection handle.
func (r *ReportRepository) collection() *mongo.Collection {
	return r.client.Collection("expenses")
}

// GetAll returns all expenses from the database that matches the filter.
func (r *ReportRepository) GetAll(ctx context.Context, filter domain.ExpenseFilter) ([]domain.Expense, error) {
	ctx, span := r.tracer.Start(ctx, "find expenses in the database")
	// span.SetAttributes(attribute.Any("filter", filter))
	defer span.End()

	// Filter expense documents.
	matchStage := bson.M{
		"$match": bson.M{
			"date": bson.M{
				"$gte": filter.From(),
				"$lte": filter.To(),
			},
		},
	}

	// Join with categories collection to get expense category.
	categoryLookupStage := bson.M{
		"$lookup": bson.M{
			"from":         "categories",
			"localField":   "categoryId",
			"foreignField": "_id",
			"as":           "category",
		},
	}

	// After the lookup take the very first match.
	addCategoryFieldStage := bson.M{
		"$addFields": bson.M{
			"category": bson.M{
				"$arrayElemAt": []interface{}{"$category", 0},
			},
		},
	}

	// Join with categories collection to get all ascendants.
	// There will be expense category as well, that needs to be filtered out.
	ascendantsLookupStage := bson.M{
		"$graphLookup": bson.M{
			"from":             "categories",
			"startWith":        "$categoryId",
			"connectFromField": "parentId",
			"connectToField":   "_id",
			"as":               "parentCategories",
		},
	}

	// Add parents to the category document, but exclude expense category.
	addAscendantsFieldStage := bson.M{
		"$addFields": bson.M{
			"category.parents": bson.M{
				"$filter": bson.M{
					"input": "$parentCategories",
					"cond": bson.M{
						"$ne": []interface{}{"$$this._id", "$categoryId"},
					},
				},
			},
		},
	}

	operations := []bson.M{
		matchStage,
		categoryLookupStage,
		addCategoryFieldStage,
		ascendantsLookupStage,
		addAscendantsFieldStage,
	}

	// span.AddEvent("start query", trace.WithAttributes(attribute.Any("filter", operations)))

	cursor, cursorErr := r.collection().Aggregate(ctx, operations)
	if cursorErr != nil {
		return nil, errors.Wrap(cursorErr, "mongodb cursor expense")
	}

	span.AddEvent("cursor iteration")

	var expenseDbModels []expenseDbModel
	if allError := cursor.All(ctx, &expenseDbModels); allError != nil {
		return nil, errors.Wrap(allError, "cursor iteration")
	}

	span.AddEvent("fetched finished", trace.WithAttributes(attribute.Int("items", len(expenseDbModels))))

	expenses := []domain.Expense{}
	for _, expenseDbModel := range expenseDbModels {
		exp, expErr := r.unmarshalExpense(expenseDbModel)
		if expErr != nil {
			return nil, expErr
		}

		expenses = append(expenses, *exp)
	}

	return expenses, nil
}

func (r ReportRepository) unmarshalExpense(expenseModel expenseDbModel) (*domain.Expense, error) {
	cat, catErr := r.unmarshalCategory(*expenseModel.Category)
	if catErr != nil {
		return nil, errors.Wrap(catErr, "unmarshal category")
	}

	exp, expErr := domain.NewExpense(expenseModel.ID.Hex(), *cat,
		expenseModel.Price, expenseModel.Currency, expenseModel.Quantity,
		expenseModel.Comment, expenseModel.Trip, expenseModel.Date)
	if expErr != nil {
		return nil, errors.Wrap(expErr, "unmarshal expense")
	}

	return exp, nil
}

func (r ReportRepository) unmarshalCategory(categoryModel categoryDbModel) (*domain.Category, error) {
	var parentID string
	if categoryModel.ParentID != nil && !categoryModel.ParentID.IsZero() {
		parentID = categoryModel.ParentID.Hex()
	}
	cat, catErr := domain.NewCategory(categoryModel.ID.Hex(), &parentID,
		categoryModel.Name, categoryModel.Icon, categoryModel.Level, categoryModel.Path)
	if catErr != nil {
		return nil, errors.Wrap(catErr, "unmarshal category")
	}

	parentCategories := make([]domain.Category, 0)
	for _, parentCat := range categoryModel.Parents {
		var parentID string
		if parentCat.ParentID != nil {
			parentID = parentCat.ParentID.Hex()
		}
		parent, parentErr := domain.NewCategory(parentCat.ID.Hex(), &parentID, parentCat.Name,
			parentCat.Icon, parentCat.Level, parentCat.Path)
		if parentErr != nil {
			return nil, errors.Wrap(parentErr, "unmarshal parent category")
		}
		parentCategories = append(parentCategories, *parent)
	}
	cat.SetParents(&parentCategories)

	return cat, nil
}
