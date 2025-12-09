package commands

import (
	"fmt"

	"github.com/livetemplate/lvt/internal/eject"
)

// Component handles the `lvt component` command group.
func Component(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printComponentHelp) {
		return nil
	}

	if len(args) < 1 {
		printComponentHelp()
		return nil
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "eject":
		return ComponentEject(subArgs)
	case "eject-template":
		return ComponentEjectTemplate(subArgs)
	case "list":
		return ComponentList(subArgs)
	default:
		return fmt.Errorf("unknown component subcommand: %s\n\nRun 'lvt component --help' for usage", subcommand)
	}
}

// ComponentList lists available components from the library.
func ComponentList(args []string) error {
	if ShowHelpIfRequested(args, printComponentListHelp) {
		return nil
	}

	fmt.Println("Available components from github.com/livetemplate/components:")
	fmt.Println()

	// Use single source of truth from eject package
	for _, c := range eject.AvailableComponents() {
		fmt.Printf("  %s\n", c.Name)
		fmt.Printf("    %s\n", c.Description)
		// Format template names with full lvt: prefix
		templates := make([]string, len(c.Templates))
		for i, t := range c.Templates {
			templates[i] = fmt.Sprintf("lvt:%s:%s:v1", c.Name, t)
		}
		fmt.Printf("    Templates: %v\n", templates)
		fmt.Println()
	}

	fmt.Println("To eject a component:")
	fmt.Println("  lvt component eject <name>")
	fmt.Println()
	fmt.Println("To eject just a template:")
	fmt.Println("  lvt component eject-template <name> <template>")

	return nil
}
