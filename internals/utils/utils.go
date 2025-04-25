package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type Envelope map[string]interface{}

type Pagination struct {
	Page     int
	PageSize int
}

func ReadYearMonthParam(r *http.Request) (int, int, error) {
	yearMonthParam := chi.URLParam(r, "year_month")
	if yearMonthParam == "" {
		return 0, 0, errors.New("invalid year_month parameter")
	}
	_, err := time.Parse("200601", yearMonthParam)
	if err != nil {
		return 0, 0, errors.New("invalid year_month parameter type")
	}

	yearStr, monthStr := yearMonthParam[:4], yearMonthParam[4:]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, errors.New("invalid year parameter type")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, 0, errors.New("invalid month parameter type")
	}
	return year, month, nil
}

func LastDayOfMonth(year int, month time.Month) time.Time {
	// Go trick: the 0th day of the next month is the last day of the current month
	return time.Date(year, month+1, 0, 23, 59, 59, 0, time.UTC)
}

func ReadSKUParam(r *http.Request) (string, error) {
	skuParam := chi.URLParam(r, "sku")
	if skuParam == "" {
		return "", errors.New("invalid year parameter")
	}
	return skuParam, nil
}
func ReadTenantParam(r *http.Request) (string, error) {
	tenantParam := chi.URLParam(r, "tenant")
	if tenantParam == "" {
		return "", errors.New("invalid year parameter")
	}
	return tenantParam, nil
}

func WriteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func CreatePagination(r *http.Request) Pagination {
	var pagination Pagination
	pagination.Page, pagination.PageSize = 1, 10

	pageStr := r.URL.Query().Get("page")         // Get the 'page' query parameter
	pageSizeStr := r.URL.Query().Get("pageSize") // Get the 'pageSize' query parameter
	if pageSizeStr != "" {
		pagination.PageSize, _ = strconv.Atoi(pageSizeStr)
		if pagination.PageSize < 1 {
			pagination.PageSize = 10
		}
	}
	if pageStr != "" {
		pagination.Page, _ = strconv.Atoi(pageStr)
		if pagination.Page < 1 {
			pagination.Page = 1
		}
	}
	return pagination
}
