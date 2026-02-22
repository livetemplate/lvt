package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livetemplate/lvt/internal/validation"
	"github.com/livetemplate/lvt/internal/validator"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/term"
)

// MCPServer starts the Model Context Protocol server for lvt
func MCPServer(args []string) error {
	// Parse flags
	hasFlags := false
	for _, arg := range args {
		switch arg {
		case "--help", "-h":
			printMCPHelp()
			return nil
		case "--setup":
			printMCPSetup()
			return nil
		case "--list-tools":
			printMCPTools()
			return nil
		case "--version", "-v":
			printMCPVersion()
			return nil
		default:
			if arg != "" && !strings.HasPrefix(arg, "-") {
				hasFlags = true
			}
		}
	}

	// Check if running in terminal (TTY)
	if isTerminal() {
		// If no flags provided and in terminal, show interactive setup
		if !hasFlags && len(args) == 0 {
			printMCPSetup()
			return nil
		}
		// Otherwise show warning
		printTTYWarning()
		return nil
	}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "lvt",
		Version: "0.1.0",
	}, nil)

	// Register tools
	registerNewTool(server)
	registerGenResourceTool(server)
	registerGenViewTool(server)
	registerGenAuthTool(server)
	registerGenSchemaTools(server)
	registerMigrationTools(server)
	registerSeedTool(server)
	registerResourceInspectTools(server)
	registerValidateTemplatesTool(server)
	registerEnvTools(server)
	registerKitsTools(server)

	// Start server with stdio transport
	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}

