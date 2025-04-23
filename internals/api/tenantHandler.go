package api

import (
	"log"
	"net/http"

	"github.com/mrangel-jr/api-billing/internals/middleware"
	"github.com/mrangel-jr/api-billing/internals/store"
	"github.com/mrangel-jr/api-billing/internals/utils"
)

type TenantHandler struct {
	// DB is the database connection
	tenantStore *store.MongoTenantStore
	logger      *log.Logger
}

func NewTenantHandler(db *store.MongoTenantStore, logger *log.Logger) *TenantHandler {
	return &TenantHandler{
		tenantStore: db,
		logger:      logger,
	}
}

func (th *TenantHandler) GetConsumeByTenant(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ReadYearMonthParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readYearMonthParam: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	tenant, err := middleware.GetUser(r)
	if err != nil {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	if tenant == "" {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	consumeParams := store.TenantQuery{
		Tenant: tenant,
		Year:   year,
		Month:  month,
	}

	tenantConsume, err := th.tenantStore.GetConsumeByTenant(consumeParams)

	if err != nil {
		th.logger.Printf("Error getConsumeByTenant %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"payload": tenantConsume})
}

func (th *TenantHandler) GetConsumeBySKU(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ReadYearMonthParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readYearMonthParam: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	sku, err := utils.ReadSKUParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readSKUParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid sku parameter"})
		return
	}

	tenant, err := middleware.GetUser(r)
	if err != nil {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	if tenant == "" {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	consumeParams := store.TenantSKUQuery{
		Tenant:     tenant,
		ProductSKU: sku,
		Year:       year,
		Month:      month,
	}

	tenantConsume, err := th.tenantStore.GetConsumeBySKU(consumeParams)

	if err != nil {
		th.logger.Printf("Error getConsumeBySKU %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"payload": tenantConsume})
}

func (th *TenantHandler) GetAllConsumes(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ReadYearMonthParam(r)
	if err != nil {
		th.logger.Printf("ERROR: readYearMonthParam: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	tenant, err := middleware.GetUser(r)
	if err != nil {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	if tenant == "" {
		th.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be authenticated"})
		return
	}

	page, pageSize := utils.Pagination(r)

	consumeParams := store.TenantQueryPagination{
		Tenant:   tenant,
		Year:     year,
		Month:    month,
		Page:     page,
		PageSize: pageSize,
	}

	tenantConsumes, err := th.tenantStore.GetAllConsumes(consumeParams)

	if err != nil {
		th.logger.Printf("Error getAllConsumes %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	dataPagination := store.SummarySKUPagination{
		Page:     page,
		PageSize: pageSize,
		Data:     tenantConsumes,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"payload": dataPagination})
}
