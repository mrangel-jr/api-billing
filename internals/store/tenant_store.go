package store

import (
	"context"
	"fmt"
	"time"

	"github.cm/mrangel-jr/api-billing/internals/utils"
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

type MongoTenantStore struct {
	db *MongoRepository
}

func NewMongoTenantStore(client *MongoRepository) *MongoTenantStore {
	return &MongoTenantStore{db: client}
}

type TenantStore interface {
	GetConsumeByTenant(year int, month int) (*UsageSummary, error)
	GetConsumeBySKU(sku string, month int, year int) (*UsageSummarySKU, error)
}

func (m *MongoTenantStore) GetConsumeByTenant(year int, month int) (*UsageSummary, error) {
	collection := m.db.GetCollection("magalu_billing_aggragation")
	fmt.Printf("Collection Name: %v\n", collection.Name())
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.LastDayOfMonth(year, time.Month(month))
	tenantName := "tenant_j"

	filterSummary := bson.D{
		{Key: "tenant.tenant", Value: tenantName},
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

func (m *MongoTenantStore) GetConsumeBySKU(sku string, year int, month int) (*UsageSummarySKU, error) {
	collection := m.db.GetCollection("magalu_billing_aggragation")
	fmt.Printf("Collection Name: %v\n", collection.Name())
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := utils.LastDayOfMonth(year, time.Month(month))
	tenantName := "tenant_j"

	ctx := context.Background()

	filterSummary := bson.D{
		{Key: "tenant.tenant", Value: tenantName},
		{Key: "tenant.product_sku", Value: sku},
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
		return nil, fmt.Errorf("no results found for SKU %s", sku)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &result, nil

}
