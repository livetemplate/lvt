package evolution

import (
	"testing"

	"github.com/livetemplate/lvt/internal/telemetry"
)

func TestClassifyError_Component(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		wantType string
		wantComp string
	}{
		{
			name:     "component subdirectory",
			file:     "components/modal/modal.go",
			wantType: "component",
			wantComp: "modal",
		},
		{
			name:     "component template",
			file:     "components/toast/templates/default.tmpl",
			wantType: "component",
			wantComp: "toast",
		},
		{
			name:     "component test file",
			file:     "components/dropdown/dropdown_test.go",
			wantType: "component",
			wantComp: "dropdown",
		},
		{
			name:     "bare component file",
			file:     "components/modal.go",
			wantType: "component",
			wantComp: "modal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := ClassifyError(telemetry.GenerationError{File: tt.file})
			if loc.Type != tt.wantType {
				t.Errorf("Type: got %q, want %q", loc.Type, tt.wantType)
			}
			if loc.Component != tt.wantComp {
				t.Errorf("Component: got %q, want %q", loc.Component, tt.wantComp)
			}
		})
	}
}

func TestClassifyError_Kit(t *testing.T) {
	loc := ClassifyError(telemetry.GenerationError{
		File: "internal/kits/system/multi/form.tmpl",
	})
	if loc.Type != "kit" {
		t.Errorf("expected type 'kit', got %q", loc.Type)
	}
	if loc.Component != "" {
		t.Errorf("expected empty component for kit, got %q", loc.Component)
	}
}

func TestClassifyError_KitWithComponentsSubdir(t *testing.T) {
	// Kit paths containing "components/" must be classified as kit, not component
	loc := ClassifyError(telemetry.GenerationError{
		File: "internal/kits/system/multi/components/form.tmpl",
	})
	if loc.Type != "kit" {
		t.Errorf("expected type 'kit' for kit components/ subdir, got %q", loc.Type)
	}
}

func TestClassifyError_Generated(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{name: "app prefix", file: "app/posts/handler.go"},
		{name: "app in middle of path", file: "cmd/server/app/handler.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := ClassifyError(telemetry.GenerationError{File: tt.file})
			if loc.Type != "generated" {
				t.Errorf("expected type 'generated' for %q, got %q", tt.file, loc.Type)
			}
		})
	}
}

func TestClassifyError_GeneratedNoFalsePositive(t *testing.T) {
	// "webapp/" and "myapp/" should NOT classify as generated
	tests := []struct {
		name string
		file string
	}{
		{name: "webapp prefix", file: "webapp/config.go"},
		{name: "myapp prefix", file: "myapp/handler.go"},
		{name: "application prefix", file: "application/foo.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := ClassifyError(telemetry.GenerationError{File: tt.file})
			if loc.Type == "generated" {
				t.Errorf("expected NOT 'generated' for %q, but got 'generated'", tt.file)
			}
		})
	}
}

func TestClassifyError_Unknown(t *testing.T) {
	loc := ClassifyError(telemetry.GenerationError{
		File: "some/other/path.go",
	})
	if loc.Type != "unknown" {
		t.Errorf("expected type 'unknown', got %q", loc.Type)
	}
}

func TestClassifyError_EmptyFile(t *testing.T) {
	loc := ClassifyError(telemetry.GenerationError{File: ""})
	if loc.Type != "unknown" {
		t.Errorf("expected type 'unknown' for empty file, got %q", loc.Type)
	}
}

func TestClassifyFix(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantType string
		wantComp string
	}{
		{
			name:     "component",
			target:   "components/modal/*.go",
			wantType: "component",
			wantComp: "modal",
		},
		{
			name:     "kit",
			target:   "internal/kits/system/multi/templates/*.tmpl",
			wantType: "kit",
			wantComp: "",
		},
		{
			name:     "generated",
			target:   "app/posts/handler.go",
			wantType: "generated",
			wantComp: "",
		},
		{
			name:     "unknown",
			target:   "config/settings.go",
			wantType: "unknown",
			wantComp: "",
		},
		{
			name:     "wildcard glob",
			target:   "components/*/handler.go",
			wantType: "unknown",
			wantComp: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := ClassifyFix(Fix{TargetFile: tt.target})
			if loc.Type != tt.wantType {
				t.Errorf("Type: got %q, want %q", loc.Type, tt.wantType)
			}
			if loc.Component != tt.wantComp {
				t.Errorf("Component: got %q, want %q", loc.Component, tt.wantComp)
			}
		})
	}
}

func TestClassifyError_NoFalseComponentMatch(t *testing.T) {
	// "custom_components/" should NOT match as a component
	tests := []struct {
		name string
		file string
	}{
		{name: "custom_components prefix", file: "vendor/custom_components/modal.go"},
		{name: "my_components prefix", file: "third_party/my_components/toast/toast.go"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc := ClassifyError(telemetry.GenerationError{File: tt.file})
			if loc.Type == "component" {
				t.Errorf("expected NOT 'component' for %q, but got 'component'", tt.file)
			}
		})
	}
}

func TestClassifyError_WindowsPathNormalization(t *testing.T) {
	// Backslash paths should be normalized and still classify correctly
	loc := ClassifyError(telemetry.GenerationError{
		File: "components\\modal\\modal.go",
	})
	if loc.Type != "component" {
		t.Errorf("expected type 'component' for Windows path, got %q", loc.Type)
	}
	if loc.Component != "modal" {
		t.Errorf("expected component 'modal', got %q", loc.Component)
	}
}

func TestClassifyError_CaseInsensitive(t *testing.T) {
	loc := ClassifyError(telemetry.GenerationError{
		File: "Components/Modal/Modal.go",
	})
	if loc.Type != "component" {
		t.Errorf("expected type 'component', got %q", loc.Type)
	}
	if loc.Component != "modal" {
		t.Errorf("expected component 'modal', got %q", loc.Component)
	}
}
