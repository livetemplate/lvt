package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/livetemplate/lvt/internal/kits"
	"github.com/livetemplate/lvt/internal/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateResource(basePath, moduleName, resourceName string, fields []parser.Field, kitName, cssFramework, styles, paginationMode string, pageSize int, editMode, parentResource string, withAuthz, searchable bool) error {
	// Defaults
	if kitName == "" {
		kitName = "multi"
	}
	if cssFramework == "" {
		cssFramework = "tailwind"
	}
	if paginationMode == "" {
		paginationMode = "infinite"
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if editMode == "" {
		editMode = "modal"
	}
	if styles == "" {
		styles = "tailwind"
	}
	validStyles := map[string]bool{"tailwind": true, "unstyled": true}
	if !validStyles[styles] {
		return fmt.Errorf("invalid styles adapter: %q (valid: tailwind, unstyled)", styles)
	}

	// appMode is the same as kit name in the new architecture
	appMode := kitName

	// Load kit using KitLoader
	kitLoader := kits.DefaultLoader()
	kit, err := kitLoader.Load(kitName)
	if err != nil {
		return fmt.Errorf("failed to load kit %q: %w", kitName, err)
	}

	// Inject CSS helpers based on CSS framework if kit doesn't have helpers
	if kit.Helpers == nil {
		if err := kit.SetHelpersForFramework(cssFramework); err != nil {
			return fmt.Errorf("failed to load CSS helpers for framework %q: %w", cssFramework, err)
		}
	}

	// Capitalize resource name and derive singular/plural forms
	resourceNameLower := strings.ToLower(resourceName)
	titleCaser := cases.Title(language.English)
	resourceName = titleCaser.String(resourceNameLower)

	// Derive singular and plural forms for struct/function names and table name
	resourceNameSingular := singularize(resourceNameLower)
	resourceNameSingularCap := titleCaser.String(resourceNameSingular)
	resourceNamePluralCap := titleCaser.String(pluralize(resourceNameSingular))
	tableName := pluralize(resourceNameSingular)

	fieldData := FieldDataFromFields(fields)

	// Read dev mode setting from .lvtrc
	devMode := ReadDevMode(basePath)

	data := ResourceData{
		PackageName:          resourceNameLower,
		ModuleName:           moduleName,
		ResourceName:         resourceName,
		ResourceNameLower:    resourceNameLower,
		ResourceNameSingular: resourceNameSingularCap,
		ResourceNamePlural:   resourceNamePluralCap,
		TableName:            tableName,
		Fields:               fieldData,
		Kit:                  kit,
		CSSFramework:         cssFramework, // Keep for backward compatibility
		DevMode:              devMode,
		PaginationMode:       paginationMode,
		PageSize:             pageSize,
		EditMode:             editMode,
		Styles:               styles,
		Searchable:           searchable,
		WithAuthz:            withAuthz,
	}
	if data.Searchable && len(data.SearchableFields()) == 0 {
		return fmt.Errorf("--searchable requires at least one string field for FTS indexing")
	}
	data.Components = ComputeComponentUsage(data)
	if data.Components.UseModal || data.Components.UseToast {
		if data.Styles == "unstyled" {
			data.StylesImportPath = "github.com/livetemplate/lvt/components/styles/unstyled"
		} else {
			data.StylesImportPath = "github.com/livetemplate/lvt/components/styles/tailwind"
		}
	}

	// Populate embedded/parent fields when --parent is specified
	if parentResource != "" {
		parentResource = strings.ToLower(parentResource)
		parentSingular := singularize(parentResource)
		data.ParentResource = parentResource
		data.ParentPackageName = parentResource
		data.ParentResourceSingular = titleCaser.String(parentSingular)
		data.IsEmbedded = true

		// Auto-detect parent reference field
		parentTable := pluralize(parentSingular)
		for _, f := range fieldData {
			if f.IsReference && f.ReferencedTable == parentTable {
				data.ParentReferenceField = f.Name
				break
			}
		}
		if data.ParentReferenceField == "" {
			return fmt.Errorf("could not find a reference field for parent table %q in child fields", parentTable)
		}
	}

	// Create resource directory
	resourceDir := filepath.Join(basePath, "app", resourceNameLower)
	if err := os.MkdirAll(resourceDir, 0755); err != nil {
		return fmt.Errorf("failed to create resource directory: %w", err)
	}

	// Embedded mode uses different templates and skips route/home injection
	if data.IsEmbedded {
		return generateEmbeddedResource(basePath, resourceDir, resourceNameLower, tableName, data, kitLoader, kitName, kit)
	}

	return generateStandaloneResource(basePath, resourceDir, resourceNameLower, tableName, moduleName, editMode, appMode, data, kitLoader, kitName, kit)
}

