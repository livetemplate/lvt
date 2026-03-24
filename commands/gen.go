package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
	"github.com/livetemplate/lvt/internal/kits"
	"github.com/livetemplate/lvt/internal/parser"
	"github.com/livetemplate/lvt/internal/telemetry"
	"github.com/livetemplate/lvt/internal/validation"
	"github.com/livetemplate/lvt/internal/validator"
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
	case "queue":
		return GenQueue(args[1:])
	case "job":
		return GenJob(args[1:])
	case "authz":
		return Authz(args[1:])
	case "api":
		return GenAPI(args[1:])
	case "task":
		return GenTask(args[1:])
	default:
		return fmt.Errorf("unknown subcommand: %s\n\nAvailable subcommands:\n  resource  Generate full CRUD resource with database\n  view      Generate view-only handler (no database)\n  schema    Generate database schema only\n  auth      Generate authentication system\n  authz     Generate role-based authorization\n  api       Generate JSON API endpoints\n  stack     Generate deployment stack configuration\n  queue     Set up background job processing (River)\n  job       Scaffold a new background job handler\n  task      Scaffold a new scheduled task\n\nRun 'lvt gen' for interactive mode", subcommand)
	}
}

func interactiveGen() error {
	fmt.Println("Usage: lvt gen <subcommand> [args...]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  resource <name> <field:type>...       Generate full CRUD with database")
	fmt.Println("  view <name>                           Generate view-only handler (no database)")
	fmt.Println("  schema <table> <field:type>...        Generate database schema only")
	fmt.Println("  auth [StructName] [table_name]        Generate authentication system")
	fmt.Println("  stack <target>                        Generate deployment stack configuration")
	fmt.Println("  queue                                 Set up background job processing (River)")
	fmt.Println("  job <name>                            Scaffold a new background job handler")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen resource posts title content:text published:bool")
	fmt.Println("  lvt gen view dashboard")
	fmt.Println("  lvt gen auth")
	fmt.Println()
	return nil
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
	parentResource := ""
	withAuthz := false
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
		} else if args[i] == "--parent" && i+1 < len(args) {
			parentResource = args[i+1]
			i++ // skip next arg
		} else if args[i] == "--with-authz" {
			withAuthz = true
		} else {
			filteredArgs = append(filteredArgs, args[i])
		}
	}

	if len(filteredArgs) < 1 {
		return fmt.Errorf("resource name required")
	}

	resourceName := filteredArgs[0]

	// Validate --with-authz prerequisites
	if withAuthz {
		if _, err := os.Stat(filepath.Join(basePath, "app", "auth")); os.IsNotExist(err) {
			return fmt.Errorf("--with-authz requires authentication. Run 'lvt gen auth' and 'lvt gen authz' first")
		}
	}

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

	// Validate --parent flag
	if parentResource != "" {
		parentResource = strings.ToLower(parentResource)
		// Check parent handler exists
		parentHandlerPath := filepath.Join(basePath, "app", parentResource, parentResource+".go")
		if _, err := os.Stat(parentHandlerPath); os.IsNotExist(err) {
			return fmt.Errorf("parent resource %q not found: %s does not exist.\nGenerate the parent first: lvt gen resource %s ...", parentResource, parentHandlerPath, parentResource)
		}
		// Check child has a reference field pointing to the parent table
		hasRef := false
		for _, f := range fields {
			if f.IsReference && f.ReferencedTable == parentResource {
				hasRef = true
				break
			}
		}
		if !hasRef {
			return fmt.Errorf("child resource must have a field referencing parent table %q (e.g., %s_id:references:%s)", parentResource, generator.Singularize(parentResource), parentResource)
		}
	}

	// Get module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	// Start telemetry capture
	collector := telemetry.NewCollector()
	defer collector.Close()
	capture := collector.StartCapture("gen resource", map[string]any{
		"resource_name":   resourceName,
		"fields":          fieldArgs,
		"kit":             kit,
		"pagination_mode": paginationMode,
		"edit_mode":       editMode,
	})
	capture.SetKit(kit) // also sets the dedicated Kit column for SQL queries; inputs has it for context

	// Detect which components this resource will use and record for telemetry
	resourceData := generator.ResourceData{Fields: generator.FieldDataFromFields(fields)}
	compUsage := generator.ComputeComponentUsage(resourceData)
	capture.RecordComponentsUsed(telemetry.ComponentsFromUsage(compUsage))

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

	styles := projectConfig.Styles
	if err := generator.GenerateResource(basePath, moduleName, resourceName, fields, kit, cssFramework, styles, paginationMode, pageSize, editMode, parentResource, withAuthz); err != nil {
		capture.RecordError(telemetry.GenerationError{Phase: "generation", Message: err.Error()})
		capture.AttributeComponentErrors() // attribute errors on failure path
		capture.Complete(false, "")
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	var validationResult *validator.ValidationResult
	if !skipValidation {
		validationResult, validationErr = runPostGenValidation(basePath)
	}
	if validationErr != nil {
		capture.RecordError(telemetry.GenerationError{
			Phase:   "validation",
			Message: validationErr.Error(),
		})
	}
	capture.AttributeComponentErrors() // attribute any captured errors to components before completing
	capture.Complete(validationErr == nil, marshalValidationResult(validationResult))

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
	if parentResource != "" {
		fmt.Printf("Embedded in parent: %s\n", parentResource)
		fmt.Printf("  app/%s/%s.go (modified)\n", parentResource, parentResource)
		fmt.Printf("  app/%s/%s.tmpl (modified)\n", parentResource, parentResource)
		fmt.Println()
		fmt.Println("No separate route — child is rendered on the parent's detail page.")
	} else if compUsage.UseUpload {
		fmt.Println("Manual route setup required (file uploads need storage):")
		fmt.Println("  Add to main.go:")
		fmt.Printf("    store := storage.NewLocalStore(\"uploads\", \"/uploads\")\n")
		fmt.Printf("    http.Handle(\"/uploads/\", http.StripPrefix(\"/uploads/\", store.FileServer()))\n")
		fmt.Printf("    http.Handle(\"/%s\", %s.Handler(queries, store))\n", resourceNameLower, resourceNameLower)
		fmt.Println()
		fmt.Println("  Add imports:")
		fmt.Println("    \"github.com/livetemplate/lvt/pkg/storage\"")
		fmt.Printf("    \"%s/app/%s\"\n", moduleName, resourceNameLower)
	} else {
		fmt.Println("Route auto-injected:")
		fmt.Printf("  http.Handle(\"/%s\", %s.Handler(queries))\n", resourceNameLower, resourceNameLower)
	}
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

	// Start telemetry capture
	collector := telemetry.NewCollector()
	defer collector.Close()
	capture := collector.StartCapture("gen view", map[string]any{
		"view_name": viewName,
		"kit":       kit,
	})
	capture.SetKit(kit) // also sets the dedicated Kit column for SQL queries; inputs has it for context

	fmt.Printf("Generating view-only handler: %s\n", viewName)
	fmt.Printf("Kit: %s\n", kit)
	fmt.Printf("CSS Framework: %s\n", cssFramework)

	if err := generator.GenerateView(basePath, moduleName, viewName, kit, cssFramework); err != nil {
		capture.RecordError(telemetry.GenerationError{Phase: "generation", Message: err.Error()})
		capture.Complete(false, "")
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	var validationResult *validator.ValidationResult
	if !skipValidation {
		validationResult, validationErr = runPostGenValidation(basePath)
	}
	if validationErr != nil {
		capture.RecordError(telemetry.GenerationError{
			Phase:   "validation",
			Message: validationErr.Error(),
		})
	}
	capture.Complete(validationErr == nil, marshalValidationResult(validationResult))

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

	// Start telemetry capture
	collector := telemetry.NewCollector()
	defer collector.Close()
	capture := collector.StartCapture("gen schema", map[string]any{
		"table_name": tableName,
		"fields":     args[1:],
		"kit":        kit,
	})
	capture.SetKit(kit) // also sets the dedicated Kit column for SQL queries; inputs has it for context

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
		capture.RecordError(telemetry.GenerationError{Phase: "generation", Message: err.Error()})
		capture.Complete(false, "")
		return err
	}

	// Post-generation validation (run before printing success banner)
	var validationErr error
	var validationResult *validator.ValidationResult
	if !skipValidation {
		validationResult, validationErr = runPostGenValidation(basePath)
	}
	if validationErr != nil {
		capture.RecordError(telemetry.GenerationError{
			Phase:   "validation",
			Message: validationErr.Error(),
		})
	}
	capture.Complete(validationErr == nil, marshalValidationResult(validationResult))

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
// sqlc generate is run. Prints the formatted result and returns both the result
// (for telemetry) and an error if validation found issues.
func runPostGenValidation(basePath string) (*validator.ValidationResult, error) {
	fmt.Println("Running validation...")
	result := validation.ValidatePostGen(context.Background(), basePath)
	fmt.Print(result.Format())
	if result.HasErrors() {
		return result, fmt.Errorf("validation failed with %d error(s)", result.ErrorCount())
	}
	return result, nil
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

		// Delegate select and file/image types to ParseFields to avoid duplication
		lowerTyp := strings.ToLower(typ)
		if lowerTyp == "select" || lowerTyp == "file" || lowerTyp == "image" {
			parsed, err := parser.ParseFields([]string{arg})
			if err != nil {
				return nil, err
			}
			fields = append(fields, parsed...)
			continue
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
			Metadata:   parser.GetFieldMetadata(typ),
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
			field.Metadata = parser.FieldMetadata{ValidateTag: "required", HTMLInputType: "text"}

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

		"avatar": "image", "photo": "image", "picture": "image",
		"logo": "image", "thumbnail": "image", "icon": "image",
		"cover": "image", "banner": "image", "headshot": "image",

		"document": "file", "attachment": "file", "resume": "file",
		"file": "file", "upload": "file", "certificate": "file",

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

// marshalValidationResult serialises a validation result to JSON for telemetry.
// Returns empty string if the result is nil.
func marshalValidationResult(result *validator.ValidationResult) string {
	if result == nil {
		return ""
	}
	b, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(b)
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
