package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Console opens an interactive database shell for the current app.
func Console(args []string) error {
	if ShowHelpIfRequested(args, printConsoleHelp) {
		return nil
	}

	dbPath := findDBPath()
	if dbPath == "" {
		return fmt.Errorf("no database found. Are you in a LiveTemplate project directory?\nExpected: app.db or DATABASE_PATH environment variable")
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found: %s\nRun 'lvt migration up' to create it", dbPath)
	}

	sqlite3Path, err := exec.LookPath("sqlite3")
	if err != nil {
		fmt.Printf("Database path: %s\n", dbPath)
		return fmt.Errorf("sqlite3 not found in PATH. Install it:\n  macOS:  brew install sqlite3\n  Ubuntu: sudo apt install sqlite3\n  Or open manually: sqlite3 %s", dbPath)
	}

	fmt.Printf("Opening %s...\n", dbPath)
	fmt.Println("Type .tables to list tables, .schema to see schema, .quit to exit")
	fmt.Println()

	cmd := exec.Command(sqlite3Path, "-header", "-column", dbPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// findDBPath locates the database file for the current project.
func findDBPath() string {
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		return path
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		path := filepath.Join(dir, "app.db")
		if _, err := os.Stat(path); err == nil {
			return path
		}
		// Stop at project boundary (go.mod indicates project root)
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func printConsoleHelp() {
	fmt.Println("Usage: lvt console  (alias: lvt db)")
	fmt.Println()
	fmt.Println("Opens an interactive SQLite database shell for the current app.")
	fmt.Println()
	fmt.Println("The database is located by:")
	fmt.Println("  1. DATABASE_PATH environment variable")
	fmt.Println("  2. app.db in the current or parent directories")
	fmt.Println()
	fmt.Println("Useful SQLite commands:")
	fmt.Println("  .tables     List all tables")
	fmt.Println("  .schema     Show CREATE TABLE statements")
	fmt.Println("  .headers on Enable column headers")
	fmt.Println("  .mode csv   Switch to CSV output")
	fmt.Println("  .quit       Exit the console")
	fmt.Println()
}
