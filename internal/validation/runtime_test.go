package validation

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/livetemplate/lvt/internal/validator"
)

// validHTTPApp is a minimal Go app that reads PORT and serves HTTP.
const validHTTPApp = `package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.ListenAndServe(":"+port, nil)
}
`

func TestRuntimeCheck_ValidApp(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", validHTTPApp)

	c := &RuntimeCheck{StartupTimeout: 10 * time.Second}
	result := c.Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}
}

func TestRuntimeCheck_BuildFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", "package main\n\nfunc main() {\n\tundefined()\n}\n")

	c := &RuntimeCheck{}
	result := c.Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for code that doesn't compile")
	}
	assertHasError(t, result, "build failed")
}

func TestRuntimeCheck_StartupTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	// App that compiles but never listens on any port.
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", `package main

import "time"

func main() {
	time.Sleep(time.Hour)
}
`)

	c := &RuntimeCheck{StartupTimeout: 1 * time.Second}
	result := c.Run(context.Background(), dir)

	if result.Valid {
		t.Error("expected invalid for app that never starts listening")
	}
	assertHasError(t, result, "did not start")
}

func TestRuntimeCheck_RouteProbing(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	// App that returns 500 on /bad and 200 on /.
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", `package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error")
	})
	http.HandleFunc("/notfound", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "not found")
	})
	http.ListenAndServe(":"+port, nil)
}
`)

	c := &RuntimeCheck{
		StartupTimeout: 10 * time.Second,
		Routes:         []string{"/", "/bad", "/notfound"},
	}
	result := c.Run(context.Background(), dir)

	// Should have errors for /bad (5xx) and warnings for /notfound (4xx).
	if result.Valid {
		t.Error("expected invalid due to 500 on /bad")
	}

	var has5xx, has4xx bool
	for _, issue := range result.Issues {
		if issue.Level == validator.LevelError && strings.Contains(issue.Message, "/bad") && strings.Contains(issue.Message, "500") {
			has5xx = true
		}
		if issue.Level == validator.LevelWarning && strings.Contains(issue.Message, "/notfound") && strings.Contains(issue.Message, "404") {
			has4xx = true
		}
	}
	if !has5xx {
		t.Error("expected error for /bad returning 500")
	}
	if !has4xx {
		t.Error("expected warning for /notfound returning 404")
	}
}

func TestRuntimeCheck_ProcessCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	// App that runs forever — we verify the process is killed after Run returns.
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", validHTTPApp)

	// We need to capture the PID. Run the check, then verify no process at that PID.
	// Since RuntimeCheck handles cleanup internally, we verify by checking
	// the binary is removed (deferred os.Remove) and the port is free.
	port, err := getFreePort()
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}

	c := &RuntimeCheck{
		StartupTimeout: 10 * time.Second,
		Port:           port,
	}
	result := c.Run(context.Background(), dir)

	if !result.Valid {
		t.Errorf("expected valid, got: %s", result.Format())
	}

	// After Run returns, the port should be free again.
	// Try to listen on it to confirm the process was cleaned up.
	verifyPortFree(t, port)

	// Binary should have been removed.
	binaryPath := fmt.Sprintf("%s/lvt-runtime-check", dir)
	if _, err := os.Stat(binaryPath); !os.IsNotExist(err) {
		t.Error("expected binary to be removed after runtime check")
	}
}

// verifyPortFree checks that nothing is listening on the given port.
func verifyPortFree(t *testing.T, port int) {
	t.Helper()
	// Small delay to allow OS to reclaim the port.
	time.Sleep(100 * time.Millisecond)

	ln, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Logf("could not create socket to verify port: %v", err)
		return
	}
	defer syscall.Close(ln)

	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], []byte{127, 0, 0, 1})
	if err := syscall.Bind(ln, sa); err != nil {
		t.Errorf("port %d still in use after runtime check cleanup", port)
	}
}

func TestRuntimeCheck_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("runtime check runs external commands")
	}

	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/testapp\n\ngo 1.21\n")
	writeFile(t, dir, "main.go", validHTTPApp)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	c := &RuntimeCheck{StartupTimeout: 5 * time.Second}
	result := c.Run(ctx, dir)

	if result.Valid {
		t.Error("expected invalid due to context cancellation")
	}
}
