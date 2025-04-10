package data

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tenants = []string{"tenant_a", "tenant_b", "tenant_c", "tenant_d", "tenant_e", "tenant_f", "tenant_g", "tenant_h", "tenant_i", "tenant_j"}

func Populate() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("magalu_billing").Collection("magalu_billing_pulse")

	// Generate and insert 1 million documents in batches
	const total = 1_000_000
	const batchSize = 10000

	fmt.Printf("Inserting %d mock documents...\n", total)

	for i := 0; i < total; i += batchSize {
		var batch []interface{}

		ini := 300

		for j := 0; j < batchSize; j++ {
			doc := bson.M{
				"tenant":      tenants[rand.Intn(len(tenants))],
				"product_sku": fmt.Sprintf("SKU-%03d", rand.Intn(100)),
				"used_amount": rand.Intn(batchSize-ini+j) + ini,
				"use_unity":   "KB",
				"created_at":  time.Now().Add(-time.Duration(rand.Intn(90)) * time.Minute), // last 90 min
			}
			batch = append(batch, doc)
		}

		if _, err := collection.InsertMany(ctx, batch); err != nil {
			panic(err)
		}

		fmt.Printf("Inserted %d...\n", i+batchSize)
	}

	fmt.Println("✅ Done!")
}
