package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/kits"
)

type AuthConfig struct {
	ModuleName          string
	StructName          string // e.g., "User", "Account", "Admin"
	TableName           string // e.g., "users", "accounts", "admin_users"
	EnablePassword      bool
	EnableMagicLink     bool
	EnableEmailConfirm  bool
	EnablePasswordReset bool
	EnableSessionsUI    bool
	EnableCSRF          bool
}

func GenerateAuth(projectRoot string, authConfig *AuthConfig) error {
	// Apply defaults if not set
	if authConfig.TableName == "" {
		authConfig.TableName = "users"
	}
	if authConfig.StructName == "" {
		authConfig.StructName = "User"
	}

	// Load project config to get the kit
	projectConfig, err := config.LoadProjectConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	kitName := projectConfig.GetKit()

	// Load kit loader
	kitLoader := kits.DefaultLoader()

	// Note: password and email utilities are imported from github.com/livetemplate/lvt/pkg
	// No need to generate shared/ directory anymore

	// Generate migration
	migrationsDir := filepath.Join(projectRoot, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	migrationFile := fmt.Sprintf("%s_create_auth_tables.sql", timestamp)
	migrationPath := filepath.Join(migrationsDir, migrationFile)

	templateContent, err := kitLoader.LoadKitTemplate(kitName, "auth/migration.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load migration template: %w", err)
	}

	funcMap := template.FuncMap{
		"singular": singularize,
	}

	tmpl, err := template.New("migration").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse migration template: %w", err)
	}

	file, err := os.Create(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, authConfig); err != nil {
		return fmt.Errorf("failed to execute migration template: %w", err)
	}

	// Append to queries.sql (or create if doesn't exist)
	queriesPath := filepath.Join(projectRoot, "database", "queries.sql")

	templateContent, err = kitLoader.LoadKitTemplate(kitName, "auth/queries.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load queries template: %w", err)
	}

	tmpl, err = template.New("queries").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse queries template: %w", err)
	}

	// Open in append mode (create if doesn't exist)
	file, err = os.OpenFile(queriesPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open queries.sql: %w", err)
	}

	// Add separator if file already has content
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to stat queries.sql: %w", err)
	}
	if stat.Size() > 0 {
		if _, err := file.WriteString("\n\n"); err != nil {
			file.Close()
			return fmt.Errorf("failed to write separator: %w", err)
		}
	}

	if err := tmpl.Execute(file, authConfig); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute queries template: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close queries.sql: %w", err)
	}

	// Append to schema.sql for sqlc (separate from migration)
	schemaPath := filepath.Join(projectRoot, "database", "schema.sql")
	schemaTemplateContent, err := kitLoader.LoadKitTemplate(kitName, "auth/schema.sql.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load schema template: %w", err)
	}

	schemaTmpl, err := template.New("schema").Delims("[[", "]]").Funcs(funcMap).Parse(string(schemaTemplateContent))
	if err != nil {
		return fmt.Errorf("failed to parse schema template: %w", err)
	}

	// Open in append mode (create if doesn't exist)
	schemaFile, err := os.OpenFile(schemaPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open schema.sql: %w", err)
	}

	// Add separator if file already has content
	schemaStat, err := schemaFile.Stat()
	if err != nil {
		schemaFile.Close()
		return fmt.Errorf("failed to stat schema.sql: %w", err)
	}
	if schemaStat.Size() > 0 {
		if _, err := schemaFile.WriteString("\n"); err != nil {
			schemaFile.Close()
			return fmt.Errorf("failed to write separator: %w", err)
		}
	}

	if err := schemaTmpl.Execute(schemaFile, authConfig); err != nil {
		schemaFile.Close()
		return fmt.Errorf("failed to execute schema template: %w", err)
	}

	if err := schemaFile.Close(); err != nil {
		return fmt.Errorf("failed to close schema.sql: %w", err)
	}

	// Generate auth handler
	authHandlerDir := filepath.Join(projectRoot, "app", "auth")
	if err := os.MkdirAll(authHandlerDir, 0755); err != nil {
		return fmt.Errorf("failed to create auth handler directory: %w", err)
	}

	// Generate handler.go
	templateContent, err = kitLoader.LoadKitTemplate(kitName, "auth/handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load handler template: %w", err)
	}

	outputPath := filepath.Join(authHandlerDir, "auth.go")
	tmpl, err = template.New("handler").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse handler template: %w", err)
	}

	file, err = os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create auth.go: %w", err)
	}

	if err := tmpl.Execute(file, authConfig); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute handler template: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close auth.go: %w", err)
	}

	// Generate template file
	templateContent, err = kitLoader.LoadKitTemplate(kitName, "auth/template.tmpl.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load template template: %w", err)
	}

	outputPath = filepath.Join(authHandlerDir, "auth.tmpl")
	tmpl, err = template.New("template").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template template: %w", err)
	}

	file, err = os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create auth.tmpl: %w", err)
	}

	if err := tmpl.Execute(file, authConfig); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute template template: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close auth.tmpl: %w", err)
	}

	// Generate middleware file
	templateContent, err = kitLoader.LoadKitTemplate(kitName, "auth/middleware.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load middleware template: %w", err)
	}

	outputPath = filepath.Join(authHandlerDir, "middleware.go")
	tmpl, err = template.New("middleware").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse middleware template: %w", err)
	}

	file, err = os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create middleware.go: %w", err)
	}

	if err := tmpl.Execute(file, authConfig); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute middleware template: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close middleware.go: %w", err)
	}

	// Generate E2E test file
	templateContent, err = kitLoader.LoadKitTemplate(kitName, "auth/e2e_test.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load e2e test template: %w", err)
	}

	outputPath = filepath.Join(authHandlerDir, "auth_e2e_test.go")
	tmpl, err = template.New("e2e_test").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse e2e test template: %w", err)
	}

	file, err = os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create auth_e2e_test.go: %w", err)
	}

	if err := tmpl.Execute(file, authConfig); err != nil {
		file.Close()
		return fmt.Errorf("failed to execute e2e test template: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close auth_e2e_test.go: %w", err)
	}

	// Update go.mod dependencies if go.mod exists
	goModPath := filepath.Join(projectRoot, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		dependencies := []string{
			"github.com/google/uuid@latest",
			"github.com/chromedp/chromedp@latest", // For E2E tests
			"github.com/livetemplate/lvt@latest",  // Auth utilities (password, email)
		}
		if authConfig.EnablePassword {
			dependencies = append(dependencies, "golang.org/x/crypto@latest")
		}
		if authConfig.EnableCSRF {
			dependencies = append(dependencies, "github.com/gorilla/csrf@latest")
		}

		if len(dependencies) > 0 {
			args := append([]string{"get"}, dependencies...)
			cmd := exec.Command("go", args...)
			cmd.Dir = projectRoot
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to update dependencies: %w\n%s", err, output)
			}
		}
	}

	// Inject auth routes into main.go
	mainGoPath := findMainGo(projectRoot)
	if mainGoPath != "" {
		// Main auth route (LiveTemplate handler)
		routes := []RouteInfo{
			{
				Path:        "/auth",
				PackageName: "auth",
				HandlerCall: "auth.Handler(queries)",
				ImportPath:  authConfig.ModuleName + "/app/auth",
			},
			// Logout route
			{
				Path:        "/auth/logout",
				PackageName: "auth",
				HandlerCall: "auth.LogoutHandler(queries)",
				ImportPath:  authConfig.ModuleName + "/app/auth",
			},
		}

		// Add magic link route if enabled
		if authConfig.EnableMagicLink {
			routes = append(routes, RouteInfo{
				Path:        "/auth/magic",
				PackageName: "auth",
				HandlerCall: "auth.MagicLinkHandler(queries)",
				ImportPath:  authConfig.ModuleName + "/app/auth",
			})
		}

		// Add password reset route if enabled
		if authConfig.EnablePasswordReset {
			routes = append(routes, RouteInfo{
				Path:        "/auth/reset",
				PackageName: "auth",
				HandlerCall: "auth.ResetPasswordHandler(queries)",
				ImportPath:  authConfig.ModuleName + "/app/auth",
			})
		}

		// Add email confirmation route if enabled
		if authConfig.EnableEmailConfirm {
			routes = append(routes, RouteInfo{
				Path:        "/auth/confirm",
				PackageName: "auth",
				HandlerCall: "auth.ConfirmEmailHandler(queries)",
				ImportPath:  authConfig.ModuleName + "/app/auth",
			})
		}

		for _, route := range routes {
			if err := InjectRoute(mainGoPath, route); err != nil {
				// Log warning but don't fail - user can add route manually
				fmt.Printf("⚠️  Could not auto-inject route %s: %v\n", route.Path, err)
				fmt.Printf("   Please add manually: http.Handle(\"%s\", auth.Handler(queries))\n", route.Path)
			}
		}

		// Register auth in home page
		if err := RegisterResource(projectRoot, "Auth", "/auth", "auth"); err != nil {
			fmt.Printf("⚠️  Could not register auth in home page: %v\n", err)
		}

		// Update home page to show login/logout buttons
		if err := updateHomeForAuth(projectRoot, authConfig); err != nil {
			fmt.Printf("⚠️  Could not update home page for auth: %v\n", err)
			fmt.Println("   You may need to manually add login/logout buttons to your home page")
		}
	}

	return nil
}

