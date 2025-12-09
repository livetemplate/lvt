package commands

import (
	"fmt"
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

	components := []struct {
		name        string
		description string
		templates   []string
	}{
		{"accordion", "Collapsible content sections", []string{"lvt:accordion:default:v1", "lvt:accordion:single:v1"}},
		{"autocomplete", "Search with suggestions", []string{"lvt:autocomplete:default:v1"}},
		{"breadcrumbs", "Navigation breadcrumb trail", []string{"lvt:breadcrumbs:default:v1"}},
		{"datatable", "Data tables with sorting/pagination", []string{"lvt:datatable:default:v1"}},
		{"datepicker", "Date selection", []string{"lvt:datepicker:single:v1", "lvt:datepicker:range:v1", "lvt:datepicker:inline:v1"}},
		{"drawer", "Slide-out panels", []string{"lvt:drawer:default:v1"}},
		{"dropdown", "Dropdown menus", []string{"lvt:dropdown:default:v1", "lvt:dropdown:searchable:v1", "lvt:dropdown:multi:v1"}},
		{"menu", "Navigation menus", []string{"lvt:menu:default:v1", "lvt:menu:nested:v1"}},
		{"modal", "Modal dialogs", []string{"lvt:modal:default:v1", "lvt:modal:confirm:v1", "lvt:modal:sheet:v1"}},
		{"popover", "Rich content popovers", []string{"lvt:popover:default:v1"}},
		{"progress", "Progress indicators", []string{"lvt:progress:default:v1", "lvt:progress:circular:v1", "lvt:progress:spinner:v1"}},
		{"rating", "Star ratings", []string{"lvt:rating:default:v1"}},
		{"skeleton", "Loading placeholders", []string{"lvt:skeleton:default:v1", "lvt:skeleton:avatar:v1", "lvt:skeleton:card:v1"}},
		{"tabs", "Tab navigation", []string{"lvt:tabs:horizontal:v1", "lvt:tabs:vertical:v1", "lvt:tabs:pills:v1"}},
		{"tagsinput", "Tag/chip input", []string{"lvt:tagsinput:default:v1"}},
		{"timeline", "Event timelines", []string{"lvt:timeline:default:v1"}},
		{"timepicker", "Time selection", []string{"lvt:timepicker:default:v1"}},
		{"toast", "Toast notifications", []string{"lvt:toast:default:v1", "lvt:toast:container:v1"}},
		{"toggle", "Toggle switches", []string{"lvt:toggle:default:v1", "lvt:toggle:checkbox:v1"}},
		{"tooltip", "Tooltips", []string{"lvt:tooltip:default:v1"}},
	}

	for _, c := range components {
		fmt.Printf("  %s\n", c.name)
		fmt.Printf("    %s\n", c.description)
		fmt.Printf("    Templates: %v\n", c.templates)
		fmt.Println()
	}

	fmt.Println("To eject a component:")
	fmt.Println("  lvt component eject <name>")
	fmt.Println()
	fmt.Println("To eject just a template:")
	fmt.Println("  lvt component eject-template <name> <template>")

	return nil
}
