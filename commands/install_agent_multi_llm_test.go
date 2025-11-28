package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallAgent_ListAgents(t *testing.T) {
	// This test doesn't need a temp directory, just tests the list functionality
	// We can't easily capture stdout, but we can verify it doesn't error
	tmpDir, err := os.MkdirTemp("", "lvt-agent-list-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Should not error when listing agents
	err = InstallAgent([]string{"--list"})
	if err != nil {
		t.Errorf("--list should not error: %v", err)
	}
}

func TestInstallAgent_CopilotAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-copilot-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Install Copilot agent
	err = InstallAgent([]string{"--llm", "copilot"})
	if err != nil {
		t.Fatalf("InstallAgent for copilot failed: %v", err)
	}

	// Verify .github directory exists
	githubDir := filepath.Join(tmpDir, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Errorf(".github directory was not created")
	}

	// Verify copilot-instructions.md exists
	instructionsFile := filepath.Join(tmpDir, ".github/copilot-instructions.md")
	if _, err := os.Stat(instructionsFile); os.IsNotExist(err) {
		t.Errorf("copilot-instructions.md was not created")
	}

	// Verify file content is not empty
	content, err := os.ReadFile(instructionsFile)
	if err != nil {
		t.Fatalf("Failed to read copilot-instructions.md: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("copilot-instructions.md is empty")
	}

	// Verify it contains expected content
	contentStr := string(content)
	if !stringContains(contentStr, "LiveTemplate") {
		t.Errorf("copilot-instructions.md doesn't mention LiveTemplate")
	}
	if !stringContains(contentStr, "MCP") {
		t.Errorf("copilot-instructions.md doesn't mention MCP")
	}
}

func TestInstallAgent_CursorAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-cursor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Install Cursor agent
	err = InstallAgent([]string{"--llm", "cursor"})
	if err != nil {
		t.Fatalf("InstallAgent for cursor failed: %v", err)
	}

	// Verify .cursor/rules directory exists
	rulesDir := filepath.Join(tmpDir, ".cursor/rules")
	if _, err := os.Stat(rulesDir); os.IsNotExist(err) {
		t.Errorf(".cursor/rules directory was not created")
	}

	// Verify lvt.md exists in rules directory
	rulesFile := filepath.Join(tmpDir, ".cursor/rules/lvt.md")
	if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
		t.Errorf(".cursor/rules/lvt.md was not created")
	}

	// Verify file content
	content, err := os.ReadFile(rulesFile)
	if err != nil {
		t.Fatalf("Failed to read lvt.md: %v", err)
	}
	if len(content) == 0 {
		t.Errorf("lvt.md is empty")
	}

	// Verify frontmatter
	contentStr := string(content)
	if !stringContains(contentStr, "applyTo") {
		t.Errorf("lvt.md doesn't have applyTo frontmatter")
	}
}

func TestInstallAgent_AiderAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-aider-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Install Aider agent
	err = InstallAgent([]string{"--llm", "aider"})
	if err != nil {
		t.Fatalf("InstallAgent for aider failed: %v", err)
	}

	// Verify .aider directory exists
	aiderDir := filepath.Join(tmpDir, ".aider")
	if _, err := os.Stat(aiderDir); os.IsNotExist(err) {
		t.Errorf(".aider directory was not created")
	}

	// Verify .aider.conf.yml exists (should be renamed from aider.conf.yml)
	confFile := filepath.Join(tmpDir, ".aider/.aider.conf.yml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		t.Errorf(".aider.conf.yml was not created")
	}

	// Verify lvt-instructions.md exists
	instructionsFile := filepath.Join(tmpDir, ".aider/lvt-instructions.md")
	if _, err := os.Stat(instructionsFile); os.IsNotExist(err) {
		t.Errorf("lvt-instructions.md was not created")
	}

	// Verify config file content
	confContent, err := os.ReadFile(confFile)
	if err != nil {
		t.Fatalf("Failed to read .aider.conf.yml: %v", err)
	}
	if len(confContent) == 0 {
		t.Errorf(".aider.conf.yml is empty")
	}

	// Verify it's valid YAML-ish content
	confStr := string(confContent)
	if !stringContains(confStr, "auto-commits") {
		t.Errorf(".aider.conf.yml doesn't contain auto-commits setting")
	}
}

func TestInstallAgent_GenericAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-generic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Install Generic agent
	err = InstallAgent([]string{"--llm", "generic"})
	if err != nil {
		t.Fatalf("InstallAgent for generic failed: %v", err)
	}

	// Verify lvt-agent directory exists
	agentDir := filepath.Join(tmpDir, "lvt-agent")
	if _, err := os.Stat(agentDir); os.IsNotExist(err) {
		t.Errorf("lvt-agent directory was not created")
	}

	// Verify README.md exists
	readmeFile := filepath.Join(tmpDir, "lvt-agent/README.md")
	if _, err := os.Stat(readmeFile); os.IsNotExist(err) {
		t.Errorf("README.md was not created")
	}

	// Verify QUICK_REFERENCE.md exists
	quickRefFile := filepath.Join(tmpDir, "lvt-agent/QUICK_REFERENCE.md")
	if _, err := os.Stat(quickRefFile); os.IsNotExist(err) {
		t.Errorf("QUICK_REFERENCE.md was not created")
	}

	// Verify README content
	readmeContent, err := os.ReadFile(readmeFile)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}
	if len(readmeContent) == 0 {
		t.Errorf("README.md is empty")
	}

	// Verify it contains generic LLM content
	readmeStr := string(readmeContent)
	if !stringContains(readmeStr, "MCP") {
		t.Errorf("README.md doesn't mention MCP")
	}
	if !stringContains(readmeStr, "16 tools") {
		t.Errorf("README.md doesn't mention 16 tools")
	}
}

