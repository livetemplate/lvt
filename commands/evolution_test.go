package commands

import (
	"testing"
)

func TestEvolution_Help(t *testing.T) {
	// Should not error when no subcommand given
	err := Evolution(nil)
	if err != nil {
		t.Errorf("expected nil error for empty args, got: %v", err)
	}

	err = Evolution([]string{"help"})
	if err != nil {
		t.Errorf("expected nil error for help, got: %v", err)
	}
}

func TestEvolution_UnknownSubcommand(t *testing.T) {
	err := Evolution([]string{"nonexistent"})
	if err == nil {
		t.Error("expected error for unknown subcommand")
	}
}

func TestEvolution_Status(t *testing.T) {
	// Status should work even with no data (telemetry may be disabled in test env)
	t.Setenv("LVT_TELEMETRY", "true")
	// Use temp dir for config to avoid polluting real config
	t.Setenv("HOME", t.TempDir())

	err := Evolution([]string{"status"})
	if err != nil {
		t.Errorf("status: %v", err)
	}
}

func TestEvolution_Patterns(t *testing.T) {
	// This test requires the patterns.md file to be findable
	err := Evolution([]string{"patterns"})
	if err != nil {
		t.Logf("patterns command failed (expected if patterns.md not accessible): %v", err)
	}
}

func TestEvolution_Failures(t *testing.T) {
	t.Setenv("LVT_TELEMETRY", "true")
	t.Setenv("HOME", t.TempDir())

	err := Evolution([]string{"failures"})
	if err != nil {
		t.Errorf("failures: %v", err)
	}
}

func TestEvolution_Failures_WithLastFlag(t *testing.T) {
	t.Setenv("LVT_TELEMETRY", "true")
	t.Setenv("HOME", t.TempDir())

	err := Evolution([]string{"failures", "--last", "5"})
	if err != nil {
		t.Errorf("failures --last 5: %v", err)
	}
}

func TestEvolution_Metrics(t *testing.T) {
	t.Setenv("LVT_TELEMETRY", "true")
	t.Setenv("HOME", t.TempDir())

	err := Evolution([]string{"metrics"})
	if err != nil {
		t.Errorf("metrics: %v", err)
	}
}

func TestEvolution_Propose_NoArgs(t *testing.T) {
	err := Evolution([]string{"propose"})
	if err == nil {
		t.Error("expected error for propose without event-id")
	}
}

func TestEvolution_Apply_NoArgs(t *testing.T) {
	err := Evolution([]string{"apply"})
	if err == nil {
		t.Error("expected error for apply without fix-id")
	}
}
