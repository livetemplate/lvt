package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newChildData() ResourceData {
	return ResourceData{
		PackageName:            "comments",
		ModuleName:             "testapp",
		ResourceName:           "Comments",
		ResourceNameLower:      "comments",
		ResourceNameSingular:   "Comment",
		ResourceNamePlural:     "Comments",
		TableName:              "comments",
		ParentResource:         "posts",
		ParentPackageName:      "posts",
		ParentResourceSingular: "Post",
		ParentReferenceField:   "post_id",
		IsEmbedded:             true,
		Fields: []FieldData{
			{Name: "post_id", GoType: "string", IsReference: true, ReferencedTable: "posts"},
			{Name: "author", GoType: "string"},
			{Name: "text", GoType: "string"},
		},
	}
}

const parentGoTemplate = `package posts

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/livetemplate/livetemplate"
	"testapp/database/models"
)

type PostsItem = models.Post

type PostsController struct {
	Queries *models.Queries
}

type PostsState struct {
	Title        string        ` + "`" + `json:"title"` + "`" + `
	EditingID    string        ` + "`" + `json:"editing_id"` + "`" + `
	EditingPosts *PostsItem    ` + "`" + `json:"editing_posts"` + "`" + `
	LastUpdated  string        ` + "`" + `json:"last_updated"` + "`" + `
}

func (c *PostsController) View(state PostsState, ctx *livetemplate.Context) (PostsState, error) {
	state.EditingID = "test-id"
	state.LastUpdated = formatTime()
	return state, nil
}

func (c *PostsController) Mount(state PostsState, _ *livetemplate.Context) (PostsState, error) {
	return c.loadPostss(state, context.Background())
}

func (c *PostsController) loadPostss(state PostsState, ctx context.Context) (PostsState, error) {
	return state, nil
}

func formatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Handler(queries *models.Queries) http.Handler {
	controller := &PostsController{
		Queries: queries,
	}

	initialState := &PostsState{
		Title:       "Posts Management",
		LastUpdated: formatTime(),
	}

	baseTmpl := livetemplate.Must(livetemplate.New("posts",
		livetemplate.WithDevMode(false),
	))
	if _, err := baseTmpl.ParseFiles("app/posts/posts.tmpl"); err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseTmpl.Handle(controller, livetemplate.AsState(initialState)).ServeHTTP(w, r)
	})
}
`

const parentTmplTemplate = `{{define "detailPage"}}
  {{if .EditingPosts}}
  <!-- View Mode -->
  <div style="display: flex;">
    <a href="/posts">Back</a>
    <button lvt-click="request_delete" lvt-data-id="{{.EditingID}}">Delete</button>
  </div>

  <!-- Detail Content -->
  <h2>Post Details</h2>

  <div style="max-width: 600px;">
    <div>
      <label>Title</label>
      <div>{{.EditingPosts.Title}}</div>
    </div>
  </div>
  {{end}}
{{end}}
`

