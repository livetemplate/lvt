// Package eject provides functionality to eject component source code
// from the components library to a user's project.
package eject

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ComponentInfo holds information about an ejectable component.
type ComponentInfo struct {
	Name        string
	Package     string // e.g., "github.com/livetemplate/components/dropdown"
	Description string
	Templates   []string
}

// AvailableComponents returns the list of components that can be ejected.
func AvailableComponents() []ComponentInfo {
	return []ComponentInfo{
		{Name: "accordion", Package: "github.com/livetemplate/components/accordion", Description: "Collapsible content sections", Templates: []string{"default", "single"}},
		{Name: "autocomplete", Package: "github.com/livetemplate/components/autocomplete", Description: "Search with suggestions", Templates: []string{"default"}},
		{Name: "breadcrumbs", Package: "github.com/livetemplate/components/breadcrumbs", Description: "Navigation breadcrumb trail", Templates: []string{"default"}},
		{Name: "datatable", Package: "github.com/livetemplate/components/datatable", Description: "Data tables with sorting/pagination", Templates: []string{"default"}},
		{Name: "datepicker", Package: "github.com/livetemplate/components/datepicker", Description: "Date selection", Templates: []string{"single", "range", "inline"}},
		{Name: "drawer", Package: "github.com/livetemplate/components/drawer", Description: "Slide-out panels", Templates: []string{"default"}},
		{Name: "dropdown", Package: "github.com/livetemplate/components/dropdown", Description: "Dropdown menus", Templates: []string{"default", "searchable", "multi"}},
		{Name: "menu", Package: "github.com/livetemplate/components/menu", Description: "Navigation menus", Templates: []string{"default", "nested"}},
		{Name: "modal", Package: "github.com/livetemplate/components/modal", Description: "Modal dialogs", Templates: []string{"default", "confirm", "sheet"}},
		{Name: "popover", Package: "github.com/livetemplate/components/popover", Description: "Rich content popovers", Templates: []string{"default"}},
		{Name: "progress", Package: "github.com/livetemplate/components/progress", Description: "Progress indicators", Templates: []string{"default", "circular", "spinner"}},
		{Name: "rating", Package: "github.com/livetemplate/components/rating", Description: "Star ratings", Templates: []string{"default"}},
		{Name: "skeleton", Package: "github.com/livetemplate/components/skeleton", Description: "Loading placeholders", Templates: []string{"default", "avatar", "card"}},
		{Name: "tabs", Package: "github.com/livetemplate/components/tabs", Description: "Tab navigation", Templates: []string{"horizontal", "vertical", "pills"}},
		{Name: "tagsinput", Package: "github.com/livetemplate/components/tagsinput", Description: "Tag/chip input", Templates: []string{"default"}},
		{Name: "timeline", Package: "github.com/livetemplate/components/timeline", Description: "Event timelines", Templates: []string{"default"}},
		{Name: "timepicker", Package: "github.com/livetemplate/components/timepicker", Description: "Time selection", Templates: []string{"default"}},
		{Name: "toast", Package: "github.com/livetemplate/components/toast", Description: "Toast notifications", Templates: []string{"default", "container"}},
		{Name: "toggle", Package: "github.com/livetemplate/components/toggle", Description: "Toggle switches", Templates: []string{"default", "checkbox"}},
		{Name: "tooltip", Package: "github.com/livetemplate/components/tooltip", Description: "Tooltips", Templates: []string{"default"}},
	}
}

