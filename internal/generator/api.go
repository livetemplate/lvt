package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/kits"
	"github.com/livetemplate/lvt/internal/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// APIData holds template data for API handler generation.
type APIData struct {
	PackageName          string
	ModuleName           string
	ResourceName         string
	ResourceNameLower    string
	ResourceNameSingular string
	ResourceNamePlural   string
	TableName            string
	Fields               []FieldData
}

// GenerateAPI generates a JSON API handler for a resource.
func GenerateAPI(basePath, moduleName, resourceName string, fields []parser.Field, kitName string) error {
	if kitName == "" {
		kitName = "multi"
	}

	projectConfig, err := config.LoadProjectConfig(basePath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	if kitName == "" {
		kitName = projectConfig.GetKit()
	}

	kitLoader := kits.DefaultLoader()

	resourceNameLower := strings.ToLower(resourceName)
	titleCaser := cases.Title(language.English)
	resourceName = titleCaser.String(resourceNameLower)

	resourceNameSingular := singularize(resourceNameLower)
	resourceNameSingularCap := titleCaser.String(resourceNameSingular)
	resourceNamePluralCap := titleCaser.String(pluralize(resourceNameSingular))
	tableName := pluralize(resourceNameSingular)

	fieldData := FieldDataFromFields(fields)

	data := APIData{
		PackageName:          "api",
		ModuleName:           moduleName,
		ResourceName:         resourceName,
		ResourceNameLower:    resourceNameLower,
		ResourceNameSingular: resourceNameSingularCap,
		ResourceNamePlural:   resourceNamePluralCap,
		TableName:            tableName,
		Fields:               fieldData,
	}

	// Create api directory
	apiDir := filepath.Join(basePath, "app", "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api directory: %w", err)
	}

	// Generate handler
	handlerTmpl, err := kitLoader.LoadKitTemplate(kitName, "api/handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load API handler template: %w", err)
	}
	handlerFile := filepath.Join(apiDir, resourceNameLower+".go")
	if err := generateAPIFile(string(handlerTmpl), data, handlerFile); err != nil {
		return fmt.Errorf("failed to generate API handler: %w", err)
	}

	// Generate test
	testTmpl, err := kitLoader.LoadKitTemplate(kitName, "api/test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load API test template: %w", err)
	}
	testFile := filepath.Join(apiDir, resourceNameLower+"_test.go")
	if err := generateAPIFile(string(testTmpl), data, testFile); err != nil {
		return fmt.Errorf("failed to generate API test: %w", err)
	}

	// Check if DB schema already exists (from gen resource)
	schemaPath := filepath.Join(basePath, "database", "schema.sql")
	schemaExists := false
	if schemaData, err := os.ReadFile(schemaPath); err == nil {
		schemaExists = strings.Contains(string(schemaData), "CREATE TABLE IF NOT EXISTS "+tableName)
	}

	if !schemaExists {
		// Generate schema, migration, and base queries using resource templates
		resourceData := ResourceData{
			PackageName:          resourceNameLower,
			ModuleName:           moduleName,
			ResourceName:         resourceName,
			ResourceNameLower:    resourceNameLower,
			ResourceNameSingular: resourceNameSingularCap,
			ResourceNamePlural:   resourceNamePluralCap,
			TableName:            tableName,
			Fields:               fieldData,
		}

		kit, err := kitLoader.Load(kitName)
		if err != nil {
			return fmt.Errorf("failed to load kit: %w", err)
		}
		if kit.Helpers == nil {
			kit.SetHelpersForFramework("tailwind")
		}

		// Generate migration
		dbDir := filepath.Join(basePath, "database")
		migrationsDir := filepath.Join(dbDir, "migrations")
		if err := os.MkdirAll(migrationsDir, 0755); err != nil {
			return fmt.Errorf("failed to create migrations directory: %w", err)
		}

		migrationTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/migration.sql.tmpl")
		if err != nil {
			return fmt.Errorf("failed to load migration template: %w", err)
		}

		timestamp := time.Now()
		var migrationPath string
		const maxRetries = 3600
		for i := 0; i < maxRetries; i++ {
			timestampStr := timestamp.Format("20060102150405")
			migrationPath = filepath.Join(migrationsDir, fmt.Sprintf("%s_create_%s.sql", timestampStr, tableName))
			matches, _ := filepath.Glob(filepath.Join(migrationsDir, timestampStr+"_*.sql"))
			if len(matches) == 0 {
				break
			}
			timestamp = timestamp.Add(1 * time.Second)
		}

		if err := generateFile(string(migrationTmpl), resourceData, migrationPath, kit); err != nil {
			return fmt.Errorf("failed to generate migration: %w", err)
		}

		schemaTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/schema.sql.tmpl")
		if err != nil {
			return fmt.Errorf("failed to load schema template: %w", err)
		}
		if err := appendToFile(string(schemaTmpl), resourceData, filepath.Join(dbDir, "schema.sql"), "\n", kit); err != nil {
			return fmt.Errorf("failed to append to schema: %w", err)
		}

		// Generate base CRUD queries
		queriesTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/queries.sql.tmpl")
		if err != nil {
			return fmt.Errorf("failed to load queries template: %w", err)
		}
		if err := appendToFile(string(queriesTmpl), resourceData, filepath.Join(dbDir, "queries.sql"), "\n", kit); err != nil {
			return fmt.Errorf("failed to append base queries: %w", err)
		}
	}

	// Append paginated API queries
	apiQueriesTmpl, err := kitLoader.LoadKitTemplate(kitName, "api/queries.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load API queries template: %w", err)
	}

	// Check if paginated queries already exist
	queriesPath := filepath.Join(basePath, "database", "queries.sql")
	queriesData, _ := os.ReadFile(queriesPath)
	paginatedName := fmt.Sprintf("List%sPaginated", resourceNamePluralCap)
	if !strings.Contains(string(queriesData), paginatedName) {
		kit, err := kitLoader.Load(kitName)
		if err != nil {
			return fmt.Errorf("failed to load kit: %w", err)
		}
		if kit.Helpers == nil {
			kit.SetHelpersForFramework("tailwind")
		}
		if err := appendToFile(string(apiQueriesTmpl), data, queriesPath, "\n", kit); err != nil {
			return fmt.Errorf("failed to append API queries: %w", err)
		}
	}

	// Inject route into main.go
	mainGoPath := findMainGo(basePath)
	if mainGoPath != "" {
		route := RouteInfo{
			Path:        "/api/v1/" + resourceNameLower + "/",
			PackageName: "api",
			HandlerCall: "api.Handler(queries)",
			ImportPath:  moduleName + "/app/api",
		}
		if err := InjectRoute(mainGoPath, route); err != nil {
			fmt.Printf("⚠️  Could not auto-inject API route: %v\n", err)
			fmt.Printf("   Add manually: http.Handle(\"/api/v1/%s/\", api.Handler(queries))\n", resourceNameLower)
		}
	}

	// Register resource for home page
	if err := RegisterResource(basePath, data.ResourceName, "/api/v1/"+resourceNameLower, "api"); err != nil {
		fmt.Printf("⚠️  Could not register API resource: %v\n", err)
	}

	return nil
}

func generateAPIFile(tmplStr string, data APIData, outPath string) error {
	funcs := template.FuncMap{
		"title":       cases.Title(language.English).String,
		"lower":       strings.ToLower,
		"upper":       strings.ToUpper,
		"camelCase":   toCamelCase,
		"singularize": singularizeForTemplate,
	}

	tmpl, err := template.New("api").Delims("[[", "]]").Funcs(funcs).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}
