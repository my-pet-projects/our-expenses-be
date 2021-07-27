package app

import (
	"context"

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
	MoveCategory   command.MoveCategoryHandlerInterface
}

// Queries struct holds available application queries.
type Queries struct {
	FindCategories     query.FindCategoriesHandlerInterface
	FindCategory       query.FindCategoryHandlerInterface
	FindCategoryUsages query.FindCategoryUsagesHandlerInterface
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
	categoryRepo := repository.NewCategoryRepo(mongoClient, logger)

	return &Application{
		Commands: Commands{
			AddCategory:    command.NewAddCategoryHandler(categoryRepo, logger),
			UpdateCategory: command.NewUpdateCategoryHandler(categoryRepo, logger),
			DeleteCategory: command.NewDeleteCategoryHandler(categoryRepo, logger),
			MoveCategory:   command.NewMoveCategoryHandler(categoryRepo, logger),
		},
		Queries: Queries{
			FindCategories:     query.NewFindCategoriesHandler(categoryRepo, logger),
			FindCategory:       query.NewFindCategoryHandler(categoryRepo, logger),
			FindCategoryUsages: query.NewFindCategoryUsagesHandler(categoryRepo, logger),
		},
		Logger: logger,
		Config: *config,
		Tracer: tracer,
	}, nil
}
