package generator

import (
	"bytes"
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

// JobsConfig holds configuration for generating the queue infrastructure.
type JobsConfig struct {
	ModuleName string
}

// JobConfig holds configuration for generating a single job handler.
type JobConfig struct {
	ModuleName   string
	JobName      string // snake_case, e.g. "send_email"
	JobNameCamel string // CamelCase, e.g. "SendEmail"
}

// GenerateQueue sets up the background job infrastructure using River.
// It creates the migration, schema, worker init file, and injects setup into main.go.
func GenerateQueue(projectRoot string, moduleName string) error {
	projectConfig, err := config.LoadProjectConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	kitName := projectConfig.GetKit()
	kitLoader := kits.DefaultLoader()

	// Check if queue already set up
	workerPath := filepath.Join(projectRoot, "app", "jobs", "worker.go")
	if _, err := os.Stat(workerPath); err == nil {
		return fmt.Errorf("job queue already set up (app/jobs/worker.go exists)")
	}

	// 1. Create migration file
	migrationsDir := filepath.Join(projectRoot, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now()
	var migrationPath string
	for i := 0; i < 3600; i++ {
		timestampStr := timestamp.Format("20060102150405")
		migrationFile := fmt.Sprintf("%s_setup_river_queue.sql", timestampStr)
		migrationPath = filepath.Join(migrationsDir, migrationFile)

		matches, err := filepath.Glob(filepath.Join(migrationsDir, timestampStr+"_*"))
		if err != nil {
			return fmt.Errorf("failed to check for existing migrations: %w", err)
		}
		if len(matches) == 0 {
			break
		}
		timestamp = timestamp.Add(1 * time.Second)
		if i == 3599 {
			return fmt.Errorf("failed to generate unique migration timestamp")
		}
	}

	if err := writeTemplateFile(kitLoader, kitName, "jobs/migration.sql.tmpl", migrationPath, nil); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	// 2. Append to schema.sql
	schemaPath := filepath.Join(projectRoot, "database", "schema.sql")
	if err := appendTemplateFile(kitLoader, kitName, "jobs/schema.sql.tmpl", schemaPath, nil); err != nil {
		return fmt.Errorf("failed to append to schema.sql: %w", err)
	}

	// 3. Create app/jobs/worker.go
	jobsDir := filepath.Join(projectRoot, "app", "jobs")
	if err := os.MkdirAll(jobsDir, 0755); err != nil {
		return fmt.Errorf("failed to create app/jobs directory: %w", err)
	}

	if err := writeTemplateFile(kitLoader, kitName, "jobs/worker_init.go.tmpl", workerPath, nil); err != nil {
		return fmt.Errorf("failed to generate worker.go: %w", err)
	}

	// 4. Add River dependencies
	goModPath := filepath.Join(projectRoot, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		dependencies := []string{
			"github.com/riverqueue/river@latest",
			"github.com/riverqueue/river/riverdriver/riversqlite@latest",
		}
		args := append([]string{"get"}, dependencies...)
		cmd := exec.Command("go", args...)
		cmd.Dir = projectRoot
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch River dependencies (run 'go mod tidy' in %s to resolve):\n%s\n", projectRoot, string(output))
		}
	}

	// 5. Inject River client setup into main.go
	mainGoPath := findMainGo(projectRoot)
	if mainGoPath != "" {
		if err := injectJobWorker(mainGoPath, moduleName); err != nil {
			return fmt.Errorf("failed to inject job worker into main.go: %w", err)
		}
	}

	return nil
}

// GenerateJob scaffolds a new job handler and registers it with the worker.
func GenerateJob(projectRoot string, moduleName string, jobName string) error {
	// Validate queue is set up
	workerPath := filepath.Join(projectRoot, "app", "jobs", "worker.go")
	if _, err := os.Stat(workerPath); os.IsNotExist(err) {
		return fmt.Errorf("job queue not set up yet. Run 'lvt gen queue' first")
	}

	projectConfig, err := config.LoadProjectConfig(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}
	kitName := projectConfig.GetKit()
	kitLoader := kits.DefaultLoader()

	jobNameCamel := toCamelCase(jobName)

	jobConfig := &JobConfig{
		ModuleName:   moduleName,
		JobName:      jobName,
		JobNameCamel: jobNameCamel,
	}

	// Check if job already exists
	jobPath := filepath.Join(projectRoot, "app", "jobs", jobName+".go")
	if _, err := os.Stat(jobPath); err == nil {
		return fmt.Errorf("job '%s' already exists (app/jobs/%s.go)", jobName, jobName)
	}

	// 1. Generate job handler file
	templateContent, err := kitLoader.LoadKitTemplate(kitName, "jobs/handler.go.tmpl")
	if err != nil {
		return fmt.Errorf("failed to load handler template: %w", err)
	}

	tmpl, err := template.New("job_handler").Delims("<<", ">>").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse handler template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, jobConfig); err != nil {
		return fmt.Errorf("failed to execute handler template: %w", err)
	}
	if err := os.WriteFile(jobPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s.go: %w", jobName, err)
	}

	// 2. Register worker in app/jobs/worker.go
	if err := injectWorkerRegistration(workerPath, jobNameCamel); err != nil {
		return fmt.Errorf("failed to register worker: %w", err)
	}

	return nil
}

