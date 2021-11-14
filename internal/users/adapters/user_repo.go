package adapters

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

const usersCollectionName string = "users"

// userDbModel defines user structure in MongoDB.
type userDbModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	Password     string             `bson:"password"`
	Token        string             `bson:"token"`
	RefreshToken string             `bson:"refreshToken"`
	CreatedAt    time.Time          `bson:"createdAt,omitempty"`
	CreatedBy    string             `bson:"createdBy,omitempty"`
	UpdatedAt    *time.Time         `bson:"updatedAt,omitempty"`
	UpdatedBy    *string            `bson:"updatedBy,omitempty"`
}

// UserRepository represents a struct to access users MongoDB collection.
type UserRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// UserRepoInterface defines a contract to persist users in the database.
type UserRepoInterface interface {
	GetOne(ctx context.Context, username string) (*domain.User, error)
	Insert(ctx context.Context, user *domain.User) (*domain.InsertResult, error)
	Update(ctx context.Context, user *domain.User) (*domain.UpdateResult, error)
}

// NewUserRepo returns a UserRepository.
func NewUserRepo(client *database.MongoClient, logger logger.LogInterface) *UserRepository {
	return &UserRepository{
		logger: logger,
		client: client,
	}
}

// collection returns collection handle.
func (r *UserRepository) collection() *mongo.Collection {
	return r.client.Collection(usersCollectionName)
}

// GetOne returns a single user from the database.
func (r *UserRepository) GetOne(ctx context.Context, username string) (*domain.User, error) {
	ctx, span := tracer.NewSpan(ctx, "find user in the database")
	tracer.AddSpanTags(span, map[string]string{"username": username})
	defer span.End()

	filter := bson.M{"username": username}
	userDbModel := userDbModel{}
	findErr := r.collection().FindOne(ctx, filter).Decode(&userDbModel)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		tracer.AddSpanError(span, findErr)
		return nil, errors.Wrap(findErr, "mongodb find user")
	}

	user, userErr := r.unmarshalUser(userDbModel)
	if userErr != nil {
		tracer.AddSpanError(span, userErr)
		return nil, errors.Wrap(userErr, "unmarshall user")
	}

	return user, nil
}

// Insert inserts a user into the database.
func (r *UserRepository) Insert(ctx context.Context, user *domain.User) (*domain.InsertResult, error) {
	ctx, span := tracer.NewSpan(ctx, "add user to the database")
	defer span.End()

	userDbModel := r.marshalUser(user)
	userDbModel.CreatedBy = domain.SystemUser
	userDbModel.CreatedAt = time.Now()

	insRes, insErr := r.collection().InsertOne(ctx, userDbModel)
	if insErr != nil {
		tracer.AddSpanError(span, insErr)
		return nil, errors.Wrap(insErr, "mongodb insert user")
	}
	objID, _ := insRes.InsertedID.(primitive.ObjectID)

	result := &domain.InsertResult{
		ID: objID.Hex(),
	}

	return result, nil
}

// Update updates a user in the database.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.UpdateResult, error) {
	ctx, span := tracer.NewSpan(ctx, "update user in the database")
	defer span.End()

	userDbModel := r.marshalUser(user)
	sysUser := domain.SystemUser
	now := time.Now()
	userDbModel.UpdatedBy = &sysUser
	userDbModel.UpdatedAt = &now

	filter := bson.M{"_id": userDbModel.ID}
	updater := bson.M{"$set": userDbModel}

	opts := options.Update().SetUpsert(false)

	updResult, updErr := r.collection().UpdateOne(ctx, filter, updater, opts)
	if updErr != nil {
		return nil, errors.Wrap(updErr, "mongodb update user")
	}

	result := &domain.UpdateResult{
		UpdateCount: int(updResult.ModifiedCount),
	}

	return result, nil
}

func (r UserRepository) marshalUser(user *domain.User) userDbModel {
	id, _ := primitive.ObjectIDFromHex(user.ID())
	return userDbModel{
		ID:           id,
		Username:     user.Username(),
		Password:     user.Password(),
		Token:        user.Token(),
		RefreshToken: user.RefreshToken(),
	}
}

func (r UserRepository) unmarshalUser(userDbModel userDbModel) (*domain.User, error) {
	user, userErr := domain.NewUser(userDbModel.ID.Hex(), userDbModel.Username, userDbModel.Password,
		userDbModel.Token, userDbModel.RefreshToken)
	return user, userErr
}
