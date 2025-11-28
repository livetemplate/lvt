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

//go:embed agent_resources/copilot
var copilotResources embed.FS

//go:embed agent_resources/cursor
var cursorResources embed.FS

//go:embed agent_resources/aider
var aiderResources embed.FS

//go:embed agent_resources/generic
var genericResources embed.FS

// AgentConfig holds configuration for different agent types
type AgentConfig struct {
	Name        string
	TargetDir   string
	ResourceFS  embed.FS
	ResourceDir string
	Description string
	NextSteps   []string
}

var agentConfigs = map[string]AgentConfig{
	"claude": {
		Name:        "Claude Code",
		TargetDir:   ".claude",
		ResourceFS:  claudeResources,
		ResourceDir: "claude_resources",
		Description: "Claude Code agent with skills and project management",
		NextSteps: []string{
			"Open Claude Code in this directory: claude",
			"Try: \"Add a posts resource with title and content\"",
			"Try: \"Generate authentication system\"",
		},
	},
	"copilot": {
		Name:        "GitHub Copilot",
		TargetDir:   ".github",
		ResourceFS:  copilotResources,
		ResourceDir: "agent_resources/copilot",
		Description: "GitHub Copilot instructions for LiveTemplate development",
		NextSteps: []string{
			"Copilot will automatically use these instructions",
			"Open a file and start coding",
			"Use @workspace to ask questions",
		},
	},
	"cursor": {
		Name:        "Cursor",
		TargetDir:   ".cursor",
		ResourceFS:  cursorResources,
		ResourceDir: "agent_resources/cursor",
		Description: "Cursor AI rules for LiveTemplate development",
		NextSteps: []string{
			"Open project in Cursor",
			"Rules apply automatically to *.go files",
			"Use Composer/Agent mode for best results",
		},
	},
	"aider": {
		Name:        "Aider",
		TargetDir:   ".aider",
		ResourceFS:  aiderResources,
		ResourceDir: "agent_resources/aider",
		Description: "Aider CLI configuration for LiveTemplate",
		NextSteps: []string{
			"Run: aider",
			"Configuration loads automatically",
			"Use lvt commands via Aider",
		},
	},
	"generic": {
		Name:        "Generic LLM",
		TargetDir:   "lvt-agent",
		ResourceFS:  genericResources,
		ResourceDir: "agent_resources/generic",
		Description: "LLM-agnostic documentation and guides",
		NextSteps: []string{
			"See lvt-agent/README.md for integration guide",
			"Adapt to your LLM's tool format",
			"Use lvt mcp-server for MCP-enabled LLMs",
		},
	},
}

