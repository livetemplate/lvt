package base

import (
	"testing"
)

func TestNewActionContext(t *testing.T) {
	data := map[string]string{
		"value": "option1",
		"index": "5",
	}

	ctx := NewActionContext("select", "dropdown-1", data)

	if ctx.Action != "select" {
		t.Errorf("expected Action 'select', got '%s'", ctx.Action)
	}

	if ctx.ComponentID != "dropdown-1" {
		t.Errorf("expected ComponentID 'dropdown-1', got '%s'", ctx.ComponentID)
	}
}

func TestNewActionContext_NilData(t *testing.T) {
	ctx := NewActionContext("toggle", "myid", nil)

	// Should not panic
	if ctx.Data("key") != "" {
		t.Error("expected empty string for missing key")
	}
}

func TestActionContext_Data(t *testing.T) {
	data := map[string]string{
		"value":   "option1",
		"label":   "Option One",
		"missing": "",
	}

	ctx := NewActionContext("select", "myid", data)

	tests := []struct {
		key      string
		expected string
	}{
		{"value", "option1"},
		{"label", "Option One"},
		{"missing", ""},
		{"nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := ctx.Data(tt.key)
			if result != tt.expected {
				t.Errorf("Data(%q) = %q, want %q", tt.key, result, tt.expected)
			}
		})
	}
}

func TestActionContext_DataInt(t *testing.T) {
	data := map[string]string{
		"count":   "42",
		"invalid": "not-a-number",
		"empty":   "",
		"float":   "3.14",
		"negative": "-10",
	}

	ctx := NewActionContext("select", "myid", data)

	tests := []struct {
		key      string
		expected int
	}{
		{"count", 42},
		{"invalid", 0},
		{"empty", 0},
		{"float", 0}, // ParseInt doesn't handle floats
		{"negative", -10},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := ctx.DataInt(tt.key)
			if result != tt.expected {
				t.Errorf("DataInt(%q) = %d, want %d", tt.key, result, tt.expected)
			}
		})
	}
}

func TestActionContext_DataFloat(t *testing.T) {
	data := map[string]string{
		"price":    "19.99",
		"integer":  "42",
		"invalid":  "not-a-number",
		"empty":    "",
		"negative": "-3.14",
	}

	ctx := NewActionContext("update", "myid", data)

	tests := []struct {
		key      string
		expected float64
	}{
		{"price", 19.99},
		{"integer", 42.0},
		{"invalid", 0.0},
		{"empty", 0.0},
		{"negative", -3.14},
		{"nonexistent", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := ctx.DataFloat(tt.key)
			if result != tt.expected {
				t.Errorf("DataFloat(%q) = %f, want %f", tt.key, result, tt.expected)
			}
		})
	}
}

func TestActionContext_DataBool(t *testing.T) {
	data := map[string]string{
		"enabled":  "true",
		"disabled": "false",
		"one":      "1",
		"zero":     "0",
		"yes":      "yes",
		"no":       "no",
		"on":       "on",
		"off":      "off",
		"True":     "True",
		"TRUE":     "TRUE",
		"empty":    "",
		"invalid":  "maybe",
	}

	ctx := NewActionContext("toggle", "myid", data)

	tests := []struct {
		key      string
		expected bool
	}{
		{"enabled", true},
		{"disabled", false},
		{"one", true},
		{"zero", false},
		{"yes", true},
		{"no", false},
		{"on", true},
		{"off", false},
		{"True", true},
		{"TRUE", true},
		{"empty", false},
		{"invalid", false},
		{"nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := ctx.DataBool(tt.key)
			if result != tt.expected {
				t.Errorf("DataBool(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestActionContext_HasData(t *testing.T) {
	data := map[string]string{
		"present": "value",
		"empty":   "",
	}

	ctx := NewActionContext("check", "myid", data)

	if !ctx.HasData("present") {
		t.Error("expected HasData('present') to be true")
	}

	if !ctx.HasData("empty") {
		t.Error("expected HasData('empty') to be true (key exists)")
	}

	if ctx.HasData("missing") {
		t.Error("expected HasData('missing') to be false")
	}
}

func TestActionContext_AllData(t *testing.T) {
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	ctx := NewActionContext("test", "myid", data)
	allData := ctx.AllData()

	// Should return a copy
	allData["key1"] = "modified"

	// Original should be unchanged
	if ctx.Data("key1") != "value1" {
		t.Error("AllData should return a copy, not the original map")
	}

	// Check all keys are present
	if len(allData) != 2 {
		t.Errorf("expected 2 keys, got %d", len(allData))
	}
}
