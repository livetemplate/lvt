package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"

	// Register both adapters so they're available for listing.
	_ "github.com/livetemplate/lvt/components/styles/tailwind"
	_ "github.com/livetemplate/lvt/components/styles/unstyled"

	"github.com/livetemplate/lvt/components/styles"
	unstyledpkg "github.com/livetemplate/lvt/components/styles/unstyled"
)

// Styles handles the "lvt styles" command and subcommands.
func Styles(args []string) error {
	if len(args) == 0 {
		printStylesHelp()
		return nil
	}

	if ShowHelpIfRequested(args, printStylesHelp) {
		return nil
	}

	switch args[0] {
	case "list":
		return stylesList()
	case "info":
		if len(args) < 2 {
			return fmt.Errorf("adapter name required: lvt styles info <name>")
		}
		return stylesInfo(args[1])
	case "scaffold":
		name := "unstyled"
		output := ""
		for i := 1; i < len(args); i++ {
			if args[i] == "--output" && i+1 < len(args) {
				output = args[i+1]
				i++
			} else {
				name = args[i]
			}
		}
		return stylesScaffold(name, output)
	default:
		return fmt.Errorf("unknown styles subcommand: %s", args[0])
	}
}

func stylesList() error {
	names := styles.Names()
	sort.Strings(names)

	defaultAdapter := styles.Default()
	defaultName := ""
	if defaultAdapter != nil {
		defaultName = defaultAdapter.Name()
	}

	fmt.Println("Registered style adapters:")
	fmt.Println()
	for _, name := range names {
		marker := "  "
		if name == defaultName {
			marker = "* "
		}
		fmt.Printf("  %s%s\n", marker, name)
	}
	fmt.Println()
	fmt.Printf("  %d adapter(s) registered (* = default)\n", len(names))
	return nil
}

func stylesInfo(name string) error {
	adapter := styles.Get(name)
	if adapter == nil {
		available := strings.Join(styles.Names(), ", ")
		return fmt.Errorf("adapter %q not found (available: %s)", name, available)
	}

	fmt.Printf("Style Adapter: %s\n", adapter.Name())
	fmt.Println()

	if name == "unstyled" {
		count := unstyledpkg.ClassCount()
		fmt.Printf("  BEM class names: %d\n", count)
		fmt.Println("  Convention: lvt-{component}__{element}--{modifier}")
		fmt.Println()
		fmt.Println("  Use 'lvt styles scaffold' to generate a CSS file with all class stubs.")
	} else if name == "tailwind" {
		fmt.Println("  Framework: Tailwind CSS")
		fmt.Println("  All component classes use Tailwind utility classes.")
		fmt.Println("  No additional CSS file needed.")
	}

	fmt.Println()
	fmt.Println("  Components: 29 style types across 20 component packages")

	return nil
}

func stylesScaffold(name, output string) error {
	if name != "unstyled" {
		return fmt.Errorf("scaffold generation is only supported for the 'unstyled' adapter")
	}

	if output == "" {
		// Write to stdout
		return unstyledpkg.GenerateCSS(os.Stdout)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	if err := unstyledpkg.GenerateCSS(f); err != nil {
		return err
	}

	fmt.Printf("CSS scaffold written to %s (%d BEM classes)\n", output, unstyledpkg.ClassCount())
	return nil
}

func printStylesHelp() {
	fmt.Println("Manage component style adapters")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lvt styles <command> [args...]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  list                      List registered style adapters")
	fmt.Println("  info <name>               Show adapter details")
	fmt.Println("  scaffold [--output file]  Generate CSS scaffold file (unstyled adapter)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  lvt styles list")
	fmt.Println("  lvt styles info tailwind")
	fmt.Println("  lvt styles scaffold --output styles.css")
}
