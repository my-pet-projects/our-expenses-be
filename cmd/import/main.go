package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"our-expenses-server/container"
	"our-expenses-server/db/repositories"
	"our-expenses-server/models"
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
	categoryRepo := repositories.ProvideCategoryRepository(&db)
	ctx := context.Background()

	deleteCount, deleteErr := categoryRepo.DeleteAll(ctx)
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

func getCategories() []models.Category {
	jsonFile, err := os.Open("cmd/import/categories.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jsonCategories []category
	_ = json.Unmarshal([]byte(byteValue), &jsonCategories)

	var categories []models.Category
	for _, jsonCategory := range jsonCategories {
		rootCategoryID := primitive.NewObjectID()
		rootCategory := models.Category{
			ID:       &rootCategoryID,
			Name:     jsonCategory.Name,
			ParentID: &primitive.ObjectID{},
			Path:     fmt.Sprintf("/%s", jsonCategory.Name),
		}
		categories = append(categories, rootCategory)

		for _, jsonSubcategory := range jsonCategory.Subcategories {
			childCategory := models.Category{
				Name:     jsonSubcategory.Name,
				ParentID: rootCategory.ID,
				Path:     strings.ToLower(fmt.Sprintf("/%s/%s", rootCategory.Name, jsonSubcategory.Name)),
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
