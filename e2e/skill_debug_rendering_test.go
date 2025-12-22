//go:build http

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDebugRenderingSkill_ReferencesValid verifies that the debug-rendering skill
// references valid files and functions in the livetemplate ecosystem.
// This test should be run when the livetemplate codebase changes to ensure
// the skill documentation stays in sync.
func TestDebugRenderingSkill_ReferencesValid(t *testing.T) {
	// Find the livetemplate monorepo root
	// This works both from main lvt dir and from worktrees
	var livetemplateRoot, clientRoot string

	// Try direct sibling first (main repo: lvt/../livetemplate)
	lvtRoot, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err, "Should resolve lvt root")

	livetemplateRoot = filepath.Join(lvtRoot, "livetemplate")
	clientRoot = filepath.Join(lvtRoot, "client")

	// If not found, try worktree path (lvt/.worktrees/xxx/../../../livetemplate)
	if _, err := os.Stat(livetemplateRoot); os.IsNotExist(err) {
		monorepoRoot, err := filepath.Abs(filepath.Join("..", "..", "..", ".."))
		require.NoError(t, err, "Should resolve monorepo root")
		livetemplateRoot = filepath.Join(monorepoRoot, "livetemplate")
		clientRoot = filepath.Join(monorepoRoot, "client")
	}

	// Skip if not in monorepo context
	if _, err := os.Stat(livetemplateRoot); os.IsNotExist(err) {
		t.Skip("Skipping: not in livetemplate monorepo context (livetemplate dir not found)")
	}
	if _, err := os.Stat(clientRoot); os.IsNotExist(err) {
		t.Skip("Skipping: not in livetemplate monorepo context (client dir not found)")
	}

	t.Run("ServerSideFilesExist", func(t *testing.T) {
		serverFiles := []string{
			"template.go",
			"dispatch.go",
			"internal/parse/parse.go",
			"internal/build/types.go",
			"internal/diff/tree_compare.go",
			"internal/diff/range_ops.go",
			"internal/observe/metrics.go",
		}

		for _, file := range serverFiles {
			fullPath := filepath.Join(livetemplateRoot, file)
			require.FileExists(t, fullPath, "Server file %s should exist", file)
		}
	})

	t.Run("ServerSideFunctionsExist", func(t *testing.T) {
		functionChecks := []struct {
			file      string
			functions []string
		}{
			{
				file:      "template.go",
				functions: []string{"buildTree", "ExecuteUpdates", "Clone"},
			},
			{
				file:      "dispatch.go",
				functions: []string{"DispatchError", "DispatchWithState"},
			},
			{
				file:      "internal/parse/parse.go",
				functions: []string{"Parse"},
			},
			{
				file:      "internal/build/types.go",
				functions: []string{"TreeNode", "RangeData"},
			},
			{
				file:      "internal/diff/tree_compare.go",
				functions: []string{"CompareTreesAndGetChangesWithPath"},
			},
		}

		for _, check := range functionChecks {
			fullPath := filepath.Join(livetemplateRoot, check.file)
			content, err := os.ReadFile(fullPath)
			require.NoError(t, err, "Should read %s", check.file)

			contentStr := string(content)
			for _, fn := range check.functions {
				require.True(t, strings.Contains(contentStr, fn),
					"File %s should contain function/type '%s'", check.file, fn)
			}
		}
	})

	t.Run("ClientSideFilesExist", func(t *testing.T) {
		clientFiles := []string{
			"livetemplate-client.ts",
			"state/tree-renderer.ts",
			"transport/websocket.ts",
			"dom/event-delegation.ts",
		}

		for _, file := range clientFiles {
			fullPath := filepath.Join(clientRoot, file)
			require.FileExists(t, fullPath, "Client file %s should exist", file)
		}
	})

	t.Run("ClientDebugVariablesExist", func(t *testing.T) {
		// Check main client file for debug variables
		clientFile := filepath.Join(clientRoot, "livetemplate-client.ts")
		content, err := os.ReadFile(clientFile)
		require.NoError(t, err, "Should read livetemplate-client.ts")

		contentStr := string(content)

		debugVars := []string{
			"__wsMessages",
			"__lastWSMessage",
			"__lvtSendCalled",
			"__lvtMessageAction",
			"__lvtSendPath",
		}

		for _, varName := range debugVars {
			require.True(t, strings.Contains(contentStr, varName),
				"Client should contain debug variable '%s'", varName)
		}
	})

	t.Run("ClientFunctionsExist", func(t *testing.T) {
		functionChecks := []struct {
			file      string
			functions []string
		}{
			{
				file:      "livetemplate-client.ts",
				functions: []string{"updateDOM", "handleWebSocketPayload"},
			},
			{
				file:      "state/tree-renderer.ts",
				functions: []string{"applyUpdate", "reconstructFromTree", "deepMergeTreeNodes"},
			},
			{
				file:      "transport/websocket.ts",
				functions: []string{"WebSocketTransport", "WebSocketManager"},
			},
		}

		for _, check := range functionChecks {
			fullPath := filepath.Join(clientRoot, check.file)
			content, err := os.ReadFile(fullPath)
			require.NoError(t, err, "Should read %s", check.file)

			contentStr := string(content)
			for _, fn := range check.functions {
				require.True(t, strings.Contains(contentStr, fn),
					"File %s should contain function/class '%s'", check.file, fn)
			}
		}
	})
}

