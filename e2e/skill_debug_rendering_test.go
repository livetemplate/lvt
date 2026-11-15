//go:build http

package e2e

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			"internal/parse/api.go",
			"internal/parse/eval.go",
			"internal/build/types.go",
			"internal/diff/tree_compare.go",
			"internal/diff/range_ops.go",
			"internal/observe/metrics.go",
		}

		for _, file := range serverFiles {
			fullPath := filepath.Join(livetemplateRoot, file)
			assert.FileExists(t, fullPath, "Server file %s should exist", file)
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
				file:      "internal/parse/api.go",
				functions: []string{"func Parse("},
			},
			{
				file:      "internal/parse/eval.go",
				functions: []string{"evalPipe"},
			},
			{
				file:      "internal/build/types.go",
				functions: []string{"TreeNode", "RangeData", "TreeMetadata"},
			},
			{
				file:      "internal/diff/tree_compare.go",
				functions: []string{"CompareTreesAndGetChangesWithPath", "ClientNeedsStatics"},
			},
		}

		for _, check := range functionChecks {
			fullPath := filepath.Join(livetemplateRoot, check.file)
			content, err := os.ReadFile(fullPath)
			if !assert.NoError(t, err, "Should read %s", check.file) {
				continue
			}

			contentStr := string(content)
			for _, fn := range check.functions {
				assert.True(t, strings.Contains(contentStr, fn),
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
			assert.FileExists(t, fullPath, "Client file %s should exist", file)
		}
	})

	// The skill points at these console entry points. The client's old
	// window.__wsMessages-style capture globals are gone, so WS delivery is
	// checked from frames instead (DevTools, or e2etest.RecordWSFrames).
	t.Run("ClientDebugEntryPointsExist", func(t *testing.T) {
		clientFile := filepath.Join(clientRoot, "livetemplate-client.ts")
		content, err := os.ReadFile(clientFile)
		require.NoError(t, err, "Should read livetemplate-client.ts")

		contentStr := string(content)

		// Patterns, not exact source, so formatting or a renamed local doesn't
		// fail the test while the entry point still works.
		entryPoints := map[string]*regexp.Regexp{
			"window.liveTemplateClient": regexp.MustCompile(`\.liveTemplateClient\s*=`),
			"getTreeState()":            regexp.MustCompile(`\bgetTreeState\s*\(`),
			"public ws":                 regexp.MustCompile(`public\s+ws\s*:`),
		}

		for name, re := range entryPoints {
			assert.True(t, re.MatchString(contentStr),
				"Client should contain debug entry point %s (/%s/)", name, re)
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
			if !assert.NoError(t, err, "Should read %s", check.file) {
				continue
			}

			contentStr := string(content)
			for _, fn := range check.functions {
				assert.True(t, strings.Contains(contentStr, fn),
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
