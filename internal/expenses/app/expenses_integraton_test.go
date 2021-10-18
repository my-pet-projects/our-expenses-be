package app

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
	"github.com/stretchr/testify/assert"
)

func setupApp(t *testing.T) *Application {
	ctx, cancel := context.WithCancel(context.Background())
	appConfig, appConfigErr := config.NewConfig()
	if appConfigErr != nil {
		t.Fatalf("Failed to create application config: %v", appConfigErr)
	}
	appConfig.Logger.Writers.FileWriter.Enabled = false
	if cfgValidErr := appConfig.Validate(); cfgValidErr != nil {
		t.Fatalf("Invalid application config: %v", cfgValidErr)
	}
	appLogger, appLoggerErr := logger.NewLogger(appConfig.Logger)
	if appLoggerErr != nil {
		t.Fatalf("Failed to create application logger: %v", appLoggerErr)
	}
	appTracer := tracer.NewTracer(appConfig.Telemetry)

	mongoClient, mongoClientErr := database.NewMongoClient(appLogger, appConfig.Database)
	if mongoClientErr != nil {
		t.Fatalf("Failed to create MongoDB client: %v", mongoClientErr)

	}
	if mongoConErr := mongoClient.OpenConnection(ctx, cancel); mongoConErr != nil {
		t.Fatalf("Failed to open MongoDB connection: %v", mongoConErr)
	}

	expensesApp, expensesAppErr := NewApplication(ctx, cancel, appConfig, appLogger, appTracer, mongoClient)
	if expensesAppErr != nil {
		t.Fatalf("Failed to instantiate Expenses application: %v", expensesAppErr)
	}

	return expensesApp
}

func TestNewServer_ReturnsServerInstanceWithSettingsFromConfig(t *testing.T) {
	// Arrange
	ctx := context.Background()
	from := time.Date(2021, 8, 2, 0, 0, 0, 0, time.UTC)
	to := time.Date(2021, 8, 2, 23, 59, 59, 999, time.UTC)
	findQuery := query.FindExpensesQuery{
		From: from,
		To:   to,
	}

	// SUT
	sut := setupApp(t)

	// Act
	res, resErr := sut.Queries.FindExpenses.Handle(ctx, findQuery)

	// Assert
	assert.NotNil(t, resErr)
	assert.NotNil(t, res)
}
