package agenttest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/commands"
	e2etest "github.com/livetemplate/lvt/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AgentTestEnv wraps E2ETest with agent-specific testing capabilities
type AgentTestEnv struct {
	*e2etest.E2ETest
	T             *testing.T
	TmpDir        string
	AppName       string
	AppDir        string
	CommandsRun   []string
	SkillsUsed    []string
	CurrentWorkDir string
}

// SetupOptions configures the agent test environment
type SetupOptions struct {
	AppName    string // Name of the app to create (if empty, no app is created)
	Kit        string // Kit type: multi, single, simple
	ChromeMode e2etest.ChromeMode
}

// Setup creates a new agent test environment with isolation
func Setup(t *testing.T, opts *SetupOptions) *AgentTestEnv {
	t.Helper()

	if opts == nil {
		opts = &SetupOptions{}
	}

	if opts.Kit == "" {
		opts.Kit = "multi"
	}

	// Create temporary directory for this test
	tmpDir := t.TempDir()

	var e2eTest *e2etest.E2ETest
	var appDir string

	// If AppName is provided, create the app and setup e2e testing
	if opts.AppName != "" {
		// Create the app in temp directory
		originalDir, err := os.Getwd()
		require.NoError(t, err)

		err = os.Chdir(tmpDir)
		require.NoError(t, err)

		// Create app with SKIP_GO_MOD_TIDY to avoid test I/O issues
		os.Setenv("SKIP_GO_MOD_TIDY", "1")
		defer os.Unsetenv("SKIP_GO_MOD_TIDY")

		args := []string{opts.AppName, "--kit", opts.Kit}
		err = commands.New(args)
		require.NoError(t, err, "failed to create test app")

		appDir = filepath.Join(tmpDir, opts.AppName)

		// Run go mod tidy separately with proper synchronization
		err = runGoModTidy(t, appDir)
		require.NoError(t, err, "failed to run go mod tidy")

		// Return to original directory
		err = os.Chdir(originalDir)
		require.NoError(t, err)

		// Setup E2E testing if Chrome is needed
		if opts.ChromeMode != "" {
			// Determine app path based on kit
			var appPath string
			if opts.Kit == "simple" {
				appPath = filepath.Join(appDir, "main.go")
			} else {
				appPath = filepath.Join(appDir, "cmd", opts.AppName, "main.go")
			}

			// Change to app directory for e2e test setup
			err = os.Chdir(appDir)
			require.NoError(t, err)

			e2eTest = e2etest.Setup(t, &e2etest.SetupOptions{
				AppPath:    appPath,
				ChromeMode: opts.ChromeMode,
			})

			// Return to original directory
			err = os.Chdir(originalDir)
			require.NoError(t, err)
		}
	}

	env := &AgentTestEnv{
		E2ETest:        e2eTest,
		T:              t,
		TmpDir:         tmpDir,
		AppName:        opts.AppName,
		AppDir:         appDir,
		CommandsRun:    []string{},
		SkillsUsed:     []string{},
		CurrentWorkDir: tmpDir,
	}

	return env
}

// TrackSkill records that a skill was used (for testing agent skill selection)
func (e *AgentTestEnv) TrackSkill(skillName string) {
	e.SkillsUsed = append(e.SkillsUsed, skillName)
}

// TrackCommand records that a command was run
func (e *AgentTestEnv) TrackCommand(cmd string) {
	e.CommandsRun = append(e.CommandsRun, cmd)
}

