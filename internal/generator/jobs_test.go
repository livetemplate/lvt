package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateQueue(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up minimal project structure
	setupTestProject(t, tmpDir)

	err := GenerateQueue(tmpDir, "testmodule")
	if err != nil {
		t.Fatalf("GenerateQueue failed: %v", err)
	}

	// Verify worker.go was created
	workerPath := filepath.Join(tmpDir, "app", "jobs", "worker.go")
	if _, err := os.Stat(workerPath); os.IsNotExist(err) {
		t.Error("app/jobs/worker.go was not created")
	}

	workerContent, err := os.ReadFile(workerPath)
	if err != nil {
		t.Fatalf("Failed to read worker.go: %v", err)
	}
	if !strings.Contains(string(workerContent), "river.NewWorkers()") {
		t.Error("worker.go missing river.NewWorkers() call")
	}
	if !strings.Contains(string(workerContent), "Register job workers below") {
		t.Error("worker.go missing registration marker comment")
	}

	// Verify migration was created
	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations dir: %v", err)
	}

	var migrationFound bool
	for _, entry := range entries {
		if strings.Contains(entry.Name(), "setup_river_queue") {
			migrationFound = true
			// Read and verify content
			content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read migration: %v", err)
			}
			if !strings.Contains(string(content), "river_job") {
				t.Error("Migration missing river_job table")
			}
			if !strings.Contains(string(content), "river_leader") {
				t.Error("Migration missing river_leader table")
			}
			if !strings.Contains(string(content), "+goose Up") {
				t.Error("Migration missing goose Up marker")
			}
		}
	}
	if !migrationFound {
		t.Error("River migration file was not created")
	}

	// Verify schema.sql was appended
	schemaPath := filepath.Join(tmpDir, "database", "schema.sql")
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema.sql: %v", err)
	}
	if !strings.Contains(string(schemaContent), "river_job") {
		t.Error("schema.sql missing river_job table")
	}
}

func TestGenerateQueueIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir)

	// First call should succeed
	if err := GenerateQueue(tmpDir, "testmodule"); err != nil {
		t.Fatalf("First GenerateQueue failed: %v", err)
	}

	// Second call should fail (already set up)
	err := GenerateQueue(tmpDir, "testmodule")
	if err == nil {
		t.Error("Expected error on second GenerateQueue call")
	}
	if err != nil && !strings.Contains(err.Error(), "already set up") {
		t.Errorf("Expected 'already set up' error, got: %v", err)
	}
}

func TestGenerateJob(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir)

	// Set up queue first
	if err := GenerateQueue(tmpDir, "testmodule"); err != nil {
		t.Fatalf("GenerateQueue failed: %v", err)
	}

	// Generate a job
	if err := GenerateJob(tmpDir, "testmodule", "send_email"); err != nil {
		t.Fatalf("GenerateJob failed: %v", err)
	}

	// Verify job handler was created
	jobPath := filepath.Join(tmpDir, "app", "jobs", "send_email.go")
	if _, err := os.Stat(jobPath); os.IsNotExist(err) {
		t.Error("app/jobs/send_email.go was not created")
	}

	jobContent, err := os.ReadFile(jobPath)
	if err != nil {
		t.Fatalf("Failed to read send_email.go: %v", err)
	}

	// Check for expected content
	checks := []string{
		"SendEmailArgs",
		"SendEmailWorker",
		`Kind() string { return "send_email" }`,
		"river.WorkerDefaults[SendEmailArgs]",
		"func (w *SendEmailWorker) Work(",
	}
	for _, check := range checks {
		if !strings.Contains(string(jobContent), check) {
			t.Errorf("send_email.go missing expected content: %s", check)
		}
	}

	// Verify worker registration was injected
	workerPath := filepath.Join(tmpDir, "app", "jobs", "worker.go")
	workerContent, err := os.ReadFile(workerPath)
	if err != nil {
		t.Fatalf("Failed to read worker.go: %v", err)
	}
	if !strings.Contains(string(workerContent), "river.AddWorker(workers, &SendEmailWorker{})") {
		t.Error("worker.go missing SendEmailWorker registration")
	}
}

