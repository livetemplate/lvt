package generator

import (
	"strings"
	"testing"
)

func TestScheduleToGo(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"@hourly", "time.Hour"},
		{"@daily", "24 * time.Hour"},
		{"@weekly", "7 * 24 * time.Hour"},
		{"@every 5m", "5 * time.Minute"},
		{"@every 1h", "1 * time.Hour"},
		{"@every 30s", "30 * time.Second"},
		{"@every 10m", "10 * time.Minute"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := scheduleToGo(tt.input)
			if got != tt.want {
				t.Errorf("scheduleToGo(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestScheduleToGo_InvalidFallback(t *testing.T) {
	got := scheduleToGo("some-invalid-schedule")
	if !strings.Contains(got, "TODO") {
		t.Errorf("expected TODO in fallback, got %q", got)
	}
}

func TestScheduleToGo_RejectsNonNumeric(t *testing.T) {
	got := scheduleToGo("@every abcm")
	if !strings.Contains(got, "TODO") {
		t.Errorf("non-numeric duration should fall back to TODO, got %q", got)
	}
}

func TestIsPositiveInt(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"1", true},
		{"5", true},
		{"100", true},
		{"0", false},
		{"-1", false},
		{"abc", false},
		{"1.5", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isPositiveInt(tt.input); got != tt.want {
				t.Errorf("isPositiveInt(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