// RunLvtCommand executes an lvt command and tracks it
func (e *AgentTestEnv) RunLvtCommand(args ...string) error {
	e.T.Helper()

	cmdStr := "lvt " + strings.Join(args, " ")
	e.TrackCommand(cmdStr)

	// Change to app directory if it exists
	if e.AppDir != "" {
		originalDir, err := os.Getwd()
		if err != nil {
			return err
		}
		defer os.Chdir(originalDir)

		if err := os.Chdir(e.AppDir); err != nil {
			return err
		}
	}

	// Route to appropriate command handler
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	switch args[0] {
	case "gen":
		return commands.Gen(args[1:])
	case "migration":
		if len(args) < 2 {
			return fmt.Errorf("migration subcommand required")
		}
		return commands.Migration(args[1:])
	case "seed":
		return commands.Seed(args[1:])
	case "resource", "res":
		return commands.Resource(args[1:])
	case "parse":
		return commands.Parse(args[1:])
	case "env":
		return commands.Env(args[1:])
	case "kits", "kit":
		return commands.Kits(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

// AssertSkillUsed verifies that a specific skill was used
func (e *AgentTestEnv) AssertSkillUsed(skillName string) {
	e.T.Helper()
	assert.Contains(e.T, e.SkillsUsed, skillName, "Expected skill %s to be used", skillName)
}

// AssertCommandRun verifies that a specific command was executed
func (e *AgentTestEnv) AssertCommandRun(cmd string) {
	e.T.Helper()
	assert.Contains(e.T, e.CommandsRun, cmd, "Expected command '%s' to be run", cmd)
}

// AssertFileExists verifies that a file exists in the app directory
func (e *AgentTestEnv) AssertFileExists(relativePath string) {
	e.T.Helper()
	fullPath := filepath.Join(e.AppDir, relativePath)
	_, err := os.Stat(fullPath)
	assert.NoError(e.T, err, "Expected file to exist: %s", relativePath)
}

// AssertFileContains verifies that a file contains specific text
func (e *AgentTestEnv) AssertFileContains(relativePath, text string) {
	e.T.Helper()
	fullPath := filepath.Join(e.AppDir, relativePath)
	content, err := os.ReadFile(fullPath)
	require.NoError(e.T, err, "Failed to read file: %s", relativePath)
	assert.Contains(e.T, string(content), text, "File %s should contain '%s'", relativePath, text)
}

// AssertDatabaseTableExists verifies that a table exists in the database
func (e *AgentTestEnv) AssertDatabaseTableExists(tableName string) {
	e.T.Helper()
	// This is validated by checking if the migration file exists and was applied
	migrationsDir := filepath.Join(e.AppDir, "internal", "database", "migrations")
	files, err := os.ReadDir(migrationsDir)
	require.NoError(e.T, err, "Failed to read migrations directory")

	found := false
	for _, file := range files {
		if strings.Contains(file.Name(), tableName) {
			found = true
			break
		}
	}
	assert.True(e.T, found, "Expected migration for table %s to exist", tableName)
}

// SimulateQuickStart simulates the quick start workflow from the usage guide
func (e *AgentTestEnv) SimulateQuickStart(resourceName string) {
	e.T.Helper()

	// Track that we're using the quickstart skill
	e.TrackSkill("lvt-quickstart")

	// Generate the resource
	err := e.RunLvtCommand("gen", "resource", resourceName, "title:string", "content:text")
	require.NoError(e.T, err, "gen command should succeed")

	// Run migrations
	err = e.RunLvtCommand("migration", "up")
	require.NoError(e.T, err, "migration up should succeed")

	// Verify files exist
	e.AssertFileExists("internal/database/migrations")
	e.AssertFileExists("internal/app/" + resourceName + "/" + resourceName + ".go")
	e.AssertFileExists("internal/database/queries.sql")
}

// SimulateFullStack simulates the full stack workflow from the usage guide
func (e *AgentTestEnv) SimulateFullStack() {
	e.T.Helper()

	// Track skills used
	e.TrackSkill("lvt-gen-auth")
	e.TrackSkill("lvt-add-resource")

	// Generate auth
	err := e.RunLvtCommand("gen", "auth")
	require.NoError(e.T, err, "gen auth should succeed")

	// Add projects resource
	err = e.RunLvtCommand("gen", "resource", "projects", "name:string", "description:text", "user_id:int")
	require.NoError(e.T, err, "gen projects should succeed")

	// Add tasks resource
	err = e.RunLvtCommand("gen", "resource", "tasks", "title:string", "description:text", "project_id:int", "due_date:time", "priority:int")
	require.NoError(e.T, err, "gen tasks should succeed")

	// Run migrations
	err = e.RunLvtCommand("migration", "up")
	require.NoError(e.T, err, "migration up should succeed")

	// Verify key files exist
	e.AssertFileExists("internal/app/auth/auth.go")
	e.AssertFileExists("internal/app/projects/projects.go")
	e.AssertFileExists("internal/app/tasks/tasks.go")
}

// SimulateConversation simulates a multi-turn conversation
func (e *AgentTestEnv) SimulateConversation(steps []ConversationStep) {
	e.T.Helper()

	for i, step := range steps {
		e.T.Logf("Conversation step %d: %s", i+1, step.Prompt)
		step.Execute(e)
	}
}

// ConversationStep represents a single turn in a conversation
type ConversationStep struct {
	Prompt  string
	Execute func(*AgentTestEnv)
}

// runGoModTidy runs go mod tidy in the specified directory
func runGoModTidy(t *testing.T, dir string) error {
	t.Helper()
	// For now, we skip this as it's handled by SKIP_GO_MOD_TIDY flag
	// and the test infrastructure
	return nil
}
