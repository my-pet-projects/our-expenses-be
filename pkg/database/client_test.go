package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewMongoClient_InvalidConfig_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Database{
		Mongo: config.MongoDB{
			Name:     "name",
			User:     "user",
			Pass:     "pass",
			URI:      "uri",
			Database: "db",
		},
	}
	logger := new(mocks.LogInterface)

	// Act
	result, err := database.NewMongoClient(logger, config)

	// Assert
	assert.Nil(t, result, "Result should be nil.")
	assert.NotNil(t, err, "Error should not be nil.")
}

func TestNewMongoClient_ReturnsMongoClient(t *testing.T) {
	t.Parallel()
	// Arrange
	config := config.Database{
		Mongo: config.MongoDB{
			Name:     "name",
			User:     "user",
			Pass:     "pass",
			URI:      "mongodb://mongodb",
			Database: "db",
		},
	}
	logger := new(mocks.LogInterface)

	// Act
	result, err := database.NewMongoClient(logger, config)

	// Assert
	assert.NotNil(t, result, "Result should not be nil.")
	assert.Nil(t, err, "Error should be nil.")
}
