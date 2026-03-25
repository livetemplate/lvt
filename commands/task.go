package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
)

// GenTask scaffolds a new scheduled task.
func GenTask(args []string) error {
	if ShowHelpIfRequested(args, printGenTaskHelp) {
		return nil
	}

	if len(args) < 1 {
		return fmt.Errorf("task name required\n\nUsage: lvt gen task <name> --schedule <interval>\n\nExamples:\n  lvt gen task cleanup --schedule @hourly\n  lvt gen task daily_report --schedule @daily")
	}

	// Parse flags
	schedule := "@hourly" // default
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		if (args[i] == "--schedule" || args[i] == "-s") && i+1 < len(args) {
			schedule = args[i+1]
			i++
		} else {
			filteredArgs = append(filteredArgs, args[i])
		}
	}

	if len(filteredArgs) < 1 {
		return fmt.Errorf("task name required\n\nUsage: lvt gen task <name> --schedule <interval>")
	}
	taskName := strings.ToLower(strings.TrimSpace(filteredArgs[0]))
	if taskName == "" {
		return fmt.Errorf("task name cannot be empty")
	}

	if err := ValidatePositionalArg(taskName, "task name"); err != nil {
		return err
	}

	// Validate schedule format
	validSchedules := map[string]bool{
		"@hourly": true, "@daily": true, "@weekly": true,
		"@every 1m": true, "@every 5m": true, "@every 10m": true, "@every 30m": true,
		"@every 1 minute": true, "@every 5 minutes": true, "@every 10 minutes": true, "@every 30 minutes": true,
	}
	if !validSchedules[schedule] && !isValidEverySchedule(schedule) {
		return fmt.Errorf("unsupported schedule %q\n\nSupported formats:\n  @hourly, @daily, @weekly\n  @every Nm (e.g., @every 5m, @every 2h, @every 30s)", schedule)
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

	if err := generator.GenerateTask(cwd, moduleName, taskName, schedule); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("✅ Scheduled task '%s' created!\n", taskName)
	fmt.Println()
	fmt.Printf("Schedule: %s\n", schedule)
	fmt.Println()
	fmt.Println("Generated files:")
	fmt.Printf("  app/jobs/%s.go    Task handler\n", taskName)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Edit app/jobs/%s.go to implement your task logic\n", taskName)
	fmt.Println("  2. The task will run automatically when the app starts")
	fmt.Println()

	return nil
}

func printGenTaskHelp() {
	fmt.Println("Usage: lvt gen task <name> [--schedule <interval>]")
	fmt.Println()
	fmt.Println("Scaffold a new scheduled/recurring task.")
	fmt.Println("Tasks run automatically on a schedule using the River job queue.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --schedule, -s    Schedule interval (default: @hourly)")
	fmt.Println()
	fmt.Println("Schedule shortcuts:")
	fmt.Println("  @hourly           Run every hour")
	fmt.Println("  @daily            Run every day")
	fmt.Println("  @weekly           Run every week")
	fmt.Println("  @every 5m         Run every 5 minutes")
	fmt.Println("  @every 30m        Run every 30 minutes")
	fmt.Println("  @every 2h         Run every 2 hours")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen task cleanup_sessions --schedule @hourly")
	fmt.Println("  lvt gen task daily_report --schedule @daily")
	fmt.Println("  lvt gen task sync_data --schedule \"@every 5m\"")
	fmt.Println()
	fmt.Println("Prerequisites:")
	fmt.Println("  Run 'lvt gen queue' first to set up the job infrastructure.")
}

func isValidEverySchedule(s string) bool {
	if !strings.HasPrefix(s, "@every ") {
		return false
	}
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return false
	}
	d := parts[1]
	if len(d) < 2 {
		return false
	}
	suffix := d[len(d)-1]
	if suffix != 'm' && suffix != 'h' && suffix != 's' {
		return false
	}
	n := d[:len(d)-1]
	for _, c := range n {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(n) > 0
}
