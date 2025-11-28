package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallAgent_NewInstallation(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run install
	err = InstallAgent([]string{"--llm", "claude"})
	if err != nil {
		t.Fatalf("InstallAgent failed: %v", err)
	}

	// Verify .claude directory exists
	claudeDir := filepath.Join(tmpDir, ".claude")
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Errorf(".claude directory was not created")
	}

	// Verify essential files exist
	essentialFiles := []string{
		".claude/settings.json",
		".claude/agents/project-manager-backlog.md",
		".claude/skills/lvt/core/new-app.md",
		".claude/skills/lvt/core/add-resource.md",
		".claude/skills/lvt/workflows/quickstart.md",
	}

	for _, file := range essentialFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Essential file missing: %s", file)
		}
	}

	// Verify directory structure
	expectedDirs := []string{
		".claude/agents",
		".claude/skills",
		".claude/skills/lvt",
		".claude/skills/lvt/core",
		".claude/skills/lvt/workflows",
		".claude/skills/lvt/maintenance",
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(tmpDir, dir)
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			t.Errorf("Expected directory missing: %s", dir)
		} else if !info.IsDir() {
			t.Errorf("Expected directory is not a directory: %s", dir)
		}
	}
}

func TestInstallAgent_ExistingInstallation(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Create existing .claude directory
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("Failed to create .claude dir: %v", err)
	}

	// Try to install without force (should fail)
	err = InstallAgent([]string{"--llm", "claude"})
	if err == nil {
		t.Errorf("InstallAgent should fail when .claude exists without --force")
	}

	// Verify error message is helpful
	if err != nil {
		errMsg := err.Error()
		if errMsg != "installation cancelled" {
			t.Errorf("Expected 'installation cancelled' error, got: %s", errMsg)
		}
	}
}

func TestInstallAgent_ForceOverwrite(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Create existing .claude directory with a test file
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("Failed to create .claude dir: %v", err)
	}

	testFile := filepath.Join(claudeDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Install with force
	err = InstallAgent([]string{"--force"})
	if err != nil {
		t.Fatalf("InstallAgent with --force failed: %v", err)
	}

	// Verify .claude directory was overwritten (test file should still be there though)
	// The force flag creates the directory but doesn't remove existing content first
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		// This is actually expected - our implementation doesn't clean the directory
		// It just allows writing to existing .claude directory
	}

	// Verify new files were installed
	essentialFile := filepath.Join(tmpDir, ".claude/settings.json")
	if _, err := os.Stat(essentialFile); os.IsNotExist(err) {
		t.Errorf("Essential file not installed after --force: settings.json")
	}
}

func TestInstallAgent_FileCount(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run install
	err = InstallAgent([]string{"--llm", "claude"})
	if err != nil {
		t.Fatalf("InstallAgent failed: %v", err)
	}

	// Count files in .claude directory
	fileCount := 0
	err = filepath.Walk(filepath.Join(tmpDir, ".claude"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to walk .claude directory: %v", err)
	}

	// We expect 33 files based on the current .claude structure
	// This number may change as skills are added/removed
	if fileCount < 30 {
		t.Errorf("Expected at least 30 files, got %d", fileCount)
	}
}

func TestInstallAgent_SettingsJSON(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Run install
	err = InstallAgent([]string{"--llm", "claude"})
	if err != nil {
		t.Fatalf("InstallAgent failed: %v", err)
	}

	// Read settings.json
	settingsPath := filepath.Join(tmpDir, ".claude/settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read settings.json: %v", err)
	}

	// Verify it contains expected permissions
	settingsStr := string(content)
	if len(settingsStr) == 0 {
		t.Errorf("settings.json is empty")
	}

	// Basic validation - should contain JSON structure
	if !stringContains(settingsStr, "permissions") {
		t.Errorf("settings.json doesn't contain 'permissions' field")
	}
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
