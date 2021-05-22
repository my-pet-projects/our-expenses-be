//+build wireinject

package container

import (
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/api"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/api/handler"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/api/router"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/infrastructure/db"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/infrastructure/db/repository"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/service/category"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/validator"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateServer(database *mongo.Database) (*api.Server, error) {
	wire.Build(
		logger.ProvideLogger,
		config.ProvideConfiguration,
		api.ProvideServer,
		router.ProvideRouter,
		validator.ProvideValidator,
		handler.ProvideCategoryController,
		category.ProvideCategoryService,
		repository.ProvideCategoryRepo,
	)

	return &api.Server{}, nil
}

func InitDatabase() (*mongo.Database, error) {
	wire.Build(config.ProvideConfiguration, logger.ProvideLogger, db.CreateMongoDBPool)

	return &mongo.Database{}, nil
}
