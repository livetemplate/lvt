package evolution

import (
	"testing"

	"github.com/livetemplate/lvt/internal/telemetry"
)

func TestClassifyError_Component(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantType  string
		wantComp  string
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
		File: "internal/kits/multi/form.tmpl",
	})
	if loc.Type != "kit" {
		t.Errorf("expected type 'kit', got %q", loc.Type)
	}
	if loc.Component != "" {
		t.Errorf("expected empty component for kit, got %q", loc.Component)
	}
}

func TestClassifyError_Generated(t *testing.T) {
	loc := ClassifyError(telemetry.GenerationError{
		File: "app/posts/handler.go",
	})
	if loc.Type != "generated" {
		t.Errorf("expected type 'generated', got %q", loc.Type)
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

func TestClassifyFix(t *testing.T) {
	fix := Fix{
		TargetFile: "components/modal/*.go",
	}
	loc := ClassifyFix(fix)
	if loc.Type != "component" {
		t.Errorf("expected type 'component', got %q", loc.Type)
	}
	if loc.Component != "modal" {
		t.Errorf("expected component 'modal', got %q", loc.Component)
	}
}

func TestClassifyFix_Kit(t *testing.T) {
	fix := Fix{
		TargetFile: "internal/kits/multi/templates/*.tmpl",
	}
	loc := ClassifyFix(fix)
	if loc.Type != "kit" {
		t.Errorf("expected type 'kit', got %q", loc.Type)
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
