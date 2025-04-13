package store

import (
	"context"

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
	return m.db.Database("magalu_billing").Collection(name)
}
