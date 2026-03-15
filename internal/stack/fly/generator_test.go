package fly

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livetemplate/lvt/internal/stack"
)

func TestGenerator_Generate(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider: stack.ProviderFly,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	}

	gen := New()
	ctx := context.Background()

	err := gen.Generate(ctx, config, deployDir)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// fly.toml should be at project root, not in deploy/
	if _, err := os.Stat(filepath.Join(tmpDir, "fly.toml")); os.IsNotExist(err) {
		t.Error("Expected fly.toml at project root")
	}
	if _, err := os.Stat(filepath.Join(deployDir, "fly.toml")); err == nil {
		t.Error("fly.toml should NOT be in deploy/ directory")
	}

	// Other files should be in deploy/
	for _, file := range []string{"Dockerfile", ".env.example", "README.md"} {
		path := filepath.Join(deployDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected deploy/%s does not exist", file)
		}
	}
}

func TestGenerator_FlyToml_Content(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider: stack.ProviderFly,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "fly.toml"))
	if err != nil {
		t.Fatalf("Failed to read fly.toml: %v", err)
	}

	checks := []struct {
		name     string
		contains string
	}{
		{"dockerfile path", `dockerfile = "deploy/Dockerfile"`},
		{"health check path", `/health/live`},
		{"force https", `force_https = true`},
		{"auto stop format", `auto_stop_machines = "stop"`},
		{"mounts for sqlite", `[mounts]`},
		{"data destination", `destination = "/data"`},
		{"initial volume size", `initial_size = "1GB"`},
		{"env database path", `DATABASE_PATH = "/data/app.db"`},
		{"sqlite memory 256mb", `memory = "256mb"`},
	}

	for _, c := range checks {
		if !bytes.Contains(content, []byte(c.contains)) {
			t.Errorf("fly.toml missing %s: expected to contain %q", c.name, c.contains)
		}
	}

	// SQLite should NOT have 512mb
	if bytes.Contains(content, []byte(`memory = "512mb"`)) {
		t.Error("SQLite fly.toml should use 256mb, not 512mb")
	}
}

func TestGenerator_Dockerfile_NoCGO(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider: stack.ProviderFly,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}

	// Must NOT contain CGO_ENABLED=1 or C compiler dependencies
	forbidden := []string{
		"CGO_ENABLED=1",
		"gcc",
		"musl-dev",
		"sqlite-dev",
		"sqlite-libs",
		"sqlc",
		"go mod tidy",
	}
	for _, f := range forbidden {
		if bytes.Contains(content, []byte(f)) {
			t.Errorf("Dockerfile should NOT contain %q (modernc.org/sqlite is pure Go)", f)
		}
	}

	// Must contain CGO_ENABLED=0
	if !bytes.Contains(content, []byte("CGO_ENABLED=0")) {
		t.Error("Dockerfile must use CGO_ENABLED=0")
	}

	// Must use pinned alpine version
	if bytes.Contains(content, []byte("alpine:latest")) {
		t.Error("Dockerfile should pin alpine version, not use :latest")
	}
	if !bytes.Contains(content, []byte("alpine:3.21")) {
		t.Error("Dockerfile should use alpine:3.21")
	}

	// Must have specific COPY directives
	if !bytes.Contains(content, []byte("COPY --from=builder /app/main .")) {
		t.Error("Dockerfile must COPY binary specifically")
	}
	if !bytes.Contains(content, []byte("COPY --from=builder /app/app ./app")) {
		t.Error("Dockerfile must COPY app/ templates specifically")
	}

	// Builder stage must ensure directories exist to prevent COPY failures
	if !bytes.Contains(content, []byte("mkdir -p database/migrations app")) {
		t.Error("Dockerfile builder stage must mkdir -p database/migrations app for safety")
	}

	// Must NOT have the broad copy-then-delete pattern
	if bytes.Contains(content, []byte("rm -rf")) {
		t.Error("Dockerfile should not use copy-then-delete pattern")
	}
}

