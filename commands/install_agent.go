package commands

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed claude_resources
var claudeResources embed.FS

// InstallAgent installs the Claude Code agent and skills into the current project
func InstallAgent(args []string) error {
	force := false
	targetDir := ".claude"

	// Parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--force" || args[i] == "-f" {
			force = true
		}
	}

	// Check if .claude directory already exists
	if _, err := os.Stat(targetDir); err == nil && !force {
		fmt.Println("⚠️  .claude directory already exists")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  - Use --force to overwrite existing installation")
		fmt.Println("  - Remove .claude manually and run again")
		fmt.Println("  - Keep your existing installation")
		return fmt.Errorf("installation cancelled")
	}

	fmt.Println("Installing Claude Code agent and skills...")
	fmt.Println()

	// Create .claude directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	// Copy embedded resources
	installed := 0
	err := fs.WalkDir(claudeResources, "claude_resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path (remove claude_resources prefix)
		relPath := strings.TrimPrefix(path, "claude_resources/")
		if relPath == "" {
			return nil // Skip root directory
		}

		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file from embedded FS
		content, err := fs.ReadFile(claudeResources, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write file to target
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		installed++
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to install agent: %w", err)
	}

	fmt.Printf("✅ Successfully installed %d files!\n", installed)
	fmt.Println()
	fmt.Println("📁 Installed to: .claude/")
	fmt.Println()
	fmt.Println("📚 What's included:")
	fmt.Println("  • 20+ skills for lvt commands and workflows")
	fmt.Println("  • Project management agent")
	fmt.Println("  • Permission settings for safe operation")
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Println("  1. Open Claude Code in this directory:")
	fmt.Println("     claude")
	fmt.Println()
	fmt.Println("  2. Try asking Claude:")
	fmt.Println("     • \"Add a posts resource with title and content\"")
	fmt.Println("     • \"Generate authentication system\"")
	fmt.Println("     • \"Create a quickstart blog app\"")
	fmt.Println()
	fmt.Println("💡 The agent will guide you through workflows and best practices!")
	fmt.Println()
	fmt.Println("📖 Learn more: docs/AGENT_USAGE_GUIDE.md")
	fmt.Println()

	return nil
}
