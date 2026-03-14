package docker

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livetemplate/lvt/internal/stack"
)

// helperGenerate creates a temp dir structure (project root + deploy subdir) and runs Generate.
// Returns (projectRoot, deployDir).
func helperGenerate(t *testing.T, config stack.StackConfig) (string, string) {
	t.Helper()
	projectRoot := t.TempDir()
	deployDir := filepath.Join(projectRoot, "deploy")
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		t.Fatalf("failed to create deploy dir: %v", err)
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	return projectRoot, deployDir
}

func TestGenerator_Generate(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	// Simple SQLite config should NOT generate docker-compose.yml
	expectedFiles := []string{
		"Dockerfile",
		".dockerignore",
		".env.example",
		"README.md",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(deployDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", file)
		}
	}

	// docker-compose.yml should NOT exist for simple SQLite config
	composePath := filepath.Join(deployDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err == nil {
		t.Error("docker-compose.yml should NOT be generated for simple SQLite config")
	}
}

func TestGenerator_Generate_SimpleNoCompose(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	composePath := filepath.Join(deployDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); err == nil {
		t.Error("docker-compose.yml should not exist for simple SQLite config without backup or redis")
	}
}

func TestGenerator_Generate_ComplexWithCompose(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabasePostgres,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	composePath := filepath.Join(deployDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Error("docker-compose.yml should exist for postgres config")
	}
}

func TestGenerator_Generate_WithLitestream(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupLitestream,
		Storage:  stack.StorageS3,
		Redis:    stack.RedisNone,
		CI:       stack.CINone,
	})

	// Check litestream.yml exists
	litestreamPath := filepath.Join(deployDir, "litestream.yml")
	if _, err := os.Stat(litestreamPath); os.IsNotExist(err) {
		t.Error("Expected litestream.yml does not exist")
	}

	// Litestream config needs compose
	composePath := filepath.Join(deployDir, "docker-compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Error("docker-compose.yml should exist for litestream config")
	}
}

func TestDockerfile_NoCGO(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}
	text := string(content)

	if strings.Contains(text, "CGO_ENABLED=1") {
		t.Error("Dockerfile should not contain CGO_ENABLED=1")
	}
	if strings.Contains(text, "gcc") {
		t.Error("Dockerfile should not reference gcc")
	}
	if strings.Contains(text, "musl-dev") {
		t.Error("Dockerfile should not reference musl-dev")
	}
	if strings.Contains(text, "sqlite-dev") {
		t.Error("Dockerfile should not reference sqlite-dev")
	}
	if strings.Contains(text, "sqlite-libs") {
		t.Error("Dockerfile should not reference sqlite-libs")
	}
	if !strings.Contains(text, "CGO_ENABLED=0") {
		t.Error("Dockerfile should contain CGO_ENABLED=0")
	}
}

func TestDockerfile_NoSqlcDownload(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}
	text := string(content)

	if strings.Contains(text, "sqlc") {
		t.Error("Dockerfile should not reference sqlc")
	}
	if strings.Contains(text, "curl") {
		t.Error("Dockerfile should not reference curl")
	}
}

func TestDockerfile_SpecificCopyDirectives(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}
	text := string(content)

	// Should have specific COPY for binary and app dir
	if !strings.Contains(text, "COPY --from=builder /app/main .") {
		t.Error("Dockerfile should contain specific COPY for binary")
	}
	if !strings.Contains(text, "COPY --from=builder /app/app ./app") {
		t.Error("Dockerfile should contain specific COPY for app directory")
	}

	// Should NOT have the broad 'COPY --from=builder /app .' pattern
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "COPY --from=builder /app ." {
			t.Error("Dockerfile should not contain broad 'COPY --from=builder /app .' directive")
		}
	}

	// Should NOT have rm -rf cleanup
	if strings.Contains(text, "rm -rf") {
		t.Error("Dockerfile should not contain rm -rf cleanup of build artifacts")
	}
}

