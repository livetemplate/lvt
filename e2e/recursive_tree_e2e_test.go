//go:build browser

package e2e_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/livetemplate/livetemplate"
	e2etest "github.com/livetemplate/lvt/testing"
)

// fileTreeNode is a self-referential tree — the canonical recursive-template case
// (a file browser / comment thread / org chart). Rendering it requires the
// {{template "treeNode" .}} call inside {{range .Children}} to be evaluated
// recursively at build time, which is the C8 capability under test.
type fileTreeNode struct {
	Name     string
	Path     string
	IsDir    bool
	Children []fileTreeNode
}

// fileTreeState is the reactive state for the recursive-tree e2e app.
type fileTreeState struct {
	Root   fileTreeNode
	NextID int
}

// fileTreeController drives the recursive-tree app. AddDeepFile appends a new file
// under the /src directory — a node two levels deep — so the click exercises a
// reactive minimal update of a NESTED branch, not just the top level.
type fileTreeController struct{}

// AddDeepFile handles lvt-on:click="addDeepFile": append a file under /src.
func (c *fileTreeController) AddDeepFile(state fileTreeState, _ *livetemplate.Context) (fileTreeState, error) {
	state.NextID++
	name := fmt.Sprintf("new-%d.go", state.NextID)
	for i := range state.Root.Children {
		if state.Root.Children[i].Path == "/src" {
			state.Root.Children[i].Children = append(state.Root.Children[i].Children, fileTreeNode{
				Name: name, Path: "/src/" + name,
			})
		}
	}
	return state, nil
}

// renameDeepNode walks to a node by path and sets its Name (leaving Path — the
// data-key — unchanged), returning true if found.
func renameDeepNode(n *fileTreeNode, path, name string) bool {
	if n.Path == path {
		n.Name = name
		return true
	}
	for i := range n.Children {
		if renameDeepNode(&n.Children[i], path, name) {
			return true
		}
	}
	return false
}

// RenameDeep handles lvt-on:click="renameDeep": rename the deepest node
// (/src/util/hash.go, three levels down) WITHOUT changing its Path. Because the
// data-key stays stable, this is a pure content edit — the diff engine must scope
// it to a nested chain of ["u", key, …] ops down to that single leaf (statics-free,
// no ancestor branch re-sent), and the client must apply it by morphing the
// existing <li> in place rather than rebuilding the branch. This is the payload +
// DOM-preservation win under test.
func (c *fileTreeController) RenameDeep(state fileTreeState, _ *livetemplate.Context) (fileTreeState, error) {
	renameDeepNode(&state.Root, "/src/util/hash.go", "hash-RENAMED.go")
	return state, nil
}

// recursiveTreeE2ETemplate is a full HTML document whose body renders a recursive
// file tree via a self-referential {{define "treeNode"}} and a button that mutates
// a deep branch. Serving a full document exercises wrapper injection around a
// recursive body — the shape a real app ships.
const recursiveTreeE2ETemplate = `<!DOCTYPE html>
<html>
<head><title>Recursive Tree E2E</title></head>
<body>
{{define "treeNode"}}<li data-key="{{.Path}}"><span class="node-name">{{.Name}}</span>{{if .IsDir}}<ul>{{range .Children}}{{template "treeNode" .}}{{end}}</ul>{{end}}</li>{{end}}
<h1>File Tree</h1>
<ul id="tree">{{template "treeNode" .Root}}</ul>
<button type="button" id="add-btn" lvt-on:click="addDeepFile">Add file to /src</button>
<button type="button" id="rename-btn" lvt-on:click="renameDeep">Rename deep file</button>
<script src="/client.js"></script>
</body>
</html>`

