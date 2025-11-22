package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/livetemplate/lvt/internal/generator"
)

func New(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("app name required")
	}

	appName := args[0]
	moduleName := appName // Default to app name
	devMode := false      // Default to production (use CDN)
	kit := "multi"        // Default kit

	// Check for flags
	for i := 1; i < len(args); i++ {
		if args[i] == "--module" && i+1 < len(args) {
			moduleName = args[i+1]
			i++ // Skip next arg
		} else if args[i] == "--dev" {
			devMode = true
		} else if args[i] == "--kit" && i+1 < len(args) {
			kit = args[i+1]
			i++ // Skip next arg
		}
	}

	// Validate kit
	validKits := map[string]bool{"multi": true, "single": true, "simple": true}
	if !validKits[kit] {
		return fmt.Errorf("invalid kit: %s (valid: multi, single, simple)", kit)
	}

	fmt.Printf("Creating new LiveTemplate app: %s\n", appName)
	fmt.Printf("Kit: %s\n", kit)
	if devMode {
		fmt.Println("Mode: Development (using local client library)")
	}

	// Check if we're inside another Go module
	isNested := false
	if _, err := os.Stat("go.mod"); err == nil {
		isNested = true
	}

	if err := generator.GenerateApp(appName, moduleName, kit, devMode); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ App created successfully!")

	if isNested {
		fmt.Println()
		fmt.Println("⚠️  Warning: Creating app inside another Go module")
		if kit == "simple" {
			fmt.Printf("   You'll need to use: GOWORK=off go run main.go\n")
		} else {
			fmt.Printf("   You'll need to use: GOWORK=off go run cmd/%s/main.go\n", appName)
		}
		fmt.Println("   For production, create apps outside Go module directories")
	}

	fmt.Println()

	// Skip go mod tidy when running under test to avoid "Test I/O incomplete" errors
	// Tests handle go mod tidy separately with proper synchronization
	if os.Getenv("SKIP_GO_MOD_TIDY") != "1" {
		// Run go mod tidy to resolve and download dependencies
		fmt.Println("Installing dependencies...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = appName
		// Disable workspace mode to prevent background processes
		cmd.Env = append(os.Environ(), "GOWORK=off")

		// Use CombinedOutput() instead of Run() to properly close pipes
		// This prevents "Test I/O incomplete" errors when running under test frameworks
		// Note: CombinedOutput captures stdout/stderr internally, so we don't set cmd.Stdout/Stderr
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("⚠️  Warning: failed to install dependencies: %v\n", err)
			if len(output) > 0 {
				fmt.Printf("   Output: %s\n", string(output))
			}
			fmt.Printf("   You can run it manually: cd %s && go mod tidy\n", appName)
		} else {
			// Print output if there was any (warnings, etc.)
			if len(output) > 0 {
				fmt.Print(string(output))
			}
			fmt.Println("✅ Dependencies installed!")
		}
	} else {
		fmt.Println("⏭️  Skipping dependency installation (test mode)")
	}

	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", appName)

	// Different instructions based on kit type
	if kit == "simple" {
		fmt.Println("  go run main.go")
		fmt.Println()
		fmt.Println("Then open http://localhost:8080 in your browser")
		fmt.Println()
		fmt.Println("Edit main.go to customize your app logic")
		fmt.Printf("Edit %s.tmpl to modify the UI\n", appName)
	} else {
		fmt.Println("  lvt gen users name:string email:string")
		fmt.Println("  lvt migration up")
		fmt.Printf("  go run cmd/%s/main.go\n", appName)
	}
	fmt.Println()

	return nil
}
