package utils_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/mrangel-jr/api-billing/internals/utils"
	"github.com/stretchr/testify/assert"
)

func SetupRequest(method, url, params, value string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(params, value) // ✅ This is the key step
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	req = req.WithContext(ctx)
	return req
}
func TestReadYearMonthParam(t *testing.T) {
	tests := []struct {
		name       string
		year_month string
		wantErr    bool
	}{
		{
			name:       "valid year_month",
			year_month: "202301",
			wantErr:    false,
		},
		{
			name:       "invalid year_month",
			year_month: "2023-01",
			wantErr:    true,
		},
		{
			name:       "invalid year_month",
			year_month: "202313",
			wantErr:    true,
		},
		{
			name:       "invalid year_month",
			year_month: "2023",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := SetupRequest(http.MethodGet, "/consume/"+tt.year_month, "year_month", tt.year_month)
			year, month, err := utils.ReadYearMonthParam(req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, year)
				assert.Equal(t, 0, month)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 2023, year)
				assert.Equal(t, 1, month)
			}
		})
	}
}

func TestCreatePagination(t *testing.T) {
	tests := []struct {
		name               string
		pagination         map[string]string
		expectedPagination utils.Pagination
	}{
		{
			name:       "valid pagination",
			pagination: map[string]string{"page": "2", "pageSize": "25"},
			expectedPagination: utils.Pagination{
				Page:     2,
				PageSize: 25,
			},
		},
		{
			name:       "valid page without pageSize",
			pagination: map[string]string{"page": "5"},
			expectedPagination: utils.Pagination{
				Page:     5,
				PageSize: 10, // default
			},
		},
		{
			name:       "valid pageSize without page",
			pagination: map[string]string{"pageSize": "50"},
			expectedPagination: utils.Pagination{
				Page:     1, // default
				PageSize: 50,
			},
		},
		{
			name:       "invalid page (zero)",
			pagination: map[string]string{"page": "0", "pageSize": "10"},
			expectedPagination: utils.Pagination{
				Page:     1, // default
				PageSize: 10,
			},
		},
		{
			name:       "invalid pageSize (zero)",
			pagination: map[string]string{"page": "2", "pageSize": "0"},
			expectedPagination: utils.Pagination{
				Page:     2,
				PageSize: 10, // default
			},
		},
		{
			name:       "invalid page and pageSize (negative)",
			pagination: map[string]string{"page": "-1", "pageSize": "-5"},
			expectedPagination: utils.Pagination{
				Page:     1,  // default
				PageSize: 10, // default
			},
		},
		{
			name:       "non-integer values",
			pagination: map[string]string{"page": "abc", "pageSize": "xyz"},
			expectedPagination: utils.Pagination{
				Page:     1,  // default
				PageSize: 10, // default
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("/consume?page=%v&pageSize=%v", tt.pagination["page"], tt.pagination["pageSize"])
			pagination := utils.CreatePagination(SetupRequest(http.MethodGet, url, "", ""))

			assert.Equal(t, tt.expectedPagination.Page, pagination.Page)
			assert.Equal(t, tt.expectedPagination.PageSize, pagination.PageSize)
		})
	}
}
