package drumkitstore

import (
	"path/filepath"
	"testing"

	"drumkit-take-home/internal/load"
)

func TestStoreSaveMergeAndReload(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "loads.json")
	store, err := New(path)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	saved := load.Load{
		ExternalTMSLoadID: "1001",
		FreightLoadID:     "canonical-1001",
		Status:            "Tendered",
		Customer:          load.Customer{ExternalTMSID: "834045", Name: "Saved Customer", RefNumber: "REF-1"},
		BillTo:            load.BillTo{Name: "Saved Billing"},
		Carrier:           load.Carrier{Name: "Saved Carrier"},
		Operator:          "Saved Operator",
		PONums:            "PO-1",
	}
	if err := store.Save("1001", saved); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	reloaded, err := New(path)
	if err != nil {
		t.Fatalf("reload New returned error: %v", err)
	}

	upstream := load.Load{
		ExternalTMSLoadID: "1001",
		FreightLoadID:     "upstream-1001",
		Status:            "Covered",
		Customer:          load.Customer{Name: "Upstream Customer"},
	}

	merged := reloaded.Merge(upstream)
	if merged.Status != "Covered" {
		t.Fatalf("expected upstream status to win, got %q", merged.Status)
	}
	if merged.FreightLoadID != saved.FreightLoadID {
		t.Fatalf("expected saved freightLoadID to be preserved, got %q", merged.FreightLoadID)
	}
	if merged.Carrier.Name != "Saved Carrier" {
		t.Fatalf("expected saved carrier fields to be restored, got %q", merged.Carrier.Name)
	}
	if merged.Operator != "Saved Operator" {
		t.Fatalf("expected saved operator to be restored, got %q", merged.Operator)
	}
}

func TestStoreMergeList(t *testing.T) {
	t.Parallel()

	store, err := New(filepath.Join(t.TempDir(), "loads.json"))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if err := store.Save("1001", load.Load{ExternalTMSLoadID: "1001", FreightLoadID: "saved-1001", Status: "Tendered"}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if err := store.Save("1002", load.Load{ExternalTMSLoadID: "1002", FreightLoadID: "saved-1002", Status: "Tendered"}); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	merged := store.MergeList([]load.Load{
		{ExternalTMSLoadID: "1001", Status: "Covered"},
		{ExternalTMSLoadID: "missing", Status: "Tendered"},
		{ExternalTMSLoadID: "1002", Status: "Covered"},
	})

	if len(merged) != 3 {
		t.Fatalf("expected 3 merged loads, got %d", len(merged))
	}
	if merged[0].FreightLoadID != "saved-1001" || merged[0].Status != "Covered" {
		t.Fatalf("expected first saved record to merge and preserve upstream status, got %+v", merged[0])
	}
	if merged[1].ExternalTMSLoadID != "missing" || merged[1].Status != "Tendered" {
		t.Fatalf("expected unmatched load to pass through unchanged, got %+v", merged[1])
	}
	if merged[2].FreightLoadID != "saved-1002" || merged[2].Status != "Covered" {
		t.Fatalf("expected second saved record to merge and preserve upstream status, got %+v", merged[2])
	}
}
