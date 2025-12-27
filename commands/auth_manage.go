package commands

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// AuthManage handles auth management subcommands (confirm, list, etc.)
func AuthManage(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand\n\nUsage:\n  lvt auth [--db <path>] confirm <email>    Confirm a user's email (for testing)\n  lvt auth [--db <path>] list               List all users")
	}

	// Parse --db flag
	var dbPath string
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--db" && i+1 < len(args) {
			dbPath = args[i+1]
			i++ // skip the path argument
		} else {
			filteredArgs = append(filteredArgs, args[i])
		}
	}

	if len(filteredArgs) == 0 {
		return fmt.Errorf("missing subcommand\n\nUsage:\n  lvt auth [--db <path>] confirm <email>    Confirm a user's email (for testing)\n  lvt auth [--db <path>] list               List all users")
	}

	subcommand := filteredArgs[0]
	subArgs := filteredArgs[1:]

	switch subcommand {
	case "confirm":
		return AuthConfirm(subArgs, dbPath)
	case "list":
		return AuthList(subArgs, dbPath)
	default:
		return fmt.Errorf("unknown auth subcommand: %s\n\nAvailable subcommands:\n  confirm <email>    Confirm a user's email (for testing)\n  list               List all users", subcommand)
	}
}

// AuthConfirm confirms a user's email in the database
func AuthConfirm(args []string, dbPath string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing email address\n\nUsage: lvt auth confirm <email>")
	}

	email := args[0]

	// Find the database file if not specified
	if dbPath == "" {
		var err error
		dbPath, err = findDatabaseFile()
		if err != nil {
			return err
		}
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Check if user exists
	var userID string
	var confirmedAt sql.NullTime
	err = db.QueryRow("SELECT id, confirmed_at FROM users WHERE email = ?", email).Scan(&userID, &confirmedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found: %s", email)
	}
	if err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	}

	// Check if already confirmed
	if confirmedAt.Valid {
		fmt.Printf("User %s is already confirmed (confirmed at: %s)\n", email, confirmedAt.Time.Format(time.RFC3339))
		return nil
	}

	// Confirm the user
	now := time.Now()
	_, err = db.Exec("UPDATE users SET confirmed_at = ?, updated_at = ? WHERE email = ?", now, now, email)
	if err != nil {
		return fmt.Errorf("failed to confirm user: %w", err)
	}

	fmt.Printf("✅ User %s confirmed successfully!\n", email)
	return nil
}

// AuthList lists all users in the database
func AuthList(args []string, dbPath string) error {
	// Find the database file if not specified
	if dbPath == "" {
		var err error
		dbPath, err = findDatabaseFile()
		if err != nil {
			return err
		}
	}

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Query users
	rows, err := db.Query("SELECT id, email, confirmed_at, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	fmt.Println("Users:")
	fmt.Println("------")
	count := 0
	for rows.Next() {
		var id, email string
		var confirmedAt sql.NullTime
		var createdAt time.Time
		if err := rows.Scan(&id, &email, &confirmedAt, &createdAt); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		status := "❌ unconfirmed"
		if confirmedAt.Valid {
			status = "✅ confirmed"
		}

		fmt.Printf("  %s (%s) - created: %s\n", email, status, createdAt.Format("2006-01-02 15:04"))
		count++
	}

	if count == 0 {
		fmt.Println("  (no users found)")
	} else {
		fmt.Printf("\nTotal: %d user(s)\n", count)
	}

	return nil
}

// findDatabaseFile looks for the app database file
func findDatabaseFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Common database file locations
	candidates := []string{
		filepath.Join(wd, "app.db"),
		filepath.Join(wd, "data.db"),
		filepath.Join(wd, "database.db"),
		filepath.Join(wd, "db", "app.db"),
		filepath.Join(wd, "data", "app.db"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("database file not found. Tried: %v\n\nMake sure you're in the project root and the app has been run at least once", candidates)
}
