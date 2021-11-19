package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
)

func TestDateRange_NewDateRange_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	to := time.Now()
	from := to.Add(-10 * 24 * time.Hour)

	// Act
	res, resErr := domain.NewDateRange(from, to)

	// Assert
	assert.NotNil(t, res)
	assert.Nil(t, resErr)
}

func TestDateRange_NewDateRange_InvalidDates_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	from := time.Now()
	to := from.Add(-1 * time.Hour)

	// Act
	res, resErr := domain.NewDateRange(from, to)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestDateRange_DatesInBetween_ReturnsDates(t *testing.T) {
	t.Parallel()
	// Arrange
	to := time.Now()
	from := to.Add(-10 * 24 * time.Hour)

	// SUT
	sut, _ := domain.NewDateRange(from, to)

	// Act
	res := sut.DatesInBetween()

	// Assert
	assert.NotNil(t, res)
	assert.Len(t, res, 10)
}
