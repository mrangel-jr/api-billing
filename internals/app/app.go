package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.cm/mrangel-jr/api-billing/internals/api"
	"github.cm/mrangel-jr/api-billing/internals/store"
)

type Application struct {
	// DB is the database connection
	Logger        *log.Logger
	TenantHandler *api.TenantHandler
	DB            *store.MongoRepository
}

func NewApplication() (*Application, error) {
	logger := log.New(log.Writer(), "api-billing: ", log.LstdFlags)
	db, err := store.OpenMongoDB(context.Background(), "mongodb://localhost:27017")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Let's here my stores
	tenantStore := store.NewMongoTenantStore(db)

	// Let's here my handlers
	tenantHandler := api.NewTenantHandler(tenantStore, logger)

	return &Application{
		Logger:        logger,
		TenantHandler: tenantHandler,
		DB:            db,
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
