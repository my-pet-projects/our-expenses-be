package database

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/config"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
)

// MongoClient provides a MongoDB client.
type MongoClient struct {
	client *mongo.Client
	logger logger.LogInterface
	config config.Database
}

// NewMongoClient creates a MongoDB client.
func NewMongoClient(log logger.LogInterface, config config.Database) (*MongoClient, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(config.Mongo.URI)
	clientOptions.Auth = &options.Credential{
		Username: config.Mongo.User,
		Password: config.Mongo.Pass,
	}
	clientOptions.SetAppName(config.Mongo.Name)
	clientOptions.SetReadConcern(readconcern.Majority())
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	clientOptions.Monitor = otelmongo.NewMonitor()

	client, clientErr := mongo.NewClient(clientOptions)
	if clientErr != nil {
		return nil, errors.Wrap(clientErr, "mongodb client")
	}

	mongoClient := &MongoClient{
		client: client,
		logger: log,
		config: config,
	}

	return mongoClient, nil
}

// OpenConnection creates a connection to MongoDB cluster.
// TODO: add graceful shutdown mechanism.
func (c MongoClient) OpenConnection(ctx context.Context, cancel context.CancelFunc) error {
	c.logger.Info(ctx, "Initializing MongoDB connection ...")

	connectCtx, connectCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectCancel()
	if connectErr := c.client.Connect(connectCtx); connectErr != nil {
		return errors.Wrap(connectErr, "connection failed")
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if pingErr := c.client.Ping(pingCtx, readpref.Primary()); pingErr != nil {
		return errors.Wrap(pingErr, "cluster ping")
	}

	c.logger.Info(ctx, "Connected to MongoDB!")

	return nil
}

// Database returns MongoDB database handler.
func (c MongoClient) Database() *mongo.Database {
	return c.client.Database(c.config.Mongo.Database)
}

// Collection returns MongoDB collection handler.
func (c MongoClient) Collection(name string) *mongo.Collection {
	// TODO: handle non-existing collections
	return c.Database().Collection(name)
}

// ListCollections returns all available MongoDB collections.
func (c MongoClient) ListCollections(ctx context.Context) ([]string, error) {
	col, colErr := c.Database().ListCollectionNames(ctx, bson.M{})
	if colErr != nil {
		return nil, colErr
	}
	return col, nil
}