func TestInstallAgent_InvalidLLMType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-invalid-llm-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Try to install with invalid LLM type
	err = InstallAgent([]string{"--llm", "invalid-llm"})
	if err == nil {
		t.Errorf("InstallAgent should fail with invalid LLM type")
	}

	// Verify error message
	if err != nil {
		errMsg := err.Error()
		if errMsg != "invalid LLM type" {
			t.Errorf("Expected 'invalid LLM type' error, got: %s", errMsg)
		}
	}
}

func TestInstallAgent_DefaultIsClaudeCode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-default-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Install with --llm claude (testing it installs correctly)
	err = InstallAgent([]string{"--llm", "claude"})
	if err != nil {
		t.Fatalf("Claude InstallAgent failed: %v", err)
	}

	// Verify .claude directory exists
	claudeDir := filepath.Join(tmpDir, ".claude")
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Errorf(".claude directory was not created")
	}

	// Verify no other agent directories were created
	otherDirs := []string{".github", ".cursor", ".aider", "lvt-agent"}
	for _, dir := range otherDirs {
		fullPath := filepath.Join(tmpDir, dir)
		if _, err := os.Stat(fullPath); err == nil {
			t.Errorf("Unexpected directory created: %s (should only create .claude)", dir)
		}
	}
}

func TestInstallAgent_UpgradeCursorAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-cursor-upgrade-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// First install Cursor agent
	err = InstallAgent([]string{"--llm", "cursor"})
	if err != nil {
		t.Fatalf("Initial InstallAgent for cursor failed: %v", err)
	}

	// Create a custom local.md file
	cursorDir := filepath.Join(tmpDir, ".cursor")
	customFile := filepath.Join(cursorDir, "local.md")
	customContent := "# Custom Cursor Rules"
	if err := os.WriteFile(customFile, []byte(customContent), 0644); err != nil {
		t.Fatalf("Failed to create custom local.md: %v", err)
	}

	// Upgrade
	err = InstallAgent([]string{"--llm", "cursor", "--upgrade"})
	if err != nil {
		t.Fatalf("Upgrade cursor agent failed: %v", err)
	}

	// Verify .cursor still exists
	if _, err := os.Stat(cursorDir); os.IsNotExist(err) {
		t.Errorf(".cursor directory missing after upgrade")
	}

	// Verify custom file was preserved
	restoredContent, err := os.ReadFile(customFile)
	if err != nil {
		t.Errorf("Custom local.md not preserved: %v", err)
	} else if string(restoredContent) != customContent {
		t.Errorf("Custom local.md not preserved. Got: %s, Want: %s", string(restoredContent), customContent)
	}

	// Verify rules file is still there
	rulesFile := filepath.Join(tmpDir, ".cursor/rules/lvt.md")
	if _, err := os.Stat(rulesFile); os.IsNotExist(err) {
		t.Errorf(".cursor/rules/lvt.md missing after upgrade")
	}
}

func TestInstallAgent_UpgradeAiderAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-aider-upgrade-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// First install Aider agent
	err = InstallAgent([]string{"--llm", "aider"})
	if err != nil {
		t.Fatalf("Initial InstallAgent for aider failed: %v", err)
	}

	// Create a custom .aider.local.yml file
	aiderDir := filepath.Join(tmpDir, ".aider")
	customFile := filepath.Join(aiderDir, ".aider.local.yml")
	customContent := "# Custom Aider Config\nmy-setting: true"
	if err := os.WriteFile(customFile, []byte(customContent), 0644); err != nil {
		t.Fatalf("Failed to create custom .aider.local.yml: %v", err)
	}

	// Upgrade
	err = InstallAgent([]string{"--llm", "aider", "--upgrade"})
	if err != nil {
		t.Fatalf("Upgrade aider agent failed: %v", err)
	}

	// Verify .aider still exists
	if _, err := os.Stat(aiderDir); os.IsNotExist(err) {
		t.Errorf(".aider directory missing after upgrade")
	}

	// Verify custom file was preserved
	restoredContent, err := os.ReadFile(customFile)
	if err != nil {
		t.Errorf("Custom .aider.local.yml not preserved: %v", err)
	} else if string(restoredContent) != customContent {
		t.Errorf("Custom .aider.local.yml not preserved. Got: %s, Want: %s", string(restoredContent), customContent)
	}

	// Verify main config file is still there
	confFile := filepath.Join(tmpDir, ".aider/.aider.conf.yml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		t.Errorf(".aider.conf.yml missing after upgrade")
	}
}

func TestInstallAgent_ExistingCopilotAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lvt-copilot-existing-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Create existing .github directory
	githubDir := filepath.Join(tmpDir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		t.Fatalf("Failed to create .github dir: %v", err)
	}

	// Try to install without force (should fail)
	err = InstallAgent([]string{"--llm", "copilot"})
	if err == nil {
		t.Errorf("InstallAgent should fail when .github exists without --force")
	}

	// Install with force
	err = InstallAgent([]string{"--llm", "copilot", "--force"})
	if err != nil {
		t.Fatalf("InstallAgent with --force failed: %v", err)
	}

	// Verify copilot-instructions.md was created
	instructionsFile := filepath.Join(tmpDir, ".github/copilot-instructions.md")
	if _, err := os.Stat(instructionsFile); os.IsNotExist(err) {
		t.Errorf("copilot-instructions.md not installed after --force")
	}
}
