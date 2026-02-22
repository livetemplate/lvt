package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
	"github.com/livetemplate/lvt/internal/kits"
	"github.com/livetemplate/lvt/internal/parser"
	"github.com/livetemplate/lvt/internal/validation"
)

func Gen(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printGenHelp) {
		return nil
	}

	if len(args) < 1 {
		// Show interactive prompt
		return interactiveGen()
	}

	// Route to subcommands
	subcommand := args[0]

	// Validate that subcommand doesn't look like a flag
	if err := ValidatePositionalArg(subcommand, "subcommand"); err != nil {
		return err
	}

	switch subcommand {
	case "resource":
		return GenResource(args[1:])
	case "view":
		return GenView(args[1:])
	case "schema":
		return GenSchema(args[1:])
	case "auth":
		return Auth(args[1:])
	case "stack":
		return GenStack(args[1:])
	default:
		return fmt.Errorf("unknown subcommand: %s\n\nAvailable subcommands:\n  resource  Generate full CRUD resource with database\n  view      Generate view-only handler (no database)\n  schema    Generate database schema only\n  auth      Generate authentication system\n  stack     Generate deployment stack configuration\n\nRun 'lvt gen' for interactive mode", subcommand)
	}
}

func interactiveGen() error {
	fmt.Println("What would you like to generate?")
	fmt.Println()
	fmt.Println("  1. Resource - Full CRUD with database (handler + template + migration + queries)")
	fmt.Println("  2. View     - UI only, no database (handler + template)")
	fmt.Println("  3. Schema   - Database tables only (migration + schema + queries)")
	fmt.Println("  4. Auth     - Authentication system (handler + middleware + migrations + E2E tests)")
	fmt.Println()
	fmt.Print("Enter your choice (1-4): ")

	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		fmt.Println("\nGenerating resource...")
		fmt.Println("You can also use: lvt gen resource <name> <field:type>...")
		fmt.Println()
		return fmt.Errorf("interactive resource generation not yet implemented - use: lvt gen resource <name> <field:type>...")
	case "2":
		fmt.Println("\nGenerating view...")
		fmt.Println("You can also use: lvt gen view <name>")
		fmt.Println()
		return fmt.Errorf("interactive view generation not yet implemented - use: lvt gen view <name>")
	case "3":
		fmt.Println("\nGenerating schema...")
		fmt.Println("You can also use: lvt gen schema <table> <field:type>...")
		fmt.Println()
		return fmt.Errorf("interactive schema generation not yet implemented - use: lvt gen schema <table> <field:type>...")
	case "4":
		fmt.Println("\nGenerating auth system...")
		fmt.Println("You can also use: lvt gen auth [StructName] [table_name] [flags...]")
		fmt.Println()
		return fmt.Errorf("interactive auth generation not yet implemented - use: lvt gen auth")
	default:
		return fmt.Errorf("invalid choice: %s", choice)
	}
}

