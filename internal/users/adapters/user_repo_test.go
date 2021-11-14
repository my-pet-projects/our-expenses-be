package adapters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/testing/mocks"
)

func TestNewUsersRepo_ReturnsRepository(t *testing.T) {
	t.Parallel()
	// Arrange
	client := &database.MongoClient{}
	log := new(mocks.LogInterface)

	// Act
	result := adapters.NewUserRepo(client, log)

	// Assert
	assert.NotNil(t, result, "Result should not be nil.")
}