func TestDockerfile_PinnedAlpine(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}
	text := string(content)

	if strings.Contains(text, "alpine:latest") {
		t.Error("Dockerfile should pin alpine version, not use alpine:latest")
	}
	if !strings.Contains(text, "alpine:3.21") {
		t.Error("Dockerfile should use alpine:3.21")
	}
	if !strings.Contains(text, "tzdata") {
		t.Error("Dockerfile should include tzdata")
	}
}

func TestDockerComposeYAML_NoDuplicateVolumes(t *testing.T) {
	tests := []struct {
		name     string
		database stack.DatabaseType
		backup   stack.BackupType
		storage  stack.StorageType
	}{
		{
			name:     "postgres only",
			database: stack.DatabasePostgres,
			backup:   stack.BackupNone,
		},
		{
			name:     "sqlite with litestream",
			database: stack.DatabaseSQLite,
			backup:   stack.BackupLitestream,
			storage:  stack.StorageS3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := tt.storage
			if storage == "" {
				storage = stack.StorageNone
			}

			_, deployDir := helperGenerate(t, stack.StackConfig{
				Provider: stack.ProviderDocker,
				Database: tt.database,
				Backup:   tt.backup,
				Redis:    stack.RedisNone,
				Storage:  storage,
				CI:       stack.CINone,
			})

			composeFile := filepath.Join(deployDir, "docker-compose.yml")
			content, err := os.ReadFile(composeFile)
			if err != nil {
				t.Fatalf("Failed to read docker-compose.yml: %v", err)
			}

			lines := strings.Split(string(content), "\n")
			topLevelVolumesCount := 0
			for _, line := range lines {
				if line == "volumes:" {
					topLevelVolumesCount++
				}
			}

			if topLevelVolumesCount > 1 {
				t.Errorf("Found %d top-level 'volumes:' declarations, expected at most 1. Content:\n%s", topLevelVolumesCount, string(content))
			}
		})
	}
}

func TestDockerCompose_CorrectBuildContext(t *testing.T) {
	_, deployDir := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabasePostgres,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(deployDir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("Failed to read docker-compose.yml: %v", err)
	}
	text := string(content)

	if strings.Contains(text, "deploy/docker/Dockerfile") {
		t.Error("docker-compose.yml should reference deploy/Dockerfile, not deploy/docker/Dockerfile")
	}
	if !strings.Contains(text, "deploy/Dockerfile") {
		t.Error("docker-compose.yml should reference deploy/Dockerfile")
	}
}

func TestGenerator_Generate_MakefileGenerated(t *testing.T) {
	projectRoot, _ := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	makefilePath := filepath.Join(projectRoot, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Fatal("Makefile should exist at project root")
	}

	content, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}
	text := string(content)

	expectedTargets := []string{"build:", "run:", "stop:", "logs:", "clean:"}
	for _, target := range expectedTargets {
		if !strings.Contains(text, target) {
			t.Errorf("Makefile should contain target %s", target)
		}
	}

	if !strings.Contains(text, "deploy/Dockerfile") {
		t.Error("Makefile should reference deploy/Dockerfile")
	}
}

func TestGenerator_Generate_MakefileSQLiteBackup(t *testing.T) {
	projectRoot, _ := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(projectRoot, "Makefile"))
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "backup:") {
		t.Error("Makefile for SQLite config should contain backup target")
	}
}

func TestGenerator_Generate_MakefileNoBackupWithoutSQLite(t *testing.T) {
	projectRoot, _ := helperGenerate(t, stack.StackConfig{
		Provider: stack.ProviderDocker,
		Database: stack.DatabasePostgres,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	})

	content, err := os.ReadFile(filepath.Join(projectRoot, "Makefile"))
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}
	text := string(content)

	if strings.Contains(text, "backup:") {
		t.Error("Makefile for non-SQLite config should NOT contain backup target")
	}
}
