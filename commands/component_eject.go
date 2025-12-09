package commands

import (
	"github.com/livetemplate/lvt/internal/eject"
)

// ComponentEject handles the `lvt component eject` command.
func ComponentEject(args []string) error {
	// Handle --help flag
	if ShowHelpIfRequested(args, printComponentEjectHelp) {
		return nil
	}

	if len(args) < 1 {
		printComponentEjectHelp()
		return nil
	}

	componentName := args[0]

	// Validate that component name doesn't look like a flag
	if err := ValidatePositionalArg(componentName, "component name"); err != nil {
		return err
	}

	// Parse options
	opts := eject.EjectOptions{
		ComponentName: componentName,
	}

	for i := 1; i < len(args); i++ {
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

	// Try to get module name
	if modName, err := eject.GetModuleName(); err == nil {
		opts.ModuleName = modName
	}

	return eject.EjectComponent(opts)
}
