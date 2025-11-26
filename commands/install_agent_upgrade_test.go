package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallAgent_Upgrade(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-upgrade-test-*")
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

	// First, install the agent
	err = InstallAgent([]string{})
	if err != nil {
		t.Fatalf("Initial InstallAgent failed: %v", err)
	}

	// Verify installation
	claudeDir := filepath.Join(tmpDir, ".claude")
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Errorf(".claude directory was not created")
	}

	// Create a custom settings.local.json
	customSettings := `{"test": "custom settings"}`
	settingsPath := filepath.Join(claudeDir, "settings.local.json")
	if err := os.WriteFile(settingsPath, []byte(customSettings), 0644); err != nil {
		t.Fatalf("Failed to create custom settings: %v", err)
	}

	// Now upgrade
	err = InstallAgent([]string{"--upgrade"})
	if err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Verify .claude still exists
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Errorf(".claude directory missing after upgrade")
	}

	// Verify custom settings were preserved
	restoredSettings, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Errorf("Custom settings not restored: %v", err)
	} else if string(restoredSettings) != customSettings {
		t.Errorf("Custom settings not preserved. Got: %s, Want: %s", string(restoredSettings), customSettings)
	}

	// Verify agent files are present
	essentialFiles := []string{
		".claude/settings.json",
		".claude/agents/project-manager-backlog.md",
		".claude/skills/lvt/core/new-app.md",
	}

	for _, file := range essentialFiles {
		fullPath := filepath.Join(tmpDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Essential file missing after upgrade: %s", file)
		}
	}
}

func TestInstallAgent_UpgradeNoExisting(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-upgrade-test-*")
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

	// Try to upgrade without existing installation
	err = InstallAgent([]string{"--upgrade"})
	if err == nil {
		t.Errorf("Upgrade should fail when no existing installation found")
	}

	// Verify error message is helpful
	if err != nil {
		errMsg := err.Error()
		if errMsg != "no existing installation to upgrade" {
			t.Errorf("Expected 'no existing installation to upgrade' error, got: %s", errMsg)
		}
	}
}

func TestInstallAgent_UpgradePreservesOnlyLocal(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "lvt-agent-upgrade-test-*")
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

	// First install
	if err := InstallAgent([]string{}); err != nil {
		t.Fatalf("Initial install failed: %v", err)
	}

	// Modify settings.json (should be overwritten)
	claudeDir := filepath.Join(tmpDir, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(`{"modified": true}`), 0644); err != nil {
		t.Fatalf("Failed to modify settings.json: %v", err)
	}

	// Create custom settings.local.json (should be preserved)
	customSettings := `{"preserved": true}`
	settingsLocalPath := filepath.Join(claudeDir, "settings.local.json")
	if err := os.WriteFile(settingsLocalPath, []byte(customSettings), 0644); err != nil {
		t.Fatalf("Failed to create settings.local.json: %v", err)
	}

	// Upgrade
	if err := InstallAgent([]string{"--upgrade"}); err != nil {
		t.Fatalf("Upgrade failed: %v", err)
	}

	// Verify settings.json was overwritten (not modified)
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read settings.json: %v", err)
	}
	if stringContains(string(settingsContent), "modified") {
		t.Errorf("settings.json should have been overwritten, but still contains modified content")
	}

	// Verify settings.local.json was preserved
	settingsLocalContent, err := os.ReadFile(settingsLocalPath)
	if err != nil {
		t.Errorf("settings.local.json should be preserved: %v", err)
	} else if string(settingsLocalContent) != customSettings {
		t.Errorf("settings.local.json not preserved. Got: %s, Want: %s", string(settingsLocalContent), customSettings)
	}
}
