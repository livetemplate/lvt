package validation

import (
	"context"
	"os"
	"path/filepath"

	"github.com/livetemplate/lvt/internal/validator"
	"golang.org/x/mod/modfile"
)

// GoModCheck validates the go.mod file in an app directory.
type GoModCheck struct{}

func (c *GoModCheck) Name() string { return "go.mod" }

func (c *GoModCheck) Run(_ context.Context, appPath string) *validator.ValidationResult {
	result := validator.NewValidationResult()
	goModPath := filepath.Join(appPath, "go.mod")

	data, err := os.ReadFile(goModPath)
	if err != nil {
		if os.IsNotExist(err) {
			result.AddError("go.mod not found", "go.mod", 0)
		} else {
			result.AddError("failed to read go.mod: "+err.Error(), "go.mod", 0)
		}
		return result
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		result.AddError("go.mod parse error: "+err.Error(), "go.mod", 0)
		return result
	}

	if f.Module == nil || f.Module.Mod.Path == "" {
		result.AddError("go.mod has no module path", "go.mod", 0)
	}

	if f.Go == nil || f.Go.Version == "" {
		result.AddWarning("go.mod has no go version directive", "go.mod", 0)
	}

	return result
}
