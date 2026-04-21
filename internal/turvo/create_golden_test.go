package turvo

import (
	"strings"
	"testing"

	"drumkit-take-home/internal/load"
)

func TestBuildCreateShipmentRequest_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputPath string
		wantPath  string
	}{
		{name: "minimal payload", inputPath: "create/create-minimal.input.json", wantPath: "create/create-minimal.want.json"},
		{name: "rich payload", inputPath: "create/create-rich.input.json", wantPath: "create/create-rich.want.json"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := readFixture[load.Load](t, tc.inputPath)
			want := readFixture[createShipmentRequest](t, tc.wantPath)

			got, err := buildCreateShipmentRequest(input)
			if err != nil {
				t.Fatalf("buildCreateShipmentRequest returned error: %v", err)
			}

			assertDeepEqualJSON(t, want, got)
		})
	}
}

func TestBuildCreateShipmentRequest_Validation(t *testing.T) {
	t.Parallel()

	base := readFixture[load.Load](t, "create/create-minimal.input.json")

	tests := []struct {
		name        string
		mutate      func(*load.Load)
		wantMessage string
	}{
		{
			name: "missing customer external id",
			mutate: func(in *load.Load) {
				in.Customer.ExternalTMSID = ""
			},
			wantMessage: "customer.externalTMSId is required",
		},
		{
			name: "missing pickup ids",
			mutate: func(in *load.Load) {
				in.Pickup.ExternalTMSID = ""
				in.Pickup.WarehouseID = ""
			},
			wantMessage: "pickup.externalTMSId is required",
		},
		{
			name: "invalid pickup time",
			mutate: func(in *load.Load) {
				in.Pickup.ApptTime = "not-a-time"
				in.Pickup.ReadyTime = ""
			},
			wantMessage: "pickup.apptTime must be a valid RFC3339 timestamp",
		},
		{
			name: "missing consignee time",
			mutate: func(in *load.Load) {
				in.Consignee.ApptTime = ""
				in.Consignee.MustDeliver = ""
			},
			wantMessage: "consignee.apptTime is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input := base
			tc.mutate(&input)

			_, err := buildCreateShipmentRequest(input)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.wantMessage) {
				t.Fatalf("expected error to contain %q, got %q", tc.wantMessage, err.Error())
			}
		})
	}
}

func TestBuildCreateShipmentRequest_Invariants(t *testing.T) {
	t.Parallel()

	input := readFixture[load.Load](t, "create/create-rich.input.json")
	got, err := buildCreateShipmentRequest(input)
	if err != nil {
		t.Fatalf("buildCreateShipmentRequest returned error: %v", err)
	}

	if got.StartDate.TimeZone != "America/New_York" {
		t.Fatalf("expected pickup timezone to default to America/New_York, got %q", got.StartDate.TimeZone)
	}
	if got.EndDate.TimeZone != "America/New_York" {
		t.Fatalf("expected delivery timezone to inherit pickup timezone when empty, got %q", got.EndDate.TimeZone)
	}
	if got.GlobalRoute[0].Location.ID != 777001 {
		t.Fatalf("expected pickup location to fall back to warehouseId, got %d", got.GlobalRoute[0].Location.ID)
	}
	if got.GlobalRoute[1].Location.ID != 888002 {
		t.Fatalf("expected consignee location to fall back to warehouseId, got %d", got.GlobalRoute[1].Location.ID)
	}
	if len(got.CustomerOrder) != 1 || len(got.CustomerOrder[0].ExternalIDs) != 5 {
		t.Fatalf("expected deduped external ids to include po, ref, external load, and freight load entries; got %+v", got.CustomerOrder)
	}
}
