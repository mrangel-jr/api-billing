package store

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db *mongo.Client
}

func OpenMongoDB(ctx context.Context, uri string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &MongoRepository{db: client}, nil
}

func (m *MongoRepository) Close() {
	if err := m.db.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}
func (m *MongoRepository) GetCollection(name string) *mongo.Collection {
	// Check if the database is connected
	dbStr := os.Getenv("MONGO_DATABASE")
	if dbStr == "" {
		log.Fatalf("DATABASE_URL not set in environment")
	}
	return m.db.Database(dbStr).Collection(name)
}