func GenResource(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printGenResourceHelp) {
		return nil
	}

	if len(args) < 1 {
		return fmt.Errorf("resource name required")
	}

	// Get current directory for project config
	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load project config
	projectConfig, err := config.LoadProjectConfig(basePath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	kit := projectConfig.GetKit()

	// Load kit manifest to get CSS framework
	loader := kits.DefaultLoader()
	kitInfo, err := loader.Load(kit)
	if err != nil {
		return fmt.Errorf("failed to load kit: %w", err)
	}
	cssFramework := kitInfo.Manifest.CSSFramework

	// Parse flags (removed --css and --mode, they're now locked in config)
	paginationMode := "infinite" // default
	pageSize := 20               // default
	editMode := "modal"          // default
	skipValidation := false
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--pagination" && i+1 < len(args) {
			paginationMode = args[i+1]
			i++ // skip next arg
		} else if args[i] == "--page-size" && i+1 < len(args) {
			if size, err := fmt.Sscanf(args[i+1], "%d", &pageSize); err != nil || size == 0 || pageSize < 1 {
				pageSize = 20 // fallback to default
			}
			i++ // skip next arg
		} else if args[i] == "--edit-mode" && i+1 < len(args) {
			editMode = args[i+1]
			i++ // skip next arg
		} else if args[i] == "--skip-validation" {
			skipValidation = true
		} else {
			filteredArgs = append(filteredArgs, args[i])
		}
	}

	if len(filteredArgs) < 1 {
		return fmt.Errorf("resource name required")
	}

	resourceName := filteredArgs[0]

	// Validate that resource name doesn't look like a flag
	if err := ValidatePositionalArg(resourceName, "resource name"); err != nil {
		return err
	}

	fieldArgs := filteredArgs[1:]

	if len(fieldArgs) == 0 {
		return fmt.Errorf("at least one field required (format: name:type)")
	}

	// Validate pagination mode
	validPaginationModes := map[string]bool{"infinite": true, "load-more": true, "prev-next": true, "numbers": true}
	if !validPaginationModes[paginationMode] {
		return fmt.Errorf("invalid pagination mode: %s (valid: infinite, load-more, prev-next, numbers)", paginationMode)
	}

	// Validate edit mode
	validEditModes := map[string]bool{"modal": true, "page": true}
	if !validEditModes[editMode] {
		return fmt.Errorf("invalid edit mode: %s (valid: modal, page)", editMode)
	}

	// Parse fields with type inference support
	fields, err := parseFieldsWithInference(fieldArgs)
	if err != nil {
		return err
	}

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	fmt.Printf("Generating CRUD resource: %s\n", resourceName)
	fmt.Printf("Kit: %s\n", kit)
	fmt.Printf("CSS Framework: %s\n", cssFramework)
	fmt.Printf("Pagination: %s (page size: %d)\n", paginationMode, pageSize)
	fmt.Printf("Edit Mode: %s\n", editMode)
	fmt.Printf("Fields: ")
	for i, f := range fields {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s:%s", f.Name, f.Type)
	}
	fmt.Println()

	if err := generator.GenerateResource(basePath, moduleName, resourceName, fields, kit, cssFramework, paginationMode, pageSize, editMode); err != nil {
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	if !skipValidation {
		validationErr = runPostGenValidation(basePath)
	}

	resourceNameLower := strings.ToLower(resourceName)

	if validationErr != nil {
		fmt.Println()
		fmt.Println("⚠️  Resource generated, but validation found issues.")
	} else {
		fmt.Println()
		fmt.Println("✅ Resource generated successfully!")
	}
	fmt.Println()
	fmt.Println("Files created:")
	fmt.Printf("  app/%s/%s.go\n", resourceNameLower, resourceNameLower)
	fmt.Printf("  app/%s/%s.tmpl\n", resourceNameLower, resourceNameLower)
	fmt.Println()
	fmt.Println("Files updated:")
	fmt.Println("  database/schema.sql")
	fmt.Println("  database/queries.sql")
	fmt.Println()
	fmt.Println("Route auto-injected:")
	fmt.Printf("  http.Handle(\"/%s\", %s.Handler(queries))\n", resourceNameLower, resourceNameLower)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run migration:")
	fmt.Println("     lvt migration up")
	fmt.Println("  2. Run your app")
	fmt.Println()

	return validationErr
}

