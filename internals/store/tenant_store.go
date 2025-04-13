package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ProductSKUInfo struct {
	ProductSKU string `json:"product_sku"`
	TotalUsed  int64  `json:"total_used"`
}

type TenantGroup struct {
	ID         int64  `json:"-" bson:"-"`
	TenantName string `json:"tenant" bson:"tenant"`
	ProductSKU string `json:"product_sku" bson:"product_sku"`
}

type Tenant struct {
	ID         string      `json:"_id" bson:"_id"`
	TotalUsed  int64       `json:"total_used" bson:"total_used"`
	UseUnity   string      `json:"use_unity" bson:"use_unity"`
	BucketTime time.Time   `json:"bucket_time" bson:"bucket_time"`
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	Tenant     TenantGroup `json:"tenant"`
}

type TenantPulse struct {
	ProductSKU string    `json:"product_sku" bson:"product_sku"`
	TotalUsed  int64     `json:"total_used" bson:"total_used"`
	UseUnity   string    `json:"use_unity" bson:"use_unity"`
	UseValue   int64     `json:"use_value" bson:"use_value"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	TenantName string    `json:"tenant" bson:"tenant"`
	ID         string    `json:"_id" bson:"_id"`
}

type MongoTenantStore struct {
	db *MongoRepository
}

func NewMongoTenantStore(client *MongoRepository) *MongoTenantStore {
	return &MongoTenantStore{db: client}
}

type TenantStore interface {
	GetConsumeBySKU(sku string, month int, year int) (*ProductSKUInfo, error)
	GetConsumeByTenant(month int, year int) (*Tenant, error)
}

func (m *MongoTenantStore) GetConsumeBySKU(sku string, month int, year int) (*[]Tenant, error) {
	collection := m.db.GetCollection("magalu_billing_aggragation")
	fmt.Printf("Collection Name: %v\n", collection.Name())
	// filter := map[string]interface{}{"tenant.product_sku": `${sku}`, "bucket_time": `${year}-${month}-01T00:00:00Z`}
	filter := bson.M{"tenant.tenant": "tenant_j"}
	// fmt.Printf("Product SKU Filter: %v\n", collection.FindOne(context.Background(), filter))
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var result []Tenant
	for cursor.Next(context.Background()) {
		var tenant Tenant
		if err := cursor.Decode(&tenant); err != nil {
			return nil, err
		}
		result = append(result, tenant)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no documents found")
	}

	return &result, nil

}

// func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
// 	tx, err := pg.db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	query := `INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
// 		VALUES ($1, $2,$3, $4, $5)
// 		RETURNING id`

// 	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, entry := range workout.Entries {
// 		query := `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps,duration_seconds, weight, notes, order_index)
// 		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
// 		RETURNING id`
// 		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return workout, nil
// }
