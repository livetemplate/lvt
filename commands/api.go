package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
	"github.com/livetemplate/lvt/internal/kits"
)

// GenAPI generates a JSON API handler for a resource.
func GenAPI(args []string) error {
	if ShowHelpIfRequested(args, printGenAPIHelp) {
		return nil
	}

	if len(args) < 1 {
		return fmt.Errorf("resource name required")
	}

	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectConfig, err := config.LoadProjectConfig(basePath)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	kit := projectConfig.GetKit()

	loader := kits.DefaultLoader()
	kitInfo, err := loader.Load(kit)
	if err != nil {
		return fmt.Errorf("failed to load kit: %w", err)
	}
	_ = kitInfo

	// Parse flags
	skipValidation := false
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--skip-validation" {
			skipValidation = true
		} else {
			filteredArgs = append(filteredArgs, args[i])
		}
	}
	_ = skipValidation

	if len(filteredArgs) < 1 {
		return fmt.Errorf("resource name required")
	}

	resourceName := filteredArgs[0]

	if err := ValidatePositionalArg(resourceName, "resource name"); err != nil {
		return err
	}

	fieldArgs := filteredArgs[1:]
	if len(fieldArgs) == 0 {
		return fmt.Errorf("at least one field required (format: name:type)")
	}

	fields, err := parseFieldsWithInference(fieldArgs)
	if err != nil {
		return err
	}

	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	fmt.Printf("Generating API resource: %s\n", resourceName)
	fmt.Printf("Fields: ")
	for i, f := range fields {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s:%s", f.Name, f.Type)
	}
	fmt.Println()

	if err := generator.GenerateAPI(basePath, moduleName, resourceName, fields, kit); err != nil {
		return err
	}

	resourceNameLower := strings.ToLower(resourceName)

	fmt.Println()
	fmt.Println("✅ API resource generated successfully!")
	fmt.Println()
	fmt.Println("Files created:")
	fmt.Printf("  app/api/%s.go\n", resourceNameLower)
	fmt.Printf("  app/api/%s_test.go\n", resourceNameLower)
	fmt.Println()
	fmt.Println("Files updated:")
	fmt.Println("  database/queries.sql (paginated queries added)")
	fmt.Println()
	fmt.Println("API endpoints:")
	fmt.Printf("  GET    /api/v1/%s        List (paginated)\n", resourceNameLower)
	fmt.Printf("  POST   /api/v1/%s        Create\n", resourceNameLower)
	fmt.Printf("  GET    /api/v1/%s/{id}   Get by ID\n", resourceNameLower)
	fmt.Printf("  PUT    /api/v1/%s/{id}   Update\n", resourceNameLower)
	fmt.Printf("  DELETE /api/v1/%s/{id}   Delete\n", resourceNameLower)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run migration:")
	fmt.Println("     lvt migration up")
	fmt.Println("  2. Regenerate sqlc code:")
	fmt.Println("     sqlc generate")
	fmt.Println("  3. Run your app")
	fmt.Println()

	return nil
}

func printGenAPIHelp() {
	fmt.Println("Usage: lvt gen api <resource> <field:type>... [--skip-validation]")
	fmt.Println()
	fmt.Println("Generates a JSON API with RESTful CRUD endpoints.")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  lvt gen api posts title content:text published:bool")
	fmt.Println()
	fmt.Println("Generated endpoints:")
	fmt.Println("  GET    /api/v1/<resource>        List (paginated)")
	fmt.Println("  POST   /api/v1/<resource>        Create")
	fmt.Println("  GET    /api/v1/<resource>/{id}   Get by ID")
	fmt.Println("  PUT    /api/v1/<resource>/{id}   Update")
	fmt.Println("  DELETE /api/v1/<resource>/{id}   Delete")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --skip-validation   Skip post-generation validation")
	fmt.Println()
}
