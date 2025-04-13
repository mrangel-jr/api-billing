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

func (th *TenantHandler) GetConsumeBySKU(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to get consume by SKU
	year, err := utils.ReadYearParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	month, err := utils.ReadMonthParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sku, err := utils.ReadSKUParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	th.logger.Printf("GetConsumeBySKU in %v/%v called with SKU: %v\n", year, month, sku)

	tenant, err := th.tenantStore.GetConsumeBySKU(sku, month, year)

	if err != nil {
		th.logger.Printf("Error getting consume by SKU: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tenant": tenant})
}
