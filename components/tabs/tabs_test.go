package tabs

import (
	"html/template"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	items := []Tab{
		{ID: "general", Label: "General"},
		{ID: "security", Label: "Security"},
		{ID: "notifications", Label: "Notifications"},
	}

	tabs := New("settings", items)

	if tabs.ID() != "settings" {
		t.Errorf("expected ID 'settings', got '%s'", tabs.ID())
	}
	if len(tabs.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(tabs.Items))
	}
	if tabs.ActiveID != "general" {
		t.Errorf("expected ActiveID 'general' (first tab), got '%s'", tabs.ActiveID)
	}
}

func TestNewWithOptions(t *testing.T) {
	items := []Tab{
		{ID: "general", Label: "General"},
		{ID: "security", Label: "Security"},
	}

	tabs := New("settings", items,
		WithActive("security"),
		WithStyled(false),
	)

	if tabs.ActiveID != "security" {
		t.Errorf("expected ActiveID 'security', got '%s'", tabs.ActiveID)
	}
	if tabs.IsStyled() {
		t.Error("expected IsStyled to be false")
	}
}

func TestNewWithEmptyItems(t *testing.T) {
	tabs := New("empty", nil)

	if tabs.ActiveID != "" {
		t.Errorf("expected empty ActiveID for empty items, got '%s'", tabs.ActiveID)
	}
}

func TestNewVertical(t *testing.T) {
	items := []Tab{
		{ID: "dashboard", Label: "Dashboard"},
	}

	tabs := NewVertical("nav", items)

	if tabs.ID() != "nav" {
		t.Errorf("expected ID 'nav', got '%s'", tabs.ID())
	}
}

func TestNewPills(t *testing.T) {
	items := []Tab{
		{ID: "all", Label: "All"},
	}

	tabs := NewPills("filter", items)

	if tabs.ID() != "filter" {
		t.Errorf("expected ID 'filter', got '%s'", tabs.ID())
	}
}

func TestTabs_SetActive(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C", Disabled: true},
	}

	tabs := New("test", items)

	tabs.SetActive("b")
	if tabs.ActiveID != "b" {
		t.Errorf("expected ActiveID 'b', got '%s'", tabs.ActiveID)
	}

	// Should not set disabled tab
	tabs.SetActive("c")
	if tabs.ActiveID != "b" {
		t.Error("SetActive should not activate disabled tab")
	}

	// Should not set non-existent tab
	tabs.SetActive("nonexistent")
	if tabs.ActiveID != "b" {
		t.Error("SetActive should not activate non-existent tab")
	}
}

func TestTabs_ActiveTab(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}

	tabs := New("test", items)

	active := tabs.ActiveTab()
	if active == nil || active.ID != "a" {
		t.Error("expected ActiveTab to return first tab")
	}

	tabs.SetActive("b")
	active = tabs.ActiveTab()
	if active == nil || active.ID != "b" {
		t.Error("expected ActiveTab to return 'b'")
	}
}

func TestTabs_ActiveTabWithEmptyItems(t *testing.T) {
	tabs := New("empty", nil)

	if tabs.ActiveTab() != nil {
		t.Error("expected ActiveTab to return nil for empty items")
	}
}

func TestTabs_IsActive(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}

	tabs := New("test", items, WithActive("a"))

	if !tabs.IsActive("a") {
		t.Error("expected 'a' to be active")
	}
	if tabs.IsActive("b") {
		t.Error("expected 'b' to not be active")
	}
}

func TestTabs_Next(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items)

	tabs.Next()
	if tabs.ActiveID != "b" {
		t.Errorf("expected 'b' after Next, got '%s'", tabs.ActiveID)
	}

	tabs.Next()
	if tabs.ActiveID != "c" {
		t.Errorf("expected 'c' after Next, got '%s'", tabs.ActiveID)
	}

	// Should wrap around
	tabs.Next()
	if tabs.ActiveID != "a" {
		t.Errorf("expected 'a' after wrap-around, got '%s'", tabs.ActiveID)
	}
}

func TestTabs_NextSkipsDisabled(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B", Disabled: true},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items)

	tabs.Next()
	if tabs.ActiveID != "c" {
		t.Errorf("expected 'c' (skipping disabled 'b'), got '%s'", tabs.ActiveID)
	}
}

func TestTabs_Previous(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items, WithActive("c"))

	tabs.Previous()
	if tabs.ActiveID != "b" {
		t.Errorf("expected 'b' after Previous, got '%s'", tabs.ActiveID)
	}

	tabs.Previous()
	if tabs.ActiveID != "a" {
		t.Errorf("expected 'a' after Previous, got '%s'", tabs.ActiveID)
	}

	// Should wrap around
	tabs.Previous()
	if tabs.ActiveID != "c" {
		t.Errorf("expected 'c' after wrap-around, got '%s'", tabs.ActiveID)
	}
}