// TestRecursiveTemplate_E2E is the black-box browser gate for C8: a recursive
// {{template}} must render its full nested structure on first load AND reactively
// apply a minimal update to a deep branch through the real published client — no
// full-page reload, no console errors. This is the end-to-end counterpart to the
// library-level render/diff unit tests in livetemplate/recursive_template*_test.go.
//
// Observability per the e2e mandate: browser console + server logs + WebSocket
// frames + rendered HTML are all captured and dumped on any failure.
func TestRecursiveTemplate_E2E(t *testing.T) {
	// Concurrency-safe server-log capture (see TestIssue414 for the log.SetOutput
	// caveat). Do NOT add t.Parallel() — global log capture is incompatible with it.
	serverLogs := e2etest.NewSafeBuffer()
	log.SetOutput(serverLogs)
	defer log.SetOutput(os.Stderr)

	controller := &fileTreeController{}
	state := &fileTreeState{
		Root: fileTreeNode{
			Name: "root", Path: "/", IsDir: true,
			Children: []fileTreeNode{
				{Name: "README.md", Path: "/README.md"},
				{Name: "src", Path: "/src", IsDir: true, Children: []fileTreeNode{
					{Name: "main.go", Path: "/src/main.go"},
					{Name: "util", Path: "/src/util", IsDir: true, Children: []fileTreeNode{
						{Name: "hash.go", Path: "/src/util/hash.go"},
					}},
				}},
			},
		},
	}

	tmpl := livetemplate.Must(livetemplate.New("recursive-e2e"))
	if _, err := tmpl.Parse(recursiveTreeE2ETemplate); err != nil {
		t.Fatalf("Parse recursive template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", tmpl.Handle(controller, livetemplate.AsState(state)))
	mux.HandleFunc("/client.js", e2etest.ServeClientLibrary)

	port, err := e2etest.GetFreePort()
	if err != nil {
		t.Fatalf("GetFreePort: %v", err)
	}
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()
	e2etest.WaitForServer(t, fmt.Sprintf("http://localhost:%d", port), 10*time.Second)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			t.Logf("Server shutdown warning: %v", err)
		}
	}()

	chromeCtx, cleanup := e2etest.SetupDockerChrome(t, 30*time.Second)
	defer cleanup()
	ctx := chromeCtx.Context

	consoleLog := e2etest.NewConsoleLogger()
	consoleLog.Start(ctx)
	wsLog := e2etest.RecordWSFrames(ctx)

	var renderedHTML string
	snapshotHTML := func() {
		if err := chromedp.Run(ctx, chromedp.OuterHTML(`html`, &renderedHTML, chromedp.ByQuery)); err != nil {
			t.Logf("snapshotHTML: %v (renderedHTML may be stale)", err)
		}
	}
	dumpDiagnostics := func(label string) {
		t.Logf("=== %s ===", label)
		logs := consoleLog.GetLogs()
		t.Logf("--- BROWSER CONSOLE (%d entries) ---", len(logs))
		for _, e := range logs {
			t.Logf("console [%s]> %s", e.Type, e.Message)
		}
		t.Logf("--- SERVER LOGS ---\n%s", serverLogs.String())
		wsFrames := wsLog.GetMessages()
		t.Logf("--- WEBSOCKET FRAMES (%d) ---", len(wsFrames))
		for _, m := range wsFrames {
			t.Logf("ws %s> %s", m.Direction, m.Data)
		}
		t.Logf("--- RENDERED HTML ---\n%s", renderedHTML)
	}

	url := e2etest.GetChromeTestURL(port)

	// Phase 1: first load must render the FULL nested tree. The deepest node
	// (/src/util/hash.go, three levels down) is the proof recursion evaluated all
	// the way — a flatten-time overflow or a truncated recursion would drop it.
	var deepNodePresent, midNodePresent bool
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#add-btn`, chromedp.ByID),
		// Wait for the framework's "WS connected" signal (data-lvt-loading cleared
		// by the client only after the WS handshake completes).
		e2etest.WaitFor(`(() => {
			const w = document.querySelector('[data-lvt-id]');
			return w && !w.hasAttribute('data-lvt-loading');
		})()`, 5*time.Second),
		chromedp.WaitVisible(`li[data-key="/src/util/hash.go"]`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/util/hash.go"]') !== null`, &deepNodePresent),
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/main.go"]') !== null`, &midNodePresent),
		chromedp.OuterHTML(`html`, &renderedHTML, chromedp.ByQuery),
	); err != nil {
		snapshotHTML()
		dumpDiagnostics("first-render / structural-check phase failed")
		t.Fatalf("chromedp.Run (first render): %v", err)
	}
	if !deepNodePresent {
		dumpDiagnostics("deepest recursive node missing on first render")
		t.Fatalf("recursion did not render the deepest node /src/util/hash.go")
	}
	if !midNodePresent {
		dumpDiagnostics("mid-level recursive node missing on first render")
		t.Fatalf("recursion did not render /src/main.go")
	}

	// Phase 2: click adds a file to the DEEP /src branch. The reactive update must
	// insert exactly that node into the existing DOM (no reload) while leaving the
	// unrelated deep node /src/util/hash.go in place.
	var newNodePresent, deepNodeStillPresent, reloaded bool
	if err := chromedp.Run(ctx,
		chromedp.Click(`#add-btn`, chromedp.ByID),
		chromedp.WaitVisible(`li[data-key="/src/new-1.go"]`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/new-1.go"]') !== null`, &newNodePresent),
		// The sibling deep node must survive the minimal update (tree not rebuilt).
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/util/hash.go"]') !== null`, &deepNodeStillPresent),
		// Sanity: the WS-driven update should not have navigated/reloaded the page.
		chromedp.Evaluate(`window.performance.getEntriesByType('navigation').length === 1`, &reloaded),
		chromedp.OuterHTML(`html`, &renderedHTML, chromedp.ByQuery),
	); err != nil {
		snapshotHTML()
		dumpDiagnostics("reactive-update phase failed")
		t.Fatalf("chromedp.Run (reactive update): %v", err)
	}
	if !newNodePresent {
		dumpDiagnostics("new node not inserted by reactive update")
		t.Fatalf("clicking Add did not insert /src/new-1.go into the recursive tree")
	}
	if !deepNodeStillPresent {
		dumpDiagnostics("deep sibling lost during reactive update")
		t.Fatalf("reactive update tore down the unrelated deep node /src/util/hash.go")
	}
	if !reloaded {
		dumpDiagnostics("page navigated/reloaded instead of a WS update")
		t.Fatalf("expected an in-place WS update, but the page performed a full navigation")
	}

	// The click must round-trip over WebSocket — no frame means delegation broke.
	if len(wsLog.GetMessages()) == 0 {
		dumpDiagnostics("no WebSocket frames captured")
		t.Fatalf("expected WebSocket frames for the reactive update, got none")
	}

	// Phase 3: a DEEP RENAME (content edit on /src/util/hash.go, three levels down,
	// keeping its Path/data-key stable) exercises the per-leaf ["u"] chain. We mark
	// two DOM nodes with a JS property BEFORE the click: the node that will be
	// renamed, and an unaffected sibling (/src/main.go). A per-leaf update morphs the
	// existing <li> in place, so BOTH markers must survive — proof the branch was not
	// rebuilt. A branch re-send would recreate these elements and drop the markers.
	wsBefore := len(wsLog.GetMessages())
	var renamedMarkerSurvived, siblingMarkerSurvived, renamedTextUpdated bool
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(`(() => {
			document.querySelector('li[data-key="/src/util/hash.go"]').__lvtMark = 'RENAMED_NODE';
			document.querySelector('li[data-key="/src/main.go"]').__lvtMark = 'SIBLING_NODE';
			return true;
		})()`, nil),
		chromedp.Click(`#rename-btn`, chromedp.ByID),
		// Wait for the renamed name to appear in the SAME data-key node.
		e2etest.WaitFor(`(() => {
			const li = document.querySelector('li[data-key="/src/util/hash.go"]');
			return li && li.querySelector('.node-name').textContent === 'hash-RENAMED.go';
		})()`, 5*time.Second),
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/util/hash.go"]').querySelector('.node-name').textContent === 'hash-RENAMED.go'`, &renamedTextUpdated),
		// The renamed <li> element itself was morphed in place (marker survives).
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/util/hash.go"]').__lvtMark === 'RENAMED_NODE'`, &renamedMarkerSurvived),
		// The unaffected sibling <li> was untouched (marker survives).
		chromedp.Evaluate(`document.querySelector('li[data-key="/src/main.go"]').__lvtMark === 'SIBLING_NODE'`, &siblingMarkerSurvived),
		chromedp.OuterHTML(`html`, &renderedHTML, chromedp.ByQuery),
	); err != nil {
		snapshotHTML()
		dumpDiagnostics("deep-rename phase failed")
		t.Fatalf("chromedp.Run (deep rename): %v", err)
	}
	if !renamedTextUpdated {
		dumpDiagnostics("deep rename did not update the node text")
		t.Fatalf("deep rename did not update /src/util/hash.go's displayed name")
	}
	if !renamedMarkerSurvived {
		dumpDiagnostics("renamed node was rebuilt, not morphed in place")
		t.Fatalf("per-leaf update must morph the existing <li> in place (marker lost → element replaced)")
	}
	if !siblingMarkerSurvived {
		dumpDiagnostics("unaffected sibling was rebuilt during deep rename")
		t.Fatalf("per-leaf update must not rebuild the unaffected sibling /src/main.go (marker lost)")
	}

	// The rename's WS update must be scoped: the frame(s) it produced must NOT carry
	// unaffected nodes (README.md, main.go, the /src/new-1.go added earlier). A branch
	// re-send would include them; a per-leaf ["u"] chain does not.
	renameFrames := wsLog.GetMessages()[wsBefore:]
	var renameServerToClient string
	for _, m := range renameFrames {
		if m.Direction == "received" || m.Direction == "recv" || m.Direction == "in" {
			renameServerToClient += m.Data
		}
	}
	if renameServerToClient == "" {
		// Direction labels vary; fall back to concatenating all new frames.
		for _, m := range renameFrames {
			renameServerToClient += m.Data
		}
	}
	if strings.Contains(renameServerToClient, "README.md") ||
		strings.Contains(renameServerToClient, "main.go") ||
		strings.Contains(renameServerToClient, "new-1.go") {
		dumpDiagnostics("deep-rename WS update re-sent unaffected nodes (not per-leaf)")
		t.Fatalf("deep-rename update must be scoped to the changed leaf; frame carried unaffected nodes:\n%s", renameServerToClient)
	}
	if !strings.Contains(renameServerToClient, "hash-RENAMED.go") {
		dumpDiagnostics("deep-rename WS frame missing the renamed value")
		t.Fatalf("expected the deep-rename WS frame to carry 'hash-RENAMED.go', got:\n%s", renameServerToClient)
	}

	// No browser-console errors (a client-side apply failure surfaces here).
	for _, e := range consoleLog.GetLogs() {
		if strings.EqualFold(e.Type, "error") {
			dumpDiagnostics("browser console error during recursive render/update")
			t.Fatalf("browser console error: %s", e.Message)
		}
	}
}
