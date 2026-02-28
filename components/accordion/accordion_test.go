package accordion

import (
	"html/template"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	items := []Item{
		{ID: "q1", Title: "Question 1", Content: "Answer 1"},
		{ID: "q2", Title: "Question 2", Content: "Answer 2"},
	}

	acc := New("faq", items)

	if acc.ID() != "faq" {
		t.Errorf("expected ID 'faq', got '%s'", acc.ID())
	}
	if len(acc.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(acc.Items))
	}
	if !acc.AllowMultiple {
		t.Error("expected AllowMultiple to be true for New()")
	}
	if acc.OpenCount() != 0 {
		t.Error("expected no items open initially")
	}
}

func TestNewWithOptions(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "Content A"},
		{ID: "b", Title: "B", Content: "Content B"},
	}

	acc := New("test", items,
		WithOpen("a"),
		WithStyled(false),
	)

	if !acc.IsOpen("a") {
		t.Error("expected 'a' to be open")
	}
	if acc.IsStyled() {
		t.Error("expected IsStyled to be false")
	}
}

func TestNewSingle(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := NewSingle("test", items)

	if acc.AllowMultiple {
		t.Error("expected AllowMultiple to be false for NewSingle()")
	}
}

func TestAccordion_Toggle(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items)

	// Open item
	acc.Toggle("a")
	if !acc.IsOpen("a") {
		t.Error("expected 'a' to be open after toggle")
	}

	// Close item
	acc.Toggle("a")
	if acc.IsOpen("a") {
		t.Error("expected 'a' to be closed after second toggle")
	}
}

func TestAccordion_ToggleDisabled(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A", Disabled: true},
	}

	acc := New("test", items)
	acc.Toggle("a")

	if acc.IsOpen("a") {
		t.Error("toggle should not open disabled item")
	}
}

func TestAccordion_ToggleSingleMode(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := NewSingle("test", items)

	acc.Toggle("a")
	if !acc.IsOpen("a") {
		t.Error("expected 'a' to be open")
	}

	acc.Toggle("b")
	if acc.IsOpen("a") {
		t.Error("expected 'a' to be closed in single mode")
	}
	if !acc.IsOpen("b") {
		t.Error("expected 'b' to be open")
	}
}

func TestAccordion_Open(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
	}

	acc := New("test", items)
	acc.Open("a")

	if !acc.IsOpen("a") {
		t.Error("expected 'a' to be open")
	}
}

func TestAccordion_OpenDisabled(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A", Disabled: true},
	}

	acc := New("test", items)
	acc.Open("a")

	if acc.IsOpen("a") {
		t.Error("Open should not open disabled item")
	}
}

func TestAccordion_Close(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
	}

	acc := New("test", items, WithOpen("a"))
	acc.Close("a")

	if acc.IsOpen("a") {
		t.Error("expected 'a' to be closed")
	}
}

func TestAccordion_OpenAll(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B", Disabled: true},
		{ID: "c", Title: "C", Content: "C"},
	}

	acc := New("test", items)
	acc.OpenAll()

	if !acc.IsOpen("a") || !acc.IsOpen("c") {
		t.Error("expected enabled items to be open")
	}
	if acc.IsOpen("b") {
		t.Error("expected disabled item to not be open")
	}
}

func TestAccordion_CloseAll(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items, WithAllOpen())
	acc.CloseAll()

	if acc.OpenCount() != 0 {
		t.Error("expected all items to be closed")
	}
}

func TestAccordion_OpenCount(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items)

	if acc.OpenCount() != 0 {
		t.Error("expected OpenCount 0 initially")
	}

	acc.Open("a")
	if acc.OpenCount() != 1 {
		t.Error("expected OpenCount 1")
	}

	acc.Open("b")
	if acc.OpenCount() != 2 {
		t.Error("expected OpenCount 2")
	}
}

func TestAccordion_ItemCount(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items)

	if acc.ItemCount() != 2 {
		t.Errorf("expected ItemCount 2, got %d", acc.ItemCount())
	}
}

func TestAccordion_GetItem(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items)

	item := acc.GetItem("a")
	if item == nil || item.ID != "a" {
		t.Error("expected to get item 'a'")
	}

	item = acc.GetItem("nonexistent")
	if item != nil {
		t.Error("expected nil for nonexistent item")
	}
}

func TestAccordion_AddItem(t *testing.T) {
	acc := New("test", nil)

	acc.AddItem(Item{ID: "a", Title: "A", Content: "A"})

	if acc.ItemCount() != 1 {
		t.Error("expected 1 item after AddItem")
	}
}

func TestAccordion_RemoveItem(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items, WithOpen("a"))

	acc.RemoveItem("a")

	if acc.ItemCount() != 1 {
		t.Error("expected 1 item after RemoveItem")
	}
	if acc.IsOpen("a") {
		t.Error("expected removed item to no longer be open")
	}
}

func TestWithAllOpen(t *testing.T) {
	items := []Item{
		{ID: "a", Title: "A", Content: "A"},
		{ID: "b", Title: "B", Content: "B"},
	}

	acc := New("test", items, WithAllOpen())

	if !acc.IsOpen("a") || !acc.IsOpen("b") {
		t.Error("expected all items to be open with WithAllOpen()")
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()

	if ts == nil {
		t.Fatal("expected non-nil TemplateSet")
	}
	if ts.Namespace != "accordion" {
		t.Errorf("expected namespace 'accordion', got '%s'", ts.Namespace)
	}
}

func TestTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Item{
		{ID: "q1", Title: "Question 1", Content: "Answer 1"},
		{ID: "q2", Title: "Question 2", Content: "Answer 2", Disabled: true},
	}

	t.Run("default template", func(t *testing.T) {
		acc := New("faq", items, WithOpen("q1"), WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:accordion:default:v1", acc)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-accordion="faq"`) {
			t.Error("expected data-accordion attribute")
		}
		if !strings.Contains(html, `lvt-click="toggle_accordion_faq"`) {
			t.Error("expected lvt-click attribute")
		}
		if !strings.Contains(html, `aria-expanded="true"`) {
			t.Error("expected aria-expanded=true for open item")
		}
		if !strings.Contains(html, "Answer 1") {
			t.Error("expected open item content to be rendered")
		}
		if strings.Contains(html, "Answer 2") {
			t.Error("expected closed item content to not be rendered")
		}
	})

	t.Run("single template", func(t *testing.T) {
		acc := NewSingle("nav", items, WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:accordion:single:v1", acc)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-single-open="true"`) {
			t.Error("expected data-single-open attribute")
		}
	})
}

func TestUnstyledTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Item{
		{ID: "a", Title: "A", Content: "Content A"},
	}

	acc := New("test", items, WithStyled(false))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:accordion:default:v1", acc)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Unstyled version should not have Tailwind classes
	if strings.Contains(html, "divide-y") {
		t.Error("unstyled template should not have Tailwind classes")
	}
	// But should still have functional attributes
	if !strings.Contains(html, `aria-expanded`) {
		t.Error("unstyled template should have aria-expanded")
	}
}

func TestItemWithIcon(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Item{
		{ID: "help", Title: "Help", Content: "Help content", Icon: "?"},
	}

	acc := New("test", items, WithStyled(true))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:accordion:default:v1", acc)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "?") {
		t.Error("expected icon in output")
	}
}