// InstallAgent installs an AI agent into the current project
func InstallAgent(args []string) error {
	force := false
	upgrade := false
	llm := "claude" // Default to Claude Code

	// Parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--force" || args[i] == "-f" {
			force = true
		} else if args[i] == "--upgrade" || args[i] == "-u" {
			upgrade = true
		} else if args[i] == "--llm" && i+1 < len(args) {
			llm = args[i+1]
			i++
		} else if args[i] == "--list" {
			listAgents()
			return nil
		}
	}

	// Validate LLM type
	config, ok := agentConfigs[llm]
	if !ok {
		fmt.Printf("❌ Unknown LLM type: %s\n\n", llm)
		fmt.Println("Available LLM types:")
		for name, cfg := range agentConfigs {
			fmt.Printf("  • %s - %s\n", name, cfg.Description)
		}
		fmt.Println()
		fmt.Println("Usage: lvt install-agent --llm <type>")
		return fmt.Errorf("invalid LLM type")
	}

	targetDir := config.TargetDir

	// Check if target directory already exists
	if _, err := os.Stat(targetDir); err == nil && !force && !upgrade {
		fmt.Printf("⚠️  %s directory already exists\n", targetDir)
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  - Use --upgrade to update agent and preserve custom settings")
		fmt.Println("  - Use --force to overwrite existing installation")
		fmt.Printf("  - Remove %s manually and run again\n", targetDir)
		fmt.Println("  - Keep your existing installation")
		return fmt.Errorf("installation cancelled")
	}

	// Handle upgrade mode
	if upgrade {
		return upgradeAgent(targetDir, config)
	}

	fmt.Printf("Installing %s agent...\n", config.Name)
	fmt.Println()

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", targetDir, err)
	}

	// Copy embedded resources
	installed := 0
	err := fs.WalkDir(config.ResourceFS, config.ResourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path (remove resource dir prefix)
		relPath := strings.TrimPrefix(path, config.ResourceDir+"/")
		if relPath == "" || relPath == config.ResourceDir {
			return nil // Skip root directory
		}

		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file from embedded FS
		content, err := fs.ReadFile(config.ResourceFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Special case: rename aider.conf.yml to .aider.conf.yml
		if filepath.Base(relPath) == "aider.conf.yml" {
			targetPath = filepath.Join(filepath.Dir(targetPath), ".aider.conf.yml")
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
	fmt.Printf("📁 Installed to: %s/\n", targetDir)
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	for i, step := range config.NextSteps {
		fmt.Printf("  %d. %s\n", i+1, step)
	}
	fmt.Println()
	fmt.Println("💡 The agent will guide you through LiveTemplate workflows!")
	fmt.Println()

	return nil
}

// upgradeAgent upgrades an existing agent installation while preserving user settings
func upgradeAgent(targetDir string, config AgentConfig) error {
	fmt.Printf("🔄 Upgrading %s agent...\n", config.Name)
	fmt.Println()

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Println("❌ No existing installation found")
		fmt.Println("   Run 'lvt install-agent' instead")
		return fmt.Errorf("no existing installation to upgrade")
	}

	// Backup user settings for Claude (settings.local.json) or other config files
	var settingsFiles = []string{"settings.local.json", ".aider.local.yml", "local.md"}
	backups := make(map[string][]byte)

	for _, filename := range settingsFiles {
		settingsPath := filepath.Join(targetDir, filename)
		if data, err := os.ReadFile(settingsPath); err == nil {
			backups[filename] = data
			fmt.Printf("📦 Backing up %s...\n", filename)
		}
	}

	// Count existing files before upgrade
	oldFileCount := 0
	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			oldFileCount++
		}
		return nil
	})

	// Remove old installation
	fmt.Println("🗑️  Removing old installation...")
	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove old installation: %w", err)
	}

	// Create fresh directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", targetDir, err)
	}

	// Install new agent files
	fmt.Println("📥 Installing updated agent...")
	installed := 0
	err := fs.WalkDir(config.ResourceFS, config.ResourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, config.ResourceDir+"/")
		if relPath == "" || relPath == config.ResourceDir {
			return nil
		}

		targetPath := filepath.Join(targetDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		content, err := fs.ReadFile(config.ResourceFS, path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Special case: rename aider.conf.yml to .aider.conf.yml
		if filepath.Base(relPath) == "aider.conf.yml" {
			targetPath = filepath.Join(filepath.Dir(targetPath), ".aider.conf.yml")
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		installed++
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to install agent: %w", err)
	}

	// Restore user settings if they existed
	if len(backups) > 0 {
		fmt.Println("♻️  Restoring your custom settings...")
		for filename, data := range backups {
			settingsPath := filepath.Join(targetDir, filename)
			if err := os.WriteFile(settingsPath, data, 0644); err != nil {
				fmt.Printf("⚠️  Warning: failed to restore %s: %v\n", filename, err)
			}
		}
	}

	// Show upgrade summary
	fmt.Println()
	fmt.Printf("✅ Successfully upgraded %s! (%d → %d files)\n", config.Name, oldFileCount, installed)
	fmt.Println()
	fmt.Println("📚 What's new:")
	fmt.Println("  • Latest documentation and improvements")
	fmt.Println("  • Updated workflows and best practices")
	fmt.Println("  • Bug fixes and enhancements")
	if len(backups) > 0 {
		fmt.Println("  • Your custom settings preserved")
	}
	fmt.Println()

	return nil
}

// listAgents lists all available agent types
func listAgents() {
	fmt.Println("Available AI agents for LiveTemplate:")
	fmt.Println()

	// Order: Claude (default), then alphabetically
	order := []string{"claude", "aider", "copilot", "cursor", "generic"}

	for _, name := range order {
		config := agentConfigs[name]
		fmt.Printf("  %s\n", name)
		fmt.Printf("    Name: %s\n", config.Name)
		fmt.Printf("    Description: %s\n", config.Description)
		fmt.Printf("    Installs to: %s/\n", config.TargetDir)
		fmt.Println()
	}

	fmt.Println("Usage:")
	fmt.Println("  lvt install-agent                    # Install Claude Code (default)")
	fmt.Println("  lvt install-agent --llm copilot      # Install GitHub Copilot")
	fmt.Println("  lvt install-agent --llm cursor       # Install Cursor")
	fmt.Println("  lvt install-agent --llm aider        # Install Aider")
	fmt.Println("  lvt install-agent --llm generic      # Install generic LLM docs")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --force, -f      Overwrite existing installation")
	fmt.Println("  --upgrade, -u    Upgrade existing installation")
	fmt.Println("  --list           Show this list")
	fmt.Println()
}
