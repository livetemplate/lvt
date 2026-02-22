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

func printComponentHelp() {
	fmt.Println("lvt component - Manage UI components from the components library")
	fmt.Println()
	fmt.Println("Usage: lvt component <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list                              List all available components")
	fmt.Println("  eject <name>                      Eject component source to project")
	fmt.Println("  eject-template <name> <template>  Eject template only (keep library logic)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt component list")
	fmt.Println("  lvt component eject dropdown")
	fmt.Println("  lvt component eject-template dropdown searchable")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printComponentListHelp() {
	fmt.Println("lvt component list - List all available components")
	fmt.Println()
	fmt.Println("Usage: lvt component list")
	fmt.Println()
	fmt.Println("Lists all components available from github.com/livetemplate/components")
	fmt.Println("with their template names and descriptions.")
	fmt.Println()
	fmt.Println("Run 'lvt component --help' for more commands.")
}

func printComponentEjectHelp() {
	fmt.Println("lvt component eject - Eject a component to your project")
	fmt.Println()
	fmt.Println("Usage: lvt component eject <name>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <name>    Component name (e.g., 'dropdown', 'tabs', 'modal')")
	fmt.Println()
	fmt.Println("This command copies the full component source code to your project,")
	fmt.Println("giving you complete control over the component's behavior and templates.")
	fmt.Println()
	fmt.Println("Files are ejected to: internal/components/<name>/")
	fmt.Println()
	fmt.Println("After ejecting, update your imports:")
	fmt.Println("  from: github.com/livetemplate/components/<name>")
	fmt.Println("  to:   yourapp/internal/components/<name>")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt component eject dropdown")
	fmt.Println("  lvt component eject modal")
	fmt.Println()
	fmt.Println("Run 'lvt component list' to see available components.")
}

func printComponentEjectTemplateHelp() {
	fmt.Println("lvt component eject-template - Eject only the template file")
	fmt.Println()
	fmt.Println("Usage: lvt component eject-template <name> <template>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <name>      Component name (e.g., 'dropdown')")
	fmt.Println("  <template>  Template variant (e.g., 'searchable', 'default')")
	fmt.Println()
	fmt.Println("This command copies only the template file to your project.")
	fmt.Println("The Go logic remains in the library and continues to update automatically.")
	fmt.Println()
	fmt.Println("Files are ejected to: internal/templates/<name>-<template>.tmpl")
	fmt.Println()
	fmt.Println("Your local template will override the library template automatically.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt component eject-template dropdown searchable")
	fmt.Println("  lvt component eject-template modal default")
	fmt.Println()
	fmt.Println("Run 'lvt component list' to see available templates.")
}

func printValidateHelp() {
	fmt.Println("lvt validate - Validate a LiveTemplate app directory")
	fmt.Println()
	fmt.Println("Usage: lvt validate [<app-path>] [--fast]")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <app-path>    Path to the app directory (default: current directory)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --fast        Skip compilation check (faster, suitable for file watchers)")
	fmt.Println()
	fmt.Println("Checks run:")
	fmt.Println("  go.mod        Validates module path and go version directive")
	fmt.Println("  templates     Parses all .tmpl files for syntax errors")
	fmt.Println("  migrations    Executes SQL migrations against an in-memory SQLite DB")
	fmt.Println("  compilation   Runs 'go build ./...' (skipped with --fast)")
	fmt.Println()
	fmt.Println("Exit codes:")
	fmt.Println("  0   All checks passed")
	fmt.Println("  1   One or more errors found")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt validate")
	fmt.Println("  lvt validate ./myapp")
	fmt.Println("  lvt validate --fast")
	fmt.Println()
	fmt.Println("Run 'lvt --help' for full documentation.")
}

func printNewComponentHelp() {
	fmt.Println("lvt new component - Scaffold a new component")
	fmt.Println()
	fmt.Println("Usage: lvt new component <name>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <name>    Component name (e.g., 'rating', 'stepper')")
	fmt.Println()
	fmt.Println("Creates a new component scaffold with:")
	fmt.Println("  - <name>.go         Component struct and constructor")
	fmt.Println("  - options.go        Functional options")
	fmt.Println("  - templates.go      Template embedding")
	fmt.Println("  - templates/default.tmpl")
	fmt.Println("  - <name>_test.go    Test file skeleton")
	fmt.Println()
	fmt.Println("Files are created in: components/<name>/")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt new component rating")
	fmt.Println("  lvt new component stepper")
	fmt.Println()
	fmt.Println("After creating, you can:")
	fmt.Println("  - Use locally in your project")
	fmt.Println("  - Submit PR to github.com/livetemplate/components")
	fmt.Println("  - Publish as separate Go module")
}
