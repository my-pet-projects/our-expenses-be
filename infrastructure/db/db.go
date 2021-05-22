package db

import (
	"context"
	"time"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/logger"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CreateMongoDBPool creates connection pool for MongoDB server.
func CreateMongoDBPool(config *config.Config, appLogger *logger.AppLogger) (*mongo.Database, error) {
	ctx := context.Background()
	appLogger.Info(ctx, "Initializing MongoDB connection ...")

	// TODO: user and pass
	// https://github.com/noelchavezsimbron/otel-api-customers-go/blob/2752135b53bcec1bd7bf0a7799fdf8525bb0554b/helpers/mongodb_helper.go
	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor("our-expenses-server")
	clientOptions.ApplyURI(config.Mongo.URI)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	client, clientError := mongo.NewClient(clientOptions)
	if clientError != nil {
		appLogger.Fatal("MongoDB client error", clientError, logger.FieldsSet{})
		return nil, clientError
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()
	if connectError := client.Connect(connectCtx); connectError != nil {
		appLogger.Fatal("MongoDB connection error", connectError, logger.FieldsSet{})
		return nil, connectError
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if pingError := client.Ping(pingCtx, readpref.Primary()); pingError != nil {
		appLogger.Fatal("MongoDB ping error", pingError, logger.FieldsSet{})
		return nil, pingError
	}

	appLogger.Info(ctx, "Connected to MongoDB!")

	// go func() {
	// 	select {
	// 	case <-connectCtx.Done():
	// 		client.Disconnect(connectCtx)
	// 	}
	// }()

	database := client.Database(config.Mongo.Database)
	collections, _ := database.ListCollectionNames(context.Background(), bson.M{})

	appLogger.Infof(ctx, "MongoDB database: %s", database.Name())
	appLogger.Infof(ctx, "Available collections: %s", collections)

	return database, nil
}

// TODO: defer code execution
// https://blog.codecentric.de/en/2020/04/golang-gin-mongodb-building-microservices-easily/
