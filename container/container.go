//+build wireinject

package container

import (
	"our-expenses-server/config"
	"our-expenses-server/db"
	"our-expenses-server/db/repositories"
	"our-expenses-server/logger"
	"our-expenses-server/validators"
	"our-expenses-server/web/api"
	"our-expenses-server/web/api/controllers"
	"our-expenses-server/web/server"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateServer(database *mongo.Database) (*server.Server, error) {
	wire.Build(
		logger.ProvideLogger,
		config.ProvideConfiguration,
		server.ProvideServer,
		api.ProvideRouter,
		controllers.ProvideCategoryController,
		repositories.ProvideCategoryRepository,
		validators.ProvideValidator)

	return &server.Server{}, nil
}

func InitDatabase() (*mongo.Database, error) {
	wire.Build(config.ProvideConfiguration, logger.ProvideLogger, db.CreateMongoDBPool)

	return &mongo.Database{}, nil
}
