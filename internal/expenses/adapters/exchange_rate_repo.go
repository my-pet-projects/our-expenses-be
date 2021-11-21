package adapters

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/expenses/domain"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/database"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/logger"
	"dev.azure.com/filimonovga/our-expenses/our-expenses-server/pkg/tracer"
)

type rateDbModel struct {
	ID    int64              `bson:"_id,omitempty"`
	Date  time.Time          `bson:"date"`
	Base  string             `bson:"base"`
	Rates map[string]float32 `bson:"rates"`
}

// ExchangeRateRepository represents a struct to access exchange rates MongoDB collection.
type ExchangeRateRepository struct {
	client *database.MongoClient
	logger logger.LogInterface
}

// ExchangeRateRepoInterface defines a contract to persist exchange rates in the database.
type ExchangeRateRepoInterface interface {
	InsertAll(ctx context.Context, rates []domain.ExchangeRate) (*domain.InsertResult, error)
	GetAll(ctx context.Context, dateRange domain.DateRange) ([]domain.ExchangeRate, error)
}

// NewExchangeRateRepo returns repository.
func NewExchangeRateRepo(client *database.MongoClient, logger logger.LogInterface) *ExchangeRateRepository {
	return &ExchangeRateRepository{
		logger: logger,
		client: client,
	}
}

// collection returns collection handle.
func (r *ExchangeRateRepository) collection() *mongo.Collection {
	return r.client.Collection("exchangeRates")
}

// InsertAll insert all rates into the database.
func (r *ExchangeRateRepository) InsertAll(
	ctx context.Context,
	rates []domain.ExchangeRate,
) (*domain.InsertResult, error) {
	ctx, span := tracer.NewSpan(ctx, "insert exchange rates to the database")
	defer span.End()

	operations := make([]mongo.WriteModel, 0)
	for _, rate := range rates {
		dbModel := r.marshalRate(rate)
		insOperation := mongo.NewUpdateOneModel()
		insOperation.SetFilter(bson.M{"date": dbModel.Date})
		insOperation.SetUpdate(bson.M{"$set": dbModel})
		insOperation.SetUpsert(true)
		operations = append(operations, insOperation)
	}

	insRes, insErr := r.collection().BulkWrite(ctx, operations)
	if insErr != nil {
		tracer.AddSpanError(span, insErr)
		return nil, errors.Wrap(insErr, "mongodb bulk write rates")
	}

	result := &domain.InsertResult{
		InsertCount: int(insRes.InsertedCount),
	}

	return result, nil
}

// GetAll fetches exchange rates from the database.
func (r *ExchangeRateRepository) GetAll(
	ctx context.Context,
	dateRange domain.DateRange,
) ([]domain.ExchangeRate, error) {
	ctx, span := tracer.NewSpan(ctx, "fetch exchange rates from the database")
	defer span.End()

	filter := bson.M{
		"date": bson.M{
			"$gte": dateRange.From(),
			"$lte": dateRange.To(),
		},
	}

	find, findErr := r.collection().Find(ctx, filter)
	if findErr != nil {
		tracer.AddSpanError(span, findErr)
		return nil, errors.Wrap(findErr, "mongo find exchange rates")
	}

	var rateDbModels []rateDbModel
	if cursorErr := find.All(ctx, &rateDbModels); cursorErr != nil {
		tracer.AddSpanError(span, cursorErr)
		return nil, errors.Wrap(cursorErr, "cursor iteration")
	}

	rates := []domain.ExchangeRate{}
	for _, rateDbModel := range rateDbModels {
		rate := r.unmarshalRate(rateDbModel)
		rates = append(rates, rate)
	}

	return rates, nil
}

func (r ExchangeRateRepository) marshalRate(exchangeRate domain.ExchangeRate) rateDbModel {
	rates := make(map[string]float32, len(exchangeRate.Rates()))
	for currency, rate := range exchangeRate.Rates() {
		rates[string(currency)], _ = rate.BigFloat().Float32()
	}
	dbModel := rateDbModel{
		ID:    exchangeRate.Date().Unix(),
		Date:  exchangeRate.Date(),
		Base:  string(exchangeRate.BaseCurrency()),
		Rates: rates,
	}
	return dbModel
}

func (r ExchangeRateRepository) unmarshalRate(rateDbModel rateDbModel) domain.ExchangeRate {
	rate := domain.NewExchageRate(rateDbModel.Date, rateDbModel.Base, rateDbModel.Rates)
	return rate
}
