package commands

import (
	"fmt"
	"os"

	"github.com/livetemplate/lvt/internal/generator"
)

// Authz generates the authorization system (role column + queries).
func Authz(args []string) error {
	if ShowHelpIfRequested(args, printGenAuthzHelp) {
		return nil
	}

	basePath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	moduleName, err := getModuleName()
	if err != nil {
		return fmt.Errorf("failed to get module name: %w (are you in a Go project?)", err)
	}

	// Default table name; could be customizable in the future
	tableName := "users"
	if len(args) > 0 {
		tableName = args[0]
	}

	fmt.Println("Generating authorization system...")

	// Validate auth exists
	if _, err := os.Stat("app/auth"); os.IsNotExist(err) {
		return fmt.Errorf("auth system not found. Run 'lvt gen auth' first")
	}

	cfg := &generator.AuthzConfig{
		ModuleName: moduleName,
		TableName:  tableName,
	}

	if err := generator.GenerateAuthz(basePath, cfg); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ Authorization system generated!")
	fmt.Println()
	fmt.Println("Files created/updated:")
	fmt.Println("  database/migrations/<timestamp>_add_user_roles.sql")
	fmt.Println("  database/queries.sql (role queries appended)")
	fmt.Println("  database/schema.sql (role column added)")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Run migration:")
	fmt.Println("     lvt migration up")
	fmt.Println("  2. Regenerate sqlc code:")
	fmt.Println("     sqlc generate")
	fmt.Println("  3. Generate resources with authorization:")
	fmt.Printf("     lvt gen resource posts title content:text --with-authz\n")
	fmt.Println("  4. Set admin role:")
	fmt.Println("     UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';")
	fmt.Println()

	return nil
}

func printGenAuthzHelp() {
	fmt.Println("Usage: lvt gen authz [table_name]")
	fmt.Println()
	fmt.Println("Generates role-based authorization for the auth system.")
	fmt.Println("Adds a 'role' column to the users table and role management queries.")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  table_name    Users table name (default: users)")
	fmt.Println()
	fmt.Println("Prerequisites:")
	fmt.Println("  Run 'lvt gen auth' first to set up the authentication system.")
	fmt.Println()
	fmt.Println("After running this command, use --with-authz on resource generation:")
	fmt.Println("  lvt gen resource posts title content:text --with-authz")
	fmt.Println()
}
