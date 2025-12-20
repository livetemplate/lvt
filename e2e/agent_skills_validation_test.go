package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/commands"
	"github.com/stretchr/testify/require"
)

// TestClaudeAgent_Installation verifies that Claude agent installation includes all documented components
func TestClaudeAgent_Installation(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	// Install Claude agent
	err = commands.InstallAgent([]string{"--llm", "claude"})
	require.NoError(t, err, "Claude agent installation should succeed")

	// Verify .claude/agents/lvt-assistant.md exists
	agentFile := filepath.Join(tmpDir, ".claude", "agents", "lvt-assistant.md")
	require.FileExists(t, agentFile, "lvt-assistant agent file should exist")

	// Verify .claude/skills/ directory exists
	skillsDir := filepath.Join(tmpDir, ".claude", "skills")
	require.DirExists(t, skillsDir, "skills directory should exist")

	// Verify settings.json exists
	settingsFile := filepath.Join(tmpDir, ".claude", "settings.json")
	require.FileExists(t, settingsFile, "settings.json should exist")
}

// TestClaudeAgent_AllSkillsExist verifies that all documented skills exist
func TestClaudeAgent_AllSkillsExist(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	// Install Claude agent
	err = commands.InstallAgent([]string{"--llm", "claude"})
	require.NoError(t, err)

	skillsDir := filepath.Join(tmpDir, ".claude", "skills")

	// List of skills documented in AGENT_USAGE_GUIDE.md
	// Each skill is a directory with a SKILL.md file
	documentedSkills := []string{
		// Core Skills (14)
		"new-app",
		"add-resource",
		"add-view",
		"add-migration",
		"gen-schema",
		"gen-auth",
		"resource-inspect",
		"manage-kits",
		"validate-templates",
		"run-and-test",
		"customize",
		"seed-data",
		"deploy",
		"manage-env",
		// Workflow Skills (3)
		"quickstart",
		"production-ready",
		"add-related-resources",
		// Maintenance Skills (4)
		"analyze",
		"suggest",
		"troubleshoot",
		"debug-rendering",
		// Meta Skill (1)
		"add-skill",
	}

	for _, skillName := range documentedSkills {
		skillDir := filepath.Join(skillsDir, skillName)
		skillFile := filepath.Join(skillDir, "SKILL.md")

		require.DirExists(t, skillDir, "Skill directory %s should exist", skillName)
		require.FileExists(t, skillFile, "Skill %s should have SKILL.md at %s", skillName, skillFile)

		// Read skill file
		content, err := os.ReadFile(skillFile)
		require.NoError(t, err, "Should be able to read skill file %s", skillName)

		// Verify it has frontmatter with name
		contentStr := string(content)
		require.Contains(t, contentStr, "---", "Skill %s should have frontmatter", skillName)
		require.Contains(t, contentStr, "name:", "Skill %s should have name in frontmatter", skillName)
	}
}

// TestClaudeAgent_AgentMetadata verifies agent frontmatter is correct
func TestClaudeAgent_AgentMetadata(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	// Install Claude agent
	err = commands.InstallAgent([]string{"--llm", "claude"})
	require.NoError(t, err)

	// Read .claude/agents/lvt-assistant.md
	agentFile := filepath.Join(tmpDir, ".claude", "agents", "lvt-assistant.md")
	content, err := os.ReadFile(agentFile)
	require.NoError(t, err, "Should be able to read agent file")

	contentStr := string(content)

	// Verify frontmatter has name: "lvt-assistant"
	require.Contains(t, contentStr, "---", "Agent should have frontmatter")
	require.Contains(t, contentStr, "name: lvt-assistant", "Agent should have correct name")
	require.Contains(t, contentStr, "description:", "Agent should have description")

	// Verify it mentions LiveTemplate
	require.Contains(t, contentStr, "LiveTemplate", "Agent should mention LiveTemplate")
}

// TestClaudeAgent_SkillInvocationSyntax verifies skill names match directory structure
func TestClaudeAgent_SkillInvocationSyntax(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	// Install Claude agent
	err = commands.InstallAgent([]string{"--llm", "claude"})
	require.NoError(t, err)

	skillsDir := filepath.Join(tmpDir, ".claude", "skills")

	// Read all skill directories
	entries, err := os.ReadDir(skillsDir)
	require.NoError(t, err, "Should be able to read skills directory")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillName := entry.Name()
		skillFile := filepath.Join(skillsDir, skillName, "SKILL.md")

		// Skip if no SKILL.md (might be other directories like docs)
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			continue
		}

		// Read skill file
		content, err := os.ReadFile(skillFile)
		require.NoError(t, err, "Should be able to read %s", skillFile)

		contentStr := string(content)

		// Parse frontmatter to get the skill name
		lines := strings.Split(contentStr, "\n")
		inFrontmatter := false
		var frontmatterName string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "---" {
				if inFrontmatter {
					break // End of frontmatter
				}
				inFrontmatter = true
				continue
			}

			if inFrontmatter && strings.HasPrefix(line, "name:") {
				frontmatterName = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
				break
			}
		}

		// Verify frontmatter name follows lvt-skill-name pattern
		expectedName := "lvt-" + skillName

		require.Equal(t, expectedName, frontmatterName,
			"Skill %s: frontmatter name '%s' should match expected format 'lvt-%s'",
			skillName, frontmatterName, skillName)
	}
}

// TestClaudeAgent_SkillCount verifies we have the expected number of skills
func TestClaudeAgent_SkillCount(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(tmpDir)
	require.NoError(t, err)

	// Install Claude agent
	err = commands.InstallAgent([]string{"--llm", "claude"})
	require.NoError(t, err)

	skillsDir := filepath.Join(tmpDir, ".claude", "skills")

	// Count skill directories with SKILL.md files
	skillCount := 0
	entries, err := os.ReadDir(skillsDir)
	require.NoError(t, err, "Should be able to read skills directory")

	for _, entry := range entries {
		if entry.IsDir() {
			skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
			if _, err := os.Stat(skillFile); err == nil {
				skillCount++
			}
		}
	}

	// AGENT_USAGE_GUIDE.md documents:
	// Core (14) + Workflow (4 - added brainstorm) + Maintenance (4) + Meta (1) = 23 skills
	require.GreaterOrEqual(t, skillCount, 23,
		"Should have at least 23 skills (documented in AGENT_USAGE_GUIDE.md), found %d", skillCount)
}
