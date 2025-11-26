package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPServer starts the Model Context Protocol server for lvt
func MCPServer(args []string) error {
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

// NewAppInput defines the input schema for lvt new
type NewAppInput struct {
	Name   string `json:"name" jsonschema:"Application name"`
	Kit    string `json:"kit,omitempty" jsonschema:"Template kit (multi, single, or simple)"`
	CSS    string `json:"css,omitempty" jsonschema:"CSS framework (tailwind, bulma, pico, or none)"`
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
	Success bool   `json:"success"`
	Message string `json:"message"`
	Files   []string `json:"files,omitempty" jsonschema:"List of generated files"`
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
		args := []string{"resource", input.Name}
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

		// List generated files
		files := []string{
			fmt.Sprintf("internal/app/%s/%s.go", input.Name, input.Name),
			fmt.Sprintf("internal/app/%s/%s.tmpl", input.Name, input.Name),
			fmt.Sprintf("internal/app/%s/%s_test.go", input.Name, input.Name),
			fmt.Sprintf("internal/app/%s/%s_ws_test.go", input.Name, input.Name),
		}

		return nil, GenResourceOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully generated %s resource with %d fields", input.Name, len(input.Fields)),
			Files:   files,
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
	Success bool   `json:"success"`
	Message string `json:"message"`
	Files   []string `json:"files,omitempty"`
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

		// Execute Gen command
		args := []string{"view", input.Name}
		if err := Gen(args); err != nil {
			return nil, GenViewOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to generate view: %v", err),
			}, nil
		}

		files := []string{
			fmt.Sprintf("internal/app/%s/%s.go", input.Name, input.Name),
			fmt.Sprintf("internal/app/%s/%s.tmpl", input.Name, input.Name),
		}

		return nil, GenViewOutput{
			Success: true,
			Message: fmt.Sprintf("Successfully generated %s view", input.Name),
			Files:   files,
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
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func registerGenAuthTool(server *mcp.Server) {
	tool := &mcp.Tool{
		Name:        "lvt_gen_auth",
		Description: "Generate a complete authentication system with sessions, password hashing, and auth handlers",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GenAuthInput) (*mcp.CallToolResult, GenAuthOutput, error) {
		args := []string{"auth"}

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

		return nil, GenAuthOutput{
			Success: true,
			Message: "Successfully generated authentication system",
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
