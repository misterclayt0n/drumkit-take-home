package turvo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func readFixture[T any](t *testing.T, relativePath string) T {
	t.Helper()

	content, err := os.ReadFile(filepath.Join("testdata", relativePath))
	if err != nil {
		t.Fatalf("read fixture %s: %v", relativePath, err)
	}

	var value T
	if err := json.Unmarshal(content, &value); err != nil {
		t.Fatalf("decode fixture %s: %v", relativePath, err)
	}

	return value
}

func assertDeepEqualJSON(t *testing.T, want, got any) {
	t.Helper()
	if reflect.DeepEqual(want, got) {
		return
	}

	wantJSON, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		t.Fatalf("marshal want json: %v", err)
	}
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("marshal got json: %v", err)
	}

	t.Fatalf("json mismatch\nwant:\n%s\n\ngot:\n%s", wantJSON, gotJSON)
}