// ValidationIssueOutput is a single validation issue in MCP responses.
type ValidationIssueOutput struct {
	Level      string `json:"level"`
	File       string `json:"file,omitempty"`
	Line       int    `json:"line,omitempty"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

// ValidationOutput is the structured validation result for MCP responses.
type ValidationOutput struct {
	Valid      bool                    `json:"valid"`
	ErrorCount int                    `json:"error_count"`
	WarnCount  int                    `json:"warning_count"`
	Issues     []ValidationIssueOutput `json:"issues,omitempty"`
}

// runMCPValidation runs structural validation (post-gen) and converts to MCP output.
// Uses PostGen engine since generated code may not compile until sqlc generate is run.
func runMCPValidation(appPath string) *ValidationOutput {
	result := validation.ValidatePostGen(context.Background(), appPath)
	return validationResultToOutput(result)
}

// validationResultToOutput converts a validator.ValidationResult to MCP output.
func validationResultToOutput(result *validator.ValidationResult) *ValidationOutput {
	out := &ValidationOutput{
		Valid:      result.Valid,
		ErrorCount: result.ErrorCount(),
		WarnCount:  result.WarningCount(),
	}
	for _, issue := range result.Issues {
		out.Issues = append(out.Issues, ValidationIssueOutput{
			Level:      string(issue.Level),
			File:       issue.File,
			Line:       issue.Line,
			Message:    issue.Message,
			Suggestion: issue.Hint,
		})
	}
	return out
}

// NewAppInput defines the input schema for lvt new
type NewAppInput struct {
	Name   string `json:"name" jsonschema:"Application name"`
	Kit    string `json:"kit,omitempty" jsonschema:"Template kit (multi, single, or simple)"`
	CSS    string `json:"css,omitempty" jsonschema:"CSS framework (tailwind or none)"`
	Module string `json:"module,omitempty" jsonschema:"Go module name (defaults to app name)"`
}

// NewAppOutput defines the output schema for lvt new
type NewAppOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AppDir  string `json:"app_dir,omitempty"`
}

func registerNewTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_new",
		Description: "Create a new LiveTemplate application with specified configuration. Returns the application directory path on success.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input NewAppInput) (*mcp.CallToolResult, NewAppOutput, error) {
		// Validate required fields
		if input.Name == "" {
			return nil, NewAppOutput{
				Success: false,
				Message: "Application name is required",
			}, nil
		}

		// Set defaults
		kit := input.Kit
		if kit == "" {
			kit = "multi" // Default kit
		}

		module := input.Module
		if module == "" {
			module = input.Name
		}

		// Build command arguments
		args := []string{input.Name, "--module", module, "--kit", kit}

		// Execute New command
		if err := New(args); err != nil {
			return nil, NewAppOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to create app: %v", err),
			}, nil
		}

		// Get absolute path of created app
		appDir, err := filepath.Abs(input.Name)
		if err != nil {
			appDir = input.Name
		}

		return nil, NewAppOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully created %s app with %s kit", input.Name, kit),
			AppDir:  appDir,
		}, nil
	}

	mcp.AddTool(server, tool, handler)
}

// GenResourceInput defines the input schema for lvt gen resource
type GenResourceInput struct {
	Name   string            `json:"name" jsonschema:"Resource name (singular, e.g. 'post' or 'user')"`
	Fields map[string]string `json:"fields" jsonschema:"Field definitions as name:type pairs"`
}

// GenResourceOutput defines the output schema
type GenResourceOutput struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Files      []string          `json:"files,omitempty" jsonschema:"List of generated files"`
	Validation *ValidationOutput `json:"validation,omitempty"`
}

func registerGenResourceTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_gen_resource",
		Description: "Generate a complete CRUD resource with database schema, handlers, templates, and tests. Valid types: string, int, bool, float, time, text, textarea.",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GenResourceInput) (*mcp.CallToolResult, GenResourceOutput, error) {
		if input.Name == "" {
			return nil, GenResourceOutput{
				Success: false,
				Message: "Resource name is required",
			}, nil
		}

		if len(input.Fields) == 0 {
			return nil, GenResourceOutput{
				Success: false,
				Message: "At least one field is required",
			}, nil
		}

		// Build command arguments: lvt gen resource <name> <field:type>...
		// Pass --skip-validation to avoid double-validation (MCP runs it explicitly).
		args := []string{"resource", input.Name, "--skip-validation"}
		for field, typ := range input.Fields {
			args = append(args, fmt.Sprintf("%s:%s", field, typ))
		}

		// Execute Gen command
		if err := Gen(args); err != nil {
			return nil, GenResourceOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate resource: %v", err),
			}, nil
		}

		// Run validation explicitly for structured results
		wd, _ := os.Getwd()
		validationResult := runMCPValidation(wd)

		// List generated files
		files := []string{
			fmt.Sprintf("app/%s/%s.go", input.Name, input.Name),
			fmt.Sprintf("app/%s/%s.tmpl", input.Name, input.Name),
			fmt.Sprintf("app/%s/%s_test.go", input.Name, input.Name),
			fmt.Sprintf("app/%s/%s_ws_test.go", input.Name, input.Name),
		}

		return nil, GenResourceOutput{
			Success:    validationResult.Valid,
			Message:    fmt.Sprintf("Successfully generated %s resource with %d fields", input.Name, len(input.Fields)),
			Files:      files,
			Validation: validationResult,
		}, nil
	}

	mcp.AddTool(server, tool, handler)
}

// GenViewInput defines the input schema for lvt gen view
type GenViewInput struct {
	Name string `json:"name" jsonschema:"View name (e.g. 'dashboard' or 'counter')"`
}

// GenViewOutput defines the output schema
type GenViewOutput struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Files      []string          `json:"files,omitempty"`
	Validation *ValidationOutput `json:"validation,omitempty"`
}

func registerGenViewTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_gen_view",
		Description: "Generate a view-only handler without database (useful for dashboards, counters, etc.)",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GenViewInput) (*mcp.CallToolResult, GenViewOutput, error) {
		if input.Name == "" {
			return nil, GenViewOutput{
				Success: false,
				Message: "View name is required",
			}, nil
		}

		// Execute Gen command with --skip-validation (MCP runs it explicitly)
		args := []string{"view", input.Name, "--skip-validation"}
		if err := Gen(args); err != nil {
			return nil, GenViewOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate view: %v", err),
			}, nil
		}

		// Run validation explicitly for structured results
		wd, _ := os.Getwd()
		validationResult := runMCPValidation(wd)

		files := []string{
			fmt.Sprintf("app/%s/%s.go", input.Name, input.Name),
			fmt.Sprintf("app/%s/%s.tmpl", input.Name, input.Name),
		}

		return nil, GenViewOutput{
			Success:    validationResult.Valid,
			Message:    fmt.Sprintf("Successfully generated %s view", input.Name),
			Files:      files,
			Validation: validationResult,
		}, nil
	}

	mcp.AddTool(server, tool, handler)
}

// GenAuthInput defines the input schema for lvt gen auth
type GenAuthInput struct {
	StructName string `json:"struct_name,omitempty" jsonschema:"Go struct name (default: User)"`
	TableName  string `json:"table_name,omitempty" jsonschema:"Database table name (default: users)"`
}

// GenAuthOutput defines the output schema
type GenAuthOutput struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Validation *ValidationOutput `json:"validation,omitempty"`
}

func registerGenAuthTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_gen_auth",
		Description: "Generate a complete authentication system with sessions, password hashing, and auth handlers",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GenAuthInput) (*mcp.CallToolResult, GenAuthOutput, error) {
		// Pass --skip-validation to avoid double-validation (MCP runs it explicitly)
		args := []string{"auth", "--skip-validation"}

		if input.StructName != "" {
			args = append(args, input.StructName)
		}
		if input.TableName != "" {
			args = append(args, input.TableName)
		}

		if err := Gen(args); err != nil {
			return nil, GenAuthOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate auth: %v", err),
			}, nil
		}

		// Run validation explicitly for structured results
		wd, _ := os.Getwd()
		validationResult := runMCPValidation(wd)

		return nil, GenAuthOutput{
			Success:    validationResult.Valid,
			Message:    "Successfully generated authentication system",
			Validation: validationResult,
		}, nil
	}

	mcp.AddTool(server, tool, handler)
}

// MigrationInput defines input for migration commands
type MigrationInput struct {
	Name string `json:"name,omitempty" jsonschema:"Migration name (only for create command)"`
}

// MigrationOutput defines output for migration commands
type MigrationOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func registerMigrationTools(server *mcp.Server) {
	// Migration Up
	upTool := &mcp.Tool{
		Name:        "lvt_migration_up",
		Description: "Run all pending database migrations",
	}
	mcp.AddTool(server, upTool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, MigrationOutput, error) {
		// Redirect stdout to capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Migration([]string{"up"})

		// Restore stdout
		w.Close()
		os.Stdout = oldStdout

		// Read captured output
		buf := make([]byte, 4096)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, MigrationOutput{
				Success: false,
				Message: fmt.Sprintf("Migration failed: %v", err),
				Details: output,
			}, nil
		}

		return nil, MigrationOutput{
			Success: true,
			Message: "Migrations applied successfully",
			Details: output,
		}, nil
	})

	// Migration Down
	downTool := &mcp.Tool{
		Name:        "lvt_migration_down",
		Description: "Rollback the last database migration",
	}
	mcp.AddTool(server, downTool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, MigrationOutput, error) {
		if err := Migration([]string{"down"}); err != nil {
			return nil, MigrationOutput{
				Success: false,
				Message: fmt.Sprintf("Rollback failed: %v", err),
			}, nil
		}

		return nil, MigrationOutput{
			Success: true,
			Message: "Migration rolled back successfully",
		}, nil
	})

	// Migration Status
	statusTool := &mcp.Tool{
		Name:        "lvt_migration_status",
		Description: "Show status of all database migrations",
	}
	mcp.AddTool(server, statusTool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, MigrationOutput, error) {
		// Capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Migration([]string{"status"})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 4096)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, MigrationOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to get status: %v", err),
			}, nil
		}

		return nil, MigrationOutput{
			Success: true,
			Message: "Migration status retrieved",
			Details: output,
		}, nil
	})

	// Migration Create
	createTool := &mcp.Tool{
		Name:        "lvt_migration_create",
		Description: "Create a new migration file with the specified name",
	}
	mcp.AddTool(server, createTool, func(ctx context.Context, req *mcp.CallToolRequest, input MigrationInput) (*mcp.CallToolResult, MigrationOutput, error) {
		if input.Name == "" {
			return nil, MigrationOutput{
				Success: false,
				Message: "Migration name is required",
			}, nil
		}

		if err := Migration([]string{"create", input.Name}); err != nil {
			return nil, MigrationOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to create migration: %v", err),
			}, nil
		}

		return nil, MigrationOutput{
			Success: true,
			Message: fmt.Sprintf("Migration '%s' created successfully", input.Name),
		}, nil
	})
}

// GenSchemaInput defines input for gen schema
type GenSchemaInput struct {
	Table  string            `json:"table" jsonschema:"Database table name"`
	Fields map[string]string `json:"fields" jsonschema:"Field definitions as name:type pairs"`
}

// GenSchemaOutput defines output for gen schema
type GenSchemaOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func registerGenSchemaTools(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_gen_schema",
		Description: "Generate database schema only (no handlers or templates)",
	}

	mcp.AddTool(server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input GenSchemaInput) (*mcp.CallToolResult, GenSchemaOutput, error) {
		if input.Table == "" {
			return nil, GenSchemaOutput{
				Success: false,
				Message: "Table name is required",
			}, nil
		}

		args := []string{"schema", input.Table}
		for field, typ := range input.Fields {
			args = append(args, fmt.Sprintf("%s:%s", field, typ))
		}

		if err := Gen(args); err != nil {
			return nil, GenSchemaOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate schema: %v", err),
			}, nil
		}

		return nil, GenSchemaOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully generated schema for %s", input.Table),
		}, nil
	})
}

// SeedInput defines input for seed command
type SeedInput struct {
	Resource string `json:"resource" jsonschema:"Resource name to seed"`
	Count    int    `json:"count,omitempty" jsonschema:"Number of records to generate (default: 10)"`
	Cleanup  bool   `json:"cleanup,omitempty" jsonschema:"Clean up existing test data before seeding"`
}

// SeedOutput defines output for seed command
type SeedOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func registerSeedTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_seed",
		Description: "Generate test data for a resource",
	}

	mcp.AddTool(server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input SeedInput) (*mcp.CallToolResult, SeedOutput, error) {
		if input.Resource == "" {
			return nil, SeedOutput{
				Success: false,
				Message: "Resource name is required",
			}, nil
		}

		args := []string{input.Resource}
		if input.Count > 0 {
			args = append(args, "--count", fmt.Sprintf("%d", input.Count))
		}
		if input.Cleanup {
			args = append(args, "--cleanup")
		}

		if err := Seed(args); err != nil {
			return nil, SeedOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to seed data: %v", err),
			}, nil
		}

		return nil, SeedOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully seeded %s", input.Resource),
		}, nil
	})
}

// ResourceInspectInput defines input for resource inspect
type ResourceInspectInput struct {
	Command  string `json:"command" jsonschema:"Command: 'list' or 'describe'"`
	Resource string `json:"resource,omitempty" jsonschema:"Resource name (for describe command)"`
}

// ResourceInspectOutput defines output
type ResourceInspectOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func registerResourceInspectTools(server *mcp.Server) {
	// List resources
	listTool := &mcp.Tool{
		Name:        "lvt_resource_list",
		Description: "List all available resources in the project",
	}
	mcp.AddTool(server, listTool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, ResourceInspectOutput, error) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Resource([]string{"list"})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 8192)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, ResourceInspectOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to list resources: %v", err),
			}, nil
		}

		return nil, ResourceInspectOutput{
			Success: true,
			Message: "Resources listed successfully",
			Details: output,
		}, nil
	})

	// Describe resource
	describeTool := &mcp.Tool{
		Name:        "lvt_resource_describe",
		Description: "Show detailed schema for a specific resource",
	}
	mcp.AddTool(server, describeTool, func(ctx context.Context, req *mcp.CallToolRequest, input ResourceInspectInput) (*mcp.CallToolResult, ResourceInspectOutput, error) {
		if input.Resource == "" {
			return nil, ResourceInspectOutput{
				Success: false,
				Message: "Resource name is required",
			}, nil
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Resource([]string{"describe", input.Resource})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 8192)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, ResourceInspectOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to describe resource: %v", err),
			}, nil
		}

		return nil, ResourceInspectOutput{
			Success: true,
			Message: fmt.Sprintf("Resource %s described successfully", input.Resource),
			Details: output,
		}, nil
	})
}

// ValidateTemplatesInput defines input
type ValidateTemplatesInput struct {
	TemplateFile string `json:"template_file" jsonschema:"Path to template file to validate"`
}

// ValidateTemplatesOutput defines output
type ValidateTemplatesOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func registerValidateTemplatesTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_validate_template",
		Description: "Validate and analyze a template file",
	}

	mcp.AddTool(server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input ValidateTemplatesInput) (*mcp.CallToolResult, ValidateTemplatesOutput, error) {
		if input.TemplateFile == "" {
			return nil, ValidateTemplatesOutput{
				Success: false,
				Message: "Template file path is required",
			}, nil
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Parse([]string{input.TemplateFile})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 8192)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, ValidateTemplatesOutput{
				Success: false,
				Message: fmt.Sprintf("Template validation failed: %v", err),
				Details: output,
			}, nil
		}

		return nil, ValidateTemplatesOutput{
			Success: true,
			Message: "Template is valid",
			Details: output,
		}, nil
	})
}

// EnvInput defines input for env commands
type EnvInput struct {
	Command string `json:"command" jsonschema:"Command: 'generate' to create .env.example"`
}

// EnvOutput defines output
type EnvOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func registerEnvTools(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_env_generate",
		Description: "Generate .env.example with detected configuration",
	}

	mcp.AddTool(server, tool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, EnvOutput, error) {
		if err := Env([]string{"generate"}); err != nil {
			return nil, EnvOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate env file: %v", err),
			}, nil
		}

		return nil, EnvOutput{
			Success: true,
			Message: ".env.example generated successfully",
		}, nil
	})
}

// KitsInput defines input for kits commands
type KitsInput struct {
	Command string `json:"command" jsonschema:"Command: 'list', 'info', 'validate', or 'create'"`
	Name    string `json:"name,omitempty" jsonschema:"Kit name (for info, validate, or create)"`
	Path    string `json:"path,omitempty" jsonschema:"Path to kit (for validate)"`
}

// KitsOutput defines output
type KitsOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func registerKitsTools(server *mcp.Server) {
	// List kits
	listTool := &mcp.Tool{
		Name:        "lvt_kits_list",
		Description: "List all available CSS framework kits",
	}
	mcp.AddTool(server, listTool, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, KitsOutput, error) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Kits([]string{"list"})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 8192)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, KitsOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to list kits: %v", err),
			}, nil
		}

		return nil, KitsOutput{
			Success: true,
			Message: "Kits listed successfully",
			Details: output,
		}, nil
	})

	// Info about kit
	infoTool := &mcp.Tool{
		Name:        "lvt_kits_info",
		Description: "Show detailed information about a specific kit",
	}
	mcp.AddTool(server, infoTool, func(ctx context.Context, req *mcp.CallToolRequest, input KitsInput) (*mcp.CallToolResult, KitsOutput, error) {
		if input.Name == "" {
			return nil, KitsOutput{
				Success: false,
				Message: "Kit name is required",
			}, nil
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := Kits([]string{"info", input.Name})

		w.Close()
		os.Stdout = oldStdout

		buf := make([]byte, 8192)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if err != nil {
			return nil, KitsOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to get kit info: %v", err),
			}, nil
		}

		return nil, KitsOutput{
			Success: true,
			Message: fmt.Sprintf("Kit %s info retrieved", input.Name),
			Details: output,
		}, nil
	})
}

// init registers the MCP server logger
func init() {
	// Suppress SDK logs in normal operation (they go to stderr which interferes with MCP protocol)
	log.SetOutput(os.Stderr)
}

// isTerminal checks if stdin/stdout are connected to a terminal
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd())) && term.IsTerminal(int(os.Stdout.Fd()))
}

// getConfigPath returns the platform-specific config file path for Claude Desktop
func getConfigPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "~/Library/Application Support/Claude/claude_desktop_config.json"
	case "windows":
		return "%APPDATA%\\Claude\\claude_desktop_config.json"
	default: // linux and others
		return "~/.config/Claude/claude_desktop_config.json"
	}
}

// printTTYWarning shows a warning when MCP server is run directly in terminal
func printTTYWarning() {
	fmt.Println("⚠️  MCP Server Warning")
	fmt.Println()
	fmt.Println("The MCP server runs via AI client configuration, not directly in terminal.")
	fmt.Println()
	fmt.Println("Quick Setup:")
	fmt.Printf("  1. Edit: %s\n", getConfigPath())
	fmt.Println("  2. Add lvt to mcpServers (see --help for JSON)")
	fmt.Println("  3. Restart AI client")
	fmt.Println()
	fmt.Println("For guided setup: lvt mcp-server --setup")
	fmt.Println("For full docs:    docs/AGENT_SETUP.md")
	fmt.Println()
}

// printMCPHelp shows comprehensive help for the MCP server
func printMCPHelp() {
	fmt.Println("LiveTemplate MCP Server")
	fmt.Println()
	fmt.Println("Provides 16 tools for AI assistants to build LiveTemplate applications.")
	fmt.Println("The server runs as a JSON-RPC service over stdio and must be configured")
	fmt.Println("in your AI client (not run directly).")
	fmt.Println()
	fmt.Println("SETUP INSTRUCTIONS:")
	fmt.Println()
	fmt.Println("1. Configure your AI client:")
	fmt.Println()

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("   Claude Desktop (macOS):")
		fmt.Println("   Edit: ~/Library/Application Support/Claude/claude_desktop_config.json")
	case "windows":
		fmt.Println("   Claude Desktop (Windows):")
		fmt.Println("   Edit: %APPDATA%\\Claude\\claude_desktop_config.json")
	default:
		fmt.Println("   Claude Desktop (Linux):")
		fmt.Println("   Edit: ~/.config/Claude/claude_desktop_config.json")
	}

	fmt.Println()
	fmt.Println("   Add this JSON:")
	fmt.Println("   {")
	fmt.Println("     \"mcpServers\": {")
	fmt.Println("       \"lvt\": {")
	fmt.Println("         \"command\": \"lvt\",")
	fmt.Println("         \"args\": [\"mcp-server\"]")
	fmt.Println("       }")
	fmt.Println("     }")
	fmt.Println("   }")
	fmt.Println()
	fmt.Println("2. Restart your AI client")
	fmt.Println()
	fmt.Println("3. Test by asking: \"List available LiveTemplate tools\"")
	fmt.Println()
	fmt.Println("AVAILABLE TOOLS (16):")
	fmt.Println()
	fmt.Println("  Generation (5):  lvt_new, lvt_gen_resource, lvt_gen_view, lvt_gen_auth,")
	fmt.Println("                   lvt_gen_schema")
	fmt.Println("  Database (4):    lvt_migration_up, lvt_migration_down, lvt_migration_status,")
	fmt.Println("                   lvt_migration_create")
	fmt.Println("  Development (7): lvt_seed, lvt_resource_list, lvt_resource_describe,")
	fmt.Println("                   lvt_validate_template, lvt_env_generate, lvt_kits_list,")
	fmt.Println("                   lvt_kits_info")
	fmt.Println()
	fmt.Println("DOCUMENTATION:")
	fmt.Println()
	fmt.Println("  Full tool docs:  docs/MCP_TOOLS.md")
	fmt.Println("  Setup guide:     docs/AGENT_SETUP.md")
	fmt.Println("  Workflows:       docs/WORKFLOWS.md")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println()
	fmt.Println("  --help, -h       Show this help message")
	fmt.Println("  --setup          Interactive setup wizard")
	fmt.Println("  --list-tools     List all tools with descriptions")
	fmt.Println("  --version, -v    Show MCP protocol version")
	fmt.Println()
	fmt.Println("NOTES:")
	fmt.Println()
	fmt.Println("  - Don't run this command directly - it's started by your AI client")
	fmt.Println("  - For project-specific setup: lvt install-agent --llm <type>")
	fmt.Println("  - For troubleshooting: docs/AGENT_SETUP.md#troubleshooting")
	fmt.Println()
}

// printMCPSetup shows interactive setup wizard
func printMCPSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  LiveTemplate MCP Server - Interactive Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	// Interactive LLM selection
	fmt.Println("Which AI assistant are you using?")
	fmt.Println()
	fmt.Println("  1. Claude Desktop / Claude Code")
	fmt.Println("  2. VS Code with Copilot Chat")
	fmt.Println("  3. Cursor AI")
	fmt.Println("  4. Aider CLI")
	fmt.Println("  5. Other MCP-compatible client")
	fmt.Println()
	fmt.Print("Enter your choice (1-5): ")

	var choice string
	fmt.Scanln(&choice)
	fmt.Println()

	switch choice {
	case "1":
		printClaudeSetup()
	case "2":
		printCopilotSetup()
	case "3":
		printCursorSetup()
	case "4":
		printAiderSetup()
	case "5":
		printGenericSetup()
	default:
		fmt.Println("Invalid choice. Showing generic setup instructions...")
		fmt.Println()
		printGenericSetup()
	}
}

// printClaudeSetup shows Claude-specific setup instructions
func printClaudeSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Claude Desktop / Claude Code Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	configPath := getConfigPath()

	fmt.Println("STEP 1: Locate Configuration File")
	fmt.Println()
	fmt.Printf("Your platform: %s\n", runtime.GOOS)
	fmt.Printf("Config file:   %s\n", configPath)
	fmt.Println()

	fmt.Println("STEP 2: Add MCP Server Configuration")
	fmt.Println()
	fmt.Println("Copy and paste this JSON into your config file:")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│ {                                                           │")
	fmt.Println("│   \"mcpServers\": {                                          │")
	fmt.Println("│     \"lvt\": {                                               │")
	fmt.Println("│       \"command\": \"lvt\",                                    │")
	fmt.Println("│       \"args\": [\"mcp-server\"]                               │")
	fmt.Println("│     }                                                       │")
	fmt.Println("│   }                                                         │")
	fmt.Println("│ }                                                           │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("NOTE: If your config file already has mcpServers, just add the")
	fmt.Println("      \"lvt\" entry inside the existing mcpServers object.")
	fmt.Println()

	fmt.Println("STEP 3: Restart Claude")
	fmt.Println()
	fmt.Println("• Claude Desktop: Quit and relaunch the application")
	fmt.Println("• Claude Code: Restart the CLI")
	fmt.Println()

	fmt.Println("STEP 4: Verify Installation")
	fmt.Println()
	fmt.Println("Ask Claude: \"List available LiveTemplate tools\"")
	fmt.Println("You should see 16 tools listed.")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// printCopilotSetup shows GitHub Copilot setup instructions
func printCopilotSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  GitHub Copilot Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("RECOMMENDED APPROACH: Use Agent Installation")
	fmt.Println()
	fmt.Println("GitHub Copilot works best with project-specific instructions")
	fmt.Println("rather than global MCP server configuration.")
	fmt.Println()

	fmt.Println("SETUP STEPS:")
	fmt.Println()
	fmt.Println("1. Install the Copilot agent:")
	fmt.Println("   $ lvt install-agent --llm copilot")
	fmt.Println()
	fmt.Println("2. Open your project in VS Code with Copilot enabled")
	fmt.Println()
	fmt.Println("3. Copilot will automatically read the instructions and")
	fmt.Println("   understand LiveTemplate commands")
	fmt.Println()

	fmt.Println("USAGE:")
	fmt.Println()
	fmt.Println("• Use @workspace to ask questions about LiveTemplate")
	fmt.Println("• Ask: \"How do I add a posts resource?\"")
	fmt.Println("• Copilot will suggest using lvt commands")
	fmt.Println()

	fmt.Println("NOTE: If your VS Code environment supports MCP servers,")
	fmt.Println("      you can configure it in VS Code settings. Check the")
	fmt.Println("      VS Code MCP documentation for details.")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// printCursorSetup shows Cursor AI setup instructions
func printCursorSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Cursor AI Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("RECOMMENDED APPROACH: Use Agent Installation")
	fmt.Println()
	fmt.Println("Cursor works best with project-specific rules rather than")
	fmt.Println("global MCP server configuration.")
	fmt.Println()

	fmt.Println("SETUP STEPS:")
	fmt.Println()
	fmt.Println("1. Install the Cursor agent:")
	fmt.Println("   $ lvt install-agent --llm cursor")
	fmt.Println()
	fmt.Println("2. Open your project in Cursor")
	fmt.Println()
	fmt.Println("3. Rules will apply automatically to *.go files")
	fmt.Println()

	fmt.Println("USAGE:")
	fmt.Println()
	fmt.Println("• Use Composer mode for best results")
	fmt.Println("• Use Agent mode for autonomous workflows")
	fmt.Println("• Ask: \"Add a blog with authentication\"")
	fmt.Println("• Cursor follows LiveTemplate patterns automatically")
	fmt.Println()

	fmt.Println("NOTE: If Cursor adds MCP server support in the future,")
	fmt.Println("      you can configure it in Cursor settings.")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// printAiderSetup shows Aider CLI setup instructions
func printAiderSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Aider CLI Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("OPTION 1: Agent Installation (Recommended)")
	fmt.Println()
	fmt.Println("1. Install the Aider agent:")
	fmt.Println("   $ lvt install-agent --llm aider")
	fmt.Println()
	fmt.Println("2. Start Aider:")
	fmt.Println("   $ aider")
	fmt.Println()
	fmt.Println("   Configuration loads automatically from .aider/.aider.conf.yml")
	fmt.Println()

	fmt.Println("OPTION 2: MCP Server (If Supported)")
	fmt.Println()
	fmt.Println("If your version of Aider supports MCP servers:")
	fmt.Println()
	fmt.Println("1. Add to .aider/.aider.conf.yml:")
	fmt.Println()
	fmt.Println("   mcp_servers:")
	fmt.Println("     - name: lvt")
	fmt.Println("       command: lvt")
	fmt.Println("       args: [mcp-server]")
	fmt.Println()
	fmt.Println("2. Start Aider:")
	fmt.Println("   $ aider")
	fmt.Println()

	fmt.Println("VERIFICATION:")
	fmt.Println()
	fmt.Println("Ask Aider: \"Add a posts resource with title and content\"")
	fmt.Println("Aider should use: lvt gen resource posts title:string content:text")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// printGenericSetup shows generic MCP setup instructions
func printGenericSetup() {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  Generic MCP Server Setup")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("For other MCP-compatible AI clients:")
	fmt.Println()

	fmt.Println("STEP 1: Check MCP Support")
	fmt.Println()
	fmt.Println("Verify your AI client supports the Model Context Protocol (MCP).")
	fmt.Println("See: https://modelcontextprotocol.io")
	fmt.Println()

	fmt.Println("STEP 2: Add MCP Server Configuration")
	fmt.Println()
	fmt.Println("Add this configuration to your client's MCP settings:")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("│ {                                                           │")
	fmt.Println("│   \"mcpServers\": {                                          │")
	fmt.Println("│     \"lvt\": {                                               │")
	fmt.Println("│       \"command\": \"lvt\",                                    │")
	fmt.Println("│       \"args\": [\"mcp-server\"]                               │")
	fmt.Println("│     }                                                       │")
	fmt.Println("│   }                                                         │")
	fmt.Println("│ }                                                           │")
	fmt.Println("└─────────────────────────────────────────────────────────────┘")
	fmt.Println()

	fmt.Println("STEP 3: Restart Your AI Client")
	fmt.Println()
	fmt.Println("Restart your AI client to load the MCP server configuration.")
	fmt.Println()

	fmt.Println("STEP 4: Verify Installation")
	fmt.Println()
	fmt.Println("Ask your AI client to list available tools.")
	fmt.Println("You should see 16 LiveTemplate tools.")
	fmt.Println()

	fmt.Println("TROUBLESHOOTING:")
	fmt.Println()
	fmt.Println("• Tools not showing? Check config file syntax (must be valid JSON)")
	fmt.Println("• Server not starting? Ensure lvt is in your PATH")
	fmt.Println("• Need help? See docs/AGENT_SETUP.md#troubleshooting")
	fmt.Println()

	fmt.Println("ALTERNATIVE: Agent Installation")
	fmt.Println()
	fmt.Println("If MCP doesn't work, try installing agent documentation:")
	fmt.Println("  $ lvt install-agent --llm generic")
	fmt.Println()
	fmt.Println("This provides complete documentation and examples for")
	fmt.Println("integrating with any LLM that can execute shell commands.")
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

// printMCPTools lists all available MCP tools
func printMCPTools() {
	fmt.Println("LiveTemplate MCP Tools (16 total)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("GENERATION TOOLS (5)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("lvt_new")
	fmt.Println("  Create a new LiveTemplate application")
	fmt.Println("  Input: name (required), kit, css, module")
	fmt.Println()
	fmt.Println("lvt_gen_resource")
	fmt.Println("  Generate a CRUD resource with database integration")
	fmt.Println("  Input: name (required), fields (object)")
	fmt.Println()
	fmt.Println("lvt_gen_view")
	fmt.Println("  Generate a view-only handler (no database)")
	fmt.Println("  Input: name (required)")
	fmt.Println()
	fmt.Println("lvt_gen_auth")
	fmt.Println("  Generate authentication system")
	fmt.Println("  Input: optional configuration flags")
	fmt.Println()
	fmt.Println("lvt_gen_schema")
	fmt.Println("  Generate database schema without UI")
	fmt.Println("  Input: table (required), fields (object)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("DATABASE TOOLS (4)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("lvt_migration_up")
	fmt.Println("  Apply all pending database migrations")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_migration_down")
	fmt.Println("  Rollback the last migration")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_migration_status")
	fmt.Println("  Show migration status (pending/applied)")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_migration_create")
	fmt.Println("  Create a new migration file")
	fmt.Println("  Input: name (required)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("DEVELOPMENT TOOLS (7)")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("lvt_seed")
	fmt.Println("  Generate test data for a resource")
	fmt.Println("  Input: resource (required), count, cleanup")
	fmt.Println()
	fmt.Println("lvt_resource_list")
	fmt.Println("  List all available resources")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_resource_describe")
	fmt.Println("  Show detailed schema for a resource")
	fmt.Println("  Input: resource (required)")
	fmt.Println()
	fmt.Println("lvt_validate_template")
	fmt.Println("  Validate template file syntax")
	fmt.Println("  Input: template_file (required)")
	fmt.Println()
	fmt.Println("lvt_env_generate")
	fmt.Println("  Generate .env.example file")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_kits_list")
	fmt.Println("  List available CSS framework kits")
	fmt.Println("  Input: none")
	fmt.Println()
	fmt.Println("lvt_kits_info")
	fmt.Println("  Show detailed information about a kit")
	fmt.Println("  Input: name (required)")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("DOCUMENTATION")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("For complete tool documentation with input/output schemas,")
	fmt.Println("examples, and best practices, see:")
	fmt.Println()
	fmt.Println("  docs/MCP_TOOLS.md")
	fmt.Println()
}

// printMCPVersion shows version information
func printMCPVersion() {
	fmt.Println("LiveTemplate MCP Server")
	fmt.Println()
	fmt.Println("Server Version:   0.1.0")
	fmt.Println("MCP Protocol:     v1.0")
	fmt.Println("Go SDK:           github.com/modelcontextprotocol/go-sdk v1.1.0")
	fmt.Println()
	fmt.Println("Compatibility:")
	fmt.Println("  • Claude Desktop (all versions)")
	fmt.Println("  • Claude Code (all versions)")
	fmt.Println("  • Any MCP-compatible client")
	fmt.Println()
	fmt.Println("Tools Available:  16")
	fmt.Println()
}
