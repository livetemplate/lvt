package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
)

// GenQueue sets up background job infrastructure using River.
func GenQueue(args []string) error {
	if ShowHelpIfRequested(args, printGenQueueHelp) {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	projectConfig, err := config.LoadProjectConfig(cwd)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	moduleName := projectConfig.Module
	if moduleName == "" {
		return fmt.Errorf("could not determine module name from project config")
	}

	if err := generator.GenerateQueue(cwd, moduleName); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ Background job queue set up successfully!")
	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Println("  app/jobs/worker.go                         Job worker registration")
	fmt.Println("  database/migrations/..._setup_river_queue.sql  River queue tables")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run 'lvt gen job <name>' to create your first job handler")
	fmt.Println("  2. Run 'go mod tidy' to fetch River dependencies")
	fmt.Println("  3. Run 'lvt migration up' to create River tables")
	fmt.Println("  4. Start your app — the job worker will start automatically")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  lvt gen job send_email")

	return nil
}

// GenJob scaffolds a new background job handler.
func GenJob(args []string) error {
	if ShowHelpIfRequested(args, printGenJobHelp) {
		return nil
	}

	if len(args) < 1 {
		return fmt.Errorf("job name required\n\nUsage: lvt gen job <name>\n\nExamples:\n  lvt gen job send_email\n  lvt gen job process_payment\n  lvt gen job generate_report")
	}

	jobName := strings.TrimSpace(args[0])
	if jobName == "" {
		return fmt.Errorf("job name cannot be empty")
	}

	if err := ValidatePositionalArg(jobName, "job name"); err != nil {
		return err
	}

	// Normalize to snake_case
	jobName = strings.ToLower(jobName)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	projectConfig, err := config.LoadProjectConfig(cwd)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	moduleName := projectConfig.Module
	if moduleName == "" {
		return fmt.Errorf("could not determine module name from project config")
	}

	if err := generator.GenerateJob(cwd, moduleName, jobName); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("✅ Job '%s' created successfully!\n", jobName)
	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Printf("  app/jobs/%s.go    Job handler\n", jobName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit app/jobs/%s.go to define your payload and logic\n", jobName)
	fmt.Println("  2. Enqueue jobs from your handlers:")
	fmt.Println()
	fmt.Printf("     // In any handler with access to riverClient:\n")
	fmt.Printf("     riverClient.Insert(ctx, jobs.%sArgs{...}, nil)\n", generator.ToCamelCase(jobName))

	return nil
}

func printGenQueueHelp() {
	fmt.Println("Usage: lvt gen queue")
	fmt.Println()
	fmt.Println("Set up background job processing infrastructure using River.")
	fmt.Println("This is a one-time setup command that creates:")
	fmt.Println()
	fmt.Println("  - Database migration for River queue tables")
	fmt.Println("  - Worker registration file (app/jobs/worker.go)")
	fmt.Println("  - River client setup in main.go")
	fmt.Println()
	fmt.Println("River (https://riverqueue.com) provides:")
	fmt.Println("  - Worker pool with configurable concurrency")
	fmt.Println("  - Retry with exponential backoff")
	fmt.Println("  - Scheduled and periodic jobs")
	fmt.Println("  - Dead letter queue (discarded jobs)")
	fmt.Println("  - Unique/deduplicated jobs")
	fmt.Println("  - Graceful shutdown")
	fmt.Println("  - SQLite and PostgreSQL support")
}

func printGenJobHelp() {
	fmt.Println("Usage: lvt gen job <name>")
	fmt.Println()
	fmt.Println("Scaffold a new background job handler.")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  name    Job name in snake_case (e.g., send_email, process_payment)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen job send_email")
	fmt.Println("  lvt gen job process_payment")
	fmt.Println("  lvt gen job generate_report")
	fmt.Println("  lvt gen job cleanup_expired_sessions")
	fmt.Println()
	fmt.Println("Prerequisites:")
	fmt.Println("  Run 'lvt gen queue' first to set up the job infrastructure.")
}
