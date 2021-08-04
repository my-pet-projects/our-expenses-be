package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	categoryDomain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	categoryRepo "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	expenseDomain "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	expenseRepo "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, _ := config.NewConfig()
	log, _ := logger.NewLogger(cfg.Logger)
	mongoClient, mongoClientErr := database.NewMongoClient(log, cfg.Database)
	if mongoClientErr != nil {
		logrus.Fatalf("failed to create mongodb client: '%+v'", mongoClientErr)
	}
	if mongoConErr := mongoClient.OpenConnection(ctx, cancel); mongoConErr != nil {
		logrus.Fatalf("failed to open mongodb connection: '%+v'", mongoClientErr)
	}
	categoryRepo := categoryRepo.NewCategoryRepo(mongoClient, log)
	expenseRepo := expenseRepo.NewExpenseRepo(mongoClient, log)

	catDelRes, catDelErr := categoryRepo.DeleteAll(ctx, categoryDomain.CategoryFilter{})
	if catDelErr != nil {
		logrus.Fatalf("Failed to delete categories: '%+v'", catDelErr)
	}
	log.Infof(ctx, "Deleted %d categories", catDelRes.DeleteCount)

	expDelRes, expDelErr := expenseRepo.DeleteAll(ctx)
	if expDelErr != nil {
		logrus.Fatalf("Failed to delete expenses: '%+v'", expDelErr)
	}
	log.Infof(ctx, "Deleted %d expenses", expDelRes.DeleteCount)

	categories, oldNewCategoriesMap := getCategories()
	expenses := getExpenses(oldNewCategoriesMap)

	log.Infof(ctx, "Inserting categories ...")
	catInsCount := 0
	for _, category := range categories {
		_, catInsErr := categoryRepo.Insert(ctx, category)
		if catInsErr != nil {
			logrus.Fatalf("Failed to insert category: '%+v'", catInsErr)
		}
		catInsCount++
	}
	log.Infof(ctx, "Inserted %d categories", catInsCount)

	log.Infof(ctx, "Inserting expenses ...")
	expInsCount := 0
	for _, expense := range expenses {
		_, expInsErr := expenseRepo.Insert(ctx, expense)
		if expInsErr != nil {
			logrus.Fatalf("Failed to insert category: '%+v'", expInsErr)
		}
		expInsCount++
	}
	log.Infof(ctx, "Inserted %d expenses", expInsCount)

	os.Exit(0)
}

func getExpenses(oldNewCategoriesMap map[string]categoryDomain.Category) []expenseDomain.Expense {
	var rawExpenses itemsExpenses
	unmarshallErr := json.Unmarshal(readFile("cmd/import/expenses.json"), &rawExpenses)
	if unmarshallErr != nil {
		fmt.Printf("failed to unmarshall expenses %v", unmarshallErr)
		os.Exit(1)
	}

	var expenses []expenseDomain.Expense
	for _, rawExpense := range rawExpenses.Items {
		newCat, ok := oldNewCategoriesMap[rawExpense.CategoryID.S]
		if !ok {
			fmt.Printf("matching category not found")
			os.Exit(1)
		}

		date, _ := time.Parse("2006-01-02", rawExpense.Date.S)

		expenseID := primitive.NewObjectID()
		expense, _ := expenseDomain.NewExpense(
			expenseID.Hex(),
			newCat.ID(),
			rawExpense.Price.S,
			rawExpense.Currency.S,
			rawExpense.Quantity.S,
			&rawExpense.Comment.S,
			date,
			time.Now(),
			nil,
		)
		expenses = append(expenses, *expense)
	}

	return expenses
}

func getCategories() ([]categoryDomain.Category, map[string]categoryDomain.Category) {
	var rawCategories []category
	unmarshallErr := json.Unmarshal(readFile("cmd/import/categories.json"), &rawCategories)
	if unmarshallErr != nil {
		fmt.Printf("failed to unmarshall categories %v", unmarshallErr)
		os.Exit(1)
	}

	var oldNewCategoriesMap = make(map[string]categoryDomain.Category)
	var categories []categoryDomain.Category
	for _, rawCategory := range rawCategories {
		rootCategoryID := primitive.NewObjectID()
		rootCategory, _ := categoryDomain.NewCategory(
			rootCategoryID.Hex(),
			rawCategory.Name,
			nil,
			fmt.Sprintf("|%s", rootCategoryID.Hex()),
			nil,
			1,
			time.Now(),
			nil,
		)

		categories = append(categories, *rootCategory)

		for _, rawSubcategory := range rawCategory.Subcategories {
			childCategoryID := primitive.NewObjectID()
			rootCategoryID := rootCategory.ID()
			childCategory, _ := categoryDomain.NewCategory(
				childCategoryID.Hex(),
				rawSubcategory.Name,
				&rootCategoryID,
				strings.ToLower(fmt.Sprintf("|%s|%s", rootCategory.ID(), childCategoryID.Hex())),
				nil,
				2,
				time.Now(),
				nil,
			)

			oldNewCategoriesMap[rawSubcategory.SubcategoryID] = *childCategory
			categories = append(categories, *childCategory)
		}
	}

	return categories, oldNewCategoriesMap
}

func readFile(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)
	return bytes
}

type category struct {
	CategoryID    string `json:"CategoryId"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	IconViewport  string `json:"iconViewport"`
	Subcategories []struct {
		ParentCategoryID   string      `json:"parentCategoryId"`
		ParentCategoryName string      `json:"parentCategoryName"`
		ParentCategoryIcon interface{} `json:"parentCategoryIcon"`
		SubcategoryID      string      `json:"subcategoryId"`
		Name               string      `json:"name"`
	} `json:"Subcategories"`
}

type itemsExpenses struct {
	Items []expense `json:"Items"`
}

type expense struct {
	ExpenseID  sStruct `json:"expenseId"`
	Price      sStruct `json:"price"`
	Quantity   sStruct `json:"quantity"`
	Currency   sStruct `json:"currency"`
	Comment    sStruct `json:"comment"`
	Date       sStruct `json:"date"`
	Trip       sStruct `json:"trip"`
	CategoryID sStruct `json:"SubcategoryId"`
}

type sStruct struct {
	S string `json:"S"`
}
