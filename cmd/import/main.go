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

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
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

	categoryRepo := repository.NewCategoryRepo(mongoClient, log)

	delRes, delErr := categoryRepo.DeleteAll(ctx, domain.CategoryFilter{})
	if delErr != nil {
		logrus.Fatalf("Failed to delete categories: '%+v'", delErr)
	}
	log.Infof(ctx, "Deleted %d categories", delRes.DeleteCount)

	categories := getCategories()

	log.Infof(ctx, "Inserting categories ...")
	insCount := 0
	for _, category := range categories {
		_, insErr := categoryRepo.Insert(ctx, category)
		if insErr != nil {
			logrus.Fatalf("Failed to insert category: '%+v'", insErr)
		}
		insCount++
	}

	log.Infof(ctx, "Inserted %d categories", insCount)

	os.Exit(0)
}

func getCategories() []domain.Category {
	jsonFile, err := os.Open("cmd/import/categories.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonCategories []category
	_ = json.Unmarshal(byteValue, &jsonCategories)

	var categories []domain.Category
	for _, jsonCategory := range jsonCategories {
		rootCategoryID := primitive.NewObjectID()
		rootCategory, _ := domain.NewCategory(
			rootCategoryID.Hex(),
			jsonCategory.Name,
			nil,
			fmt.Sprintf("|%s", rootCategoryID.Hex()),
			1,
			time.Now(),
			nil,
		)

		categories = append(categories, *rootCategory)

		for _, jsonSubcategory := range jsonCategory.Subcategories {
			childCategoryID := primitive.NewObjectID()
			rootCategoryID := rootCategory.ID()
			childCategory, _ := domain.NewCategory(
				childCategoryID.Hex(),
				jsonSubcategory.Name,
				&rootCategoryID,
				strings.ToLower(fmt.Sprintf("|%s|%s", rootCategory.ID(), childCategoryID.Hex())),
				2,
				time.Now(),
				nil,
			)

			categories = append(categories, *childCategory)
		}
	}

	return categories
}

type category struct {
	CategoryID    string `json:"categoryId"`
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	IconViewport  string `json:"iconViewport"`
	Subcategories []struct {
		ParentCategoryID   string      `json:"parentCategoryId"`
		ParentCategoryName string      `json:"parentCategoryName"`
		ParentCategoryIcon interface{} `json:"parentCategoryIcon"`
		SubcategoryID      string      `json:"subcategoryId"`
		Name               string      `json:"name"`
	} `json:"subcategories"`
}
