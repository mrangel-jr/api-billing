package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrangel-jr/api-billing/internals/api"
	"github.com/mrangel-jr/api-billing/internals/middleware"
	"github.com/mrangel-jr/api-billing/internals/store"
)

type Application struct {
	// DB is the database connection
	Logger        *log.Logger
	TenantHandler *api.TenantHandler
	Middleware    middleware.AuthMiddleware
	DB            *store.MongoRepository
}

func NewApplication() (*Application, error) {
	logger := log.New(log.Writer(), "api-billing: ", log.LstdFlags)
	dbConnStr := os.Getenv("MONGO_URL")
	if dbConnStr == "" {
		logger.Fatalf("MONGO_URL not set in environment")
	}
	db, err := store.OpenMongoDB(context.Background(), dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Let's here my stores
	tenantStore := store.NewMongoTenantStore(db)

	// Let's here my handlers
	tenantHandler := api.NewTenantHandler(tenantStore, logger)
	middlewareHandler := middleware.AuthMiddleware{}

	return &Application{
		Logger:        logger,
		TenantHandler: tenantHandler,
		Middleware:    middlewareHandler,
		DB:            db,
	}, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
