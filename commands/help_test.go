package commands

import (
	"testing"
)

func TestValidatePositionalArg(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		argName string
		wantErr bool
	}{
		// Valid names
		{"valid app name", "myapp", "app name", false},
		{"valid name with hyphen in middle", "my-app", "app name", false},
		{"valid name with underscore", "my_app", "app name", false},
		{"valid name with number", "app123", "app name", false},

		// Flag-like arguments (should error)
		{"double dash help", "--help", "app name", true},
		{"double dash short help", "-h", "app name", true},
		{"unknown double dash flag", "--typo", "app name", true},
		{"unknown single dash flag", "-x", "app name", true},
		{"double dash module flag", "--module", "app name", true},
		{"single dash verbose", "-v", "app name", true},

		// Edge cases
		{"empty string", "", "app name", false}, // empty is not flag-like, let command handle it
		{"just hyphen", "-", "app name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositionalArg(tt.arg, tt.argName)
			if tt.wantErr && err == nil {
				t.Errorf("ValidatePositionalArg(%q, %q) should return error", tt.arg, tt.argName)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidatePositionalArg(%q, %q) unexpected error: %v", tt.arg, tt.argName, err)
			}
		})
	}
}

func TestShowHelpIfRequested(t *testing.T) {
	helpCalled := false
	mockHelp := func() {
		helpCalled = true
	}

	tests := []struct {
		name       string
		args       []string
		wantHelp   bool
		wantReturn bool
	}{
		{"no args", []string{}, false, false},
		{"--help flag", []string{"--help"}, true, true},
		{"-h flag", []string{"-h"}, true, true},
		{"--help with other args before", []string{"myapp", "--help"}, true, true},
		{"--help with other args after", []string{"--help", "myapp"}, true, true},
		{"no help flag", []string{"myapp", "--module", "foo"}, false, false},
		{"similar but not help", []string{"--helping"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helpCalled = false
			result := ShowHelpIfRequested(tt.args, mockHelp)

			if result != tt.wantReturn {
				t.Errorf("ShowHelpIfRequested(%v) returned %v, want %v", tt.args, result, tt.wantReturn)
			}

			if helpCalled != tt.wantHelp {
				t.Errorf("ShowHelpIfRequested(%v) helpCalled=%v, want %v", tt.args, helpCalled, tt.wantHelp)
			}
		})
	}
}

func TestNewRejectsFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{"--help shows help, no error", []string{"--help"}, false, ""},
		{"-h shows help, no error", []string{"-h"}, false, ""},
		{"--typo rejected", []string{"--typo"}, true, "looks like a flag"},
		{"-x rejected", []string{"-x"}, true, "looks like a flag"},
		{"--module without value rejected", []string{"--module"}, true, "looks like a flag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't actually run New() because it creates directories
			// Instead, test the validation logic directly

			// First check if help is requested
			if ShowHelpIfRequested(tt.args, func() {}) {
				if tt.wantErr {
					t.Errorf("Expected error for %v but help was shown instead", tt.args)
				}
				return
			}

			// Then check if first arg is valid
			if len(tt.args) > 0 {
				err := ValidatePositionalArg(tt.args[0], "app name")
				if tt.wantErr && err == nil {
					t.Errorf("Expected error for %v", tt.args)
				}
				if !tt.wantErr && err != nil {
					t.Errorf("Unexpected error for %v: %v", tt.args, err)
				}
			}
		})
	}
}
