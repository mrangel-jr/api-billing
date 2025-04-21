package api

import (
	"log"
	"net/http"
	"strconv"

	// "strconv"

	"github.cm/mrangel-jr/api-billing/internals/middleware"
	"github.cm/mrangel-jr/api-billing/internals/store"
	"github.cm/mrangel-jr/api-billing/internals/utils"
)

type TenantHandler struct {
	// DB is the database connection
	tenantStore *store.MongoTenantStore
	logger      *log.Logger
}

type SummarySKUPagination struct {
	Page     int
	PageSize int
	Data     *[]store.UsageSummarySKU
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

	tenantConsume, err := th.tenantStore.GetConsumeByTenant(year, month)

	if err != nil {
		th.logger.Printf("Error getConsumeByTenant %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tenant": tenantConsume})
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

	tenantConsume, err := th.tenantStore.GetConsumeBySKU(sku, year, month)

	if err != nil {
		th.logger.Printf("Error getConsumeBySKU %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"Tenant": tenantConsume})
}

func (th *TenantHandler) GetAllConsumes(w http.ResponseWriter, r *http.Request) {
	year, month, err := utils.ReadYearMonthParam(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tenant, err := middleware.GetUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Optional: Handle pagination if needed
	page, pageSize := 1, 10

	pageStr := r.URL.Query().Get("page")         // Get the 'page' query parameter
	pageSizeStr := r.URL.Query().Get("pageSize") // Get the 'pageSize' query parameter
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid pageSize parameter", http.StatusBadRequest)
			return
		}
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	th.logger.Printf("GetConsumeBySKU %v in %v/%v \n", tenant, year, month)

	th.logger.Printf("Pagination -> Page: %v PageSize: %v \n", page, pageSize)

	tenantConsumes, err := th.tenantStore.GetAllConsumes(tenant, year, month, page, pageSize)

	if err != nil {
		th.logger.Printf("Error getConsumeBySKU %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pagination := SummarySKUPagination{
		Page:     page,
		PageSize: pageSize,
		Data:     tenantConsumes,
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"tenant": pagination})
}
