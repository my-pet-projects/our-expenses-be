package app

import (
	"context"

	"github.com/pkg/errors"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/command"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/app/query"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/categories/repository"
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
	AddCategory    command.AddCategoryHandlerInterface
	UpdateCategory command.UpdateCategoryHandlerInterface
	DeleteCategory command.DeleteCategoryHandlerInterface
}

// Queries struct holds available application queries.
type Queries struct {
	FindCategories query.FindCategoriesHandlerInterface
	FindCategory   query.FindCategoryHandlerInterface
}

// NewApplication returns application instance.
func NewApplication(ctx context.Context, cancel context.CancelFunc) (*Application, error) {
	cfg, cfgErr := config.NewConfig()
	if cfgErr != nil {
		return nil, errors.Wrap(cfgErr, "create application config")
	}
	if cfgValidErr := cfg.Validate(); cfgValidErr != nil {
		return nil, errors.Wrap(cfgValidErr, "validate application config")
	}

	log, logErr := logger.NewLogger(cfg.Logger)
	if logErr != nil {
		return nil, errors.Wrap(logErr, "create logger")
	}

	tracer := tracer.NewTracer(cfg.Telemetry)

	log.Info(ctx, "Application starting ...")

	mongoClient, mongoClientErr := database.NewMongoClient(log, cfg.Database)
	if mongoClientErr != nil {
		return nil, errors.Wrap(mongoClientErr, "create mongodb client")
	}
	if mongoConErr := mongoClient.OpenConnection(ctx, cancel); mongoConErr != nil {
		return nil, errors.Wrap(mongoConErr, "open mongodb connection")
	}

	categoryRepo := repository.NewCategoryRepo(mongoClient, log)

	return &Application{
		Commands: Commands{
			AddCategory:    command.NewAddCategoryHandler(categoryRepo, log),
			UpdateCategory: command.NewUpdateCategoryHandler(categoryRepo, log),
			DeleteCategory: command.NewDeleteCategoryHandler(categoryRepo, log),
		},
		Queries: Queries{
			FindCategories: query.NewFindCategoriesHandler(categoryRepo, log),
			FindCategory:   query.NewFindCategoryHandler(categoryRepo, log),
		},
		Logger: log,
		Config: *cfg,
		Tracer: tracer,
	}, nil
}