func generateEmbeddedResource(basePath, resourceDir, resourceNameLower, tableName string, data ResourceData, kitLoader *kits.KitLoader, kitName string, kit *kits.KitInfo) error {
	// Load embedded-specific templates
	handlerTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/embedded_handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded handler template: %w", err)
	}

	templateTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/embedded_template.tmpl.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded template: %w", err)
	}

	// Use embedded queries (adds filtered-by-parent query)
	queriesTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/embedded_queries.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read embedded queries template: %w", err)
	}

	// Migration and schema are the same as standalone
	migrationTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/migration.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read migration template: %w", err)
	}

	schemaTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/schema.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read schema template: %w", err)
	}

	// Generate embedded handler
	if err := generateFile(string(handlerTmpl), data, filepath.Join(resourceDir, resourceNameLower+".go"), kit); err != nil {
		return fmt.Errorf("failed to generate embedded handler: %w", err)
	}

	// Generate embedded template
	tmplPath := filepath.Join(resourceDir, resourceNameLower+".tmpl")
	if err := generateFile(string(templateTmpl), data, tmplPath, kit); err != nil {
		return fmt.Errorf("failed to generate embedded template: %w", err)
	}
	if err := ValidateTemplate(tmplPath); err != nil {
		return err
	}

	// Generate migration
	dbDir := filepath.Join(basePath, "database")
	migrationsDir := filepath.Join(dbDir, "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}
	timestamp := time.Now()
	var migrationPath string
	for {
		timestampStr := timestamp.Format("20060102150405")
		migrationPath = filepath.Join(migrationsDir, fmt.Sprintf("%s_create_%s.sql", timestampStr, tableName))
		matches, _ := filepath.Glob(filepath.Join(migrationsDir, timestampStr+"_*.sql"))
		if len(matches) == 0 {
			break
		}
		timestamp = timestamp.Add(1 * time.Second)
	}
	if err := generateFile(string(migrationTmpl), data, migrationPath, kit); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	// Append to schema.sql
	if err := appendToFile(string(schemaTmpl), data, filepath.Join(dbDir, "schema.sql"), "\n", kit); err != nil {
		return fmt.Errorf("failed to append to schema: %w", err)
	}

	// Append to queries.sql (embedded queries include filtered-by-parent)
	if err := appendToFile(string(queriesTmpl), data, filepath.Join(dbDir, "queries.sql"), "\n", kit); err != nil {
		return fmt.Errorf("failed to append to queries: %w", err)
	}

	// Inject child into parent handler and template
	parentGoPath := filepath.Join(basePath, "app", data.ParentPackageName, data.ParentPackageName+".go")
	parentTmplPath := filepath.Join(basePath, "app", data.ParentPackageName, data.ParentPackageName+".tmpl")
	if err := InjectEmbeddedChild(parentGoPath, data); err != nil {
		return fmt.Errorf("failed to inject child into parent handler: %w", err)
	}
	if err := InjectEmbeddedChildTemplate(parentTmplPath, data); err != nil {
		return fmt.Errorf("failed to inject child into parent template: %w", err)
	}

	// Skip route injection and home page registration (child is rendered on parent's page)
	return nil
}

