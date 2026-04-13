package cliutil

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/walnut1024/efi-cli/internal/model"
)

func TestApplySortDescendingNumeric(t *testing.T) {
	cfg := OutputConfig{Sort: "pct"}
	records := []model.Record{
		{"code": "a", "pct": 1.2},
		{"code": "b", "pct": 3.4},
		{"code": "c", "pct": -0.5},
	}

	got, err := cfg.ApplySort(records)
	if err != nil {
		t.Fatalf("ApplySort returned error: %v", err)
	}

	codes := []string{got[0]["code"].(string), got[1]["code"].(string), got[2]["code"].(string)}
	want := []string{"b", "a", "c"}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("unexpected order: got %v want %v", codes, want)
	}
}

func TestApplySortAscendingNumeric(t *testing.T) {
	cfg := OutputConfig{Sort: "pct", Asc: true}
	records := []model.Record{
		{"code": "a", "pct": 1.2},
		{"code": "b", "pct": 3.4},
		{"code": "c", "pct": -0.5},
	}

	got, err := cfg.ApplySort(records)
	if err != nil {
		t.Fatalf("ApplySort returned error: %v", err)
	}

	codes := []string{got[0]["code"].(string), got[1]["code"].(string), got[2]["code"].(string)}
	want := []string{"c", "a", "b"}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("unexpected order: got %v want %v", codes, want)
	}
}

func TestApplySortPlacesNilLast(t *testing.T) {
	cfg := OutputConfig{Sort: "pct", Desc: true}
	records := []model.Record{
		{"code": "a", "pct": nil},
		{"code": "b", "pct": 3.4},
		{"code": "c", "pct": 1.0},
	}

	got, err := cfg.ApplySort(records)
	if err != nil {
		t.Fatalf("ApplySort returned error: %v", err)
	}

	codes := []string{got[0]["code"].(string), got[1]["code"].(string), got[2]["code"].(string)}
	want := []string{"b", "c", "a"}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("unexpected order: got %v want %v", codes, want)
	}
}

func TestApplySortRejectsConflictingDirectionFlags(t *testing.T) {
	cfg := OutputConfig{Sort: "pct", Asc: true, Desc: true}
	_, err := cfg.ApplySort([]model.Record{{"code": "a", "pct": 1}})
	if err == nil {
		t.Fatal("expected error when --asc and --desc are both set")
	}
}

func TestValidateSortRejectsDirectionWithoutField(t *testing.T) {
	cfg := OutputConfig{Asc: true}
	if err := cfg.ValidateSort(); err == nil {
		t.Fatal("expected error when sort direction is set without --sort")
	}
}

func TestFormatWithSchemaOrExitFiltersJSONFields(t *testing.T) {
	cfg := OutputConfig{Format: "json", Fields: "name,code"}
	records := []model.Record{
		{"code": "600000", "name": "浦发银行", "pct": 1.2},
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe error: %v", err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	cfg.FormatWithSchemaOrExit(records, nil)
	_ = w.Close()
	os.Stdout = oldStdout
	out := <-done

	if strings.Contains(out, "\"pct\"") {
		t.Fatalf("unexpected pct field in output: %s", out)
	}
	if !strings.Contains(out, "\"name\"") || !strings.Contains(out, "\"code\"") {
		t.Fatalf("expected filtered fields in output: %s", out)
	}
}
