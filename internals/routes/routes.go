package routes

import (
	"github.cm/mrangel-jr/api-billing/internals/app"
	"github.com/go-chi/chi"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/consume/{year_month}/{sku}", app.TenantHandler.GetConsumeBySKU)
	r.Get("/consume/{year_month}", app.TenantHandler.GetConsumeByTenant)

	return r
}