// FindComponent returns component info by name, or nil if not found.
func FindComponent(name string) *ComponentInfo {
	for _, c := range AvailableComponents() {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// EjectOptions configures the eject operation.
type EjectOptions struct {
	// ComponentName is the name of the component to eject.
	ComponentName string

	// DestDir is the destination directory (default: internal/components/<name>).
	DestDir string

	// ModuleName is the current project's Go module name.
	ModuleName string

	// Force overwrites existing files.
	Force bool
}

// EjectComponent copies the full component source to the user's project.
func EjectComponent(opts EjectOptions) error {
	comp := FindComponent(opts.ComponentName)
	if comp == nil {
		return fmt.Errorf("unknown component: %s\nRun 'lvt component list' to see available components", opts.ComponentName)
	}

	// Determine destination directory
	destDir := opts.DestDir
	if destDir == "" {
		destDir = filepath.Join("internal", "components", opts.ComponentName)
	}

	// Check if destination exists
	if !opts.Force {
		if _, err := os.Stat(destDir); err == nil {
			return fmt.Errorf("destination already exists: %s\nUse --force to overwrite", destDir)
		}
	}

	// Find the component source in Go module cache
	srcDir, err := findModulePath(comp.Package)
	if err != nil {
		return fmt.Errorf("component not found in module cache: %v\nRun 'go get %s' first", err, comp.Package)
	}

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Copy all files
	files, err := copyDir(srcDir, destDir)
	if err != nil {
		return fmt.Errorf("failed to copy files: %v", err)
	}

	// Update package imports in copied files if module name is provided
	if opts.ModuleName != "" {
		newPkg := opts.ModuleName + "/internal/components/" + opts.ComponentName
		for _, file := range files {
			if strings.HasSuffix(file, ".go") {
				if err := updateImports(file, comp.Package, newPkg); err != nil {
					fmt.Printf("Warning: failed to update imports in %s: %v\n", file, err)
				}
			}
		}
	}

	fmt.Printf("✅ Ejected %s to %s\n", opts.ComponentName, destDir)
	fmt.Println()
	fmt.Println("Ejected files:")
	for _, f := range files {
		rel, _ := filepath.Rel(destDir, f)
		fmt.Printf("  - %s\n", rel)
	}
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Update your imports from:\n")
	fmt.Printf("       %s\n", comp.Package)
	fmt.Printf("     to:\n")
	if opts.ModuleName != "" {
		fmt.Printf("       %s/internal/components/%s\n", opts.ModuleName, opts.ComponentName)
	} else {
		fmt.Printf("       yourmodule/internal/components/%s\n", opts.ComponentName)
	}
	fmt.Printf("  2. Customize the component as needed\n")

	return nil
}

// EjectTemplateOptions configures template-only ejection.
type EjectTemplateOptions struct {
	// ComponentName is the name of the component.
	ComponentName string

	// TemplateName is the template variant to eject.
	TemplateName string

	// DestDir is the destination directory (default: internal/templates).
	DestDir string

	// Force overwrites existing files.
	Force bool
}

// EjectTemplate copies only the template file to the user's project.
func EjectTemplate(opts EjectTemplateOptions) error {
	comp := FindComponent(opts.ComponentName)
	if comp == nil {
		return fmt.Errorf("unknown component: %s\nRun 'lvt component list' to see available components", opts.ComponentName)
	}

	// Validate template name
	found := false
	for _, t := range comp.Templates {
		if t == opts.TemplateName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown template '%s' for component '%s'\nAvailable: %v", opts.TemplateName, opts.ComponentName, comp.Templates)
	}

	// Determine destination
	destDir := opts.DestDir
	if destDir == "" {
		destDir = filepath.Join("internal", "templates")
	}

	destFile := filepath.Join(destDir, fmt.Sprintf("%s-%s.tmpl", opts.ComponentName, opts.TemplateName))

	// Check if destination exists
	if !opts.Force {
		if _, err := os.Stat(destFile); err == nil {
			return fmt.Errorf("destination already exists: %s\nUse --force to overwrite", destFile)
		}
	}

	// Find the component source in Go module cache
	srcDir, err := findModulePath(comp.Package)
	if err != nil {
		return fmt.Errorf("component not found in module cache: %v\nRun 'go get %s' first", err, comp.Package)
	}

	// Find the template file
	srcFile := filepath.Join(srcDir, "templates", opts.TemplateName+".tmpl")
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", srcFile)
	}

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Copy the template file
	if err := copyFile(srcFile, destFile); err != nil {
		return fmt.Errorf("failed to copy template: %v", err)
	}

	fmt.Printf("✅ Ejected template to %s\n", destFile)
	fmt.Println()
	fmt.Println("The template will automatically override the library version.")
	fmt.Println("The Go logic remains in the library and updates with 'go get -u'.")

	return nil
}

// findModulePath locates a package in the Go module cache.
func findModulePath(pkg string) (string, error) {
	// Use 'go list' to find the module directory
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", pkg)
	output, err := cmd.Output()
	if err != nil {
		// Try without the module flag for packages
		cmd = exec.Command("go", "list", "-f", "{{.Dir}}", pkg)
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("package not found: %s", pkg)
		}
	}

	dir := strings.TrimSpace(string(output))
	if dir == "" {
		return "", fmt.Errorf("empty module path for: %s", pkg)
	}

	return dir, nil
}

// copyDir recursively copies a directory and returns list of copied files.
func copyDir(src, dst string) ([]string, error) {
	var copied []string

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip hidden files and test files
		if strings.HasPrefix(filepath.Base(path), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		destPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		if err := copyFile(path, destPath); err != nil {
			return err
		}
		copied = append(copied, destPath)

		return nil
	})

	return copied, err
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// updateImports updates import paths in a Go file.
func updateImports(file, oldPkg, newPkg string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Simple string replacement for imports
	updated := strings.ReplaceAll(string(content), oldPkg, newPkg)

	return os.WriteFile(file, []byte(updated), 0644)
}

// GetModuleName attempts to get the current Go module name from go.mod.
func GetModuleName() (string, error) {
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