func TestInjectEmbeddedChild_HandlerModifications(t *testing.T) {
	tmpDir := t.TempDir()

	parentGoPath := filepath.Join(tmpDir, "posts.go")
	if err := os.WriteFile(parentGoPath, []byte(parentGoTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	childData := newChildData()
	if err := InjectEmbeddedChild(parentGoPath, childData); err != nil {
		t.Fatalf("InjectEmbeddedChild failed: %v", err)
	}

	content, err := os.ReadFile(parentGoPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(content)

	// Verify import was added
	if !strings.Contains(src, `"testapp/app/comments"`) {
		t.Error("expected child import to be added")
	}

	// Verify state field was added
	if !strings.Contains(src, `Comments *comments.EmbeddedState`) {
		t.Error("expected Comments field to be added to PostsState")
	}

	// Verify controller field was added
	if !strings.Contains(src, `CommentsCtrl *comments.EmbeddedController`) {
		t.Error("expected CommentsCtrl field to be added to PostsController")
	}

	// Verify controller initialization
	if !strings.Contains(src, `controller.CommentsCtrl = comments.NewEmbeddedController(queries)`) {
		t.Error("expected child controller initialization")
	}

	// Verify initial state includes child state
	if !strings.Contains(src, `Comments: &comments.EmbeddedState{}`) {
		t.Error("expected child state initialization")
	}

	// Verify ParseFiles includes child template
	if !strings.Contains(src, `"app/comments/comments.tmpl"`) {
		t.Error("expected child template to be added to ParseFiles")
	}

	// Verify forwarding methods
	for _, method := range []string{"CommentAdd", "CommentEdit", "CommentUpdate", "CommentDelete", "CommentCancelEdit"} {
		if !strings.Contains(src, method) {
			t.Errorf("expected forwarding method %s to be generated", method)
		}
	}

	// Verify View method loads child data
	if !strings.Contains(src, `c.CommentsCtrl.Load(state.Comments`) {
		t.Error("expected child loading in View method")
	}
}

func TestInjectEmbeddedChild_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	parentGoPath := filepath.Join(tmpDir, "posts.go")
	if err := os.WriteFile(parentGoPath, []byte(parentGoTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	childData := newChildData()

	// Inject twice
	if err := InjectEmbeddedChild(parentGoPath, childData); err != nil {
		t.Fatalf("first injection failed: %v", err)
	}
	if err := InjectEmbeddedChild(parentGoPath, childData); err != nil {
		t.Fatalf("second injection failed: %v", err)
	}

	content, err := os.ReadFile(parentGoPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(content)

	// Verify no duplication
	count := strings.Count(src, `Comments *comments.EmbeddedState`)
	if count != 1 {
		t.Errorf("expected exactly 1 Comments field, got %d", count)
	}

	count = strings.Count(src, `CommentsCtrl *comments.EmbeddedController`)
	if count != 1 {
		t.Errorf("expected exactly 1 CommentsCtrl field, got %d", count)
	}
}

func TestInjectEmbeddedChildTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	parentTmplPath := filepath.Join(tmpDir, "posts.tmpl")
	if err := os.WriteFile(parentTmplPath, []byte(parentTmplTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	childData := newChildData()
	if err := InjectEmbeddedChildTemplate(parentTmplPath, childData); err != nil {
		t.Fatalf("InjectEmbeddedChildTemplate failed: %v", err)
	}

	content, err := os.ReadFile(parentTmplPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(content)

	// Verify child section template call was added
	if !strings.Contains(src, `{{template "comments:section" .Comments}}`) {
		t.Error("expected child section template call to be inserted")
	}
}

func TestInjectEmbeddedChildTemplate_Idempotent(t *testing.T) {
	tmpDir := t.TempDir()

	parentTmplPath := filepath.Join(tmpDir, "posts.tmpl")
	if err := os.WriteFile(parentTmplPath, []byte(parentTmplTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	childData := newChildData()

	// Inject twice
	if err := InjectEmbeddedChildTemplate(parentTmplPath, childData); err != nil {
		t.Fatalf("first injection failed: %v", err)
	}
	if err := InjectEmbeddedChildTemplate(parentTmplPath, childData); err != nil {
		t.Fatalf("second injection failed: %v", err)
	}

	content, err := os.ReadFile(parentTmplPath)
	if err != nil {
		t.Fatal(err)
	}
	src := string(content)

	count := strings.Count(src, `{{template "comments:section" .Comments}}`)
	if count != 1 {
		t.Errorf("expected exactly 1 child section call, got %d", count)
	}
}

func TestInjectEmbeddedChild_MissingParentFile(t *testing.T) {
	childData := newChildData()
	err := InjectEmbeddedChild("/nonexistent/posts.go", childData)
	if err == nil {
		t.Error("expected error for missing parent file")
	}
}

func TestInjectEmbeddedChildTemplate_MissingFile(t *testing.T) {
	childData := newChildData()
	err := InjectEmbeddedChildTemplate("/nonexistent/posts.tmpl", childData)
	if err == nil {
		t.Error("expected error for missing template file")
	}
}

func TestFindClosingBrace(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		startIdx int
		want     int
	}{
		{
			name:     "simple struct",
			src:      "type Foo struct { bar int }",
			startIdx: 0,
			want:     26, // index of the closing }
		},
		{
			name:     "nested braces",
			src:      "func() { if true { x++ } }",
			startIdx: 0,
			want:     25, // index of the outer closing }
		},
		{
			name:     "no closing brace",
			src:      "type Foo struct {",
			startIdx: 0,
			want:     -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findClosingBrace(tt.src, tt.startIdx)
			if got != tt.want {
				t.Errorf("findClosingBrace() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNonReferenceFields(t *testing.T) {
	data := ResourceData{
		ParentReferenceField: "post_id",
		Fields: []FieldData{
			{Name: "post_id", GoType: "string", IsReference: true, ReferencedTable: "posts"},
			{Name: "author", GoType: "string"},
			{Name: "text", GoType: "string"},
		},
	}

	nonRef := data.NonReferenceFields()
	if len(nonRef) != 2 {
		t.Errorf("expected 2 non-reference fields, got %d", len(nonRef))
	}
	for _, f := range nonRef {
		if f.Name == "post_id" {
			t.Error("expected post_id to be excluded from non-reference fields")
		}
	}
}

func TestNonReferenceFields_NoParent(t *testing.T) {
	data := ResourceData{
		Fields: []FieldData{
			{Name: "name", GoType: "string"},
			{Name: "email", GoType: "string"},
		},
	}

	nonRef := data.NonReferenceFields()
	if len(nonRef) != 2 {
		t.Errorf("expected all 2 fields when no parent ref, got %d", len(nonRef))
	}
}
