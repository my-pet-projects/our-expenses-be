package app

import (
	"context"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/adapters"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

// Application provides an application.
type Application struct {
	Commands Commands
	Queries  Queries
	Logger   logger.LogInterface
	Config   config.Config
	Tracer   tracer.TraceInterface
}

// Commands struct holds available application commands.
type Commands struct {
	AddExpense         command.AddExpenseHandlerInterface
	FetchExchangeRates command.FetchExchangeRatesHandlerInterface
}

// Queries struct holds available application queries.
type Queries struct {
	FindExpenses query.FindExpensesHandlerInterface
	FindCategory query.FindExpenseCategoryHandlerInterface
}

// NewApplication returns application instance.
func NewApplication(
	ctx context.Context,
	cancel context.CancelFunc,
	config *config.Config,
	logger *logger.Logger,
	tracer *tracer.Tracer,
	mongoClient *database.MongoClient,
) (*Application, error) {
	rateConfig := adapters.NewExchangeRateFetcherConfig()
	expenseRepo := adapters.NewExpenseRepo(mongoClient, logger)
	reportRepo := adapters.NewReportRepo(mongoClient, logger)
	categoryRepo := adapters.NewCategoryRepo(mongoClient, logger)
	rateRepo := adapters.NewExchangeRateRepo(mongoClient, logger)
	rateFetcher := adapters.NewExchangeRateFetcher(rateConfig)

	return &Application{
		Commands: Commands{
			AddExpense:         command.NewAddExpenseHandler(expenseRepo, logger),
			FetchExchangeRates: command.NewFetchExchangeRatesHandler(rateFetcher, rateRepo, logger),
		},
		Queries: Queries{
			FindExpenses: query.NewFindExpensesHandler(reportRepo, logger),
			FindCategory: query.NewFindCategoryHandler(categoryRepo, logger),
		},
		Logger: logger,
		Config: *config,
		Tracer: tracer,
	}, nil
}
