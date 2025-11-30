package commands

import (
	"fmt"
	"strings"
)

// ShowHelpIfRequested checks if args contain --help or -h flags.
// If found, it calls the provided help function and returns true.
// Callers should return early when this returns true.
func ShowHelpIfRequested(args []string, helpFunc func()) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			helpFunc()
			return true
		}
	}
	return false
}

// ValidatePositionalArg ensures a positional argument doesn't look like a flag.
// Returns an error if the argument starts with "-" (likely a typo or unknown flag).
func ValidatePositionalArg(arg, argName string) error {
	if strings.HasPrefix(arg, "-") {
		return fmt.Errorf("invalid %s: %q looks like a flag\n\nDid you mean to use a flag? Run with --help for usage", argName, arg)
	}
	return nil
}

// Help functions for each command

func printNewHelp() {
	fmt.Println("lvt new - Create a new LiveTemplate application")
	fmt.Println()
	fmt.Println("Usage: lvt new <app-name> [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --module <name>   Go module name (defaults to app name)")
	fmt.Println("  --kit <kit>       Template kit: multi, single, simple (default: multi)")
	fmt.Println("  --dev             Use local development mode")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printGenHelp() {
	fmt.Println("lvt gen - Generate code for resources, views, schemas, or auth")
	fmt.Println()
	fmt.Println("Usage: lvt gen <subcommand> [args...]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  resource <name> <field:type>...   Generate full CRUD with database")
	fmt.Println("  view <name>                       Generate view-only handler (no database)")
	fmt.Println("  schema <table> <field:type>...    Generate database schema only")
	fmt.Println("  auth [StructName] [table_name]    Generate authentication system")
	fmt.Println("  stack <provider>                  Generate deployment stack")
	fmt.Println()
	fmt.Println("Run 'lvt gen <subcommand> --help' for subcommand-specific help.")
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printGenResourceHelp() {
	fmt.Println("lvt gen resource - Generate a CRUD resource with database integration")
	fmt.Println()
	fmt.Println("Usage: lvt gen resource <name> <field:type>...")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <name>          Resource name (singular, e.g., 'post', 'user')")
	fmt.Println("  <field:type>    Field definitions (type optional, defaults to string)")
	fmt.Println()
	fmt.Println("Types: string, int, bool, float, time, text, textarea")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen resource posts title content:text published:bool")
	fmt.Println("  lvt gen resource users name email age:int")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printGenViewHelp() {
	fmt.Println("lvt gen view - Generate a view-only handler (no database)")
	fmt.Println()
	fmt.Println("Usage: lvt gen view <name>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <name>    View name (e.g., 'dashboard', 'counter')")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen view dashboard")
	fmt.Println("  lvt gen view counter")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printGenSchemaHelp() {
	fmt.Println("lvt gen schema - Generate database schema only (no handlers/templates)")
	fmt.Println()
	fmt.Println("Usage: lvt gen schema <table> <field:type>...")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <table>         Table name")
	fmt.Println("  <field:type>    Field definitions")
	fmt.Println()
	fmt.Println("Types: string, int, bool, float, time, text")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt gen schema products name price:float quantity:int")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printGenStackHelp() {
	fmt.Println("lvt gen stack - Generate deployment stack configuration")
	fmt.Println()
	fmt.Println("Usage: lvt gen stack <provider>")
	fmt.Println()
	fmt.Println("Providers:")
	fmt.Println("  fly        Fly.io deployment configuration")
	fmt.Println("  docker     Docker/Docker Compose configuration")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printParseHelp() {
	fmt.Println("lvt parse - Validate and analyze a template file")
	fmt.Println()
	fmt.Println("Usage: lvt parse <template-file>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <template-file>    Path to .tmpl file to validate")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printSeedHelp() {
	fmt.Println("lvt seed - Generate test data for a resource")
	fmt.Println()
	fmt.Println("Usage: lvt seed <resource> [options]")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <resource>    Resource name to seed")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --count N     Number of records to generate (default: 10)")
	fmt.Println("  --cleanup     Remove existing test data before seeding")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt seed posts --count 50")
	fmt.Println("  lvt seed users --cleanup")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printResourceHelp() {
	fmt.Println("lvt resource - Inspect resources and schemas")
	fmt.Println()
	fmt.Println("Usage: lvt resource <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list              List all available resources")
	fmt.Println("  describe <name>   Show detailed schema for a resource")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printMigrationHelp() {
	fmt.Println("lvt migration - Manage database migrations")
	fmt.Println()
	fmt.Println("Usage: lvt migration <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up              Run pending migrations")
	fmt.Println("  down            Rollback last migration")
	fmt.Println("  status          Show migration status")
	fmt.Println("  create <name>   Create new migration file")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printKitsHelp() {
	fmt.Println("lvt kits - Manage CSS framework kits")
	fmt.Println()
	fmt.Println("Usage: lvt kits <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list              List all available kits")
	fmt.Println("  info <kit>        Show detailed information about a kit")
	fmt.Println("  create <name>     Create a new custom kit")
	fmt.Println("  validate <path>   Validate a kit implementation")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printEnvHelp() {
	fmt.Println("lvt env - Manage environment variables")
	fmt.Println()
	fmt.Println("Usage: lvt env <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  generate          Generate .env.example with detected config")
	fmt.Println("  list              List configured environment variables")
	fmt.Println("  set <key> <val>   Set an environment variable")
	fmt.Println("  unset <key>       Unset an environment variable")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}
