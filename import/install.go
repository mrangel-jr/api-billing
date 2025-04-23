package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mrangel-jr/api-billing/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	data.Populate()
	// MongoDB connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("magalu_billing")
	sourceCollection := db.Collection("magalu_billing_pulse")
	targetCollection := db.Collection("magalu_billing_aggregation")
	err = runAggregation(ctx, sourceCollection, targetCollection)
	if err != nil {
		log.Println("Aggregation error:", err)
	}

}

func runAggregation(ctx context.Context, source *mongo.Collection, target *mongo.Collection) error {
	// Determine time window
	now := time.Now().UTC()
	start := now.Truncate(30 * time.Minute)
	end := start.Add(30 * time.Minute)

	fmt.Printf("Aggregating between %s and %s\n", start, end)

	// Build aggregation pipeline
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "created_at", Value: bson.D{
				{Key: "$gte", Value: start},
				{Key: "$lt", Value: end},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "tenant", Value: "$tenant"},
				{Key: "product_sku", Value: "$product_sku"},
			}},
			{Key: "total_used", Value: bson.D{{Key: "$sum", Value: "$used_amount"}}},
		}}},
	}

	// Execute aggregation
	cursor, err := source.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// Prepare results to insert
	var docs []interface{}
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return err
		}
		docs = append(docs, bson.M{
			"tenant":      doc["_id"],
			"total_used":  doc["total_used"],
			"bucket_time": start,
			"created_at":  time.Now().UTC(),
		})
	}

	if len(docs) == 0 {
		fmt.Println("No data to insert")
		return nil
	}

	// Insert summary
	_, err = target.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Inserted %d summary documents\n", len(docs))
	return nil
}
