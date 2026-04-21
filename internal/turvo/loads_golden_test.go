package turvo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"drumkit-take-home/internal/load"
)

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
