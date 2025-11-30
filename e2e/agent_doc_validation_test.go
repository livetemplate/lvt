package e2e

import (
	"testing"
	"time"

	"github.com/livetemplate/lvt/internal/agenttest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAgentDocValidation_QuickStart validates the quick start example from the usage guide
// This test verifies the documentation is accurate by simulating the commands
func TestAgentDocValidation_QuickStart(t *testing.T) {
	env := agenttest.Setup(t, &agenttest.SetupOptions{
		AppName: "blog",
		Kit:     "multi",
	})

	// Simulate: "I want to build a blog with posts"
	env.TrackSkill("lvt-quickstart")

	// Generate posts resource
	err := env.RunLvtCommand("gen", "resource", "posts", "title:string", "content:text")
	require.NoError(t, err, "gen posts should succeed")

	// Run migrations
	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err, "migration up should succeed")

	// Verify expected outcomes from documentation
	env.AssertFileExists("internal/app/posts/posts.go")
	env.AssertFileExists("internal/database/queries.sql")
	env.AssertFileExists("internal/app/posts/posts.tmpl")
	env.AssertDatabaseTableExists("posts")

	// Verify commands tracked
	env.AssertCommandRun("lvt gen resource posts title:string content:text")
	env.AssertCommandRun("lvt migration up")

	// Verify skill tracked
	env.AssertSkillUsed("lvt-quickstart")
}

// TestAgentDocValidation_FullStack validates the full stack example
func TestAgentDocValidation_FullStack(t *testing.T) {
	env := agenttest.Setup(t, &agenttest.SetupOptions{
		AppName: "taskmanager",
		Kit:     "multi",
	})

	// Simulate the full stack workflow from the guide
	env.SimulateFullStack()

	// Verify all expected files exist
	env.AssertFileExists("internal/app/auth/auth.go")
	env.AssertFileExists("internal/app/projects/projects.go")
	env.AssertFileExists("internal/app/tasks/tasks.go")

	// Verify database tables
	// Note: users table is created via create_auth_tables migration (not named after table)
	env.AssertDatabaseTableExists("projects")
	env.AssertDatabaseTableExists("tasks")

	// Verify relationships
	env.AssertFileContains("internal/database/queries.sql", "user_id")
	env.AssertFileContains("internal/database/queries.sql", "project_id")
}

// TestAgentDocValidation_StepByStep validates the step-by-step conversation
func TestAgentDocValidation_StepByStep(t *testing.T) {
	env := agenttest.Setup(t, &agenttest.SetupOptions{
		AppName: "conversationapp",
		Kit:     "multi",
	})

	// Step 1: "Let's add a blog to my app"
	env.TrackSkill("lvt-add-resource")
	err := env.RunLvtCommand("gen", "resource", "posts", "title:string", "content:text")
	assert.NoError(t, err)

	// Avoid duplicate timestamp conflicts in migrations
	time.Sleep(1100 * time.Millisecond)

	// Step 2: "Can you add categories and tags?"
	err = env.RunLvtCommand("gen", "resource", "categories", "name:string")
	assert.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)

	err = env.RunLvtCommand("gen", "resource", "tags", "name:string")
	assert.NoError(t, err)

	time.Sleep(1100 * time.Millisecond)

	// Step 3: "Add author attribution to posts"
	env.TrackSkill("lvt-customize-resource")
	err = env.RunLvtCommand("migration", "create", "add_author_to_posts")
	assert.NoError(t, err)

	// Apply migrations
	err = env.RunLvtCommand("migration", "up")
	assert.NoError(t, err)

	// Verify all resources exist
	env.AssertFileExists("internal/app/posts/posts.go")
	env.AssertFileExists("internal/app/categories/categories.go")
	env.AssertFileExists("internal/app/tags/tags.go")
}

