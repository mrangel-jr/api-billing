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

func ReadIDParam(r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}
	return id, nil
}

func ReadYearParam(r *http.Request) (int, error) {
	yearParam := chi.URLParam(r, "year")
	if yearParam == "" {
		return 0, errors.New("invalid year parameter")
	}
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return 0, errors.New("invalid year parameter type")
	}

	actualYear := time.Now().Year()

	if year < 1900 || year > actualYear {
		return 0, errors.New("invalid year parameter value")
	}
	return year, nil
}

func ReadMonthParam(r *http.Request) (int, error) {
	monthParam := chi.URLParam(r, "month")
	if monthParam == "" {
		return 0, errors.New("invalid year parameter")
	}
	month, err := strconv.Atoi(monthParam)
	if err != nil {
		return 0, errors.New("invalid month parameter type")
	}

	if month < 1 || month > 12 {
		return 0, errors.New("invalid month parameter value")
	}
	return month, nil
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
