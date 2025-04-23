package routes

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mrangel-jr/api-billing/internals/app"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.JWTMiddleware)
		r.Use(middleware.Logger)

		r.Get("/consume/{year_month}/{sku}", app.TenantHandler.GetConsumeBySKU)
		r.Get("/consume/{year_month}/summary", app.TenantHandler.GetConsumeByTenant)
		r.Get("/consume/{year_month}", app.TenantHandler.GetAllConsumes)

	})
	r.Get("/health", app.HealthCheck)

	return r
}