// TestDebugRenderingSkill_SkillFileExists verifies the skill file exists and has proper format
func TestDebugRenderingSkill_SkillFileExists(t *testing.T) {
	// Find skill file relative to test
	skillFile := filepath.Join("..", "commands", "claude_resources", "skills", "debug-rendering", "SKILL.md")

	require.FileExists(t, skillFile, "debug-rendering skill file should exist")

	content, err := os.ReadFile(skillFile)
	require.NoError(t, err, "Should read skill file")

	contentStr := string(content)

	// Verify frontmatter
	require.Contains(t, contentStr, "---", "Skill should have frontmatter")
	require.Contains(t, contentStr, "name: lvt-debug-rendering", "Skill should have correct name")
	require.Contains(t, contentStr, "description:", "Skill should have description")
	require.Contains(t, contentStr, "category: maintenance", "Skill should be in maintenance category")

	// Verify key sections exist
	sections := []string{
		"## Quick Symptom Lookup",
		"## Rendering Pipeline Overview",
		"## Priority Issue Workflows",
		"## Server-Side Debugging Guide",
		"## Client-Side Debugging Guide",
		"## Debug Commands & Techniques",
		"## Source Code References",
		"## Maintenance",
	}

	for _, section := range sections {
		require.Contains(t, contentStr, section, "Skill should contain section: %s", section)
	}

	// Verify priority issues are documented
	priorityIssues := []string{
		"Template Not Updating",
		"Partial/Broken Renders",
		"Action Dispatch Errors",
	}

	for _, issue := range priorityIssues {
		require.Contains(t, contentStr, issue, "Skill should document priority issue: %s", issue)
	}
}

// TestDebugRenderingSkill_PipelinePhases verifies the 5-phase pipeline is documented
func TestDebugRenderingSkill_PipelinePhases(t *testing.T) {
	skillFile := filepath.Join("..", "commands", "claude_resources", "skills", "debug-rendering", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	require.NoError(t, err, "Should read skill file")

	contentStr := string(content)

	// Verify all 5 phases are documented
	phases := []string{
		"Phase 1: Parse",
		"Phase 2: Build",
		"Phase 3: Diff",
		"Phase 4: Render",
		"Phase 5: Send",
	}

	for _, phase := range phases {
		require.Contains(t, contentStr, phase, "Skill should document %s", phase)
	}
}
