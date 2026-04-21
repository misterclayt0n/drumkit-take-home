package turvo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"drumkit-take-home/internal/load"
)

func TestFlexibleInt_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    flexibleInt
		wantErr bool
	}{
		{name: "integer", input: `1`, want: 1},
		{name: "float whole", input: `1.0`, want: 1},
		{name: "string whole", input: `"2.0"`, want: 2},
		{name: "null", input: `null`, want: 0},
		{name: "fractional", input: `1.5`, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got flexibleInt
			err := json.Unmarshal([]byte(tc.input), &got)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %d, got %d", tc.want, got)
			}
		})
	}
}

func TestMapShipmentToLoad_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputPath string
		wantPath  string
	}{
		{name: "route detail", inputPath: "loads/detail-route.input.json", wantPath: "loads/detail-route.want.json"},
		{name: "global route fallback", inputPath: "loads/detail-global.input.json", wantPath: "loads/detail-global.want.json"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := readFixture[shipmentDetail](t, tc.inputPath)
			want := readFixture[load.Load](t, tc.wantPath)

			got := mapShipmentToLoad(input)

			assertDeepEqualJSON(t, want, got)
		})
	}
}

func TestMapShipmentToLoad_Invariants(t *testing.T) {
	t.Parallel()

	detail := readFixture[shipmentDetail](t, "loads/detail-route.input.json")
	got := mapShipmentToLoad(detail)

	if got.ExternalTMSLoadID != "1001" {
		t.Fatalf("expected externalTMSLoadID to map from shipment id, got %q", got.ExternalTMSLoadID)
	}
	if got.RateData.ProfitPercent != 25 {
		t.Fatalf("expected profitPercent 25, got %v", got.RateData.ProfitPercent)
	}
	if got.Specifications.MinTempFahrenheit != 39.2 || got.Specifications.MaxTempFahrenheit != 39.2 {
		t.Fatalf("expected temperature conversion from celsius to fahrenheit, got min=%v max=%v", got.Specifications.MinTempFahrenheit, got.Specifications.MaxTempFahrenheit)
	}
	if got.InPalletCount != 10 || got.OutPalletCount != 10 {
		t.Fatalf("expected pallet counts to be derived from handling units, got in=%d out=%d", got.InPalletCount, got.OutPalletCount)
	}
	if got.TotalWeight != 1200 || got.BillableWeight != 1200 {
		t.Fatalf("expected total and billable weight to be derived from items, got total=%v billable=%v", got.TotalWeight, got.BillableWeight)
	}
	if got.Operator != "Jane Operator" {
		t.Fatalf("expected operator to come from first active contributor, got %q", got.Operator)
	}
}

func TestListShipmentSummaryPage_UsesTurvoFiltersAndPagination(t *testing.T) {
	t.Parallel()

	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"Status":"OK","details":{"pagination":{"start":20,"pageSize":20,"totalRecordsInPage":1,"moreAvailable":true},"shipments":[{"id":1,"status":{"code":{"value":"Tendered"}}}]}}`)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, apiKey: "test", httpClient: server.Client(), accessToken: "token", expiresAt: time.Now().Add(time.Hour)}
	_, err := client.listShipmentSummaryPage(context.Background(), 20, 20, load.ListParams{
		Status:               "Tendered",
		CustomerID:           "834045",
		PickupDateSearchFrom: "2026-05-01",
		PickupDateSearchTo:   "2026-05-02",
	})
	if err != nil {
		t.Fatalf("listShipmentSummaryPage returned error: %v", err)
	}

	checks := []string{
		"start=20",
		"pageSize=20",
		"status%5Beq%5D=2101",
		"customerId%5Beq%5D=834045",
		"pickupDate%5Bgte%5D=2026-05-01T00%3A00%3A00Z",
		"pickupDate%5Blte%5D=2026-05-02T23%3A59%3A59Z",
	}
	for _, check := range checks {
		if !strings.Contains(query, check) {
			t.Fatalf("expected query %q to contain %q", query, check)
		}
	}
}

func TestPaginationFromSummaryPage_UsesLowerBoundWhenMoreAvailable(t *testing.T) {
	t.Parallel()

	got := paginationFromSummaryPage(2, 20, 20, 20, true)
	if got.Total != 41 {
		t.Fatalf("expected lower-bound total 41, got %d", got.Total)
	}
	if got.Pages != 3 {
		t.Fatalf("expected pages 3, got %d", got.Pages)
	}

	got = paginationFromSummaryPage(3, 20, 40, 5, false)
	if got.Total != 45 {
		t.Fatalf("expected exact total 45 on last page, got %d", got.Total)
	}
	if got.Pages != 3 {
		t.Fatalf("expected exact pages 3 on last page, got %d", got.Pages)
	}
}

func TestGetShipmentDetailWithRetry_RetriesAndSucceeds(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if r.URL.Path != "/shipments/123" {
			http.NotFound(w, r)
			return
		}
		if attempts < 3 {
			http.Error(w, "upstream unavailable", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"Status":"OK","details":{"id":123}}`)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL,
		apiKey:      "test",
		httpClient:  server.Client(),
		accessToken: "token",
		expiresAt:   time.Now().Add(time.Hour),
	}

	detail, skipped, err := client.getShipmentDetailWithRetry(context.Background(), 123)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if skipped {
		t.Fatal("expected shipment detail not to be skipped")
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if detail.ID != 123 {
		t.Fatalf("expected detail id 123, got %d", detail.ID)
	}
}

func TestGetShipmentDetailWithRetry_SkipsAfterRetriableFailures(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "still failing", http.StatusBadGateway)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL,
		apiKey:      "test",
		httpClient:  server.Client(),
		accessToken: "token",
		expiresAt:   time.Now().Add(time.Hour),
	}

	_, skipped, err := client.getShipmentDetailWithRetry(context.Background(), 999)
	if err != nil {
		t.Fatalf("expected no error on skipped shipment, got %v", err)
	}
	if !skipped {
		t.Fatal("expected shipment detail to be skipped after retries")
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestGetShipmentDetailWithRetry_DoesNotRetryNonRetriableStatus(t *testing.T) {
	var attempts int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		baseURL:     server.URL,
		apiKey:      "test",
		httpClient:  server.Client(),
		accessToken: "token",
		expiresAt:   time.Now().Add(time.Hour),
	}

	_, skipped, err := client.getShipmentDetailWithRetry(context.Background(), 500)
	if err == nil {
		t.Fatal("expected error for non-retriable status")
	}
	if skipped {
		t.Fatal("did not expect shipment detail to be skipped on first non-retriable error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", attempts)
	}
}
