package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewReportRepo_ReturnsRepository(t *testing.T) {
	// Arrange
	client := &database.MongoClient{}
	log := new(mocks.LogInterface)

	// Act
	result := repository.NewReportRepo(client, log)

	// Assert
	assert.NotNil(t, result, "Result should not be nil.")
}
