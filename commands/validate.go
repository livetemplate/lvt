package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/livetemplate/lvt/internal/validation"
)

// Validate runs validation checks against a LiveTemplate app directory and
// prints a formatted report. It exits with code 1 if any errors are found.
func Validate(args []string) error {
	if ShowHelpIfRequested(args, printValidateHelp) {
		return nil
	}

	appPath := "."
	fast := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--fast":
			fast = true
		default:
			if len(args[i]) > 0 && args[i][0] == '-' {
				return fmt.Errorf("unknown flag: %s\nRun 'lvt validate --help' for usage", args[i])
			}
			if appPath != "." {
				return fmt.Errorf("unexpected argument: %q\nUsage: lvt validate [<app-path>] [--fast]", args[i])
			}
			appPath = args[i]
		}
	}

	fmt.Printf("Validating: %s\n", appPath)
	if fast {
		fmt.Println("Mode: fast (skipping compilation check)")
	}
	fmt.Println()

	var engine *validation.Engine
	if fast {
		engine = validation.NewEngine(
			validation.WithCheck(&validation.GoModCheck{}),
			validation.WithCheck(&validation.TemplateCheck{}),
			validation.WithCheck(&validation.MigrationCheck{}),
		)
	} else {
		engine = validation.DefaultEngine()
	}

	result := engine.Run(context.Background(), appPath)
	fmt.Print(result.Format())

	if !result.Valid {
		os.Exit(1)
	}

	return nil
}
