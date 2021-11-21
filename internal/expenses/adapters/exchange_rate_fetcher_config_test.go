package adapters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
)

func TestExchangeRateFetcherConfig_NewExchangeRateFetcherConfig_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Act
	result := adapters.NewExchangeRateFetcherConfig()

	// Assert
	assert.NotNil(t, result)
}
