package api

import (
	"log"
	"net/http"

	"github.cm/mrangel-jr/api-billing/internals/store"
	"github.cm/mrangel-jr/api-billing/internals/utils"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	th.logger.Printf("GetConsumeByTenant in %v/%v \n", year, month)

	tenant, err := th.tenantStore.GetConsumeByTenant(year, month)

	if err != nil {
		th.logger.Printf("Error getConsumeByTenant %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tenant": tenant})
}

func (th *TenantHandler) GetConsumeBySKU(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ReadYearMonthParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sku, err := utils.ReadSKUParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	th.logger.Printf("GetConsumeBySKU %v in %v/%v \n", sku, year, month)

	tenant, err := th.tenantStore.GetConsumeBySKU(sku, year, month)

	if err != nil {
		th.logger.Printf("Error getConsumeBySKU %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"Tenant": tenant})
}