func TestTabs_PreviousSkipsDisabled(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B", Disabled: true},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items, WithActive("c"))

	tabs.Previous()
	if tabs.ActiveID != "a" {
		t.Errorf("expected 'a' (skipping disabled 'b'), got '%s'", tabs.ActiveID)
	}
}

func TestTabs_AddTab(t *testing.T) {
	tabs := New("test", nil)

	tabs.AddTab(Tab{ID: "a", Label: "A"})

	if len(tabs.Items) != 1 {
		t.Error("expected 1 item after AddTab")
	}
	if tabs.ActiveID != "a" {
		t.Error("expected first added tab to become active")
	}

	tabs.AddTab(Tab{ID: "b", Label: "B"})

	if len(tabs.Items) != 2 {
		t.Error("expected 2 items after second AddTab")
	}
	if tabs.ActiveID != "a" {
		t.Error("expected ActiveID to remain 'a'")
	}
}

func TestTabs_RemoveTab(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items, WithActive("b"))

	tabs.RemoveTab("c")
	if len(tabs.Items) != 2 {
		t.Error("expected 2 items after RemoveTab")
	}
	if tabs.ActiveID != "b" {
		t.Error("expected ActiveID to remain 'b'")
	}

	// Remove active tab - should activate first tab
	tabs.RemoveTab("b")
	if len(tabs.Items) != 1 {
		t.Error("expected 1 item after removing active tab")
	}
	if tabs.ActiveID != "a" {
		t.Errorf("expected ActiveID to be 'a' after removing active tab, got '%s'", tabs.ActiveID)
	}
}

func TestTabs_RemoveNonExistent(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
	}

	tabs := New("test", items)
	tabs.RemoveTab("nonexistent")

	if len(tabs.Items) != 1 {
		t.Error("expected items to remain unchanged")
	}
}

func TestTabs_TabCount(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B"},
	}

	tabs := New("test", items)

	if tabs.TabCount() != 2 {
		t.Errorf("expected TabCount 2, got %d", tabs.TabCount())
	}
}

func TestTabs_EnabledTabCount(t *testing.T) {
	items := []Tab{
		{ID: "a", Label: "A"},
		{ID: "b", Label: "B", Disabled: true},
		{ID: "c", Label: "C"},
	}

	tabs := New("test", items)

	if tabs.EnabledTabCount() != 2 {
		t.Errorf("expected EnabledTabCount 2, got %d", tabs.EnabledTabCount())
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()

	if ts == nil {
		t.Fatal("expected non-nil TemplateSet")
	}
	if ts.Namespace != "tabs" {
		t.Errorf("expected namespace 'tabs', got '%s'", ts.Namespace)
	}
}

func TestTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Tab{
		{ID: "general", Label: "General"},
		{ID: "security", Label: "Security", Badge: "3"},
		{ID: "disabled", Label: "Disabled", Disabled: true},
	}

	t.Run("horizontal template", func(t *testing.T) {
		tabs := New("settings", items, WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:tabs:horizontal:v1", tabs)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `data-tabs="settings"`) {
			t.Error("expected data-tabs attribute")
		}
		if !strings.Contains(html, `role="tablist"`) {
			t.Error("expected role=tablist")
		}
		if !strings.Contains(html, `lvt-click="select_tab_settings"`) {
			t.Error("expected lvt-click attribute")
		}
		if !strings.Contains(html, `aria-selected="true"`) {
			t.Error("expected aria-selected=true for active tab")
		}
		if !strings.Contains(html, "3") {
			t.Error("expected badge text")
		}
	})

	t.Run("vertical template", func(t *testing.T) {
		tabs := NewVertical("nav", items, WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:tabs:vertical:v1", tabs)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, `aria-orientation="vertical"`) {
			t.Error("expected aria-orientation=vertical")
		}
	})

	t.Run("pills template", func(t *testing.T) {
		tabs := NewPills("filter", items, WithStyled(true))

		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "lvt:tabs:pills:v1", tabs)
		if err != nil {
			t.Fatalf("failed to execute template: %v", err)
		}

		html := buf.String()
		if !strings.Contains(html, "rounded-full") {
			t.Error("expected pill styling class")
		}
	})
}

func TestUnstyledTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Tab{
		{ID: "a", Label: "A"},
	}

	tabs := New("test", items, WithStyled(false))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:tabs:horizontal:v1", tabs)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Unstyled version should not have Tailwind classes
	if strings.Contains(html, "border-b-2") {
		t.Error("unstyled template should not have Tailwind classes")
	}
	// But should still have functional attributes
	if !strings.Contains(html, `role="tablist"`) {
		t.Error("unstyled template should have role=tablist")
	}
}

func TestTabWithIcon(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	items := []Tab{
		{ID: "home", Label: "Home", Icon: "🏠"},
	}

	tabs := New("test", items, WithStyled(true))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:tabs:horizontal:v1", tabs)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, "🏠") {
		t.Error("expected icon in output")
	}
}