// updateHomeForAuth modifies the home page handler and template to show login/logout buttons
func updateHomeForAuth(projectRoot string, authConfig *AuthConfig) error {
	// Update home.go to accept queries and check auth state
	if err := updateHomeHandler(projectRoot, authConfig); err != nil {
		return fmt.Errorf("failed to update home handler: %w", err)
	}

	// Update home.tmpl to show login/logout buttons
	if err := updateHomeTemplate(projectRoot, authConfig); err != nil {
		return fmt.Errorf("failed to update home template: %w", err)
	}

	// Update main.go to pass queries to home.Handler
	if err := updateMainGoHomeHandler(projectRoot); err != nil {
		return fmt.Errorf("failed to update main.go home handler: %w", err)
	}

	return nil
}

// updateMainGoHomeHandler updates the home.Handler() call in main.go to pass queries
func updateMainGoHomeHandler(projectRoot string) error {
	mainGoPath := findMainGo(projectRoot)
	if mainGoPath == "" {
		return fmt.Errorf("could not find main.go")
	}

	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	mainContent := string(content)

	// Check if already updated
	if strings.Contains(mainContent, "home.Handler(queries)") {
		return nil // Already updated
	}

	// Update home.Handler() to home.Handler(queries)
	mainContent = strings.Replace(mainContent, "home.Handler()", "home.Handler(queries)", 1)

	if err := os.WriteFile(mainGoPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	return nil
}

// updateHomeHandler modifies home.go to check auth state and pass it to the template
func updateHomeHandler(projectRoot string, authConfig *AuthConfig) error {
	homeGoPath := filepath.Join(projectRoot, "app", "home", "home.go")

	content, err := os.ReadFile(homeGoPath)
	if err != nil {
		return fmt.Errorf("failed to read home.go: %w", err)
	}

	homeContent := string(content)

	// Check if already updated
	if strings.Contains(homeContent, "IsLoggedIn") {
		return nil // Already updated
	}

	// Add auth import
	authImport := fmt.Sprintf("\t\"%s/app/auth\"", authConfig.ModuleName)
	if !strings.Contains(homeContent, authImport) {
		// Find the import block and add auth import
		importEnd := strings.Index(homeContent, "\n)")
		if importEnd == -1 {
			return fmt.Errorf("could not find import block")
		}
		homeContent = homeContent[:importEnd] + "\n" + authImport + homeContent[importEnd:]
	}

	// Add models import if not present
	modelsImport := fmt.Sprintf("\t\"%s/database/models\"", authConfig.ModuleName)
	if !strings.Contains(homeContent, modelsImport) {
		importEnd := strings.Index(homeContent, "\n)")
		if importEnd == -1 {
			return fmt.Errorf("could not find import block")
		}
		homeContent = homeContent[:importEnd] + "\n" + modelsImport + homeContent[importEnd:]
	}

	// Add IsLoggedIn and UserEmail fields to HomeState
	stateFields := `	IsLoggedIn   bool       ` + "`json:\"is_logged_in\"`" + `
	UserEmail    string     ` + "`json:\"user_email\"`"

	// Find the LastUpdated field and add after it
	lastUpdatedIdx := strings.Index(homeContent, "LastUpdated  string")
	if lastUpdatedIdx != -1 {
		// Find the end of that line (after the json tag)
		lineEnd := strings.Index(homeContent[lastUpdatedIdx:], "\n")
		if lineEnd != -1 {
			insertPos := lastUpdatedIdx + lineEnd
			homeContent = homeContent[:insertPos] + "\n" + stateFields + homeContent[insertPos:]
		}
	}

	// Update Handler function signature to accept queries
	homeContent = strings.Replace(homeContent,
		"func Handler() http.Handler {",
		"func Handler(queries *models.Queries) http.Handler {",
		1)

	// Add auth controller and wrap the handler to check auth state
	// Find the return statement and replace the handler logic
	oldReturn := "return tmpl.Handle(controller, livetemplate.AsState(initialState))"
	if strings.Contains(homeContent, oldReturn) {
		newHandler := `// Create auth controller to check login state
	authController := auth.NewUserController(queries, nil, "")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clone state for this request
		state := *initialState

		// Check if user is logged in
		if user, err := authController.GetCurrentUser(r); err == nil && user != nil {
			state.IsLoggedIn = true
			state.UserEmail = user.Email
		}

		tmpl.Handle(controller, livetemplate.AsState(&state)).ServeHTTP(w, r)
	})`
		homeContent = strings.Replace(homeContent, oldReturn, newHandler, 1)
	}

	if err := os.WriteFile(homeGoPath, []byte(homeContent), 0644); err != nil {
		return fmt.Errorf("failed to write home.go: %w", err)
	}

	return nil
}

// updateHomeTemplate modifies home.tmpl to show login/logout buttons
func updateHomeTemplate(projectRoot string, authConfig *AuthConfig) error {
	homeTmplPath := filepath.Join(projectRoot, "app", "home", "home.tmpl")

	content, err := os.ReadFile(homeTmplPath)
	if err != nil {
		return fmt.Errorf("failed to read home.tmpl: %w", err)
	}

	tmplContent := string(content)

	// Check if already updated
	if strings.Contains(tmplContent, "IsLoggedIn") {
		return nil // Already updated
	}

	// Find {{define "content"}} and add auth buttons after it
	contentDefine := `{{define "content"}}`
	contentIdx := strings.Index(tmplContent, contentDefine)
	if contentIdx == -1 {
		return fmt.Errorf("could not find content template definition")
	}

	authButtons := `
  <!-- Auth buttons -->
  <div style="display: flex; justify-content: flex-end; gap: 1rem; align-items: center; margin-bottom: 1rem;">
    {{if .IsLoggedIn}}
      <a href="/dashboard" style="padding: 0.5rem 1rem; background: #059669; color: white; border-radius: 4px; text-decoration: none;">Dashboard</a>
      <span style="color: #666;">{{.UserEmail}}</span>
      <a href="/auth/logout" style="padding: 0.5rem 1rem; background: #dc2626; color: white; border-radius: 4px; text-decoration: none;">Logout</a>
    {{else}}
      <a href="/dashboard" style="padding: 0.5rem 1rem; background: #6b7280; color: white; border-radius: 4px; text-decoration: none;">Dashboard (protected)</a>
      <a href="/auth" style="padding: 0.5rem 1rem; background: #4f46e5; color: white; border-radius: 4px; text-decoration: none;">Login</a>
    {{end}}
  </div>
`

	insertPos := contentIdx + len(contentDefine)
	tmplContent = tmplContent[:insertPos] + authButtons + tmplContent[insertPos:]

	if err := os.WriteFile(homeTmplPath, []byte(tmplContent), 0644); err != nil {
		return fmt.Errorf("failed to write home.tmpl: %w", err)
	}

	return nil
}