func generateStandaloneResource(basePath, resourceDir, resourceNameLower, tableName, moduleName, editMode, appMode string, data ResourceData, kitLoader *kits.KitLoader, kitName string, kit *kits.KitInfo) error {
	// Read templates using kit loader (checks project kits, user kits, then embedded)
	handlerTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read handler template: %w", err)
	}

	// Load main template based on mode
	var templateTmpl []byte
	if appMode == "multi" {
		componentNames := []string{
			"layout.tmpl",
			"form.tmpl",
			"toolbar.tmpl",
			"table.tmpl",
			"pagination.tmpl",
			"search.tmpl",
			"stats.tmpl",
			"sort.tmpl",
			"detail.tmpl",
		}

		var fullTemplate string
		for _, compName := range componentNames {
			compTmpl, err := kitLoader.LoadKitComponent(kitName, compName)
			if err != nil {
				return fmt.Errorf("failed to load component %s: %w", compName, err)
			}
			fullTemplate += string(compTmpl) + "\n\n"
		}

		mainTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/template_components.tmpl.tmpl")
		if err != nil {
			return fmt.Errorf("failed to load main template: %w", err)
		}
		fullTemplate += string(mainTmpl)

		templateTmpl = []byte(fullTemplate)
	} else {
		templateTmpl, err = kitLoader.LoadKitTemplate(kitName, "resource/template.tmpl.tmpl")
		if err != nil {
			return fmt.Errorf("failed to load template: %w", err)
		}
	}

	queriesTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/queries.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read queries template: %w", err)
	}

	testTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read test template: %w", err)
	}

	migrationTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/migration.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read migration template: %w", err)
	}

	schemaTmpl, err := kitLoader.LoadKitTemplate(kitName, "resource/schema.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to read schema template: %w", err)
	}

	// Generate handler
	if err := generateFile(string(handlerTmpl), data, filepath.Join(resourceDir, resourceNameLower+".go"), kit); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	// Generate template and validate it parses correctly
	tmplPath := filepath.Join(resourceDir, resourceNameLower+".tmpl")
	if err := generateFile(string(templateTmpl), data, tmplPath, kit); err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}
	if err := ValidateTemplate(tmplPath); err != nil {
		return err
	}

	// Generate migration file
	dbDir := filepath.Join(basePath, "database")
	migrationsDir := filepath.Join(dbDir, "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now()
	migrationFilename := ""
	migrationPath := ""
	for {
		timestampStr := timestamp.Format("20060102150405")
		migrationFilename = fmt.Sprintf("%s_create_%s.sql", timestampStr, tableName)
		migrationPath = filepath.Join(migrationsDir, migrationFilename)
		matches, _ := filepath.Glob(filepath.Join(migrationsDir, timestampStr+"_*.sql"))
		if len(matches) == 0 {
			break
		}
		timestamp = timestamp.Add(1 * time.Second)
	}
	if err := generateFile(string(migrationTmpl), data, migrationPath, kit); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	if err := appendToFile(string(schemaTmpl), data, filepath.Join(dbDir, "schema.sql"), "\n", kit); err != nil {
		return fmt.Errorf("failed to append to schema: %w", err)
	}

	if err := appendToFile(string(queriesTmpl), data, filepath.Join(dbDir, "queries.sql"), "\n", kit); err != nil {
		return fmt.Errorf("failed to append to queries: %w", err)
	}

	// Generate consolidated test file (E2E + WebSocket)
	if err := generateFile(string(testTmpl), data, filepath.Join(resourceDir, resourceNameLower+"_test.go"), kit); err != nil {
		return fmt.Errorf("failed to generate test: %w", err)
	}

	// Inject router registration into main.go
	// When file uploads are used, skip auto-injection because the handler
	// requires a storage.Store parameter that must be declared in main.go first.
	mainGoPath := findMainGo(basePath)
	if mainGoPath != "" && !data.Components.UseUpload {
		handlerCall := resourceNameLower + ".Handler(queries)"

		routes := []RouteInfo{
			{
				Path:        "/" + resourceNameLower,
				PackageName: resourceNameLower,
				HandlerCall: handlerCall,
				ImportPath:  moduleName + "/app/" + resourceNameLower,
			},
		}

		if editMode == "page" {
			routes = append(routes, RouteInfo{
				Path:        "/" + resourceNameLower + "/",
				PackageName: resourceNameLower,
				HandlerCall: handlerCall,
				ImportPath:  moduleName + "/app/" + resourceNameLower,
			})
		}

		for _, route := range routes {
			if err := InjectRoute(mainGoPath, route); err != nil {
				fmt.Printf("⚠️  Could not auto-inject route %s: %v\n", route.Path, err)
				fmt.Printf("   Please add manually: http.Handle(\"%s\", %s)\n",
					route.Path, handlerCall)
			}
		}
	}

	// Register resource for home page
	if err := RegisterResource(basePath, data.ResourceName, "/"+resourceNameLower, "resource"); err != nil {
		fmt.Printf("⚠️  Could not register resource in home page: %v\n", err)
	}

	return nil
}

