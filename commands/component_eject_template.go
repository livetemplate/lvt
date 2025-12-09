package commands

import (
	"fmt"

	"github.com/livetemplate/lvt/internal/eject"
)

// ComponentEjectTemplate handles the `lvt component eject-template` command.
func ComponentEjectTemplate(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printComponentEjectTemplateHelp) {
		return nil
	}

	if len(args) < 2 {
		if len(args) == 1 {
			// Show available templates for this component
			comp := eject.FindComponent(args[0])
			if comp == nil {
				return fmt.Errorf("unknown component: %s\nRun 'lvt component list' to see available components", args[0])
			}
			fmt.Printf("Available templates for %s:\n", comp.Name)
			for _, t := range comp.Templates {
				fmt.Printf("  - %s\n", t)
			}
			fmt.Println()
			fmt.Printf("Usage: lvt component eject-template %s <template>\n", comp.Name)
			return nil
		}
		printComponentEjectTemplateHelp()
		return nil
	}

	componentName := args[0]
	templateName := args[1]

	// Validate that names don't look like flags
	if err := ValidatePositionalArg(componentName, "component name"); err != nil {
		return err
	}
	if err := ValidatePositionalArg(templateName, "template name"); err != nil {
		return err
	}

	// Parse options
	opts := eject.EjectTemplateOptions{
		ComponentName: componentName,
		TemplateName:  templateName,
	}

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--force", "-f":
			opts.Force = true
		case "--dest", "-d":
			if i+1 < len(args) {
				opts.DestDir = args[i+1]
				i++
			}
		}
	}

	return eject.EjectTemplate(opts)
}