// writeTemplateFile loads a kit template, executes it, and writes the result atomically.
func writeTemplateFile(kitLoader *kits.KitLoader, kitName, templatePath, outputPath string, data interface{}) error {
	content, err := kitLoader.LoadKitTemplate(kitName, templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

// appendTemplateFile loads a kit template and appends it to the output file.
func appendTemplateFile(kitLoader *kits.KitLoader, kitName, templatePath, outputPath string, data interface{}) error {
	content, err := kitLoader.LoadKitTemplate(kitName, templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() > 0 {
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}

	return tmpl.Execute(file, data)
}

// injectJobWorker injects River client setup into main.go.
func injectJobWorker(mainGoPath string, moduleName string) error {
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	mainStr := string(content)

	// Check if already injected
	if strings.Contains(mainStr, "river.NewClient") {
		return nil // Already injected
	}

	// Find injection point: after appCtx creation (needed by River client)
	lines := strings.Split(mainStr, "\n")
	var result []string
	injected := false

	for _, line := range lines {
		result = append(result, line)

		// Inject after appCtx creation — River needs the context for Start/Stop
		if !injected && strings.Contains(line, "appCtx, appCancel := context.WithCancel") {
			// Find the next line (defer appCancel()) and include it
			continue
		}
		if !injected && strings.Contains(line, "defer appCancel()") {
			riverSetup := []string{
				"",
				"\t// Background job processing (River)",
				"\triverDB, err := sql.Open(\"sqlite\", dbPath+\"?_pragma=journal_mode(WAL)\")",
				"\tif err != nil {",
				"\t\tslog.Error(\"Failed to open River database\", \"error\", err)",
				"\t\tos.Exit(1)",
				"\t}",
				"\triverDB.SetMaxOpenConns(1)",
				"\tdefer riverDB.Close()",
				"",
				"\tjobWorkers := jobs.SetupWorkers()",
				"\triverClient, err := river.NewClient(riversqlite.New(riverDB), &river.Config{",
				"\t\tQueues: map[string]river.QueueConfig{",
				"\t\t\triver.QueueDefault: {MaxWorkers: 100},",
				"\t\t},",
				"\t\tWorkers: jobWorkers,",
				"\t})",
				"\tif err != nil {",
				"\t\tslog.Error(\"Failed to create River client\", \"error\", err)",
				"\t\tos.Exit(1)",
				"\t}",
				"\tif err := riverClient.Start(appCtx); err != nil {",
				"\t\tslog.Error(\"Failed to start job workers\", \"error\", err)",
				"\t\tos.Exit(1)",
				"\t}",
				"\tdefer func() {",
				"\t\tstopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)",
				"\t\tdefer stopCancel()",
				"\t\t_ = riverClient.Stop(stopCtx)",
				"\t}()",
				"\tjobs.SetClient(riverClient)",
				"\tslog.Info(\"Background job worker started\")",
			}
			result = append(result, riverSetup...)
			injected = true
		}
	}

	if !injected {
		return fmt.Errorf("could not find injection point in main.go (expected 'appCtx, appCancel := context.WithCancel')")
	}

	// Inject imports using existing helper
	resultStr := strings.Join(result, "\n")

	imports := []string{
		"\t\"database/sql\"",
		"\t\"time\"",
		fmt.Sprintf("\t\"%s/app/jobs\"", moduleName),
		"\t\"github.com/riverqueue/river\"",
		"\t\"github.com/riverqueue/river/riverdriver/riversqlite\"",
	}
	for _, imp := range imports {
		resultStr, err = injectImport(resultStr, imp)
		if err != nil {
			return fmt.Errorf("failed to inject import %s: %w", imp, err)
		}
	}

	return os.WriteFile(mainGoPath, []byte(resultStr), 0644)
}

// injectWorkerRegistration adds a river.AddWorker call to app/jobs/worker.go.
func injectWorkerRegistration(workerPath string, jobNameCamel string) error {
	content, err := os.ReadFile(workerPath)
	if err != nil {
		return fmt.Errorf("failed to read worker.go: %w", err)
	}

	workerStr := string(content)

	// Check if already registered
	registrationLine := fmt.Sprintf("river.AddWorker(workers, &%sWorker{})", jobNameCamel)
	if strings.Contains(workerStr, registrationLine) {
		return nil // Already registered
	}

	// Find the marker comment and insert after it
	marker := "// Register job workers below (added by `lvt gen job`)"
	idx := strings.Index(workerStr, marker)
	if idx == -1 {
		return fmt.Errorf("could not find registration marker in worker.go")
	}

	insertPos := idx + len(marker)
	registration := "\n\t" + registrationLine

	workerStr = workerStr[:insertPos] + registration + workerStr[insertPos:]

	return os.WriteFile(workerPath, []byte(workerStr), 0644)
}