func generateFile(tmplStr string, data interface{}, outPath string, kit *kits.KitInfo) error {
	// Merge base funcMap with kit helpers
	funcs := make(template.FuncMap)
	for k, v := range funcMap {
		funcs[k] = v
	}

	// Use kit helpers if provided, otherwise fallback to static CSSHelpers() for backward compatibility
	if kit != nil && kit.Helpers != nil {
		// Get kit-specific helpers using the CSSHelpers interface
		// Note: Old templates pass framework parameter (e.g., [[csscdn .CSSFramework]])
		// but kit helpers don't need it since they're already kit-specific
		// We accept the parameter but ignore it for backward compatibility
		funcs["csscdn"] = func(args ...interface{}) string { return kit.Helpers.CSSCDN() }
		funcs["containerClass"] = func(args ...interface{}) string { return kit.Helpers.ContainerClass() }
		funcs["sectionClass"] = func(args ...interface{}) string { return kit.Helpers.SectionClass() }
		funcs["boxClass"] = func(args ...interface{}) string { return kit.Helpers.BoxClass() }
		funcs["titleClass"] = func(args ...interface{}) string { return kit.Helpers.TitleClass(1) }
		funcs["subtitleClass"] = func(args ...interface{}) string { return kit.Helpers.SubtitleClass() }
		funcs["fieldClass"] = func(args ...interface{}) string { return kit.Helpers.FieldClass() }
		funcs["labelClass"] = func(args ...interface{}) string { return kit.Helpers.LabelClass() }
		funcs["controlClass"] = func(args ...interface{}) string { return "" } // Not in interface, return empty
		funcs["inputClass"] = func(args ...interface{}) string { return kit.Helpers.InputClass() }
		funcs["inputErrorClass"] = func(args ...interface{}) string { return kit.Helpers.InputClass() } // Fallback to inputClass
		funcs["selectClass"] = func(args ...interface{}) string { return kit.Helpers.SelectClass() }
		funcs["selectWrapperClass"] = func(args ...interface{}) string { return "" } // Not in interface
		funcs["checkboxClass"] = func(args ...interface{}) string { return kit.Helpers.CheckboxClass() }
		funcs["textareaClass"] = func(args ...interface{}) string { return kit.Helpers.TextareaClass() }
		funcs["buttonClass"] = func(framework string, args ...interface{}) string {
			// Template calls: buttonClass .CSSFramework "variant"
			// framework is .CSSFramework (unused), args[0] is the actual variant
			variant := "primary"
			if len(args) > 0 {
				if v, ok := args[0].(string); ok {
					variant = v
				}
			}
			return kit.Helpers.ButtonClass(variant)
		}
		funcs["buttonGroupClass"] = func(args ...interface{}) string { return kit.Helpers.ButtonGroupClass() }
		funcs["formClass"] = func(args ...interface{}) string { return kit.Helpers.FormClass() }
		funcs["tableClass"] = func(args ...interface{}) string { return kit.Helpers.TableClass() }
		funcs["tableContainerClass"] = func(args ...interface{}) string { return kit.Helpers.TableContainerClass() }
		funcs["theadClass"] = func(args ...interface{}) string { return kit.Helpers.TheadClass() }
		funcs["thClass"] = func(args ...interface{}) string { return kit.Helpers.ThClass() }
		funcs["tbodyClass"] = func(args ...interface{}) string { return kit.Helpers.TbodyClass() }
		funcs["trClass"] = func(args ...interface{}) string { return kit.Helpers.TrClass() }
		funcs["tdClass"] = func(args ...interface{}) string { return kit.Helpers.TdClass() }
		funcs["textClass"] = func(args ...interface{}) string { return kit.Helpers.TextClass("") }
		funcs["textMutedClass"] = func(args ...interface{}) string { return kit.Helpers.TextMutedClass() }
		funcs["textPrimaryClass"] = func(args ...interface{}) string { return kit.Helpers.TextPrimaryClass() }
		funcs["textDangerClass"] = func(args ...interface{}) string { return kit.Helpers.TextDangerClass() }
		funcs["paginationClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationClass() }
		funcs["paginationButtonClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("") }
		funcs["paginationActiveClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("active") }
		funcs["paginationCurrentClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("current") }
		funcs["paginationInfoClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationListClass() }
		funcs["helpTextClass"] = func(args ...interface{}) string { return kit.Helpers.TextMutedClass() } // Use textMuted as fallback
		funcs["errorClass"] = func(args ...interface{}) string { return kit.Helpers.TextDangerClass() }   // Use textDanger as fallback
		funcs["loadingClass"] = func(args ...interface{}) string { return kit.Helpers.LoadingClass() }
		funcs["codeClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["listClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["listItemClass"] = func(args ...interface{}) string { return "" } // Not in interface
		funcs["linkClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["needsWrapper"] = func(args ...interface{}) bool { return kit.Helpers.NeedsWrapper() }
		funcs["needsArticle"] = func(args ...interface{}) bool { return kit.Helpers.NeedsArticle() }
		funcs["needsTableWrapper"] = func(args ...interface{}) bool {
			// Use TableContainerClass as a proxy for whether table wrapper is needed
			return kit.Helpers.TableContainerClass() != ""
		}
		funcs["tableWrapperClass"] = func(args ...interface{}) string { return kit.Helpers.TableContainerClass() }
		// Add dict helper from kit
		funcs["dict"] = kit.Helpers.Dict
		funcs["until"] = kit.Helpers.Until
		funcs["add"] = kit.Helpers.Add
	} else {
		// Fallback to static CSS helpers for backward compatibility
		for k, v := range CSSHelpers() {
			funcs[k] = v
		}
	}

	// Use custom delimiters to avoid conflicts with Go template syntax in the generated files
	tmpl, err := template.New("template").Delims("[[", "]]").Funcs(funcs).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func appendToFile(tmplStr string, data interface{}, outPath, separator string, kit *kits.KitInfo) error {
	// Merge base funcMap with kit helpers
	funcs := make(template.FuncMap)
	for k, v := range funcMap {
		funcs[k] = v
	}

	// Use kit helpers if provided (same logic as generateFile)
	if kit != nil && kit.Helpers != nil {
		// Accept optional framework parameter for backward compatibility
		funcs["csscdn"] = func(args ...interface{}) string { return kit.Helpers.CSSCDN() }
		funcs["containerClass"] = func(args ...interface{}) string { return kit.Helpers.ContainerClass() }
		funcs["sectionClass"] = func(args ...interface{}) string { return kit.Helpers.SectionClass() }
		funcs["boxClass"] = func(args ...interface{}) string { return kit.Helpers.BoxClass() }
		funcs["titleClass"] = func(args ...interface{}) string { return kit.Helpers.TitleClass(1) }
		funcs["subtitleClass"] = func(args ...interface{}) string { return kit.Helpers.SubtitleClass() }
		funcs["fieldClass"] = func(args ...interface{}) string { return kit.Helpers.FieldClass() }
		funcs["labelClass"] = func(args ...interface{}) string { return kit.Helpers.LabelClass() }
		funcs["controlClass"] = func(args ...interface{}) string { return "" } // Not in interface, return empty
		funcs["inputClass"] = func(args ...interface{}) string { return kit.Helpers.InputClass() }
		funcs["inputErrorClass"] = func(args ...interface{}) string { return kit.Helpers.InputClass() } // Fallback to inputClass
		funcs["selectClass"] = func(args ...interface{}) string { return kit.Helpers.SelectClass() }
		funcs["selectWrapperClass"] = func(args ...interface{}) string { return "" } // Not in interface
		funcs["checkboxClass"] = func(args ...interface{}) string { return kit.Helpers.CheckboxClass() }
		funcs["textareaClass"] = func(args ...interface{}) string { return kit.Helpers.TextareaClass() }
		funcs["buttonClass"] = func(framework string, args ...interface{}) string {
			// Template calls: buttonClass .CSSFramework "variant"
			// framework is .CSSFramework (unused), args[0] is the actual variant
			variant := "primary"
			if len(args) > 0 {
				if v, ok := args[0].(string); ok {
					variant = v
				}
			}
			return kit.Helpers.ButtonClass(variant)
		}
		funcs["buttonGroupClass"] = func(args ...interface{}) string { return kit.Helpers.ButtonGroupClass() }
		funcs["formClass"] = func(args ...interface{}) string { return kit.Helpers.FormClass() }
		funcs["tableClass"] = func(args ...interface{}) string { return kit.Helpers.TableClass() }
		funcs["tableContainerClass"] = func(args ...interface{}) string { return kit.Helpers.TableContainerClass() }
		funcs["theadClass"] = func(args ...interface{}) string { return kit.Helpers.TheadClass() }
		funcs["thClass"] = func(args ...interface{}) string { return kit.Helpers.ThClass() }
		funcs["tbodyClass"] = func(args ...interface{}) string { return kit.Helpers.TbodyClass() }
		funcs["trClass"] = func(args ...interface{}) string { return kit.Helpers.TrClass() }
		funcs["tdClass"] = func(args ...interface{}) string { return kit.Helpers.TdClass() }
		funcs["textClass"] = func(args ...interface{}) string { return kit.Helpers.TextClass("") }
		funcs["textMutedClass"] = func(args ...interface{}) string { return kit.Helpers.TextMutedClass() }
		funcs["textPrimaryClass"] = func(args ...interface{}) string { return kit.Helpers.TextPrimaryClass() }
		funcs["textDangerClass"] = func(args ...interface{}) string { return kit.Helpers.TextDangerClass() }
		funcs["paginationClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationClass() }
		funcs["paginationButtonClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("") }
		funcs["paginationActiveClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("active") }
		funcs["paginationCurrentClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationButtonClass("current") }
		funcs["paginationInfoClass"] = func(args ...interface{}) string { return kit.Helpers.PaginationListClass() }
		funcs["helpTextClass"] = func(args ...interface{}) string { return kit.Helpers.TextMutedClass() } // Use textMuted as fallback
		funcs["errorClass"] = func(args ...interface{}) string { return kit.Helpers.TextDangerClass() }   // Use textDanger as fallback
		funcs["loadingClass"] = func(args ...interface{}) string { return kit.Helpers.LoadingClass() }
		funcs["codeClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["listClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["listItemClass"] = func(args ...interface{}) string { return "" } // Not in interface
		funcs["linkClass"] = func(args ...interface{}) string { return "" }     // Not in interface
		funcs["needsWrapper"] = func(args ...interface{}) bool { return kit.Helpers.NeedsWrapper() }
		funcs["needsArticle"] = func(args ...interface{}) bool { return kit.Helpers.NeedsArticle() }
		funcs["needsTableWrapper"] = func(args ...interface{}) bool {
			return kit.Helpers.TableContainerClass() != ""
		}
		funcs["tableWrapperClass"] = func(args ...interface{}) string { return kit.Helpers.TableContainerClass() }
		funcs["dict"] = kit.Helpers.Dict
		funcs["until"] = kit.Helpers.Until
		funcs["add"] = kit.Helpers.Add
	} else {
		for k, v := range CSSHelpers() {
			funcs[k] = v
		}
	}

	// Use custom delimiters to avoid conflicts with Go template syntax in the generated files
	tmpl, err := template.New("template").Delims("[[", "]]").Funcs(funcs).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Open file for appending (create if doesn't exist)
	f, err := os.OpenFile(outPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Write separator and content
	if _, err := f.WriteString(separator); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	return nil
}

// findMainGo finds the main.go file in cmd/* directory
func findMainGo(basePath string) string {
	// Try to find cmd/*/main.go
	cmdDir := filepath.Join(basePath, "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			mainGoPath := filepath.Join(cmdDir, entry.Name(), "main.go")
			if _, err := os.Stat(mainGoPath); err == nil {
				return mainGoPath
			}
		}
	}

	return ""
}

// Singularize is the exported version of singularize for use by other packages.
func Singularize(word string) string {
	return singularize(word)
}

// singularize handles basic English singularization
func singularize(word string) string {
	// Common irregular plurals (reverse map)
	irregulars := map[string]string{
		"people":   "person",
		"children": "child",
		"teeth":    "tooth",
		"feet":     "foot",
		"men":      "man",
		"women":    "woman",
		"mice":     "mouse",
	}
	if singular, ok := irregulars[word]; ok {
		return singular
	}

	// Words ending in ies -> y (e.g., categories -> category)
	if strings.HasSuffix(word, "ies") && len(word) > 3 {
		return word[:len(word)-3] + "y"
	}

	// Words ending in ses, xes, zes -> remove es (e.g., boxes -> box)
	if strings.HasSuffix(word, "ses") || strings.HasSuffix(word, "xes") || strings.HasSuffix(word, "zes") {
		return word[:len(word)-2]
	}

	// Words ending in ches, shes -> remove es (e.g., watches -> watch)
	if strings.HasSuffix(word, "ches") || strings.HasSuffix(word, "shes") {
		return word[:len(word)-2]
	}

	// Words ending in s -> remove s (e.g., users -> user)
	if strings.HasSuffix(word, "s") && len(word) > 1 {
		return word[:len(word)-1]
	}

	// Already singular
	return word
}

// pluralize handles basic English pluralization rules
func pluralize(word string) string {
	// If already ends in 's', return as-is
	if strings.HasSuffix(word, "s") {
		return word
	}

	// Common irregular plurals
	irregulars := map[string]string{
		"person": "people",
		"child":  "children",
		"tooth":  "teeth",
		"foot":   "feet",
		"man":    "men",
		"woman":  "women",
		"mouse":  "mice",
	}
	if plural, ok := irregulars[word]; ok {
		return plural
	}

	// Words ending in consonant + y -> ies
	if len(word) >= 2 && word[len(word)-1] == 'y' {
		preceding := word[len(word)-2]
		if preceding != 'a' && preceding != 'e' && preceding != 'i' && preceding != 'o' && preceding != 'u' {
			return word[:len(word)-1] + "ies"
		}
	}

	// Words ending in s, x, z, ch, sh -> es
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") ||
		strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}

	// Default: just add s
	return word + "s"
}
