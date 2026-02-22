package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/livetemplate/lvt/internal/config"
	"github.com/livetemplate/lvt/internal/generator"
)

type AuthFlags struct {
	NoPassword      bool
	NoMagicLink     bool
	NoEmailConfirm  bool
	NoPasswordReset bool
	NoSessionsUI    bool
	NoCSRF          bool
}

func Auth(args []string) error {
	flags := &AuthFlags{}
	skipValidation := false
	var structName, tableName string
	var positionalArgs []string

	// Separate flags from positional arguments
	for _, arg := range args {
		switch arg {
		case "--no-password":
			flags.NoPassword = true
		case "--no-magic-link":
			flags.NoMagicLink = true
		case "--no-email-confirm":
			flags.NoEmailConfirm = true
		case "--no-password-reset":
			flags.NoPasswordReset = true
		case "--no-sessions-ui":
			flags.NoSessionsUI = true
		case "--no-csrf":
			flags.NoCSRF = true
		case "--skip-validation":
			skipValidation = true
		default:
			if !startsWithDash(arg) {
				positionalArgs = append(positionalArgs, arg)
			}
		}
	}

	// Parse positional arguments: [struct_name] [table_name]
	if len(positionalArgs) > 0 {
		structName = positionalArgs[0]
	}
	if len(positionalArgs) > 1 {
		tableName = positionalArgs[1]
	}

	// Default to User/users if not specified
	if structName == "" {
		structName = "User"
	}
	if tableName == "" {
		tableName = pluralizeNoun(structName)
	}

	// Validate flags
	if flags.NoPassword && flags.NoMagicLink {
		return errors.New("at least one authentication method (password or magic-link) must be enabled")
	}

	// Get project root
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Load project config
	cfg, err := config.LoadProjectConfig(wd)
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	// Create generator config
	genConfig := &generator.AuthConfig{
		ModuleName:          cfg.Module,
		StructName:          structName,
		TableName:           tableName,
		EnablePassword:      !flags.NoPassword,
		EnableMagicLink:     !flags.NoMagicLink,
		EnableEmailConfirm:  !flags.NoEmailConfirm,
		EnablePasswordReset: !flags.NoPasswordReset,
		EnableSessionsUI:    !flags.NoSessionsUI,
		EnableCSRF:          !flags.NoCSRF,
	}

	// Generate auth files
	fmt.Println("Generating authentication system...")
	if err := generator.GenerateAuth(wd, genConfig); err != nil {
		return fmt.Errorf("failed to generate auth: %w", err)
	}

	fmt.Println("✅ Authentication system generated successfully!")
	fmt.Println("\n📁 Generated files:")
	fmt.Println("  - app/auth/auth.go          (handler with all auth flows)")
	fmt.Println("  - app/auth/auth.tmpl        (LiveTemplate UI)")
	fmt.Println("  - app/auth/middleware.go    (route protection middleware)")
	fmt.Println("  - app/auth/auth_e2e_test.go (E2E tests with chromedp)")
	fmt.Println("  - database/migrations/      (auth tables migration)")
	fmt.Println("  - database/queries.sql      (auth SQL queries)")
	fmt.Println("\n📦 Dependencies added:")
	fmt.Println("  - github.com/livetemplate/lvt/pkg/password (bcrypt utilities)")
	fmt.Println("  - github.com/livetemplate/lvt/pkg/email    (email sender interface)")

	// Post-generation validation (before interactive prompts).
	// Unlike gen.go which defers the error to show the full file listing,
	// auth returns immediately on validation failure because the interactive
	// resource-protection prompts below depend on a healthy app state.
	if !skipValidation {
		if err := runPostGenValidation(wd); err != nil {
			return err
		}
	}

	// Check for existing resources to protect
	resources, err := generator.ReadResources(wd)
	if err != nil {
		fmt.Printf("⚠️  Could not read resources: %v\n", err)
	} else if len(resources) > 0 {
		// Filter out auth and home from protectable resources
		var protectableResources []generator.ResourceEntry
		for _, r := range resources {
			if r.Name != "Auth" && r.Name != "Home" && r.Type == "resource" {
				protectableResources = append(protectableResources, r)
			}
		}

		if len(protectableResources) > 0 {
			fmt.Println("\n🔒 Resource Protection")
			fmt.Println("Would you like to protect any resources with authentication?")
			fmt.Println("(Users will need to log in to access protected resources)")
			fmt.Println()

			// Show available resources
			for i, r := range protectableResources {
				fmt.Printf("  %d. %s (%s)\n", i+1, r.Name, r.Path)
			}
			fmt.Println()
			fmt.Println("Enter your selection:")
			fmt.Println("  - Numbers separated by commas (e.g., 1,2,3)")
			fmt.Println("  - 'all' to protect all resources")
			fmt.Println("  - 'none' to skip (you can protect resources later)")
			fmt.Println()
			fmt.Print("Your choice: ")

			reader := bufio.NewReader(os.Stdin)
			choice, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("⚠️  Could not read input: %v\n", err)
			} else {
				choice = strings.TrimSpace(strings.ToLower(choice))

				var selectedResources []generator.ResourceEntry

				if choice == "all" {
					selectedResources = protectableResources
				} else if choice != "none" && choice != "" {
					// Parse comma-separated numbers, preventing duplicates
					parts := strings.Split(choice, ",")
					seen := make(map[int]bool)
					for _, part := range parts {
						part = strings.TrimSpace(part)
						if num, err := strconv.Atoi(part); err == nil {
							if num >= 1 && num <= len(protectableResources) {
								if !seen[num] {
									seen[num] = true
									selectedResources = append(selectedResources, protectableResources[num-1])
								}
							} else {
								fmt.Printf("⚠️  Invalid selection: %d (out of range)\n", num)
							}
						} else {
							fmt.Printf("⚠️  Invalid selection: %s (not a number)\n", part)
						}
					}
				}

				if len(selectedResources) > 0 {
					fmt.Printf("\nProtecting %d resource(s)...\n", len(selectedResources))
					if err := generator.ProtectResources(wd, genConfig.ModuleName, selectedResources); err != nil {
						fmt.Printf("⚠️  Could not protect resources: %v\n", err)
						fmt.Println("   You can manually wrap routes with authController.RequireAuth()")
					} else {
						fmt.Println("✅ Resources protected!")
						for _, r := range selectedResources {
							fmt.Printf("   - %s (%s)\n", r.Name, r.Path)
						}
					}
				}
			}
		}
	}

	fmt.Println("\n📝 Next steps:")
	fmt.Println("  1. Run migrations:")
	fmt.Println("     lvt migration up")
	fmt.Println("\n  2. Generate sqlc code:")
	fmt.Println("     sqlc generate")
	fmt.Println("\n  3. Configure email sender (see github.com/livetemplate/lvt/pkg/email)")
	fmt.Println("\n  4. Run E2E tests (requires Docker):")
	fmt.Println("     go test ./app/auth -run TestAuthE2E -v")
	fmt.Println("\n💡 Tip: Check app/auth/auth.go for complete usage examples!")

	return nil
}

// startsWithDash checks if a string starts with a dash (flag indicator)
func startsWithDash(s string) bool {
	return strings.HasPrefix(s, "-")
}

// pluralizeNoun converts a singular noun to plural (simple English rules)
func pluralizeNoun(word string) string {
	if word == "" {
		return ""
	}

	lower := strings.ToLower(word)

	// Handle special cases
	if lower == "user" {
		return "users"
	}
	if lower == "account" {
		return "accounts"
	}
	if lower == "admin" {
		return "admins"
	}
	if lower == "member" {
		return "members"
	}

	// General rules
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") ||
		strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}

	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		// Check if preceded by consonant
		secondLast := lower[len(lower)-2]
		if !isVowel(secondLast) {
			return lower[:len(lower)-1] + "ies"
		}
	}

	return lower + "s"
}

// isVowel checks if a character is a vowel
func isVowel(c byte) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}
