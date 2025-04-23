package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mrangel-jr/api-billing/internals/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsageSummary struct {
	Tenant    string `bson:"_id"`
	TotalUsed int64  `bson:"total_used"`
}
type UsageSummarySKU struct {
	ID struct {
		Tenant     string `bson:"tenant"`
		ProductSKU string `bson:"product_sku"`
	} `bson:"_id"`
	TotalUsed int64 `bson:"total_used"`
}

type SummarySKUPagination struct {
	Page     int
	PageSize int
	Data     *[]UsageSummarySKU
}

type TenantQuery struct {
	Tenant string
	Year   int
	Month  int
}
type TenantSKUQuery struct {
	Tenant     string
	ProductSKU string
	Year       int
	Month      int
}

type TenantQueryPagination struct {
	Tenant   string
	Year     int
	Month    int
	Page     int
	PageSize int
}

type MongoTenantStore struct {
	db         *MongoRepository
	collection string
}

func NewMongoTenantStore(client *MongoRepository) *MongoTenantStore {
	magalu_collection := os.Getenv("MONGO_COLLECTION_AGGREGATION")
	if magalu_collection == "" {
		log.Fatalf("MONGO_COLLECTION_AGGREGATION not set in environment")
	}
	return &MongoTenantStore{db: client, collection: magalu_collection}
}

type TenantStore interface {
	GetConsumeByTenant(consumeParams TenantQuery) (*UsageSummary, error)
	GetConsumeBySKU(consumeParams TenantSKUQuery) (*UsageSummarySKU, error)
	GetAllConsumes(consumeParams TenantQueryPagination) (*[]UsageSummarySKU, error)
}

func (m *MongoTenantStore) GetConsumeByTenant(consumeParams TenantQuery) (*UsageSummary, error) {
	collection := m.db.GetCollection(m.collection)
	fmt.Printf("Collection Name: %v\n", collection.Name())
	startDate := time.Date(consumeParams.Year, time.Month(consumeParams.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.LastDayOfMonth(consumeParams.Year, time.Month(consumeParams.Month))

	filterSummary := bson.D{
		{Key: "tenant.tenant", Value: consumeParams.Tenant},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	groupStage := bson.D{
		{Key: "_id", Value: "$tenant.tenant"},
		{Key: "total_used", Value: bson.D{
			{Key: "$sum", Value: "$total_used"},
		}},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterSummary}},
		{{Key: "$group", Value: groupStage}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var result UsageSummary
	if cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil

}

func (m *MongoTenantStore) GetConsumeBySKU(consumeParams TenantSKUQuery) (*UsageSummarySKU, error) {
	collection := m.db.GetCollection(m.collection)
	startDate := time.Date(consumeParams.Year, time.Month(consumeParams.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.LastDayOfMonth(consumeParams.Year, time.Month(consumeParams.Month))

	ctx := context.Background()

	filterSummary := bson.D{
		{Key: "tenant.tenant", Value: consumeParams.Tenant},
		{Key: "tenant.product_sku", Value: consumeParams.ProductSKU},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	groupStage := bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "tenant", Value: "$tenant.tenant"},
			{Key: "product_sku", Value: "$tenant.product_sku"},
		}},
		{Key: "total_used", Value: bson.D{
			{Key: "$sum", Value: "$total_used"},
		}},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterSummary}},
		{{Key: "$group", Value: groupStage}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result UsageSummarySKU

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no results found for SKU %s", consumeParams.ProductSKU)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil

}

func (m *MongoTenantStore) GetAllConsumes(consumeParams TenantQueryPagination) (*[]UsageSummarySKU, error) {
	collection := m.db.GetCollection(m.collection)
	startDate := time.Date(consumeParams.Year, time.Month(consumeParams.Month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.LastDayOfMonth(consumeParams.Year, time.Month(consumeParams.Month))
	skip := (consumeParams.Page - 1) * consumeParams.PageSize

	ctx := context.Background()

	filterSummary := bson.D{
		{Key: "tenant.tenant", Value: consumeParams.Tenant},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	groupStage := bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "tenant", Value: "$tenant.tenant"},
			{Key: "product_sku", Value: "$tenant.product_sku"},
		}},
		{Key: "total_used", Value: bson.D{
			{Key: "$sum", Value: "$total_used"},
		}},
		{Key: "count", Value: bson.D{
			{Key: "$sum", Value: 1},
		}},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filterSummary}},
		{{Key: "$group", Value: groupStage}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: consumeParams.PageSize}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []UsageSummarySKU

	for cursor.Next(ctx) {
		var item UsageSummarySKU
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no results found for tenant %s", consumeParams.Tenant)
	}

	return &result, nil

}
