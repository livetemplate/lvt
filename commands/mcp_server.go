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
	registerMigrationTools(server)

	// Start server with stdio transport
	ctx := context.Background()
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}

// NewAppInput defines the input schema for lvt new
type NewAppInput struct {
	Name   string `json:"name" jsonschema:"required,description=Application name"`
	Kit    string `json:"kit,omitempty" jsonschema:"enum=multi|single|simple,description=Template kit: multi (Tailwind multi-page), single (Tailwind SPA), simple (Pico CSS)"`
	CSS    string `json:"css,omitempty" jsonschema:"enum=tailwind|bulma|pico|none,description=CSS framework to use"`
	Module string `json:"module,omitempty" jsonschema:"description=Go module name (defaults to app name)"`
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
	Name   string            `json:"name" jsonschema:"required,description=Resource name (singular, e.g., 'post', 'user')"`
	Fields map[string]string `json:"fields" jsonschema:"required,description=Field definitions as name:type pairs (e.g., {'title': 'string', 'content': 'text', 'published_at': 'time'})"`
}

// GenResourceOutput defines the output schema
type GenResourceOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Files   []string `json:"files,omitempty" jsonschema:"description=List of generated files"`
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
	Name string `json:"name" jsonschema:"required,description=View name (e.g., 'dashboard', 'counter')"`
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
	StructName string `json:"struct_name,omitempty" jsonschema:"description=Go struct name (default: User)"`
	TableName  string `json:"table_name,omitempty" jsonschema:"description=Database table name (default: users)"`
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
	Name string `json:"name,omitempty" jsonschema:"description=Migration name (only for create command)"`
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

// init registers the MCP server logger
func init() {
	// Suppress SDK logs in normal operation (they go to stderr which interferes with MCP protocol)
	log.SetOutput(os.Stderr)
}