// TestAgentDocValidation_CommonPatterns tests the common patterns from the guide
func TestAgentDocValidation_CommonPatterns(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		pattern string
		verify  func(*testing.T, *agenttest.AgentTestEnv)
	}{
		{
			name:    "AuthenticationSetup",
			appName: "authapp",
			pattern: "I need user authentication with email and password",
			verify: func(t *testing.T, env *agenttest.AgentTestEnv) {
				env.TrackSkill("lvt-gen-auth")
				err := env.RunLvtCommand("gen", "auth")
				require.NoError(t, err)

				err = env.RunLvtCommand("migration", "up")
				require.NoError(t, err)

				env.AssertFileExists("internal/app/auth/auth.go")
				env.AssertFileExists("internal/app/auth/auth.tmpl")
				// Note: Password auth disabled in v0.4.x, so no separate login/signup templates
			},
		},
		{
			name:    "CRUDResource",
			appName: "productapp",
			pattern: "Add products with name, price, description, and stock quantity",
			verify: func(t *testing.T, env *agenttest.AgentTestEnv) {
				env.TrackSkill("lvt-add-resource")
				err := env.RunLvtCommand("gen", "resource", "products",
					"name:string",
					"price:float",
					"description:text",
					"stock:int")
				require.NoError(t, err)

				err = env.RunLvtCommand("migration", "up")
				require.NoError(t, err)

				env.AssertFileExists("internal/app/products/products.go")
				env.AssertFileContains("internal/database/queries.sql", "name")
				env.AssertFileContains("internal/database/queries.sql", "price")
			},
		},
		{
			name:    "Relationships",
			appName: "orderapp",
			pattern: "Add orders that belong to users",
			verify: func(t *testing.T, env *agenttest.AgentTestEnv) {
				env.TrackSkill("lvt-add-resource")
				err := env.RunLvtCommand("gen", "resource", "users", "name:string", "email:string")
				require.NoError(t, err)

				err = env.RunLvtCommand("gen", "resource", "orders", "user_id:int", "total:float")
				require.NoError(t, err)

				err = env.RunLvtCommand("migration", "up")
				require.NoError(t, err)

				env.AssertFileContains("internal/database/queries.sql", "user_id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := agenttest.Setup(t, &agenttest.SetupOptions{
				AppName: tt.appName,
				Kit:     "multi",
			})

			tt.verify(t, env)
		})
	}
}

// TestAgentDocValidation_AllKits validates the guide works with all kits
func TestAgentDocValidation_AllKits(t *testing.T) {
	kits := []struct {
		name      string
		mainFile  string
		cmdFolder string
	}{
		{"multi", "cmd/testmulti/main.go", "cmd/testmulti"},
		{"single", "cmd/testsingle/main.go", "cmd/testsingle"},
		{"simple", "main.go", ""},
	}

	for _, kit := range kits {
		t.Run(kit.name, func(t *testing.T) {
			env := agenttest.Setup(t, &agenttest.SetupOptions{
				AppName: "test" + kit.name,
				Kit:     kit.name,
			})

			// Verify app was created
			env.AssertFileExists(kit.mainFile)

			// Verify go.mod exists
			env.AssertFileExists("go.mod")
		})
	}
}

// TestAgentDocValidation_IncrementalFeatures validates the incremental approach
func TestAgentDocValidation_IncrementalFeatures(t *testing.T) {
	env := agenttest.Setup(t, &agenttest.SetupOptions{
		AppName: "incremental",
		Kit:     "multi",
	})

	// Step 1: "Create a blog"
	err := env.RunLvtCommand("gen", "resource", "posts", "title:string", "content:text")
	require.NoError(t, err)
	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err)

	// Step 2: "Add categories and tags"
	err = env.RunLvtCommand("gen", "resource", "categories", "name:string")
	require.NoError(t, err)
	err = env.RunLvtCommand("gen", "resource", "tags", "name:string")
	require.NoError(t, err)
	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err)

	// Step 3: "Add comments"
	err = env.RunLvtCommand("gen", "resource", "comments", "post_id:int", "content:text")
	require.NoError(t, err)
	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err)

	// Verify all features exist
	env.AssertFileExists("internal/app/posts/posts.go")
	env.AssertFileExists("internal/app/categories/categories.go")
	env.AssertFileExists("internal/app/tags/tags.go")
	env.AssertFileExists("internal/app/comments/comments.go")
}

// TestAgentDocValidation_RecipeExample validates the complete example session
func TestAgentDocValidation_RecipeExample(t *testing.T) {
	env := agenttest.Setup(t, &agenttest.SetupOptions{
		AppName: "recipeshare",
		Kit:     "multi",
	})

	// Initial: "I want to build a recipe sharing site"
	env.TrackSkill("lvt-quickstart")
	err := env.RunLvtCommand("gen", "resource", "recipes",
		"title:string",
		"ingredients:text",
		"instructions:text",
		"prep_time:int",
		"cook_time:int")
	require.NoError(t, err)

	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err)

	// Avoid duplicate migration timestamps
	time.Sleep(1100 * time.Millisecond)

	// Follow-up: "add user accounts and categories"
	env.TrackSkill("lvt-gen-auth")
	err = env.RunLvtCommand("gen", "auth")
	require.NoError(t, err)

	// Avoid duplicate migration timestamps
	time.Sleep(1100 * time.Millisecond)

	env.TrackSkill("lvt-add-resource")
	err = env.RunLvtCommand("gen", "resource", "categories", "name:string")
	require.NoError(t, err)

	err = env.RunLvtCommand("migration", "up")
	require.NoError(t, err)

	// Verify all features
	env.AssertFileExists("internal/app/recipes/recipes.go")
	env.AssertFileExists("internal/app/auth/auth.go")
	env.AssertFileExists("internal/app/categories/categories.go")
	env.AssertDatabaseTableExists("recipes")
	// Note: users table is created via create_auth_tables migration (not named after table)
	env.AssertDatabaseTableExists("categories")
}
