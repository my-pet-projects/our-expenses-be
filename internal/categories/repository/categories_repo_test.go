package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewCategoriesRepo_ReturnsRepository(t *testing.T) {
	// Arrange
	client := &database.MongoClient{}
	log := new(mocks.LogInterface)

	// Act
	result := repository.NewCategoryRepo(client, log)

	// Assert
	assert.NotNil(t, result, "Result should not be nil.")
}
