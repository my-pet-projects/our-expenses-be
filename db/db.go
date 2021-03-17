package db

import (
	"context"
	"fmt"
	"our-expenses-server/config"
	"our-expenses-server/logger"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CreateMongoDBPool creates connection pool for MongoDB server.
func CreateMongoDBPool(config *config.Config, appLogger *logger.AppLogger) (*mongo.Database, error) {
	appLogger.Info("Initializing MongoDB connection ...", logger.Fields{})

	clientOptions := options.Client().ApplyURI(config.Mongo.URI)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	client, clientError := mongo.NewClient(clientOptions)
	if clientError != nil {
		appLogger.Fatal("MongoDB client error", clientError, logger.Fields{})
		return nil, clientError
	}

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()
	if connectError := client.Connect(connectCtx); connectError != nil {
		appLogger.Fatal("MongoDB connection error", connectError, logger.Fields{})
		return nil, connectError
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if pingError := client.Ping(pingCtx, readpref.Primary()); pingError != nil {
		appLogger.Fatal("MongoDB ping error", pingError, logger.Fields{})
		return nil, pingError
	}

	appLogger.Info("Connected to MongoDB!", logger.Fields{})

	// go func() {
	// 	select {
	// 	case <-connectCtx.Done():
	// 		client.Disconnect(connectCtx)
	// 	}
	// }()

	database := client.Database(config.Mongo.Database)
	collections, _ := database.ListCollectionNames(context.Background(), bson.M{})

	appLogger.Info(fmt.Sprintf("MongoDB database: %s", database.Name()), logger.Fields{})
	appLogger.Info(fmt.Sprintf("Available collections: %s", collections), logger.Fields{})

	return database, nil
}

// TODO: defer code execution
// https://blog.codecentric.de/en/2020/04/golang-gin-mongodb-building-microservices-easily/
