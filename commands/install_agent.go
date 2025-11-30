package commands

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		Description: "Claude Code agent for LiveTemplate development (22 skills)",
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
	llm := ""
	hasLLMFlag := false

	// Parse flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--force" || args[i] == "-f" {
			force = true
		} else if args[i] == "--upgrade" || args[i] == "-u" {
			upgrade = true
		} else if args[i] == "--llm" && i+1 < len(args) {
			llm = args[i+1]
			hasLLMFlag = true
			i++
		} else if args[i] == "--list" {
			listAgents()
			return nil
		}
	}

	// If no --llm flag provided
	if !hasLLMFlag {
		// If upgrade mode, try to auto-detect installed agent
		if upgrade {
			detected := detectInstalledAgent()
			if detected != "" {
				llm = detected
			} else {
				return fmt.Errorf("no existing installation to upgrade")
			}
		} else {
			// Show interactive menu for new installations
			selectedLLM, err := selectAgentInteractive()
			if err != nil {
				return err
			}
			llm = selectedLLM
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
		action, err := handleExistingInstallation(targetDir, config.Name)
		if err != nil {
			return err
		}

		switch action {
		case "upgrade":
			upgrade = true
		case "force":
			force = true
		case "cancel":
			return fmt.Errorf("installation cancelled")
		}
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

	// Write version marker file
	versionFile := filepath.Join(targetDir, ".version")
	currentVersion := getCurrentVersion()
	if err := os.WriteFile(versionFile, []byte(currentVersion+"\n"), 0644); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("⚠️  Warning: Could not write version file: %v\n", err)
	}

	fmt.Printf("✅ Successfully installed %d files!\n", installed)
	fmt.Println()
	fmt.Printf("📁 Installed to: %s/\n", targetDir)
	fmt.Printf("📦 Version: %s\n", currentVersion)
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

	// Get version information
	oldVersion := getInstalledVersion(targetDir)
	newVersion := getCurrentVersion()

	fmt.Println("📊 Version Information:")
	fmt.Printf("   Current: %s\n", oldVersion)
	fmt.Printf("   New:     %s\n", newVersion)
	fmt.Println()

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
	var oldTimestamp time.Time
	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			oldFileCount++
			if oldTimestamp.IsZero() || info.ModTime().After(oldTimestamp) {
				oldTimestamp = info.ModTime()
			}
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

	// Write version marker file
	versionFile := filepath.Join(targetDir, ".version")
	if err := os.WriteFile(versionFile, []byte(newVersion+"\n"), 0644); err != nil {
		fmt.Printf("⚠️  Warning: Could not write version file: %v\n", err)
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
	fmt.Printf("✅ Successfully upgraded %s! (%s → %s)\n", config.Name, oldVersion, newVersion)
	fmt.Printf("   Files: %d → %d\n", oldFileCount, installed)
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

// getCurrentVersion returns the current lvt version
func getCurrentVersion() string {
	// Try to read VERSION file
	data, err := os.ReadFile("VERSION")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	return "unknown"
}

// detectInstalledAgent tries to auto-detect which agent is installed
func detectInstalledAgent() string {
	for name, config := range agentConfigs {
		if _, err := os.Stat(config.TargetDir); err == nil {
			return name
		}
	}
	return ""
}

// getInstalledVersion tries to determine the version of an installed agent
func getInstalledVersion(targetDir string) string {
	// Try to find a version marker file
	versionFile := filepath.Join(targetDir, ".version")
	if data, err := os.ReadFile(versionFile); err == nil {
		return strings.TrimSpace(string(data))
	}

	// Fallback: use modification time of directory
	if info, err := os.Stat(targetDir); err == nil {
		modTime := info.ModTime()
		return fmt.Sprintf("installed %s", modTime.Format("2006-01-02"))
	}

	return "unknown"
}

// handleExistingInstallation shows an interactive menu for handling existing installations
func handleExistingInstallation(targetDir, agentName string) (string, error) {
	fmt.Printf("⚠️  %s agent already installed\n", agentName)
	fmt.Println()
	fmt.Println("What would you like to do?")
	fmt.Println()
	fmt.Println("  [1] Upgrade (update agent, preserve custom settings)")
	fmt.Println("  [2] Overwrite (fresh install, removes custom settings)")
	fmt.Println("  [3] Cancel (keep existing installation)")
	fmt.Println()
	fmt.Print("Enter your choice (1-3): ")

	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > 3 {
		fmt.Println()
		fmt.Println("❌ Invalid choice")
		return "cancel", fmt.Errorf("invalid choice")
	}

	fmt.Println()

	switch choice {
	case 1:
		return "upgrade", nil
	case 2:
		fmt.Println("⚠️  Warning: This will remove all existing files including custom settings!")
		fmt.Print("Type 'yes' to confirm: ")
		var confirm string
		fmt.Scanf("%s", &confirm)
		fmt.Println()
		if confirm == "yes" {
			return "force", nil
		}
		return "cancel", fmt.Errorf("overwrite cancelled")
	case 3:
		fmt.Println("Installation cancelled - keeping existing files")
		return "cancel", fmt.Errorf("installation cancelled")
	default:
		return "cancel", fmt.Errorf("invalid choice")
	}
}

// selectAgentInteractive shows an interactive menu to select an agent
func selectAgentInteractive() (string, error) {
	fmt.Println("✨ LiveTemplate Agent Installation")
	fmt.Println()
	fmt.Println("Select your AI assistant:")
	fmt.Println()

	// Order for display
	order := []string{"claude", "copilot", "cursor", "aider", "generic"}

	for i, name := range order {
		config := agentConfigs[name]
		fmt.Printf("  [%d] %s\n", i+1, config.Name)
		fmt.Printf("      %s\n", config.Description)
		fmt.Println()
	}

	fmt.Print("Enter your choice (1-5): ")

	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > 5 {
		fmt.Println()
		fmt.Println("❌ Invalid choice")
		fmt.Println()
		fmt.Println("💡 Tip: You can also use: lvt install-agent --llm <type>")
		fmt.Println("   Available types: claude, copilot, cursor, aider, generic")
		fmt.Println()
		return "", fmt.Errorf("invalid choice")
	}

	selectedLLM := order[choice-1]
	fmt.Println()
	fmt.Printf("Installing %s agent...\n", agentConfigs[selectedLLM].Name)
	fmt.Println()

	return selectedLLM, nil
}

// listAgents lists all available agent types
func listAgents() {
	fmt.Println("Available AI agents for LiveTemplate:")
	fmt.Println()

	// Order: Claude, then alphabetically
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
	fmt.Println("  lvt install-agent                    # Interactive menu")
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