func TestGenerator_Generate_WithLitestream(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider: stack.ProviderFly,
		Database: stack.DatabaseSQLite,
		Backup:   stack.BackupLitestream,
		Storage:  stack.StorageS3,
		Redis:    stack.RedisNone,
		CI:       stack.CINone,
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check litestream.yml exists
	if _, err := os.Stat(filepath.Join(deployDir, "litestream.yml")); os.IsNotExist(err) {
		t.Error("Expected litestream.yml does not exist")
	}

	// Check Dockerfile contains Litestream with checksum verification
	content, err := os.ReadFile(filepath.Join(deployDir, "Dockerfile"))
	if err != nil {
		t.Fatalf("Failed to read Dockerfile: %v", err)
	}

	// Litestream path must include wget in apk install
	if !bytes.Contains(content, []byte("apk --no-cache add ca-certificates tzdata wget")) {
		t.Error("Dockerfile with Litestream must install wget for downloading Litestream binary")
	}

	amd64Checksum := "eb75a3de5cab03875cdae9f5f539e6aedadd66607003d9b1e7a9077948818ba0"
	arm64Checksum := "9585f5a508516bd66af2b2376bab4de256a5ef8e2b73ec760559e679628f2d59"
	if !bytes.Contains(content, []byte(amd64Checksum)) {
		t.Error("Dockerfile does not contain expected amd64 SHA256 checksum for Litestream v0.3.13")
	}
	if !bytes.Contains(content, []byte(arm64Checksum)) {
		t.Error("Dockerfile does not contain expected arm64 SHA256 checksum for Litestream v0.3.13")
	}
	if !bytes.Contains(content, []byte("sha256sum -c -")) {
		t.Error("Dockerfile does not contain SHA256 verification command")
	}
	if !bytes.Contains(content, []byte(`CMD ["litestream", "replicate"`)) {
		t.Error("Dockerfile should use litestream replicate in exec form as CMD")
	}
	if !bytes.Contains(content, []byte("COPY deploy/litestream.yml /etc/litestream.yml")) {
		t.Error("Dockerfile should COPY deploy/litestream.yml to /etc/litestream.yml")
	}
}

func TestGenerator_Generate_Postgres(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider: stack.ProviderFly,
		Database: stack.DatabasePostgres,
		Backup:   stack.BackupNone,
		Redis:    stack.RedisNone,
		Storage:  stack.StorageNone,
		CI:       stack.CINone,
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// fly.toml should exist at project root
	flyToml, err := os.ReadFile(filepath.Join(tmpDir, "fly.toml"))
	if err != nil {
		t.Fatalf("Expected fly.toml at project root: %v", err)
	}

	// Should NOT have mounts section (postgres doesn't need local storage)
	if bytes.Contains(flyToml, []byte("[mounts]")) {
		t.Error("PostgreSQL config should not have [mounts] section")
	}

	// Should NOT have DATABASE_PATH
	if bytes.Contains(flyToml, []byte("DATABASE_PATH")) {
		t.Error("PostgreSQL config should not have DATABASE_PATH")
	}

	// PostgreSQL should get 512mb memory
	if !bytes.Contains(flyToml, []byte(`memory = "512mb"`)) {
		t.Error("PostgreSQL fly.toml should use 512mb memory")
	}
}

func TestGenerator_Generate_MultiRegion(t *testing.T) {
	tmpDir := t.TempDir()
	deployDir := filepath.Join(tmpDir, "deploy")

	config := stack.StackConfig{
		Provider:    stack.ProviderFly,
		Database:    stack.DatabaseSQLite,
		Backup:      stack.BackupNone,
		Redis:       stack.RedisNone,
		Storage:     stack.StorageNone,
		CI:          stack.CINone,
		MultiRegion: true,
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "fly.toml"))
	if err != nil {
		t.Fatalf("Expected fly.toml at project root: %v", err)
	}

	if !bytes.Contains(content, []byte("Multi-region")) {
		t.Error("Multi-region config should contain multi-region comment")
	}
}

func TestGenerator_Generate_WithProjectDir(t *testing.T) {
	// Test the explicit ProjectDir path (vs. the filepath.Dir fallback).
	// deployDir is intentionally two levels below projectRoot so that
	// filepath.Dir(deployDir) != projectRoot, verifying that the generator
	// uses config.ProjectDir rather than deriving the root itself.
	tmpDir := t.TempDir()
	projectRoot := filepath.Join(tmpDir, "myproject")
	deployDir := filepath.Join(projectRoot, "nested", "deploy")

	if err := os.MkdirAll(projectRoot, 0755); err != nil {
		t.Fatalf("mkdir project root: %v", err)
	}

	config := stack.StackConfig{
		Provider:   stack.ProviderFly,
		Database:   stack.DatabaseSQLite,
		Backup:     stack.BackupNone,
		Redis:      stack.RedisNone,
		Storage:    stack.StorageNone,
		CI:         stack.CINone,
		ProjectDir: projectRoot, // explicitly set
	}

	gen := New()
	if err := gen.Generate(context.Background(), config, deployDir); err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// fly.toml should be at ProjectDir, not filepath.Dir(deployDir)
	if _, err := os.Stat(filepath.Join(projectRoot, "fly.toml")); os.IsNotExist(err) {
		t.Error("Expected fly.toml at ProjectDir (project root)")
	}

	// It should NOT be at filepath.Dir(deployDir) which is nested/
	nestedDir := filepath.Dir(deployDir)
	if nestedDir != projectRoot {
		if _, err := os.Stat(filepath.Join(nestedDir, "fly.toml")); err == nil {
			t.Error("fly.toml should NOT be at filepath.Dir(outputDir) when ProjectDir is set")
		}
	}
}
