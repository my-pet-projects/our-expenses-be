package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	categoriesApp "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app"
	categoriesPorts "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/ports"
	expensesApp "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/app"
	expensesPorts "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/ports"
	usersApp "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/app"
	usersPorts "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/ports"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/server"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	appConfig, appConfigErr := config.NewConfig()
	if appConfigErr != nil {
		log.Printf("Failed to create application config: %v", appConfigErr)
		os.Exit(1)
	}
	if cfgValidErr := appConfig.Validate(); cfgValidErr != nil {
		log.Printf("Invalid application config: %v", cfgValidErr)
		os.Exit(1)
	}
	appLogger, appLoggerErr := logger.NewLogger(appConfig.Logger)
	if appLoggerErr != nil {
		log.Printf("Failed to create application logger: %v", appLoggerErr)
		os.Exit(1)
	}
	appTracer := tracer.NewTracer(appConfig.Telemetry)
	defer appTracer.Shutdown()

	appLogger.Info(ctx, "Application starting ...")

	mongoClient, mongoClientErr := database.NewMongoClient(appLogger, appConfig.Database)
	if mongoClientErr != nil {
		appLogger.Error(ctx, "Failed to create MongoDB client!", mongoClientErr)
		os.Exit(1)
	}
	if mongoConErr := mongoClient.OpenConnection(ctx, cancel); mongoConErr != nil {
		appLogger.Error(ctx, "Failed to open MongoDB connection!", mongoConErr)
		os.Exit(1)
	}

	categoriesApp, categoriesAppErr := categoriesApp.NewApplication(ctx, cancel, appConfig, appLogger, appTracer, mongoClient)
	if categoriesAppErr != nil {
		appLogger.Error(ctx, "Failed to instantiate Categories application!", categoriesAppErr)
		os.Exit(1)
	}
	expensesApp, expensesAppErr := expensesApp.NewApplication(ctx, cancel, appConfig, appLogger, appTracer, mongoClient)
	if expensesAppErr != nil {
		appLogger.Error(ctx, "Failed to instantiate Expenses application!", expensesAppErr)
		os.Exit(1)
	}
	usersApp, usersAppErr := usersApp.NewApplication(ctx, cancel, appConfig, appLogger, appTracer, mongoClient)
	if usersAppErr != nil {
		appLogger.Error(ctx, "Failed to instantiate Users application!", usersAppErr)
		os.Exit(1)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		signal := <-signChan
		categoriesApp.Logger.Warnf(ctx, "Received interuption signal: %s", signal)
		cancel()
	}()

	server := server.NewServer(categoriesApp.Logger, categoriesApp.Config.Server,
		func(router *echo.Echo) {
			categoriesPorts.RegisterHandlersWithBaseURL(router, categoriesPorts.NewHTTPServer(categoriesApp), "api")
			expensesPorts.RegisterHandlersWithBaseURL(router, expensesPorts.NewHTTPServer(expensesApp), "api")
			usersPorts.RegisterHandlersWithBaseURL(router, usersPorts.NewHTTPServer(usersApp), "api")
		})
	if err := server.Start(ctx); err != nil {
		categoriesApp.Logger.Error(ctx, "Failed to start HTTP server", err)
		os.Exit(1)
	}

	categoriesApp.Tracer.Shutdown()
	categoriesApp.Logger.Info(ctx, "Application shutdown")

	defer os.Exit(0)
}
