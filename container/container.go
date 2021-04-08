//+build wireinject

package container

import (
	"our-expenses-server/api"
	"our-expenses-server/api/handler"
	"our-expenses-server/api/router"
	"our-expenses-server/config"
	"our-expenses-server/infrastructure/db"
	"our-expenses-server/infrastructure/db/repository"
	"our-expenses-server/logger"
	"our-expenses-server/service/category"
	"our-expenses-server/validator"

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