func GenView(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printGenViewHelp) {
		return nil
	}

	// Parse --skip-validation flag before checking positional args,
	// otherwise `lvt gen view --skip-validation` panics on args[0].
	skipValidation := false
	var filteredArgs []string
	for _, arg := range args {
		if arg == "--skip-validation" {
			skipValidation = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	args = filteredArgs
	_ = skipValidation // used later

	if len(args) < 1 {
		return fmt.Errorf("view name required")
	}

	// Get current directory for project config
	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load project config
	projectConfig, err := config.LoadProjectConfig(basePath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	kit := projectConfig.GetKit()

	// Load kit manifest to get CSS framework
	loader := kits.DefaultLoader()
	kitInfo, err := loader.Load(kit)
	if err != nil {
		return fmt.Errorf("failed to load kit: %w", err)
	}
	cssFramework := kitInfo.Manifest.CSSFramework

	viewName := args[0]

	// Validate that view name doesn't look like a flag
	if err := ValidatePositionalArg(viewName, "view name"); err != nil {
		return err
	}

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	fmt.Printf("Generating view-only handler: %s\n", viewName)
	fmt.Printf("Kit: %s\n", kit)
	fmt.Printf("CSS Framework: %s\n", cssFramework)

	if err := generator.GenerateView(basePath, moduleName, viewName, kit, cssFramework); err != nil {
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	if !skipValidation {
		validationErr = runPostGenValidation(basePath)
	}

	viewNameLower := strings.ToLower(viewName)

	if validationErr != nil {
		fmt.Println()
		fmt.Println("⚠️  View generated, but validation found issues.")
	} else {
		fmt.Println()
		fmt.Println("✅ View generated successfully!")
	}
	fmt.Println()
	fmt.Println("Files created:")
	fmt.Printf("  app/%s/%s.go\n", viewNameLower, viewNameLower)
	fmt.Printf("  app/%s/%s.tmpl\n", viewNameLower, viewNameLower)
	fmt.Printf("  app/%s/%s_test.go\n", viewNameLower, viewNameLower)
	fmt.Println()
	fmt.Println("Route auto-injected:")
	fmt.Printf("  http.Handle(\"/%s\", %s.Handler())\n", viewNameLower, viewNameLower)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Customize handler: app/%s/%s.go\n", viewNameLower, viewNameLower)
	fmt.Printf("  2. Edit template: app/%s/%s.tmpl\n", viewNameLower, viewNameLower)
	fmt.Println("  3. Run your app")
	fmt.Println()

	return validationErr
}

func GenSchema(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printGenSchemaHelp) {
		return nil
	}

	// Parse --skip-validation flag before checking positional args,
	// otherwise `lvt gen schema --skip-validation` panics on args[0].
	skipValidation := false
	var filteredArgs []string
	for _, arg := range args {
		if arg == "--skip-validation" {
			skipValidation = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	args = filteredArgs
	_ = skipValidation // used later

	if len(args) < 1 {
		return fmt.Errorf("table name required")
	}

	// Get current directory for project config
	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Load project config
	projectConfig, err := config.LoadProjectConfig(basePath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	kit := projectConfig.GetKit()

	// Load kit manifest to get CSS framework
	loader := kits.DefaultLoader()
	kitInfo, err := loader.Load(kit)
	if err != nil {
		return fmt.Errorf("failed to load kit: %w", err)
	}
	cssFramework := kitInfo.Manifest.CSSFramework

	tableName := args[0]

	// Validate that table name doesn't look like a flag
	if err := ValidatePositionalArg(tableName, "table name"); err != nil {
		return err
	}

	fieldArgs := args[1:]

	if len(fieldArgs) == 0 {
		return fmt.Errorf("at least one field required (format: name:type)")
	}

	// Parse fields with type inference support
	fields, err := parseFieldsWithInference(fieldArgs)
	if err != nil {
		return err
	}

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	fmt.Printf("Generating database schema: %s\n", tableName)
	fmt.Printf("Kit: %s\n", kit)
	fmt.Printf("Fields: ")
	for i, f := range fields {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s:%s", f.Name, f.Type)
	}
	fmt.Println()

	if err := generator.GenerateSchema(basePath, moduleName, tableName, fields, kit, cssFramework); err != nil {
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	if !skipValidation {
		validationErr = runPostGenValidation(basePath)
	}

	tableNameLower := strings.ToLower(tableName)

	if validationErr != nil {
		fmt.Println()
		fmt.Println("⚠️  Schema generated, but validation found issues.")
	} else {
		fmt.Println()
		fmt.Println("✅ Schema generated successfully!")
	}
	fmt.Println()
	fmt.Println("Files created/updated:")
	fmt.Println("  database/migrations/<timestamp>_create_" + tableNameLower + ".sql")
	fmt.Println("  database/schema.sql (updated)")
	fmt.Println("  database/queries.sql (updated)")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run migration:")
	fmt.Println("     lvt migration up")
	fmt.Println("  2. Use generated types in your handlers")
	fmt.Println()

	return validationErr
}

// runPostGenValidation runs structural validation (go.mod, templates, migrations)
// after code generation. It skips compilation because the app may not compile until
// sqlc generate is run. Prints the formatted result and returns an error if found.
//
// TODO: accept context.Context so Ctrl+C propagates to validation.
// Structural checks are fast today so context.Background() is acceptable.
func runPostGenValidation(basePath string) error {
	fmt.Println("Running validation...")
	result := validation.ValidatePostGen(context.Background(), basePath)
	fmt.Print(result.Format())
	if result.HasErrors() {
		return fmt.Errorf("validation failed with %d error(s)", result.ErrorCount())
	}
	return nil
}

func parseFieldsWithInference(fieldArgs []string) ([]parser.Field, error) {
	// Try parsing with type inference first
	fields := make([]parser.Field, 0, len(fieldArgs))

	for _, arg := range fieldArgs {
		var name, typ string

		// Check if it contains ":"
		if strings.Contains(arg, ":") {
			// Explicit type - use normal parser
			parts := strings.SplitN(arg, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid field format: %s (expected name:type)", arg)
			}
			name = strings.TrimSpace(parts[0])
			typ = strings.TrimSpace(parts[1])
		} else {
			// No type specified - infer from name
			name = strings.TrimSpace(arg)
			typ = inferTypeForDirectMode(name)
		}

		// Map to Go and SQL types
		goType, sqlType, isTextarea, err := parser.MapType(typ)
		if err != nil {
			return nil, fmt.Errorf("field '%s': %w", name, err)
		}

		// Create field with reference metadata
		field := parser.Field{
			Name:       name,
			Type:       typ,
			GoType:     goType,
			SQLType:    sqlType,
			IsTextarea: isTextarea,
		}

		// Parse reference metadata if it's a reference type
		if strings.HasPrefix(strings.ToLower(typ), "references:") {
			parts := strings.Split(typ, ":")
			if len(parts) < 2 {
				return nil, fmt.Errorf("field '%s': invalid references syntax, expected 'references:table_name'", name)
			}

			field.IsReference = true
			field.ReferencedTable = parts[1]
			field.OnDelete = "CASCADE" // Default

			// Check for custom on_delete action
			if len(parts) > 2 {
				action := strings.ToUpper(parts[2])
				switch action {
				case "CASCADE", "SET NULL", "RESTRICT", "NO ACTION", "SET_NULL":
					if action == "SET_NULL" {
						action = "SET NULL"
					}
					field.OnDelete = action
				default:
					return nil, fmt.Errorf("field '%s': invalid ON DELETE action '%s'", name, parts[2])
				}
			}
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func inferTypeForDirectMode(fieldName string) string {
	lower := strings.ToLower(fieldName)

	// Exact matches for common field names
	exactMatches := map[string]string{
		"name": "string", "email": "string", "title": "string",
		"description": "text", "content": "text", "body": "text",
		"username": "string", "password": "string", "token": "string",
		"url": "string", "slug": "string", "path": "string",
		"address": "string", "city": "string", "state": "string",
		"country": "string", "phone": "string", "status": "string",

		"age": "int", "count": "int", "quantity": "int",
		"views": "int", "likes": "int", "shares": "int",
		"year": "int", "rating": "int",

		"price": "float", "amount": "float", "total": "float",
		"latitude": "float", "longitude": "float",

		"enabled": "bool", "active": "bool", "visible": "bool",
		"published": "bool", "deleted": "bool", "featured": "bool",

		"created_at": "time", "updated_at": "time", "deleted_at": "time",
		"published_at": "time", "expires_at": "time",
	}

	if t, ok := exactMatches[lower]; ok {
		return t
	}

	// Pattern matching for suffixes/prefixes
	if strings.HasSuffix(lower, "_at") || strings.HasSuffix(lower, "_date") ||
		strings.HasSuffix(lower, "_time") || strings.HasSuffix(lower, "date") {
		return "time"
	}

	if strings.HasPrefix(lower, "is_") || strings.HasPrefix(lower, "has_") ||
		strings.HasPrefix(lower, "can_") || strings.HasPrefix(lower, "should_") {
		return "bool"
	}

	if strings.HasSuffix(lower, "_count") || strings.HasSuffix(lower, "_number") ||
		strings.HasSuffix(lower, "_id") || strings.HasSuffix(lower, "id") {
		return "int"
	}

	if strings.HasSuffix(lower, "_price") || strings.HasSuffix(lower, "_amount") ||
		strings.HasSuffix(lower, "_total") || strings.HasSuffix(lower, "price") {
		return "float"
	}

	// Default to string
	return "string"
}

func getModuleName() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	return "", fmt.Errorf("module name not found in go.mod")
}