func TestGenerateJobWithoutQueue(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir)

	// Try to generate job without queue setup
	err := GenerateJob(tmpDir, "testmodule", "send_email")
	if err == nil {
		t.Error("Expected error when generating job without queue")
	}
	if err != nil && !strings.Contains(err.Error(), "Run 'lvt gen queue' first") {
		t.Errorf("Expected 'Run lvt gen queue first' error, got: %v", err)
	}
}

func TestGenerateJobDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir)

	if err := GenerateQueue(tmpDir, "testmodule"); err != nil {
		t.Fatalf("GenerateQueue failed: %v", err)
	}

	// First job should succeed
	if err := GenerateJob(tmpDir, "testmodule", "send_email"); err != nil {
		t.Fatalf("First GenerateJob failed: %v", err)
	}

	// Duplicate should fail
	err := GenerateJob(tmpDir, "testmodule", "send_email")
	if err == nil {
		t.Error("Expected error on duplicate job")
	}
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestGenerateMultipleJobs(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestProject(t, tmpDir)

	if err := GenerateQueue(tmpDir, "testmodule"); err != nil {
		t.Fatalf("GenerateQueue failed: %v", err)
	}

	jobs := []string{"send_email", "process_payment", "generate_report"}
	for _, job := range jobs {
		if err := GenerateJob(tmpDir, "testmodule", job); err != nil {
			t.Fatalf("GenerateJob(%s) failed: %v", job, err)
		}
	}

	// Verify all jobs registered in worker.go
	workerContent, err := os.ReadFile(filepath.Join(tmpDir, "app", "jobs", "worker.go"))
	if err != nil {
		t.Fatalf("Failed to read worker.go: %v", err)
	}

	expectedRegistrations := []string{
		"river.AddWorker(workers, &SendEmailWorker{})",
		"river.AddWorker(workers, &ProcessPaymentWorker{})",
		"river.AddWorker(workers, &GenerateReportWorker{})",
	}
	for _, reg := range expectedRegistrations {
		if !strings.Contains(string(workerContent), reg) {
			t.Errorf("worker.go missing registration: %s", reg)
		}
	}
}

func TestInjectJobWorker(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal main.go that matches the generated template structure
	mainGoContent := `package main

import (
	"context"
	"log/slog"
	"os"

	"testmodule/database"
)

func main() {
	dbPath := "app.db"
	_, err := database.InitDB(dbPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.CloseDB()

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Routes go here
	slog.Info("Server starting")
}
`
	mainGoPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	if err := injectJobWorker(mainGoPath, "testmodule"); err != nil {
		t.Fatalf("injectJobWorker failed: %v", err)
	}

	result, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to read modified main.go: %v", err)
	}
	resultStr := string(result)

	// Verify River setup was injected
	checks := []string{
		"river.NewClient",
		"riversqlite.New",
		"jobs.SetupWorkers()",
		"riverClient.Start(appCtx)",
		"riverClient.Stop(stopCtx)",
		"jobs.SetClient(riverClient)",
		"\"database/sql\"",
		"\"github.com/riverqueue/river\"",
		"\"testmodule/app/jobs\"",
	}
	for _, check := range checks {
		if !strings.Contains(resultStr, check) {
			t.Errorf("main.go missing expected content: %s", check)
		}
	}

	// Verify idempotency — second call should be no-op
	if err := injectJobWorker(mainGoPath, "testmodule"); err != nil {
		t.Fatalf("Second injectJobWorker call failed: %v", err)
	}

	// Verify no duplicate injection
	count := strings.Count(resultStr, "river.NewClient")
	if count != 1 {
		t.Errorf("Expected 1 river.NewClient occurrence, got %d", count)
	}
}

// setupTestProject creates a minimal project structure for testing.
func setupTestProject(t *testing.T, dir string) {
	t.Helper()

	// Create database directory
	dbDir := filepath.Join(dir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	// Create migrations directory
	migrationsDir := filepath.Join(dbDir, "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create schema.sql
	schemaPath := filepath.Join(dbDir, "schema.sql")
	if err := os.WriteFile(schemaPath, []byte("-- existing schema\n"), 0644); err != nil {
		t.Fatalf("Failed to create schema.sql: %v", err)
	}

	// Create .lvtrc (project config)
	lvtConfig := "kit=multi\nmodule=testmodule\n"
	if err := os.WriteFile(filepath.Join(dir, ".lvtrc"), []byte(lvtConfig), 0644); err != nil {
		t.Fatalf("Failed to create .lvtrc: %v", err)
	}
}
