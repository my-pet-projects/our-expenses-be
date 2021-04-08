package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"our-expenses-server/config"
	"our-expenses-server/container"
	"our-expenses-server/entity"
	"our-expenses-server/infrastructure/db/repository"
	"our-expenses-server/logger"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	mongoDb, mongoError := container.InitDatabase()
	if mongoError != nil {
		logrus.Fatalf("Failed establish MongoDB connection: '%s'", mongoError)
	}
	var db mongo.Database = *mongoDb

	configConfig, _ := config.ProvideConfiguration()
	appLogger, _ := logger.ProvideLogger(configConfig)
	categoryRepo := repository.ProvideCategoryRepo(appLogger, &db)

	ctx := context.Background()

	deleteCount, deleteErr := categoryRepo.DeleteAll(ctx, entity.CategoryFilter{})
	if deleteErr != nil {
		logrus.Fatalf("Failed to delete categories: '%s'", deleteErr)
	}
	logrus.Infof("Deleted %s categories", strconv.FormatInt(int64(deleteCount), 10))

	categories := getCategories()
	for _, category := range categories {
		categoryRepo.Insert(ctx, &category)
	}

	os.Exit(0)
}

func getCategories() []entity.Category {
	jsonFile, err := os.Open("cmd/import/categories.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonCategories []category
	_ = json.Unmarshal([]byte(byteValue), &jsonCategories)

	var categories []entity.Category
	for _, jsonCategory := range jsonCategories {
		rootCategoryID := primitive.NewObjectID()
		rootCategory := entity.Category{
			ID:       rootCategoryID,
			Name:     jsonCategory.Name,
			ParentID: primitive.ObjectID{},
			Path:     fmt.Sprintf("|%s", rootCategoryID.Hex()),
			Level:    1,
		}
		categories = append(categories, rootCategory)

		for _, jsonSubcategory := range jsonCategory.Subcategories {
			childCategoryID := primitive.NewObjectID()
			childCategory := entity.Category{
				ID:       childCategoryID,
				Name:     jsonSubcategory.Name,
				ParentID: rootCategory.ID,
				Path:     strings.ToLower(fmt.Sprintf("|%s|%s", rootCategory.ID.Hex(), childCategoryID.Hex())),
				Level:    2,
			}
			categories = append(categories, childCategory)
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
